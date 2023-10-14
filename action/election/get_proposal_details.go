package election

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
	return GetProposalDetailsResponse{}, nil
}
