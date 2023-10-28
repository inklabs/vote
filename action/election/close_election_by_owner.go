package election

import (
	"context"

	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/pkg/clock"

	"github.com/inklabs/vote/event"
	"github.com/inklabs/vote/internal/electionrepository"
)

type CloseElectionByOwner struct {
	ID         string
	ElectionID string
}

type closeElectionByOwnerHandler struct {
	repository electionrepository.Repository
	clock      clock.Clock
}

func NewCloseElectionByOwnerHandler(repository electionrepository.Repository, clock clock.Clock) *closeElectionByOwnerHandler {
	return &closeElectionByOwnerHandler{
		repository: repository,
		clock:      clock,
	}
}

func (h *closeElectionByOwnerHandler) On(ctx context.Context, cmd CloseElectionByOwner, eventRaiser cqrs.EventRaiser, logger cqrs.AsyncCommandLogger) error {
	election, err := h.repository.GetElection(ctx, cmd.ElectionID)
	if err != nil {
		logger.LogError("election not found: %s", cmd.ElectionID)
		return err
	}

	// TODO: tabulate results, and persist to storage
	selectedAt := int(h.clock.Now().Unix())
	winningProposalID := "todo"

	election.IsClosed = true
	election.ClosedAt = selectedAt
	election.SelectedAt = selectedAt
	election.WinningProposalID = winningProposalID

	err = h.repository.SaveElection(ctx, election)
	if err != nil {
		return err
	}

	logger.LogInfo("Closing election with winner: %s", winningProposalID)

	eventRaiser.Raise(event.ElectionWinnerWasSelected{
		ElectionID:        cmd.ElectionID,
		WinningProposalID: winningProposalID,
		SelectedAt:        selectedAt,
	})

	return nil
}
