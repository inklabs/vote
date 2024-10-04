package inmemoryrepo

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/inklabs/cqrs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/inklabs/vote/internal/electionrepository"
	"github.com/inklabs/vote/pkg/sleep"
)

const instrumentationName = "github.com/inklabs/vote/internal/electionrepository/in-memory"

var tracer = otel.Tracer(instrumentationName)

type inMemoryElectionRepository struct {
	mux sync.RWMutex

	// elections key by electionID
	elections map[string]electionrepository.Election

	// proposals key by proposalID
	proposals map[string]electionrepository.Proposal

	// votes key by electionID
	votes map[string][]electionrepository.Vote
}

func New() *inMemoryElectionRepository {
	return &inMemoryElectionRepository{
		elections: make(map[string]electionrepository.Election),
		proposals: make(map[string]electionrepository.Proposal),
		votes:     make(map[string][]electionrepository.Vote),
	}
}

func (r *inMemoryElectionRepository) SaveElection(ctx context.Context, election electionrepository.Election) error {
	_, span := tracer.Start(ctx, "db.save-election")
	defer span.End()

	r.mux.Lock()
	defer r.mux.Unlock()

	sleep.Rand(2 * time.Millisecond)

	r.elections[election.ElectionID] = election

	return nil
}

func (r *inMemoryElectionRepository) GetElection(ctx context.Context, electionID string) (electionrepository.Election, error) {
	_, span := tracer.Start(ctx, "db.get-election")
	defer span.End()

	r.mux.RLock()
	defer r.mux.RUnlock()

	sleep.Rand(1 * time.Millisecond)

	if election, ok := r.elections[electionID]; ok {
		return election, nil
	}

	err := electionrepository.NewErrElectionNotFound(electionID)
	recordSpanError(span, err)

	return electionrepository.Election{}, err
}

func (r *inMemoryElectionRepository) SaveProposal(ctx context.Context, proposal electionrepository.Proposal) error {
	_, span := tracer.Start(ctx, "db.save-election")
	defer span.End()

	r.mux.Lock()
	defer r.mux.Unlock()

	sleep.Rand(2 * time.Millisecond)

	if _, ok := r.elections[proposal.ElectionID]; !ok {
		err := electionrepository.NewErrElectionNotFound(proposal.ElectionID)
		recordSpanError(span, err)

		return err
	}

	r.proposals[proposal.ProposalID] = proposal

	return nil
}

func (r *inMemoryElectionRepository) GetProposal(ctx context.Context, proposalID string) (electionrepository.Proposal, error) {
	_, span := tracer.Start(ctx, "db.get-proposal")
	defer span.End()

	r.mux.RLock()
	defer r.mux.RUnlock()

	sleep.Rand(1 * time.Millisecond)

	if proposal, ok := r.proposals[proposalID]; ok {
		return proposal, nil
	}

	err := electionrepository.NewErrProposalNotFound(proposalID)
	recordSpanError(span, err)

	return electionrepository.Proposal{}, err
}

func (r *inMemoryElectionRepository) SaveVote(ctx context.Context, vote electionrepository.Vote) error {
	_, span := tracer.Start(ctx, "db.save-vote")
	defer span.End()

	r.mux.Lock()
	defer r.mux.Unlock()

	sleep.Rand(2 * time.Millisecond)

	if _, ok := r.elections[vote.ElectionID]; !ok {
		err := electionrepository.NewErrElectionNotFound(vote.ElectionID)
		recordSpanError(span, err)

		return err
	}

	for _, proposalID := range vote.RankedProposalIDs {
		if proposal, ok := r.proposals[proposalID]; ok {
			if proposal.ElectionID != vote.ElectionID {
				err := electionrepository.NewErrInvalidElectionProposal(proposal.ProposalID, vote.ElectionID)
				recordSpanError(span, err)

				return err
			}
		} else {
			err := electionrepository.NewErrProposalNotFound(proposalID)
			recordSpanError(span, err)

			return err
		}
	}

	r.votes[vote.ElectionID] = append(r.votes[vote.ElectionID], vote)

	return nil
}

func (r *inMemoryElectionRepository) GetVotes(ctx context.Context, electionID string) ([]electionrepository.Vote, error) {
	_, span := tracer.Start(ctx, "db.get-votes")
	defer span.End()

	r.mux.RLock()
	defer r.mux.RUnlock()

	sleep.Rand(1 * time.Millisecond)

	if votes, ok := r.votes[electionID]; ok {
		return votes, nil
	}

	err := electionrepository.NewErrElectionNotFound(electionID)
	recordSpanError(span, err)

	return nil, err
}

func (r *inMemoryElectionRepository) ListOpenElections(ctx context.Context, page, itemsPerPage int, sortBy, sortDirection *string) (int, []electionrepository.Election, error) {
	_, span := tracer.Start(ctx, "db.list-open-elections")
	defer span.End()

	r.mux.RLock()
	defer r.mux.RUnlock()

	sleep.Rand(2 * time.Millisecond)

	var openElections []electionrepository.Election

	for _, election := range r.elections {
		if !election.IsClosed {
			openElections = append(openElections, election)
		}
	}

	sortElections(openElections, sortBy, sortDirection)

	totalResults := len(openElections)
	return totalResults, pageEntity(openElections, page, itemsPerPage), nil
}

func (r *inMemoryElectionRepository) ListProposals(ctx context.Context, electionID string, page, itemsPerPage int) (int, []electionrepository.Proposal, error) {
	_, span := tracer.Start(ctx, "db.list-proposals")
	defer span.End()

	r.mux.RLock()
	defer r.mux.RUnlock()

	sleep.Rand(2 * time.Millisecond)

	if _, ok := r.elections[electionID]; !ok {
		err := electionrepository.NewErrElectionNotFound(electionID)
		recordSpanError(span, err)

		return 0, nil, err
	}

	var proposals []electionrepository.Proposal

	for _, proposal := range r.proposals {
		if proposal.ElectionID == electionID {
			proposals = append(proposals, proposal)
		}
	}

	sort.Slice(proposals, func(i, j int) bool {
		return proposals[i].ProposedAt < proposals[j].ProposedAt
	})

	totalResults := len(proposals)
	return totalResults, pageEntity(proposals, page, itemsPerPage), nil
}

func sortElections(elections []electionrepository.Election, by, direction *string) {
	sortBy, sortDirection := cqrs.DefaultSort(by, direction, "CommencedAt", "ascending")

	var sortFunction func(i, j int) bool

	switch sortBy {
	case "Name":
		if sortDirection == "ascending" {
			sortFunction = func(i, j int) bool {
				return elections[i].Name < elections[j].Name
			}
		} else {
			sortFunction = func(i, j int) bool {
				return elections[i].Name > elections[j].Name
			}
		}
	case "CommencedAt":
		if sortDirection == "ascending" {
			sortFunction = func(i, j int) bool {
				return elections[i].CommencedAt < elections[j].CommencedAt
			}
		} else {
			sortFunction = func(i, j int) bool {
				return elections[i].CommencedAt > elections[j].CommencedAt
			}
		}
	}

	sort.Slice(elections, sortFunction)
}

func pageEntity[T any](entities []T, page, itemsPerPage int) []T {
	startIndex := (page - 1) * itemsPerPage
	endIndex := startIndex + itemsPerPage

	if startIndex >= len(entities) {
		return nil
	}

	if endIndex > len(entities) {
		endIndex = len(entities)
	}

	return entities[startIndex:endIndex]
}

func recordSpanError(span trace.Span, err error) {
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
}
