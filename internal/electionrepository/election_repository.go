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
	WinningProposalID string
	IsClosed          bool
	CommencedAt       int
	ClosedAt          int
	SelectedAt        int
}

type Repository interface {
	SaveElection(ctx context.Context, election Election) error
	GetElection(ctx context.Context, electionID string) (Election, error)
}

var ErrElectionNotFound = fmt.Errorf("election not found")
