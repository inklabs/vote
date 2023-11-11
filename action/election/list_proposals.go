package election

import (
	"context"

	"github.com/inklabs/cqrs"

	"github.com/inklabs/vote/internal/electionrepository"
)

type ListProposals struct {
	ElectionID   string
	Page         *int
	ItemsPerPage *int
}

type ListProposalsResponse struct {
	Proposals []Proposal
}

type Proposal struct {
	ElectionID  string
	ProposalID  string
	OwnerUserID string
	Name        string
	Description string
	ProposedAt  int
}

type listProposalsHandler struct {
	repository electionrepository.Repository
}

func NewListProposalsHandler(repository electionrepository.Repository) *listProposalsHandler {
	return &listProposalsHandler{
		repository: repository,
	}
}

func (h *listProposalsHandler) On(ctx context.Context, query ListProposals) (ListProposalsResponse, error) {

	page, itemsPerPage := cqrs.DefaultPagination(query.Page, query.ItemsPerPage, electionrepository.DefaultItemsPerPage)

	repoProposals, err := h.repository.ListProposals(ctx,
		query.ElectionID,
		page,
		itemsPerPage,
	)
	if err != nil {
		return ListProposalsResponse{}, err
	}

	proposals := make([]Proposal, len(repoProposals))
	for i := range repoProposals {
		proposals[i] = ToProposal(repoProposals[i])
	}

	return ListProposalsResponse{
		Proposals: proposals,
	}, nil
}

func ToProposal(proposal electionrepository.Proposal) Proposal {
	return Proposal{
		ElectionID:  proposal.ElectionID,
		OwnerUserID: proposal.OwnerUserID,
		ProposalID:  proposal.ProposalID,
		Name:        proposal.Name,
		Description: proposal.Description,
		ProposedAt:  proposal.ProposedAt,
	}
}
