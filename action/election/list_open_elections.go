package election

import (
	"context"

	"github.com/inklabs/cqrs"
	"go.opentelemetry.io/otel/attribute"

	"github.com/inklabs/vote/internal/electionrepository"
)

type ListOpenElections struct {
	Page          *int
	ItemsPerPage  *int
	SortBy        *string
	SortDirection *string
}

func (q ListOpenElections) ValidationRules() cqrs.ValidationRuleMap {
	return cqrs.ValidationRuleMap{
		"SortBy": cqrs.OptionalValidValues(
			"Name",
			"CommencedAt",
		),
		"SortDirection": cqrs.OptionalValidSortDirection(),
		"Page":          cqrs.OptionalValidMinRange(1),
		"ItemsPerPage":  cqrs.OptionalValidRange(1, 10),
	}
}

type ListOpenElectionsResponse struct {
	OpenElections []OpenElection
}

type OpenElection struct {
	ElectionID      string
	OrganizerUserID string
	Name            string
	Description     string
	CommencedAt     int
}

type listOpenElectionsHandler struct {
	repository electionrepository.Repository
}

func NewListOpenElectionsHandler(repository electionrepository.Repository) *listOpenElectionsHandler {
	return &listOpenElectionsHandler{
		repository: repository,
	}
}

func (h *listOpenElectionsHandler) On(ctx context.Context, query ListOpenElections) (ListOpenElectionsResponse, error) {
	ctx, span := tracer.Start(ctx, "vote.list-open-elections")
	defer span.End()

	page, itemsPerPage := cqrs.DefaultPagination(query.Page, query.ItemsPerPage, electionrepository.DefaultItemsPerPage)
	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("itemsPerPage", itemsPerPage),
	)

	elections, err := h.repository.ListOpenElections(ctx,
		page,
		itemsPerPage,
		query.SortBy,
		query.SortDirection,
	)
	if err != nil {
		return ListOpenElectionsResponse{}, err
	}

	return ListOpenElectionsResponse{
		OpenElections: ToOpenElections(elections),
	}, nil
}

func ToOpenElections(elections []electionrepository.Election) []OpenElection {
	openElections := make([]OpenElection, len(elections))
	for i := range elections {
		openElections[i] = ToOpenElection(elections[i])
	}
	return openElections
}

func ToOpenElection(election electionrepository.Election) OpenElection {
	return OpenElection{
		ElectionID:      election.ElectionID,
		OrganizerUserID: election.OrganizerUserID,
		Name:            election.Name,
		Description:     election.Description,
		CommencedAt:     election.CommencedAt,
	}
}
