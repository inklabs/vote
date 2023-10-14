package election

type GetElectionResults struct {
	ElectionID string
}

type GetElectionResultsResponse struct {
	ElectionID        string
	WinningProposalID string
	SelectedAt        int
}

type getElectionResultsHandler struct{}

func NewGetElectionResultsHandler() *getElectionResultsHandler {
	return &getElectionResultsHandler{}
}

func (h *getElectionResultsHandler) On(query GetElectionResults) (GetElectionResultsResponse, error) {
	return GetElectionResultsResponse{}, nil
}
