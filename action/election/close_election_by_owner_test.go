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

func TestCloseElectionByOwner(t *testing.T) {
	t.Run("closes election and raises event", func(t *testing.T) {
		// Given
		app := votetest.NewTestApp(t)
		const (
			electionID        = "c14490e2-c27f-4dd3-89e9-d8a0f17341f1"
			winningProposalID = "todo"
		)

		ctx := cqrstest.TimeoutContext(t)
		election1 := electionrepository.Election{
			ElectionID:      electionID,
			OrganizerUserID: "1b207fbf-9797-4bfa-91e3-6b5eef1b9fc0",
			Name:            "Election Name",
			Description:     "Election Description",
			OccurredAt:      0,
		}
		require.NoError(t, app.ElectionRepository.SaveElection(ctx, election1))

		commandID := "4f4442af-a4b0-43d7-acc7-f83a6fd1220c"
		command := election.CloseElectionByOwner{
			ID:         commandID,
			ElectionID: electionID,
		}

		// When
		response, err := app.EnqueueCommand(command)

		// Then
		require.NoError(t, err)
		assert.Equal(t, &cqrs.AsyncCommandResponse{
			ID:            commandID,
			Status:        "QUEUED",
			HasBeenQueued: true,
		}, response)

		assert.Equal(t, event.ElectionWinnerWasSelected{
			ElectionID:        electionID,
			WinningProposalID: winningProposalID,
			OccurredAt:        2,
		}, app.EventDispatcher.GetEvent(0))

		status, err := app.AsyncCommandStore.GetAsyncCommandStatus(ctx, commandID)
		require.NoError(t, err)
		assert.True(t, status.IsFinished)
		assert.True(t, status.IsSuccess)

		logs, err := app.AsyncCommandStore.GetAsyncCommandLogs(ctx, commandID)
		require.NoError(t, err)
		assert.Equal(t, []cqrs.AsyncCommandLog{
			{
				Type:           cqrs.CommandLogInfo,
				CreatedAtMicro: 3000000,
				Message:        "Closing election with winner: " + winningProposalID,
			},
		}, logs)

		actualElection, err := app.ElectionRepository.GetElection(ctx, electionID)
		require.NoError(t, err)
		assert.Equal(t, electionrepository.Election{
			ElectionID:        electionID,
			OrganizerUserID:   election1.OrganizerUserID,
			Name:              election1.Name,
			Description:       election1.Description,
			OccurredAt:        0,
			WinningProposalID: winningProposalID,
			IsClosed:          true,
			ClosedAt:          2,
		}, actualElection)
	})

	t.Run("errors when election not found", func(t *testing.T) {
		// Given
		app := votetest.NewTestApp(t)
		const electionID = "a1255efd-3917-458e-83d3-24a7cc7e0b40"

		ctx := cqrstest.TimeoutContext(t)
		commandID := "4a9bc117-8ac9-4ad4-bee1-17ac1f044bf4"
		command := election.CloseElectionByOwner{
			ID:         commandID,
			ElectionID: electionID,
		}

		// When
		response, err := app.EnqueueCommand(command)

		// Then
		require.NoError(t, err)
		assert.Equal(t, &cqrs.AsyncCommandResponse{
			ID:            commandID,
			Status:        "QUEUED",
			HasBeenQueued: true,
		}, response)
		assert.Empty(t, app.EventDispatcher.GetEvents())

		status, err := app.AsyncCommandStore.GetAsyncCommandStatus(ctx, commandID)
		require.NoError(t, err)
		assert.True(t, status.IsFinished)
		assert.False(t, status.IsSuccess)

		logs, err := app.AsyncCommandStore.GetAsyncCommandLogs(ctx, commandID)
		require.NoError(t, err)
		assert.Equal(t, []cqrs.AsyncCommandLog{
			{
				Type:           cqrs.CommandLogError,
				CreatedAtMicro: 3000000,
				Message:        "election not found: " + electionID,
			},
		}, logs)
	})
}
