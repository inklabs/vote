package postgresrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/lib/pq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/inklabs/vote/internal/electionrepository"
)

const instrumentationName = "github.com/inklabs/vote/internal/electionrepository/postgres"

var tracer = otel.Tracer(instrumentationName)

type postgresRepository struct {
	db *sql.DB
}

func NewFromDB(db *sql.DB) (*postgresRepository, error) {
	r := &postgresRepository{
		db: db,
	}

	return r, nil
}

func NewFromConfig(config Config) (*postgresRepository, error) {
	db, err := NewDB(config)
	if err != nil {
		return nil, err
	}

	return NewFromDB(db)
}

func (r *postgresRepository) SaveElection(ctx context.Context, election electionrepository.Election) error {
	_, span := tracer.Start(ctx, "db.save-election")
	defer span.End()

	sqlStatement := `INSERT INTO election (
						ElectionID,
						OrganizerUserID,
						Name,
						Description,
						WinningProposalID,
						IsClosed,
						CommencedAt,
						ClosedAt,
						SelectedAt
                     ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
                     ON CONFLICT (ElectionID)
					 DO UPDATE SET
					     Name = EXCLUDED.Name,
					     Description = EXCLUDED.Description,
					     WinningProposalID = EXCLUDED.WinningProposalID,
					     IsClosed = EXCLUDED.IsClosed,
					     ClosedAt = EXCLUDED.ClosedAt,
					     SelectedAt = EXCLUDED.SelectedAt`

	_, err := r.db.ExecContext(ctx, sqlStatement,
		election.ElectionID,
		election.OrganizerUserID,
		election.Name,
		election.Description,
		election.WinningProposalID,
		election.IsClosed,
		election.CommencedAt,
		election.ClosedAt,
		election.SelectedAt,
	)
	if err != nil {
		recordSpanError(span, err)
		return fmt.Errorf("unable to save election: %w", err)
	}

	return nil
}

func (r *postgresRepository) GetElection(ctx context.Context, electionID string) (electionrepository.Election, error) {
	_, span := tracer.Start(ctx, "db.get-election")
	defer span.End()

	sqlStatement := `SELECT
						ElectionID,
						OrganizerUserID,
						Name,
						Description,
						WinningProposalID,
						IsClosed,
						CommencedAt,
						ClosedAt,
						SelectedAt
                     FROM election
                     WHERE ElectionID = $1`

	var election electionrepository.Election

	row := r.db.QueryRowContext(ctx, sqlStatement, electionID)
	if row.Err() != nil {
		err := fmt.Errorf("unable to get election: %w", row.Err())
		recordSpanError(span, err)
		return election, err
	}

	err := row.Scan(
		&election.ElectionID,
		&election.OrganizerUserID,
		&election.Name,
		&election.Description,
		&election.WinningProposalID,
		&election.IsClosed,
		&election.CommencedAt,
		&election.ClosedAt,
		&election.SelectedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return election, electionrepository.NewErrElectionNotFound(electionID)
		}

		err = fmt.Errorf("unable to get election data: %w", err)
		recordSpanError(span, err)
		return election, err
	}

	return election, nil
}

func (r *postgresRepository) SaveProposal(ctx context.Context, proposal electionrepository.Proposal) error {
	_, span := tracer.Start(ctx, "db.save-election")
	defer span.End()

	sqlStatement := `INSERT INTO proposal (
                      	ProposalID,
						ElectionID,
						OwnerUserID,
						Name,
						Description,
						ProposedAt
                     ) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, sqlStatement,
		proposal.ProposalID,
		proposal.ElectionID,
		proposal.OwnerUserID,
		proposal.Name,
		proposal.Description,
		proposal.ProposedAt,
	)
	if err != nil {
		recordSpanError(span, err)
		return fmt.Errorf("unable to save proposal: %w", err)
	}

	return nil
}

func (r *postgresRepository) GetProposal(ctx context.Context, proposalID string) (electionrepository.Proposal, error) {
	_, span := tracer.Start(ctx, "db.get-proposal")
	defer span.End()

	sqlStatement := `SELECT
						ProposalID,
						ElectionID,
						OwnerUserID,
						Name,
						Description,
						ProposedAt
                     FROM proposal
                     WHERE ProposalID = $1`

	var proposal electionrepository.Proposal

	row := r.db.QueryRowContext(ctx, sqlStatement, proposalID)
	if row.Err() != nil {
		err := fmt.Errorf("unable to get proposal: %w", row.Err())
		recordSpanError(span, err)
		return proposal, err
	}

	err := row.Scan(
		&proposal.ProposalID,
		&proposal.ElectionID,
		&proposal.OwnerUserID,
		&proposal.Name,
		&proposal.Description,
		&proposal.ProposedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return proposal, electionrepository.NewErrProposalNotFound(proposalID)
		}
		err = fmt.Errorf("unable to get proposal data: %w", err)
		recordSpanError(span, err)
		return proposal, err
	}

	return proposal, nil
}

func (r *postgresRepository) SaveVote(ctx context.Context, vote electionrepository.Vote) error {
	_, span := tracer.Start(ctx, "db.save-vote")
	defer span.End()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("unable to create transaction: %w", err)
		recordSpanError(span, err)
		return err
	}

	sqlStatement := `INSERT INTO vote (
                      	VoteID,
						ElectionID,
						UserID,
						SubmittedAt
                     ) VALUES ($1, $2, $3, $4)`

	_, err = tx.ExecContext(ctx, sqlStatement,
		vote.VoteID,
		vote.ElectionID,
		vote.UserID,
		vote.SubmittedAt,
	)
	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Code == "23503" && pqError.Constraint == "vote_electionid_fkey" {
				err = electionrepository.NewErrElectionNotFound(vote.ElectionID)
				recordSpanError(span, err)
				_ = tx.Rollback()
				return err
			}
		}
		recordSpanError(span, err)
		_ = tx.Rollback()
		return fmt.Errorf("unable to save vote: %w", err)
	}

	err = r.saveRankedProposals(ctx, tx, vote)
	if err != nil {
		recordSpanError(span, err)
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("unable to commit transaction: %w", err)
		recordSpanError(span, err)
		return err
	}

	return nil
}

func (r *postgresRepository) saveRankedProposals(ctx context.Context, tx *sql.Tx, vote electionrepository.Vote) error {
	var valueStrings []string
	var valueArgs []interface{}

	i := 0
	for _, proposalID := range vote.RankedProposalIDs {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs,
			vote.VoteID,
			proposalID,
			vote.ElectionID,
			i,
		)
		i++
	}

	sqlStatement := fmt.Sprintf(
		"INSERT INTO vote_ranked_proposal (VoteID, ProposalID, ElectionID, Position) VALUES %s",
		strings.Join(valueStrings, ","))

	_, err := tx.ExecContext(ctx, sqlStatement, valueArgs...)
	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Code == "23503" && pqError.Constraint == "vote_ranked_proposal_proposalid_fkey" {
				proposalID := uuidRegex.FindString(pqError.Detail)
				err = electionrepository.NewErrProposalNotFound(proposalID)
				return err
			}
			if pqError.Code == "23503" && pqError.Constraint == "vote_ranked_proposal_proposalid_electionid_fkey" {
				uuids := uuidRegex.FindAllString(pqError.Detail, -1)
				var proposalID, electionID string
				if len(uuids) == 2 {
					proposalID = uuids[0]
					electionID = uuids[1]
				}
				err = electionrepository.NewErrInvalidElectionProposal(proposalID, electionID)
				return err
			}
		}
		return fmt.Errorf("unable to save ranked proposals: %w", err)
	}

	return nil
}

func (r *postgresRepository) GetVotes(ctx context.Context, electionID string) ([]electionrepository.Vote, error) {
	_, span := tracer.Start(ctx, "db.get-votes")
	defer span.End()

	sqlStatement := `SELECT
						v.VoteID,
						v.ElectionID,
						v.UserID,
						ARRAY_REMOVE(ARRAY_AGG(vrp.ProposalID ORDER BY vrp.Position), NULL),
						v.SubmittedAt
                     FROM vote AS v
                     LEFT JOIN vote_ranked_proposal AS vrp ON vrp.VoteID = v.VoteID
                     WHERE v.ElectionID = $1
                     GROUP BY v.VoteID`

	rows, err := r.db.QueryContext(ctx, sqlStatement, electionID)
	if err != nil {
		err = fmt.Errorf("unable to get votes: %w", err)
		recordSpanError(span, err)
		return nil, err
	}

	var votes []electionrepository.Vote

	for rows.Next() {
		var vote electionrepository.Vote

		err = rows.Scan(
			&vote.VoteID,
			&vote.ElectionID,
			&vote.UserID,
			pq.Array(&vote.RankedProposalIDs),
			&vote.SubmittedAt,
		)
		if err != nil {
			err = fmt.Errorf("unable to get vote data: %w", err)
			recordSpanError(span, err)
			return nil, err
		}

		votes = append(votes, vote)
	}

	if rows.Err() != nil {
		err = fmt.Errorf("unable to get votes: %w", rows.Err())
		recordSpanError(span, err)
		return nil, err
	}

	return votes, nil
}

func (r *postgresRepository) ListOpenElections(ctx context.Context, page, itemsPerPage int, sortBy, sortDirection *string) (int, []electionrepository.Election, error) {
	_, span := tracer.Start(ctx, "db.list-open-elections")
	defer span.End()

	orderBy := getOrderBy(sortBy, sortDirection, "CommencedAt", "ASC")
	limit, offset := getLimitOffset(page, itemsPerPage)

	sqlStatement := `SELECT
						ElectionID,
						OrganizerUserID,
						Name,
						Description,
						WinningProposalID,
						IsClosed,
						CommencedAt,
						ClosedAt,
						SelectedAt,
						count(*) OVER()
                     FROM election
					 WHERE IsClosed = FALSE
                     ` + orderBy + `
                     LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, sqlStatement, limit, offset)
	if err != nil {
		err = fmt.Errorf("unable to list open elections: %w", err)
		recordSpanError(span, err)
		return 0, nil, err
	}

	var elections []electionrepository.Election
	var totalResults int

	for rows.Next() {
		var election electionrepository.Election

		err = rows.Scan(
			&election.ElectionID,
			&election.OrganizerUserID,
			&election.Name,
			&election.Description,
			&election.WinningProposalID,
			&election.IsClosed,
			&election.CommencedAt,
			&election.ClosedAt,
			&election.SelectedAt,
			&totalResults,
		)
		if err != nil {
			err = fmt.Errorf("unable to get election data: %w", err)
			recordSpanError(span, err)
			return 0, nil, err
		}

		elections = append(elections, election)
	}

	if rows.Err() != nil {
		err = fmt.Errorf("unable to get elections: %w", rows.Err())
		recordSpanError(span, err)
		return 0, nil, err
	}

	return totalResults, elections, nil
}

func (r *postgresRepository) ListProposals(ctx context.Context, electionID string, page, itemsPerPage int) (int, []electionrepository.Proposal, error) {
	_, span := tracer.Start(ctx, "db.list-proposals")
	defer span.End()

	limit, offset := getLimitOffset(page, itemsPerPage)

	sqlStatement := `SELECT
						ProposalID,
						ElectionID,
						OwnerUserID,
						Name,
						Description,
						ProposedAt,
						count(*) OVER()
                     FROM proposal
                     WHERE electionID = $1
                     ORDER BY ProposedAt ASC
                     LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, sqlStatement, electionID, limit, offset)
	if err != nil {
		err = fmt.Errorf("unable to list proposals: %w", err)
		recordSpanError(span, err)
		return 0, nil, err
	}

	var proposals []electionrepository.Proposal
	var totalResults int

	for rows.Next() {
		var proposal electionrepository.Proposal

		err = rows.Scan(
			&proposal.ProposalID,
			&proposal.ElectionID,
			&proposal.OwnerUserID,
			&proposal.Name,
			&proposal.Description,
			&proposal.ProposedAt,
			&totalResults,
		)
		if err != nil {
			err = fmt.Errorf("unable to get election data: %w", err)
			recordSpanError(span, err)
			return 0, nil, err
		}

		proposals = append(proposals, proposal)
	}

	if rows.Err() != nil {
		err = fmt.Errorf("unable to get proposals: %w", rows.Err())
		recordSpanError(span, err)
		return 0, nil, err
	}

	return totalResults, proposals, nil
}

func NewDB(config Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DataSourceName())
	if err != nil {
		return nil, fmt.Errorf("unable to open DB connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to DB: %w", err)
	}

	return db, nil
}

func (r *postgresRepository) InitDB(ctx context.Context) error {
	sqlStatements := []string{
		`CREATE TABLE IF NOT EXISTS election (
			ElectionID TEXT PRIMARY KEY,
			OrganizerUserID TEXT,
            Name TEXT,
            Description TEXT,
            WinningProposalID TEXT,
            IsClosed BOOLEAN,
            CommencedAt BIGINT,
            ClosedAt BIGINT,
            SelectedAt BIGINT
		);`,
		`CREATE TABLE IF NOT EXISTS proposal (
			ProposalID TEXT PRIMARY KEY,
			ElectionID TEXT REFERENCES election (ElectionID),
			OwnerUserID TEXT,
            Name TEXT,
            Description TEXT,
            ProposedAt BIGINT,
    		CONSTRAINT unique_proposal_election UNIQUE (ProposalID, ElectionID)
		);`,
		`CREATE TABLE IF NOT EXISTS vote (
			VoteID TEXT PRIMARY KEY,
			ElectionID TEXT REFERENCES election (ElectionID),
			UserID TEXT,
    		SubmittedAt BIGINT,
    		CONSTRAINT unique_vote_election UNIQUE (VoteID, ElectionID)
		);`,
		`CREATE TABLE IF NOT EXISTS vote_ranked_proposal (
			VoteID TEXT REFERENCES vote (VoteID),
			ProposalID TEXT REFERENCES proposal (ProposalID),
    		ElectionID TEXT REFERENCES election (ElectionID),
			Position SMALLINT,
    		PRIMARY KEY (VoteID, ProposalID),
	        FOREIGN KEY (VoteID, ElectionID) REFERENCES vote (VoteID, ElectionID),
		    FOREIGN KEY (ProposalID, ElectionID) REFERENCES proposal (ProposalID, ElectionID)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_proposal_election_id ON proposal(ElectionID);`,
		`CREATE INDEX IF NOT EXISTS idx_vote_election_id ON vote(ElectionID);`,
	}

	for _, statement := range sqlStatements {
		_, err := r.db.ExecContext(ctx, statement)
		if err != nil {
			return err
		}
	}

	return nil
}

func getLimitOffset(page int, itemsPerPage int) (int, int) {
	offset := (itemsPerPage * page) - itemsPerPage
	return itemsPerPage, offset
}

func getOrderBy(sortBy, sortDirection *string, defaultSort, defaultDirection string) string {
	validDirection := defaultDirection == "ASC" || defaultDirection == "DESC"
	if !validDirection {
		return ""
	}

	if sortBy == nil {
		return fmt.Sprintf("ORDER BY %s %s", defaultSort, defaultDirection)
	}

	direction := defaultDirection
	if *sortDirection == "ascending" {
		direction = "ASC"
	} else if *sortDirection == "descending" {
		direction = "DESC"
	}

	return fmt.Sprintf("ORDER BY %s %s", *sortBy, direction)
}

func recordSpanError(span trace.Span, err error) {
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
}

var uuidRegex = regexp.MustCompile(`[a-f0-9\-]{36}`)
