package election_test

import (
	"testing"

	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/cqrstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inklabs/vote/action/election"
	"github.com/inklabs/vote/internal/electionrepository"
	"github.com/inklabs/vote/votetest"
)

func TestListProposals(t *testing.T) {
	// Given
	app := votetest.NewTestApp(t)
	election1 := electionrepository.Election{
		ElectionID:      "56320258-9b45-45c2-bb8b-4e5a204bbf23",
		OrganizerUserID: "2aa3897d-fd75-4b05-831b-f250d72984ba",
		Name:            "Election Name 1",
		Description:     "Election Description 1",
		CommencedAt:     1,
	}
	proposal1 := electionrepository.Proposal{
		ElectionID:  election1.ElectionID,
		ProposalID:  "6ee4e705-f6cc-4251-a306-af7fd005ea8a",
		OwnerUserID: "27543010-e150-474b-bd4c-6b51311d4bee",
		Name:        "Proposal Name 1",
		Description: "Proposal Description 1",
		ProposedAt:  1,
	}
	proposal2 := electionrepository.Proposal{
		ElectionID:  election1.ElectionID,
		ProposalID:  "6d7ee247-7b70-4f4e-a72e-9e5d825a1e54",
		OwnerUserID: "d05c0aaf-8c6b-4534-952a-529df8f5b0ee",
		Name:        "Proposal Name 2",
		Description: "Proposal Description 2",
		ProposedAt:  2,
	}
	proposal3 := electionrepository.Proposal{
		ElectionID:  election1.ElectionID,
		ProposalID:  "83f4163c-cac3-4a06-9f73-d3db42f34ae3",
		OwnerUserID: "7afa7ab3-4eee-49d0-805a-f1db74d23f7d",
		Name:        "Proposal Name 3",
		Description: "Proposal Description 3",
		ProposedAt:  3,
	}

	proposalDTO1 := election.ToProposal(proposal1)
	proposalDTO2 := election.ToProposal(proposal2)
	proposalDTO3 := election.ToProposal(proposal3)

	ctx := cqrstest.TimeoutContext(t)
	require.NoError(t, app.ElectionRepository.SaveElection(ctx, election1))
	require.NoError(t, app.ElectionRepository.SaveProposal(ctx, proposal1))
	require.NoError(t, app.ElectionRepository.SaveProposal(ctx, proposal2))
	require.NoError(t, app.ElectionRepository.SaveProposal(ctx, proposal3))

	t.Run("return proposals with default pagination", func(t *testing.T) {
		// Given
		query := election.ListProposals{
			ElectionID: election1.ElectionID,
		}

		// When
		response, err := app.ExecuteQuery(query)

		// Then
		require.NoError(t, err)
		assert.Equal(t, election.ListProposalsResponse{
			Proposals: []election.Proposal{
				proposalDTO1,
				proposalDTO2,
				proposalDTO3,
			},
		}, response)
	})

	t.Run("first page", func(t *testing.T) {
		// Given
		query := election.ListProposals{
			ElectionID:   election1.ElectionID,
			Page:         cqrs.Int(1),
			ItemsPerPage: cqrs.Int(2),
		}

		// When
		response, err := app.ExecuteQuery(query)

		// Then
		require.NoError(t, err)
		assert.Equal(t, election.ListProposalsResponse{
			Proposals: []election.Proposal{
				proposalDTO1,
				proposalDTO2,
			},
		}, response)
	})
}
