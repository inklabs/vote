package election

import (
	"context"

	"github.com/inklabs/cqrs"
	"go.opentelemetry.io/otel/attribute"

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
	ctx, span := tracer.Start(ctx, "vote.list-proposals")
	defer span.End()

	page, itemsPerPage := cqrs.DefaultPagination(query.Page, query.ItemsPerPage, electionrepository.DefaultItemsPerPage)

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("itemsPerPage", itemsPerPage),
	)

	proposals, err := h.repository.ListProposals(ctx,
		query.ElectionID,
		page,
		itemsPerPage,
	)
	if err != nil {
		return ListProposalsResponse{}, err
	}

	return ListProposalsResponse{
		Proposals: ToProposals(proposals),
	}, nil
}

func ToProposals(repoProposals []electionrepository.Proposal) []Proposal {
	proposals := make([]Proposal, len(repoProposals))
	for i := range repoProposals {
		proposals[i] = ToProposal(repoProposals[i])
	}
	return proposals
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
