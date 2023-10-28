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

func TestCommenceElection(t *testing.T) {
	app := votetest.NewTestApp(t)

	t.Run("saves to repository and raises event", func(t *testing.T) {
		// Given
		const electionID = "6c194c91-bb68-4933-a6ba-7c5867a5f54d"
		command := election.CommenceElection{
			ElectionID:      electionID,
			OrganizerUserID: "73adf147-ce92-4c9f-9f9c-5464210e68da",
			Name:            "Election Name",
			Description:     "Election Description",
		}

		// When
		response, err := app.ExecuteCommand(command)

		// Then
		require.NoError(t, err)
		assert.Equal(t, &cqrs.CommandResponse{
			Status: "OK",
		}, response)
		assert.Equal(t, event.ElectionHasCommenced{
			ElectionID:      electionID,
			OrganizerUserID: command.OrganizerUserID,
			Name:            command.Name,
			Description:     command.Description,
			OccurredAt:      0,
		}, app.EventDispatcher.GetEvent(0))

		ctx := cqrstest.TimeoutContext(t)
		actualElection, err := app.ElectionRepository.GetElection(ctx, electionID)
		require.NoError(t, err)
		assert.Equal(t, electionrepository.Election{
			ElectionID:        electionID,
			OrganizerUserID:   command.OrganizerUserID,
			Name:              command.Name,
			Description:       command.Description,
			CommencedAt:       0,
			WinningProposalID: "",
			IsClosed:          false,
			ClosedAt:          0,
		}, actualElection)
	})
}
