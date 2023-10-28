package election_test

import (
	"testing"

	"github.com/inklabs/cqrs/cqrstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inklabs/vote/action/election"
	"github.com/inklabs/vote/internal/electionrepository"
	"github.com/inklabs/vote/votetest"
)

func TestGetProposalDetails(t *testing.T) {
	t.Run("returns proposal details", func(t *testing.T) {
		// Given
		app := votetest.NewTestApp(t)
		const (
			electionID = "246f38ca-382c-4b8e-85a6-2c05d25093a2"
			proposalID = "e0476e03-6c6e-4ab9-9c3b-c20970664b63"
		)
		election1 := electionrepository.Election{
			ElectionID:      electionID,
			OrganizerUserID: "439a3234-60ad-4afc-9f95-5db239102a38",
			Name:            "Election Name",
			Description:     "Election Description",
		}
		proposal1 := electionrepository.Proposal{
			ElectionID:  electionID,
			ProposalID:  proposalID,
			OwnerUserID: "67b2c7b7-173f-4cb8-9f06-299cc345fd50",
			Name:        "Proposal Name",
			Description: "Proposal Description",
			ProposedAt:  1,
		}

		ctx := cqrstest.TimeoutContext(t)
		require.NoError(t, app.ElectionRepository.SaveElection(ctx, election1))
		require.NoError(t, app.ElectionRepository.SaveProposal(ctx, proposal1))

		query := election.GetProposalDetails{
			ProposalID: proposalID,
		}

		// When
		response, err := app.ExecuteQuery(query)

		// Then
		require.NoError(t, err)
		assert.Equal(t, election.GetProposalDetailsResponse{
			ElectionID:  proposal1.ElectionID,
			ProposalID:  proposal1.ProposalID,
			OwnerUserID: proposal1.OwnerUserID,
			Name:        proposal1.Name,
			Description: proposal1.Description,
			ProposedAt:  proposal1.ProposedAt,
		}, response)
	})

	t.Run("errors", func(t *testing.T) {
		t.Run("when proposal not found", func(t *testing.T) {
			// Given
			app := votetest.NewTestApp(t)
			query := election.GetProposalDetails{
				ProposalID: "e0476e03-6c6e-4ab9-9c3b-c20970664b63",
			}

			// When
			_, err := app.ExecuteQuery(query)

			// Then
			require.Equal(t, err, electionrepository.ErrProposalNotFound)
		})
	})
}
