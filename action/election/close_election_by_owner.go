package election

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/pkg/clock"

	"github.com/inklabs/vote/event"
	"github.com/inklabs/vote/internal/authorization"
	"github.com/inklabs/vote/internal/electionrepository"
	"github.com/inklabs/vote/internal/rcv"
	"github.com/inklabs/vote/pkg/sleep"
)

// CloseElectionByOwner is an asynchronous command that closes an election and
// calculates a winner by using the Ranked Choice Voting (RCV) electoral system.
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

func (h *closeElectionByOwnerHandler) Verify(ctx authorization.Context, cmd CloseElectionByOwner) error {
	election, err := h.repository.GetElection(ctx.Context(), cmd.ElectionID)
	if err != nil {
		return err
	}

	if ctx.UserID() != election.OrganizerUserID {
		log.Printf("user %s does not match election organizer user %s", ctx.UserID(), election.OrganizerUserID)
		return cqrs.ErrAccessDenied
	}

	return nil
}

func (h *closeElectionByOwnerHandler) On(ctx context.Context, cmd CloseElectionByOwner, eventRaiser cqrs.EventRaiser, logger cqrs.AsyncCommandLogger) error {
	ctx, span := tracer.Start(ctx, "vote.close-election-by-owner")
	defer span.End()

	election, err := h.repository.GetElection(ctx, cmd.ElectionID)
	if err != nil {
		logger.LogError("election not found: %s", cmd.ElectionID)
		return err
	}

	winningProposalID, err := h.getWinningProposalID(ctx, cmd.ElectionID, logger)
	if err != nil {
		logger.LogError("unable to get winning proposal")
		err = fmt.Errorf("unable to get winning proposal: %w", err)
		cqrs.RecordSpanError(span, err)
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

func simulateProcessing(logger cqrs.AsyncCommandLogger, totalToProcess int) {
	logger.SetTotalToProcess(totalToProcess)

	sleepDuration := 5 * time.Second / time.Duration(totalToProcess)

	for i := 0; i < totalToProcess; i++ {
		logger.IncrementTotalProcessed()

		if totalToProcess < 10 || i%(totalToProcess/10) == 0 {
			logger.Flush()
			sleep.Sleep(sleepDuration)
		}
	}

	logger.Flush()
}

func (h *closeElectionByOwnerHandler) getWinningProposalID(ctx context.Context, electionID string, logger cqrs.AsyncCommandLogger) (string, error) {
	votes, err := h.repository.GetVotes(ctx, electionID)
	if err != nil {
		return "", err
	}

	if len(votes) == 0 {
		logger.LogError("no votes found for election")
		return "", ErrNoVotesFound
	}

	simulateProcessing(logger, len(votes))

	tabulator := rcv.NewSingleWinner(toRankedProposalVotes(votes))
	winningProposalID, err := tabulator.GetWinningProposal()
	if err != nil {
		if errors.Is(err, rcv.ErrWinnerNotFound) {
			logger.LogError("winner not found")
		}
		return "", err
	}

	return winningProposalID, nil
}

func toRankedProposalVotes(votes []electionrepository.Vote) rcv.Ballots {
	var rankedProposalVotes rcv.Ballots

	for _, vote := range votes {
		proposalIDs := append([]string{}, vote.RankedProposalIDs...)
		rankedProposalVotes = append(rankedProposalVotes, proposalIDs)
	}

	return rankedProposalVotes
}

var ErrNoVotesFound = errors.New("no votes found for election")
