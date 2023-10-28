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
	ProposedAt  int
}

type Vote struct {
	ElectionID        string
	UserID            string
	RankedProposalIDs []string
}

type Repository interface {
	SaveElection(ctx context.Context, election Election) error
	GetElection(ctx context.Context, electionID string) (Election, error)
	SaveProposal(ctx context.Context, proposal Proposal) error
	GetProposal(ctx context.Context, proposalID string) (Proposal, error)
	GetProposals(ctx context.Context, electionID string) ([]Proposal, error)
	SaveVote(ctx context.Context, vote Vote) error
	GetVotes(ctx context.Context, electionID string) ([]Vote, error)
}

var ErrElectionNotFound = fmt.Errorf("election not found")
var ErrProposalNotFound = fmt.Errorf("proposal not found")
var ErrInvalidElectionProposal = fmt.Errorf("invalid election")
