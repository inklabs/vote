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

func TestCastVote(t *testing.T) {
	t.Run("saves to repository and raises event", func(t *testing.T) {
		// Given
		app := votetest.NewTestApp(t)
		ctx := app.GetAuthenticatedUserContext()
		const (
			electionID  = "eda792cc-9c16-4497-b21b-e64caa0e7629"
			proposalID1 = "fe420ed3-d56e-419c-a242-522ba89f92a2"
			proposalID2 = "6ade8319-ff51-4ea0-94bc-4b21b40cf872"
			proposalID3 = "c7a3936f-c5d0-4f56-a85e-c1544934c00e"
			ownerUserID = "a75f86b8-4454-4faa-af9e-19274264f621"
			voteID      = "9e08f651-f02f-414c-bf3e-1148a8c87e0c"
		)
		election1 := electionrepository.Election{
			ElectionID:      electionID,
			OrganizerUserID: "09dce1e9-568a-4fb2-945d-0ee9b95f5b04",
			Name:            "Election Name",
			Description:     "Election Description",
		}
		proposal1 := electionrepository.Proposal{
			ElectionID:  electionID,
			ProposalID:  proposalID1,
			OwnerUserID: ownerUserID,
			Name:        "Proposal Name 1",
			Description: "Proposal Description 1",
		}
		proposal2 := electionrepository.Proposal{
			ElectionID:  electionID,
			ProposalID:  proposalID2,
			OwnerUserID: ownerUserID,
			Name:        "Proposal Name 2",
			Description: "Proposal Description 2",
		}
		proposal3 := electionrepository.Proposal{
			ElectionID:  electionID,
			ProposalID:  proposalID3,
			OwnerUserID: ownerUserID,
			Name:        "Proposal Name 3",
			Description: "Proposal Description 3",
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
			VoteID:            voteID,
			ElectionID:        electionID,
			UserID:            ownerUserID,
			RankedProposalIDs: rankedProposalIDs,
		}

		// When
		response, err := app.ExecuteCommand(ctx, command)

		// Then
		require.NoError(t, err)
		assert.Equal(t, &cqrs.CommandResponse{
			Status: "OK",
		}, response)
		assert.Equal(t, event.VoteWasCast{
			VoteID:            voteID,
			ElectionID:        electionID,
			UserID:            ownerUserID,
			RankedProposalIDs: rankedProposalIDs,
			OccurredAt:        0,
		}, app.EventDispatcher.GetEvent(0))

		actualVotes, err := app.ElectionRepository.GetVotes(ctx, electionID)
		require.NoError(t, err)
		assert.Equal(t, []electionrepository.Vote{
			{
				VoteID:            voteID,
				ElectionID:        electionID,
				UserID:            ownerUserID,
				RankedProposalIDs: rankedProposalIDs,
			},
		}, actualVotes)
	})

	t.Run("errors", func(t *testing.T) {
		t.Run("when election not found", func(t *testing.T) {
			// Given
			app := votetest.NewTestApp(t)
			ctx := app.GetAuthenticatedUserContext()
			command := election.CastVote{
				ElectionID: "8493f7d9-1080-42a7-ae08-5d42d06941de",
				UserID:     "19a0abe7-fdd7-49f8-a01f-4fc0e0f12480",
				RankedProposalIDs: []string{
					"f6cec37d-cb86-4798-ae7d-51fd3cfc075b",
				},
			}

			// When
			_, err := app.ExecuteCommand(ctx, command)

			// Then
			require.Equal(t, electionrepository.NewErrElectionNotFound(command.ElectionID), err)
			assert.Empty(t, app.EventDispatcher.GetEvents())
		})

		t.Run("when proposal not found", func(t *testing.T) {
			// Given
			const unknownProposalID = "306cf23d-8196-4742-aca8-4bf9f43cd301"
			app := votetest.NewTestApp(t)
			ctx := app.GetAuthenticatedUserContext()
			election1 := electionrepository.Election{
				ElectionID:      "a771df96-3957-48d0-bcd3-63e3ab73ac75",
				OrganizerUserID: "3c51a70e-14cc-4cbb-b2dc-f58317470729",
				Name:            "Election Name",
				Description:     "Election Description",
			}
			require.NoError(t, app.ElectionRepository.SaveElection(ctx, election1))

			command := election.CastVote{
				ElectionID: election1.ElectionID,
				UserID:     "19a0abe7-fdd7-49f8-a01f-4fc0e0f12480",
				RankedProposalIDs: []string{
					unknownProposalID,
				},
			}

			// When
			_, err := app.ExecuteCommand(ctx, command)

			// Then
			require.Equal(t, electionrepository.NewErrProposalNotFound(unknownProposalID), err)
			assert.Empty(t, app.EventDispatcher.GetEvents())
		})

		t.Run("with proposal from another election", func(t *testing.T) {
			// Given
			app := votetest.NewTestApp(t)
			ctx := app.GetAuthenticatedUserContext()
			election1 := electionrepository.Election{
				ElectionID:      "dce69b68-aaa4-4602-88c6-0790c13c73b4",
				OrganizerUserID: "0c57ef2b-f0e1-40ef-a95e-82df8da6ad4e",
				Name:            "Election Name 1",
				Description:     "Election Description 1",
			}
			proposal1 := electionrepository.Proposal{
				ElectionID:  election1.ElectionID,
				ProposalID:  "841a61bd-6b8a-45ad-bd68-d01b7132e3b8",
				OwnerUserID: "2c3bfc60-ad8c-4f70-bb2b-94a2b9d98464",
				Name:        "Proposal Name 1",
				Description: "Proposal Description 1",
			}
			election2 := electionrepository.Election{
				ElectionID:      "cb14b1f7-47b3-4fbf-a554-79f4f5420a85",
				OrganizerUserID: "843b7424-0512-44a8-9f9b-f87a2d736475",
				Name:            "Election Name 2",
				Description:     "Election Description 2",
			}
			require.NoError(t, app.ElectionRepository.SaveElection(ctx, election1))
			require.NoError(t, app.ElectionRepository.SaveElection(ctx, election2))
			require.NoError(t, app.ElectionRepository.SaveProposal(ctx, proposal1))

			command := election.CastVote{
				ElectionID: election2.ElectionID,
				UserID:     "80d3c0bf-c523-459a-af47-36a2f3a643dd",
				RankedProposalIDs: []string{
					proposal1.ProposalID,
				},
			}

			// When
			_, err := app.ExecuteCommand(ctx, command)

			// Then
			expectedErr := electionrepository.NewErrInvalidElectionProposal(proposal1.ProposalID, election2.ElectionID)
			require.Equal(t, expectedErr, err)
			assert.Empty(t, app.EventDispatcher.GetEvents())
		})
	})
}
