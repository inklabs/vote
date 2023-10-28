package election

import (
	"context"

	"github.com/inklabs/vote/internal/electionrepository"
)

type GetProposalDetails struct {
	ProposalID string
}

type GetProposalDetailsResponse struct {
	ElectionID  string
	ProposalID  string
	OwnerUserID string
	Name        string
	Description string
	ProposedAt  int
}

type getProposalDetailsHandler struct {
	repository electionrepository.Repository
}

func NewGetProposalDetailsHandler(repository electionrepository.Repository) *getProposalDetailsHandler {
	return &getProposalDetailsHandler{
		repository: repository,
	}
}

func (h *getProposalDetailsHandler) On(ctx context.Context, query GetProposalDetails) (GetProposalDetailsResponse, error) {
	proposal, err := h.repository.GetProposal(ctx, query.ProposalID)
	if err != nil {
		return GetProposalDetailsResponse{}, err
	}

	return GetProposalDetailsResponse{
		ElectionID:  proposal.ElectionID,
		ProposalID:  proposal.ProposalID,
		OwnerUserID: proposal.OwnerUserID,
		Name:        proposal.Name,
		Description: proposal.Description,
		ProposedAt:  proposal.ProposedAt,
	}, nil
}
