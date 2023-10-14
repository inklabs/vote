package election

import (
	"context"

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

	rankedProposalIDs := append([]string{}, cmd.RankedProposalIDs...)

	eventRaiser.Raise(event.VoteWasCast{
		ElectionID:        cmd.ElectionID,
		UserID:            cmd.UserID,
		RankedProposalIDs: rankedProposalIDs,
		OccurredAt:        0,
	})

	return nil
}
