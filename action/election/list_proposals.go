package election

type ListProposals struct {
	ElectionID string
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
}

type listProposalsHandler struct{}

func NewListProposalsHandler() *listProposalsHandler {
	return &listProposalsHandler{}
}

func (h *listProposalsHandler) On(query ListProposals) (ListProposalsResponse, error) {
	return ListProposalsResponse{
		Proposals: []Proposal{
			{
				ElectionID:  "ec938e42-2009-403c-9ffb-71ae3c709f7d",
				ProposalID:  "261ef094-5a3d-47fd-aa4b-ee3c335e9f84",
				OwnerUserID: "8123e8ee-ab86-4743-92fe-15797ab873b4",
				Name:        "Lorem Ipsum",
				Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod",
			},
			{
				ElectionID:  "ec938e42-2009-403c-9ffb-71ae3c709f7d",
				ProposalID:  "abd18e54-0f1e-4e71-a02b-15e2011c6169",
				OwnerUserID: "8123e8ee-ab86-4743-92fe-15797ab873b4",
				Name:        "Ut enim",
				Description: "Ut enim ad minim veniam, quis nostrud exercitation ullamco",
			},
		},
	}, nil
}
