package electionrepository

import (
	"context"
	"fmt"
)

const DefaultItemsPerPage = 10

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
	VoteID            string
	ElectionID        string
	UserID            string
	RankedProposalIDs []string
	SubmittedAt       int
}

type Repository interface {
	SaveElection(ctx context.Context, election Election) error
	GetElection(ctx context.Context, electionID string) (Election, error)
	SaveProposal(ctx context.Context, proposal Proposal) error
	GetProposal(ctx context.Context, proposalID string) (Proposal, error)
	SaveVote(ctx context.Context, vote Vote) error
	GetVotes(ctx context.Context, electionID string) ([]Vote, error)
	ListOpenElections(ctx context.Context, page, itemsPerPage int, sortBy, sortDirection *string) (int, []Election, error)
	ListProposals(ctx context.Context, electionID string, page, itemsPerPage int) (int, []Proposal, error)
}

type ErrElectionNotFound struct {
	electionID string
}

func NewErrElectionNotFound(electionID string) *ErrElectionNotFound {
	return &ErrElectionNotFound{electionID: electionID}
}

func (e ErrElectionNotFound) Error() string {
	return fmt.Sprintf("election (%s) not found", e.electionID)
}

type ErrProposalNotFound struct {
	proposalID string
}

func NewErrProposalNotFound(proposalID string) *ErrProposalNotFound {
	return &ErrProposalNotFound{proposalID: proposalID}
}

func (e ErrProposalNotFound) Error() string {
	return fmt.Sprintf("proposal (%s) not found", e.proposalID)
}

type ErrInvalidElectionProposal struct {
	proposalID string
	electionID string
}

func NewErrInvalidElectionProposal(proposalID, electionID string) *ErrInvalidElectionProposal {
	return &ErrInvalidElectionProposal{
		proposalID: proposalID,
		electionID: electionID,
	}
}

func (e ErrInvalidElectionProposal) Error() string {
	return fmt.Sprintf("invalid proposal (%s) for wrong election (%s)", e.proposalID, e.electionID)
}
