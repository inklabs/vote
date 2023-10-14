package election

import (
	"context"
	"time"

	"github.com/inklabs/cqrs"

	"github.com/inklabs/vote/event"
)

type MakeProposal struct {
	ElectionID  string
	ProposalID  string
	OwnerUserID string
	Name        string
	Description string
}

type makeProposalHandler struct{}

func NewMakeProposalHandler() *makeProposalHandler {
	return &makeProposalHandler{}
}

func (h *makeProposalHandler) On(_ context.Context, cmd MakeProposal, eventRaiser cqrs.EventRaiser) error {
	// TODO: save proposal details to storage

	eventRaiser.Raise(event.ProposalWasMade{
		ElectionID:  cmd.ElectionID,
		ProposalID:  cmd.ProposalID,
		OwnerUserID: cmd.OwnerUserID,
		Name:        cmd.Name,
		Description: cmd.Description,
		OccurredAt:  int(time.Now().Unix()),
	})

	return nil
}
