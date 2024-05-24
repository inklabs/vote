package election_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inklabs/vote/action/election"
	"github.com/inklabs/vote/internal/electionrepository"
	"github.com/inklabs/vote/votetest"
)

func TestGetElectionResults(t *testing.T) {
	t.Run("returns election results", func(t *testing.T) {
		// Given
		app := votetest.NewTestApp(t)
		ctx := app.GetAuthenticatedUserContext()
		const (
			electionID        = "ef18565e-eba3-43ed-964e-40d872568f0a"
			winningProposalID = "35d414ea-4b5f-430a-9f57-ef48bce34ef2"
		)

		election1 := electionrepository.Election{
			ElectionID:        electionID,
			OrganizerUserID:   "1b207fbf-9797-4bfa-91e3-6b5eef1b9fc0",
			Name:              "Election Name",
			Description:       "Election Description",
			WinningProposalID: winningProposalID,
			IsClosed:          true,
			CommencedAt:       0,
			SelectedAt:        1,
			ClosedAt:          1,
		}
		require.NoError(t, app.ElectionRepository.SaveElection(ctx, election1))

		query := election.GetElectionResults{
			ElectionID: electionID,
		}

		// When
		response, err := app.ExecuteQuery(ctx, query)

		// Then
		require.NoError(t, err)
		assert.Equal(t, election.GetElectionResultsResponse{
			ElectionID:        electionID,
			WinningProposalID: winningProposalID,
			SelectedAt:        1,
		}, response)
	})

	t.Run("errors when election not found", func(t *testing.T) {
		// Given
		app := votetest.NewTestApp(t)
		ctx := app.GetAuthenticatedUserContext()
		query := election.GetElectionResults{
			ElectionID: "574af1df-542a-4644-8977-a5c6b1e0b26a",
		}

		// When
		_, err := app.ExecuteQuery(ctx, query)

		// Then
		require.Equal(t, err, electionrepository.ErrElectionNotFound)
	})
}
