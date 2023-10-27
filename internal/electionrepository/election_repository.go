package electionrepository

import (
	"context"
	"fmt"
)

type Election struct {
	ElectionID        string
	OrganizerUserID   string
	Name              string
	Description       string
	OccurredAt        int
	WinningProposalID string
	IsClosed          bool
	ClosedAt          int
}

type Repository interface {
	SaveElection(ctx context.Context, election Election) error
	GetElection(ctx context.Context, electionID string) (Election, error)
}

var ErrElectionNotFound = fmt.Errorf("election not found")
