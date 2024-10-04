package election

import (
	"context"

	"github.com/inklabs/vote/internal/electionrepository"
)

type GetElection struct {
	ElectionID string
}

type GetElectionResponse struct {
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

type getElectionHandler struct {
	repository electionrepository.Repository
}

func NewGetElectionHandler(repository electionrepository.Repository) *getElectionHandler {
	return &getElectionHandler{
		repository: repository,
	}
}

func (h *getElectionHandler) On(ctx context.Context, query GetElection) (GetElectionResponse, error) {
	election, err := h.repository.GetElection(ctx, query.ElectionID)
	if err != nil {
		return GetElectionResponse{}, err
	}

	return GetElectionResponse{
		ElectionID:        election.ElectionID,
		OrganizerUserID:   election.OrganizerUserID,
		Name:              election.Name,
		Description:       election.Description,
		WinningProposalID: election.WinningProposalID,
		IsClosed:          election.IsClosed,
		CommencedAt:       election.CommencedAt,
		ClosedAt:          election.ClosedAt,
		SelectedAt:        election.SelectedAt,
	}, nil
}
