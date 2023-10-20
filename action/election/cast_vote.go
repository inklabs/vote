package election

import (
	"context"
	"time"

	"github.com/inklabs/cqrs"

	"github.com/inklabs/vote/event"
)

type CastVote struct {
	ElectionID        string
	UserID            string
	RankedProposalIDs []string
}

type castVoteHandler struct{}

func NewCastVoteHandler() *castVoteHandler {
	return &castVoteHandler{}
}

func (h *castVoteHandler) On(_ context.Context, cmd CastVote, eventRaiser cqrs.EventRaiser) error {
	// TODO: save vote details to storage

	eventRaiser.Raise(event.VoteWasCast{
		ElectionID:        cmd.ElectionID,
		UserID:            cmd.UserID,
		RankedProposalIDs: append([]string(nil), cmd.RankedProposalIDs...),
		OccurredAt:        int(time.Now().Unix()),
	})

	return nil
}
