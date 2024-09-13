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

type MakeProposal struct {
	ElectionID  string
	ProposalID  string
	OwnerUserID string
	Name        string
	Description string
}

type makeProposalHandler struct {
	repository electionrepository.Repository
	clock      clock.Clock
}

func NewMakeProposalHandler(repository electionrepository.Repository, clock clock.Clock) *makeProposalHandler {
	return &makeProposalHandler{
		repository: repository,
		clock:      clock,
	}
}

func (h *makeProposalHandler) On(ctx context.Context, cmd MakeProposal, eventRaiser cqrs.EventRaiser) error {
	proposedAt := int(h.clock.Now().Unix())

	sleep.Rand(2 * time.Millisecond)

	err := h.repository.SaveProposal(ctx, electionrepository.Proposal{
		ElectionID:  cmd.ElectionID,
		ProposalID:  cmd.ProposalID,
		OwnerUserID: cmd.OwnerUserID,
		Name:        cmd.Name,
		Description: cmd.Description,
		ProposedAt:  proposedAt,
	})
	if err != nil {
		return err
	}

	eventRaiser.Raise(event.ProposalWasMade{
		ElectionID:  cmd.ElectionID,
		ProposalID:  cmd.ProposalID,
		OwnerUserID: cmd.OwnerUserID,
		Name:        cmd.Name,
		Description: cmd.Description,
		ProposedAt:  proposedAt,
	})

	return nil
}
