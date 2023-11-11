package election

import (
	"context"

	"github.com/inklabs/cqrs"

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

	page, itemsPerPage := cqrs.DefaultPagination(query.Page, query.ItemsPerPage, electionrepository.DefaultItemsPerPage)

	elections, err := h.repository.ListOpenElections(ctx,
		page,
		itemsPerPage,
		query.SortBy,
		query.SortDirection,
	)
	if err != nil {
		return ListOpenElectionsResponse{}, err
	}

	openElections := make([]OpenElection, len(elections))
	for i := range elections {
		openElections[i] = ToOpenElection(elections[i])
	}

	return ListOpenElectionsResponse{
		OpenElections: openElections,
	}, nil
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
