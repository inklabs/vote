package election

import (
	"time"
)

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
	return GetElectionResultsResponse{
		ElectionID:        query.ElectionID,
		WinningProposalID: "250a63e3-97f6-452f-8557-9f85c5dc054f",
		SelectedAt:        int(time.Now().Unix()),
	}, nil
}
