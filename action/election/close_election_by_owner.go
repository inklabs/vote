package election

import (
	"context"
	"time"

	"github.com/inklabs/cqrs"

	"github.com/inklabs/vote/event"
)

type CloseElectionByOwner struct {
	ID         string
	ElectionID string
}

type closeElectionByOwnerHandler struct{}

func NewCloseElectionByOwnerHandler() *closeElectionByOwnerHandler {
	return &closeElectionByOwnerHandler{}
}

func (h *closeElectionByOwnerHandler) On(_ context.Context, cmd CloseElectionByOwner, eventRaiser cqrs.EventRaiser, logger cqrs.AsyncCommandLogger) error {
	// TODO: tabulate results, and persist to storage
	winningProposalID := "todo"

	eventRaiser.Raise(event.ElectionWinnerWasSelected{
		ElectionID:        cmd.ElectionID,
		WinningProposalID: winningProposalID,
		OccurredAt:        int(time.Now().Unix()),
	})

	return nil
}
