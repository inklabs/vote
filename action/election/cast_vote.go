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

// CastVote casts a ballot for a given ElectionID. RankedProposalIDs contains the
// ranked candidates in order of preference: first, second, third and so forth. If your
// first choice doesn’t have a chance to win, your ballot counts for your next choice.
type CastVote struct {
	VoteID            string
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

	sleep.Rand(2 * time.Millisecond)

	err := h.repository.SaveVote(ctx, electionrepository.Vote{
		VoteID:            cmd.VoteID,
		ElectionID:        cmd.ElectionID,
		UserID:            cmd.UserID,
		RankedProposalIDs: append([]string{}, cmd.RankedProposalIDs...),
		SubmittedAt:       occurredAt,
	})
	if err != nil {
		return err
	}

	eventRaiser.Raise(event.VoteWasCast{
		VoteID:            cmd.VoteID,
		ElectionID:        cmd.ElectionID,
		UserID:            cmd.UserID,
		RankedProposalIDs: append([]string{}, cmd.RankedProposalIDs...),
		OccurredAt:        occurredAt,
	})

	return nil
}
