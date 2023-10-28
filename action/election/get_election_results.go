package election

import (
	"context"

	"github.com/inklabs/vote/internal/electionrepository"
)

type GetElectionResults struct {
	ElectionID string
}

type GetElectionResultsResponse struct {
	ElectionID        string
	WinningProposalID string
	SelectedAt        int
}

type getElectionResultsHandler struct {
	repository electionrepository.Repository
}

func NewGetElectionResultsHandler(repository electionrepository.Repository) *getElectionResultsHandler {
	return &getElectionResultsHandler{
		repository: repository,
	}
}

func (h *getElectionResultsHandler) On(ctx context.Context, query GetElectionResults) (GetElectionResultsResponse, error) {
	election, err := h.repository.GetElection(ctx, query.ElectionID)
	if err != nil {
		return GetElectionResultsResponse{}, err
	}

	return GetElectionResultsResponse{
		ElectionID:        election.ElectionID,
		WinningProposalID: election.WinningProposalID,
		SelectedAt:        election.SelectedAt,
	}, nil
}
