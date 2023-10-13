package election

import (
	"context"
	"time"

	"github.com/inklabs/cqrs"

	"github.com/inklabs/vote/event"
)

type CommenceElection struct {
	ElectionID      string
	OrganizerUserID string
	Name            string
	Description     string
}

type commenceElectionHandler struct{}

func NewCommenceElectionHandler() *commenceElectionHandler {
	return &commenceElectionHandler{}
}

func (h *commenceElectionHandler) On(_ context.Context, cmd CommenceElection, eventRaiser cqrs.EventRaiser) error {
	// TODO: Save election details to storage

	eventRaiser.Raise(event.ElectionHasCommenced{
		ElectionID:      cmd.ElectionID,
		OrganizerUserID: cmd.OrganizerUserID,
		Name:            cmd.Name,
		Description:     cmd.Description,
		OccurredAt:      int(time.Now().Unix()),
	})
	return nil
}
