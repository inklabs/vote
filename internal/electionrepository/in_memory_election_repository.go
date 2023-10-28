package electionrepository

import (
	"context"
	"sync"
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

func (r *inMemoryElectionRepository) SaveElection(_ context.Context, election Election) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.elections[election.ElectionID] = election

	return nil
}

func (r *inMemoryElectionRepository) GetElection(_ context.Context, electionID string) (Election, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	if election, ok := r.elections[electionID]; ok {
		return election, nil
	}

	return Election{}, ErrElectionNotFound
}

func (r *inMemoryElectionRepository) SaveProposal(_ context.Context, proposal Proposal) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, ok := r.elections[proposal.ElectionID]; !ok {
		return ErrElectionNotFound
	}

	r.proposals[proposal.ProposalID] = proposal

	return nil
}

func (r *inMemoryElectionRepository) GetProposal(_ context.Context, electionID string) (Proposal, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	if proposal, ok := r.proposals[electionID]; ok {
		return proposal, nil
	}

	return Proposal{}, ErrProposalNotFound
}

func (r *inMemoryElectionRepository) GetProposals(_ context.Context, electionID string) ([]Proposal, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	if _, ok := r.elections[electionID]; !ok {
		return nil, ErrElectionNotFound
	}

	var proposals []Proposal
	for _, proposal := range r.proposals {
		if proposal.ElectionID == electionID {
			proposals = append(proposals, proposal)
		}
	}

	return proposals, nil
}

func (r *inMemoryElectionRepository) SaveVote(_ context.Context, vote Vote) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, ok := r.elections[vote.ElectionID]; !ok {
		return ErrElectionNotFound
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

func (r *inMemoryElectionRepository) GetVotes(_ context.Context, electionID string) ([]Vote, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	if votes, ok := r.votes[electionID]; ok {
		return votes, nil
	}

	return nil, ErrElectionNotFound
}
