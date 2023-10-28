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

type Proposal struct {
	ElectionID  string
	ProposalID  string
	OwnerUserID string
	Name        string
	Description string
	OccurredAt  int
}

type Repository interface {
	SaveElection(ctx context.Context, election Election) error
	GetElection(ctx context.Context, electionID string) (Election, error)
	SaveProposal(ctx context.Context, proposal Proposal) error
	GetProposals(ctx context.Context, electionID string) ([]Proposal, error)
}

var ErrElectionNotFound = fmt.Errorf("election not found")
