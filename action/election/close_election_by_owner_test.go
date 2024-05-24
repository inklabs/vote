package election_test

import (
	"testing"

	"github.com/inklabs/cqrs"
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
		ctx := app.GetAuthenticatedUserContext()
		const electionID = "c14490e2-c27f-4dd3-89e9-d8a0f17341f1"

		election1 := electionrepository.Election{
			ElectionID:      electionID,
			OrganizerUserID: app.RegularUserID,
			Name:            "Election Name",
			Description:     "Election Description",
			CommencedAt:     0,
		}
		proposal1 := electionrepository.Proposal{
			ElectionID:  electionID,
			ProposalID:  "e156a821-d2b2-404b-b8eb-53f052cbe30d",
			OwnerUserID: "d0adb8db-b56e-4f53-8e4a-4e6cac0cb95b",
			Name:        "Proposal Name",
			Description: "Proposal Description",
			ProposedAt:  0,
		}
		vote1 := electionrepository.Vote{
			ElectionID:        electionID,
			UserID:            "fa465d85-ad59-49ca-8ae4-9be7c88c6ef1",
			RankedProposalIDs: []string{proposal1.ProposalID},
		}
		require.NoError(t, app.ElectionRepository.SaveElection(ctx, election1))
		require.NoError(t, app.ElectionRepository.SaveProposal(ctx, proposal1))
		require.NoError(t, app.ElectionRepository.SaveVote(ctx, vote1))

		winningProposalID := proposal1.ProposalID
		commandID := "4f4442af-a4b0-43d7-acc7-f83a6fd1220c"
		command := election.CloseElectionByOwner{
			ID:         commandID,
			ElectionID: electionID,
		}

		// When
		response, err := app.EnqueueCommand(ctx, command)

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
			SelectedAt:        2,
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
			WinningProposalID: winningProposalID,
			IsClosed:          true,
			CommencedAt:       0,
			ClosedAt:          2,
			SelectedAt:        2,
		}, actualElection)
	})

	t.Run("errors", func(t *testing.T) {
		t.Run("when election not found during authorization", func(t *testing.T) {
			// Given
			app := votetest.NewTestApp(t)
			ctx := app.GetAuthenticatedUserContext()
			const electionID = "a1255efd-3917-458e-83d3-24a7cc7e0b40"

			commandID := "4a9bc117-8ac9-4ad4-bee1-17ac1f044bf4"
			command := election.CloseElectionByOwner{
				ID:         commandID,
				ElectionID: electionID,
			}

			// When
			_, err := app.EnqueueCommand(ctx, command)

			// Then
			require.Equal(t, electionrepository.ErrElectionNotFound, err)
		})

		t.Run("when not authorized", func(t *testing.T) {
			// Given
			app := votetest.NewTestApp(t)
			ctx := app.GetAuthenticatedUserContext()
			const electionID = "c225599f-4c1e-4bce-91e8-84f09bfb663e"

			election1 := electionrepository.Election{
				ElectionID:      electionID,
				OrganizerUserID: "53293c94-dc72-4beb-8a1f-de9ad5f67329",
				Name:            "Election Name",
				Description:     "Election Description",
				CommencedAt:     0,
			}
			require.NoError(t, app.ElectionRepository.SaveElection(ctx, election1))

			commandID := "2d191963-1e65-49f8-9eef-db0571c48daf"
			command := election.CloseElectionByOwner{
				ID:         commandID,
				ElectionID: electionID,
			}

			// When
			_, err := app.EnqueueCommand(ctx, command)

			// Then
			require.Equal(t, cqrs.ErrAccessDenied, err)
		})
	})
}
