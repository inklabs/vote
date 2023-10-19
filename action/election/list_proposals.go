package election

type ListProposals struct {
	ElectionID string
}

type ListProposalsResponse struct {
	// TODO: Add support for complex slices
}

type listProposalsHandler struct{}

func NewListProposalsHandler() *listProposalsHandler {
	return &listProposalsHandler{}
}

func (h *listProposalsHandler) On(query ListProposals) (ListProposalsResponse, error) {
	return ListProposalsResponse{}, nil
}
