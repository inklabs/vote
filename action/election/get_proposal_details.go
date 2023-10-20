package election

import (
	"time"
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

type getProposalDetailsHandler struct{}

func NewGetProposalDetailsHandler() *getProposalDetailsHandler {
	return &getProposalDetailsHandler{}
}

func (h *getProposalDetailsHandler) On(query GetProposalDetails) (GetProposalDetailsResponse, error) {
	return GetProposalDetailsResponse{
		ElectionID:  "073b1fdb-7af4-4d4d-b2de-50d8d64f8a15",
		ProposalID:  query.ProposalID,
		OwnerUserID: "b8192901-6384-4474-ba8f-531941348033",
		Name:        "Lorem Ipsum",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod",
		ProposedAt:  int(time.Now().Unix()),
	}, nil
}
