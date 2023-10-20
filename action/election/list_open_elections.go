package election

import (
	"time"
)

type ListOpenElections struct{}

type ListOpenElectionsResponse struct {
	OpenElections []OpenElection
}

type OpenElection struct {
	ElectionID      string
	OrganizerUserID string
	Name            string
	Description     string
	OpenedAt        int
}

type listOpenElectionsHandler struct{}

func NewListOpenElectionsHandler() *listOpenElectionsHandler {
	return &listOpenElectionsHandler{}
}

func (h *listOpenElectionsHandler) On(query ListOpenElections) (ListOpenElectionsResponse, error) {
	return ListOpenElectionsResponse{
		OpenElections: []OpenElection{
			{
				ElectionID:      "0aa295fb-151e-48d1-87a5-241b6403728f",
				OrganizerUserID: "df53cef2-18d1-47a2-81ba-e852e97c662e",
				Name:            "Lorem Ipsum",
				Description:     "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod",
				OpenedAt:        int(time.Now().Unix()),
			},
			{
				ElectionID:      "3d8bd133-f5fd-45ac-acb4-be9d56de3ba1",
				OrganizerUserID: "e6c718dc-d307-4f13-8844-a4e8cde67713",
				Name:            "Ut enim",
				Description:     "Ut enim ad minim veniam, quis nostrud exercitation ullamco",
				OpenedAt:        int(time.Now().Unix()),
			},
		},
	}, nil
}
