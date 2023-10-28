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

func TestCastVote(t *testing.T) {
	t.Run("saves to repository and raises event", func(t *testing.T) {
		// Given
		app := votetest.NewTestApp(t)
		const (
			electionID  = "eda792cc-9c16-4497-b21b-e64caa0e7629"
			proposalID1 = "fe420ed3-d56e-419c-a242-522ba89f92a2"
			proposalID2 = "6ade8319-ff51-4ea0-94bc-4b21b40cf872"
			proposalID3 = "c7a3936f-c5d0-4f56-a85e-c1544934c00e"
			ownerUserID = "a75f86b8-4454-4faa-af9e-19274264f621"
		)
		ctx := cqrstest.TimeoutContext(t)
		election1 := electionrepository.Election{
			ElectionID:      electionID,
			OrganizerUserID: "09dce1e9-568a-4fb2-945d-0ee9b95f5b04",
			Name:            "Election Name",
			Description:     "Election Description",
			CommencedAt:     0,
		}
		proposal1 := electionrepository.Proposal{
			ElectionID:  electionID,
			ProposalID:  proposalID1,
			OwnerUserID: ownerUserID,
			Name:        "Proposal Name 1",
			Description: "Proposal Description 1",
			OccurredAt:  0,
		}
		proposal2 := electionrepository.Proposal{
			ElectionID:  electionID,
			ProposalID:  proposalID2,
			OwnerUserID: ownerUserID,
			Name:        "Proposal Name 2",
			Description: "Proposal Description 2",
			OccurredAt:  0,
		}
		proposal3 := electionrepository.Proposal{
			ElectionID:  electionID,
			ProposalID:  proposalID3,
			OwnerUserID: ownerUserID,
			Name:        "Proposal Name 3",
			Description: "Proposal Description 3",
			OccurredAt:  0,
		}
		require.NoError(t, app.ElectionRepository.SaveElection(ctx, election1))
		require.NoError(t, app.ElectionRepository.SaveProposal(ctx, proposal1))
		require.NoError(t, app.ElectionRepository.SaveProposal(ctx, proposal2))
		require.NoError(t, app.ElectionRepository.SaveProposal(ctx, proposal3))

		rankedProposalIDs := []string{
			proposalID2,
			proposalID3,
			proposalID1,
		}
		command := election.CastVote{
			ElectionID:        electionID,
			UserID:            ownerUserID,
			RankedProposalIDs: rankedProposalIDs,
		}

		// When
		response, err := app.ExecuteCommand(command)

		// Then
		require.NoError(t, err)
		assert.Equal(t, &cqrs.CommandResponse{
			Status: "OK",
		}, response)
		assert.Equal(t, event.VoteWasCast{
			ElectionID:        electionID,
			UserID:            ownerUserID,
			RankedProposalIDs: rankedProposalIDs,
			OccurredAt:        0,
		}, app.EventDispatcher.GetEvent(0))

		actualVotes, err := app.ElectionRepository.GetVotes(ctx, electionID)
		require.NoError(t, err)
		assert.Equal(t, []electionrepository.Vote{
			{
				ElectionID:        electionID,
				UserID:            ownerUserID,
				RankedProposalIDs: rankedProposalIDs,
			},
		}, actualVotes)
	})
}
