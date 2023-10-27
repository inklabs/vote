package electionrepository

import (
	"context"
	"sync"
)

type inMemoryElectionRepository struct {
	mux       sync.RWMutex
	elections map[string]Election
}

func NewInMemory() *inMemoryElectionRepository {
	return &inMemoryElectionRepository{
		elections: make(map[string]Election),
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
