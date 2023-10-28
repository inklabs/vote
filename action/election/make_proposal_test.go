package election_test

import (
	"testing"

	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/cqrstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inklabs/vote/action/election"
	"github.com/inklabs/vote/event"
	"github.com/inklabs/vote/internal/electionrepository"
	"github.com/inklabs/vote/votetest"
)

func TestMakeProposal(t *testing.T) {
	t.Run("saves proposal", func(t *testing.T) {
		// Given
		app := votetest.NewTestApp(t)
		const (
			electionID = "130fd0f1-5872-447b-8938-ed97d8df082c"
			proposalID = "324031bd-6071-45de-bce9-b69ea902d4c2"
		)

		ctx := cqrstest.TimeoutContext(t)
		election1 := electionrepository.Election{
			ElectionID:      electionID,
			OrganizerUserID: "b916d870-ea41-4ca0-a6dd-e1b06ba33105",
			Name:            "Election Name",
			Description:     "Election Description",
			CommencedAt:     0,
		}
		require.NoError(t, app.ElectionRepository.SaveElection(ctx, election1))

		command := election.MakeProposal{
			ElectionID:  electionID,
			ProposalID:  proposalID,
			OwnerUserID: "5ee66ff6-060b-4722-8574-0b298628b3be",
			Name:        "Proposal Name",
			Description: "Proposal Description",
		}

		// When
		response, err := app.ExecuteCommand(command)

		// Then
		require.NoError(t, err)
		assert.Equal(t, &cqrs.CommandResponse{
			Status: "OK",
		}, response)
		assert.Equal(t, event.ProposalWasMade{
			ElectionID:  electionID,
			ProposalID:  proposalID,
			OwnerUserID: command.OwnerUserID,
			Name:        command.Name,
			Description: command.Description,
			OccurredAt:  0,
		}, app.EventDispatcher.GetEvent(0))

		actualElection, err := app.ElectionRepository.GetProposals(ctx, electionID)
		require.NoError(t, err)
		assert.Equal(t, []electionrepository.Proposal{
			{
				ElectionID:  electionID,
				ProposalID:  proposalID,
				OwnerUserID: command.OwnerUserID,
				Name:        command.Name,
				Description: command.Description,
				ProposedAt:  0,
			},
		}, actualElection)
	})
}
