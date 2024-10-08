package election

import (
	"context"
	"time"

	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/pkg/clock"

	"github.com/inklabs/vote/event"
	"github.com/inklabs/vote/internal/electionrepository"
	"github.com/inklabs/vote/pkg/sleep"
)

// CommenceElection instantiates a new open election that is ready for proposals and voting.
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

	sleep.Rand(2 * time.Millisecond)

	err := h.repository.SaveElection(ctx, electionrepository.Election{
		ElectionID:      cmd.ElectionID,
		OrganizerUserID: cmd.OrganizerUserID,
		Name:            cmd.Name,
		Description:     cmd.Description,
		CommencedAt:     occurredAt,
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
