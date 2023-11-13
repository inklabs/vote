package election

import (
	"context"
	"errors"

	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/pkg/clock"

	"github.com/inklabs/vote/event"
	"github.com/inklabs/vote/internal/electionrepository"
	"github.com/inklabs/vote/internal/tabulation"
)

type CloseElectionByOwner struct {
	ID         string
	ElectionID string
}

type closeElectionByOwnerHandler struct {
	repository electionrepository.Repository
	clock      clock.Clock
}

func NewCloseElectionByOwnerHandler(
	repository electionrepository.Repository,
	clock clock.Clock,
) *closeElectionByOwnerHandler {
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

	winningProposalID, err := h.getWinningProposalID(ctx, cmd.ElectionID, logger)
	if err != nil {
		logger.LogError("unable to get winning proposal")
		return err
	}

	selectedAt := int(h.clock.Now().Unix())
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

func (h *closeElectionByOwnerHandler) getWinningProposalID(ctx context.Context, electionID string, logger cqrs.AsyncCommandLogger) (string, error) {
	votes, err := h.repository.GetVotes(ctx, electionID)
	if err != nil {
		return "", err
	}

	tabulator := tabulation.NewRankedChoice(toRankedProposalVotes(votes))
	winningProposalID, err := tabulator.GetWinningProposal()
	if err != nil {
		if errors.Is(err, tabulation.ErrWinnerNotFound) {
			logger.LogError("winner not found")
		}
		return "", err
	}

	return winningProposalID, nil
}

func toRankedProposalVotes(votes []electionrepository.Vote) tabulation.Ballots {
	var rankedProposalVotes tabulation.Ballots

	for _, vote := range votes {
		proposalIDs := append([]string{}, vote.RankedProposalIDs...)
		rankedProposalVotes = append(rankedProposalVotes, proposalIDs)
	}

	return rankedProposalVotes
}
