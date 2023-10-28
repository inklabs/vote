package election

import (
	"context"

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
	occurredAt := int(h.clock.Now().Unix())

	err := h.repository.SaveVote(ctx, electionrepository.Vote{
		ElectionID:        cmd.ElectionID,
		UserID:            cmd.UserID,
		RankedProposalIDs: append([]string(nil), cmd.RankedProposalIDs...),
	})
	if err != nil {
		return err
	}

	eventRaiser.Raise(event.VoteWasCast{
		ElectionID:        cmd.ElectionID,
		UserID:            cmd.UserID,
		RankedProposalIDs: append([]string(nil), cmd.RankedProposalIDs...),
		OccurredAt:        occurredAt,
	})

	return nil
}
