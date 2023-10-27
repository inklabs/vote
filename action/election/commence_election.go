package election

import (
	"context"

	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/pkg/clock"

	"github.com/inklabs/vote/event"
	"github.com/inklabs/vote/internal/electionrepository"
)

type CommenceElection struct {
	ElectionID      string
	OrganizerUserID string
	Name            string
	Description     string
}

type commenceElectionHandler struct {
	repository electionrepository.Repository
	clock      clock.Clock
}

func NewCommenceElectionHandler(repository electionrepository.Repository, clock clock.Clock) *commenceElectionHandler {
	return &commenceElectionHandler{
		repository: repository,
		clock:      clock,
	}
}

func (h *commenceElectionHandler) On(ctx context.Context, cmd CommenceElection, eventRaiser cqrs.EventRaiser) error {
	occurredAt := int(h.clock.Now().Unix())

	err := h.repository.SaveElection(ctx, electionrepository.Election{
		ElectionID:      cmd.ElectionID,
		OrganizerUserID: cmd.OrganizerUserID,
		Name:            cmd.Name,
		Description:     cmd.Description,
		OccurredAt:      occurredAt,
	})
	if err != nil {
		return err
	}

	eventRaiser.Raise(event.ElectionHasCommenced{
		ElectionID:      cmd.ElectionID,
		OrganizerUserID: cmd.OrganizerUserID,
		Name:            cmd.Name,
		Description:     cmd.Description,
		OccurredAt:      occurredAt,
	})

	return nil
}
