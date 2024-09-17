package electionrepository

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/inklabs/cqrs"
	"go.opentelemetry.io/otel"

	"github.com/inklabs/vote/pkg/sleep"
)

const instrumentationName = "github.com/inklabs/vote/internal/electionrepository/in-memory"

var (
	tracer = otel.Tracer(instrumentationName)
)

type inMemoryElectionRepository struct {
	mux sync.RWMutex

	// elections key by electionID
	elections map[string]Election

	// proposals key by proposalID
	proposals map[string]Proposal

	// votes key by electionID
	votes map[string][]Vote
}

func NewInMemory() *inMemoryElectionRepository {
	return &inMemoryElectionRepository{
		elections: make(map[string]Election),
		proposals: make(map[string]Proposal),
		votes:     make(map[string][]Vote),
	}
}

func (r *inMemoryElectionRepository) SaveElection(ctx context.Context, election Election) error {
	_, span := tracer.Start(ctx, "db.save-election")
	defer span.End()

	r.mux.Lock()
	defer r.mux.Unlock()

	sleep.Rand(2 * time.Millisecond)

	r.elections[election.ElectionID] = election

	return nil
}

func (r *inMemoryElectionRepository) GetElection(ctx context.Context, electionID string) (Election, error) {
	_, span := tracer.Start(ctx, "db.get-election")
	defer span.End()

	r.mux.RLock()
	defer r.mux.RUnlock()

	sleep.Rand(1 * time.Millisecond)

	if election, ok := r.elections[electionID]; ok {
		return election, nil
	}

	return Election{}, NewErrElectionNotFound(electionID)
}

func (r *inMemoryElectionRepository) SaveProposal(ctx context.Context, proposal Proposal) error {
	_, span := tracer.Start(ctx, "db.save-election")
	defer span.End()

	r.mux.Lock()
	defer r.mux.Unlock()

	sleep.Rand(2 * time.Millisecond)

	if _, ok := r.elections[proposal.ElectionID]; !ok {
		return NewErrElectionNotFound(proposal.ElectionID)
	}

	r.proposals[proposal.ProposalID] = proposal

	return nil
}

func (r *inMemoryElectionRepository) GetProposal(ctx context.Context, electionID string) (Proposal, error) {
	_, span := tracer.Start(ctx, "db.get-proposal")
	defer span.End()

	r.mux.RLock()
	defer r.mux.RUnlock()

	sleep.Rand(1 * time.Millisecond)

	if proposal, ok := r.proposals[electionID]; ok {
		return proposal, nil
	}

	return Proposal{}, ErrProposalNotFound
}

func (r *inMemoryElectionRepository) GetProposals(ctx context.Context, electionID string) ([]Proposal, error) {
	_, span := tracer.Start(ctx, "db.get-proposals")
	defer span.End()

	r.mux.RLock()
	defer r.mux.RUnlock()

	sleep.Rand(1 * time.Millisecond)

	if _, ok := r.elections[electionID]; !ok {
		return nil, NewErrElectionNotFound(electionID)
	}

	var proposals []Proposal
	for _, proposal := range r.proposals {
		if proposal.ElectionID == electionID {
			proposals = append(proposals, proposal)
		}
	}

	return proposals, nil
}

func (r *inMemoryElectionRepository) SaveVote(ctx context.Context, vote Vote) error {
	_, span := tracer.Start(ctx, "db.save-vote")
	defer span.End()

	r.mux.Lock()
	defer r.mux.Unlock()

	sleep.Rand(2 * time.Millisecond)

	if _, ok := r.elections[vote.ElectionID]; !ok {
		return NewErrElectionNotFound(vote.ElectionID)
	}

	for _, proposalID := range vote.RankedProposalIDs {
		if proposal, ok := r.proposals[proposalID]; ok {
			if proposal.ElectionID != vote.ElectionID {
				return ErrInvalidElectionProposal
			}
		} else {
			return ErrProposalNotFound
		}
	}

	r.votes[vote.ElectionID] = append(r.votes[vote.ElectionID], vote)

	return nil
}

func (r *inMemoryElectionRepository) GetVotes(ctx context.Context, electionID string) ([]Vote, error) {
	_, span := tracer.Start(ctx, "db.get-votes")
	defer span.End()

	r.mux.RLock()
	defer r.mux.RUnlock()

	sleep.Rand(1 * time.Millisecond)

	if votes, ok := r.votes[electionID]; ok {
		return votes, nil
	}

	return nil, NewErrElectionNotFound(electionID)
}

func (r *inMemoryElectionRepository) ListOpenElections(ctx context.Context, page, itemsPerPage int, sortBy, sortDirection *string) ([]Election, error) {
	_, span := tracer.Start(ctx, "db.list-open-elections")
	defer span.End()

	r.mux.RLock()
	defer r.mux.RUnlock()

	sleep.Rand(2 * time.Millisecond)

	var openElections []Election

	for _, election := range r.elections {
		if !election.IsClosed {
			openElections = append(openElections, election)
		}
	}

	sortElections(openElections, sortBy, sortDirection)

	return pageEntity(openElections, page, itemsPerPage), nil
}

func (r *inMemoryElectionRepository) ListProposals(ctx context.Context, electionID string, page, itemsPerPage int) ([]Proposal, error) {
	_, span := tracer.Start(ctx, "db.list-proposals")
	defer span.End()

	r.mux.RLock()
	defer r.mux.RUnlock()

	sleep.Rand(2 * time.Millisecond)

	if _, ok := r.elections[electionID]; !ok {
		return nil, NewErrElectionNotFound(electionID)
	}

	var proposals []Proposal

	for _, proposal := range r.proposals {
		if proposal.ElectionID == electionID {
			proposals = append(proposals, proposal)
		}
	}

	sort.Slice(proposals, func(i, j int) bool {
		return proposals[i].ProposedAt < proposals[j].ProposedAt
	})

	return pageEntity(proposals, page, itemsPerPage), nil
}

func sortElections(elections []Election, by, direction *string) {
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
