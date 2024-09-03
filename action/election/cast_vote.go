package election

import (
	"context"
	"time"

	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/pkg/clock"

	"github.com/inklabs/vote/event"
	"github.com/inklabs/vote/internal/electionrepository"
)

type CastVote struct {
	ElectionID        string
	UserID            string
	RankedProposalIDs []string
}

type castVoteHandler struct {
	repository electionrepository.Repository
	clock      clock.Clock
}

func NewCastVoteHandler(repository electionrepository.Repository, clock clock.Clock) *castVoteHandler {
	return &castVoteHandler{
		repository: repository,
		clock:      clock,
	}
}

func (h *castVoteHandler) On(ctx context.Context, cmd CastVote, eventRaiser cqrs.EventRaiser) error {
	ctx, span := tracer.Start(ctx, "vote.cast-vote")
	defer span.End()

	occurredAt := int(h.clock.Now().Unix())

	err := h.repository.SaveVote(ctx, electionrepository.Vote{
		ElectionID:        cmd.ElectionID,
		UserID:            cmd.UserID,
		RankedProposalIDs: append([]string{}, cmd.RankedProposalIDs...),
	})
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Millisecond)

	eventRaiser.Raise(event.VoteWasCast{
		ElectionID:        cmd.ElectionID,
		UserID:            cmd.UserID,
		RankedProposalIDs: append([]string{}, cmd.RankedProposalIDs...),
		OccurredAt:        occurredAt,
	})

	return nil
}
