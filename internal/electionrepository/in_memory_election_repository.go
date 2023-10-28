package electionrepository

import (
	"context"
	"sync"
)

type inMemoryElectionRepository struct {
	mux sync.RWMutex

	// elections key by electionID
	elections map[string]Election

	// proposals key by electionID
	proposals map[string][]Proposal
}

func NewInMemory() *inMemoryElectionRepository {
	return &inMemoryElectionRepository{
		elections: make(map[string]Election),
		proposals: make(map[string][]Proposal),
	}
}

func (i *inMemoryElectionRepository) SaveElection(_ context.Context, election Election) error {
	i.mux.Lock()
	defer i.mux.Unlock()

	i.elections[election.ElectionID] = election

	return nil
}

func (i *inMemoryElectionRepository) GetElection(_ context.Context, electionID string) (Election, error) {
	i.mux.RLock()
	defer i.mux.RUnlock()

	if election, ok := i.elections[electionID]; ok {
		return election, nil
	}

	return Election{}, ErrElectionNotFound
}

func (i *inMemoryElectionRepository) SaveProposal(ctx context.Context, proposal Proposal) error {
	i.mux.Lock()
	defer i.mux.Unlock()

	if _, ok := i.elections[proposal.ElectionID]; !ok {
		return ErrElectionNotFound
	}

	i.proposals[proposal.ElectionID] = append(i.proposals[proposal.ElectionID], proposal)

	return nil
}

func (i *inMemoryElectionRepository) GetProposals(_ context.Context, electionID string) ([]Proposal, error) {
	i.mux.RLock()
	defer i.mux.RUnlock()

	if proposals, ok := i.proposals[electionID]; ok {
		return proposals, nil
	}

	return nil, ErrElectionNotFound
}
