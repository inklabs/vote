package election_test

import (
	"testing"

	"github.com/inklabs/cqrs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inklabs/vote/action/election"
	"github.com/inklabs/vote/internal/electionrepository"
	"github.com/inklabs/vote/votetest"
)

func TestListOpenElections(t *testing.T) {
	// Given
	app := votetest.NewTestApp(t)
	ctx := app.GetAuthenticatedUserContext()
	election1 := electionrepository.Election{
		ElectionID:      "00fa69d7-9e49-449e-9281-e5a6db476e33",
		OrganizerUserID: "76574368-caa8-478b-9764-a7f1e0fa4662",
		Name:            "Election Name 1",
		Description:     "Election Description 1",
		CommencedAt:     1,
	}
	election2 := electionrepository.Election{
		ElectionID:      "e3b09d7a-85e8-4736-b41f-e859ec6a77ab",
		OrganizerUserID: "7e616235-ff17-4fde-a4ef-fc0ec29f4fc0",
		Name:            "Election Name 2",
		Description:     "Election Description 2",
		CommencedAt:     2,
	}
	election3 := electionrepository.Election{
		ElectionID:      "ca3a38cc-e2cf-4c5e-aa57-b9bea1b37faa",
		OrganizerUserID: "5e77628e-9d76-46b4-b46c-cfef8f5f6945",
		Name:            "Election Name 3",
		Description:     "Election Description 3",
		CommencedAt:     3,
	}

	openElection1 := election.ToOpenElection(election1)
	openElection2 := election.ToOpenElection(election2)
	openElection3 := election.ToOpenElection(election3)

	require.NoError(t, app.ElectionRepository.SaveElection(ctx, election1))
	require.NoError(t, app.ElectionRepository.SaveElection(ctx, election2))
	require.NoError(t, app.ElectionRepository.SaveElection(ctx, election3))

	t.Run("returns open elections with default pagination and sorting", func(t *testing.T) {
		// Given
		query := election.ListOpenElections{}

		// When
		response, err := app.ExecuteQuery(ctx, query)

		// Then
		require.NoError(t, err)
		assert.Equal(t, election.ListOpenElectionsResponse{
			OpenElections: []election.OpenElection{
				openElection1,
				openElection2,
				openElection3,
			},
		}, response)
	})

	t.Run("sorted by name descending", func(t *testing.T) {
		// Given
		query := election.ListOpenElections{
			SortBy:        cqrs.String("Name"),
			SortDirection: cqrs.String("descending"),
		}

		// When
		response, err := app.ExecuteQuery(ctx, query)

		// Then
		require.NoError(t, err)
		assert.Equal(t, election.ListOpenElectionsResponse{
			OpenElections: []election.OpenElection{
				openElection3,
				openElection2,
				openElection1,
			},
		}, response)
	})

	t.Run("sorted by CommencedAt descending", func(t *testing.T) {
		// Given
		query := election.ListOpenElections{
			SortBy:        cqrs.String("CommencedAt"),
			SortDirection: cqrs.String("descending"),
		}

		// When
		response, err := app.ExecuteQuery(ctx, query)

		// Then
		require.NoError(t, err)
		assert.Equal(t, election.ListOpenElectionsResponse{
			OpenElections: []election.OpenElection{
				openElection3,
				openElection2,
				openElection1,
			},
		}, response)
	})

	t.Run("first page default sort", func(t *testing.T) {
		// Given
		query := election.ListOpenElections{
			Page:         cqrs.Int(1),
			ItemsPerPage: cqrs.Int(2),
		}

		// When
		response, err := app.ExecuteQuery(ctx, query)

		// Then
		require.NoError(t, err)
		assert.Equal(t, election.ListOpenElectionsResponse{
			OpenElections: []election.OpenElection{
				openElection1,
				openElection2,
			},
		}, response)
	})

	t.Run("second page default sort", func(t *testing.T) {
		// Given
		query := election.ListOpenElections{
			Page:         cqrs.Int(2),
			ItemsPerPage: cqrs.Int(2),
		}

		// When
		response, err := app.ExecuteQuery(ctx, query)

		// Then
		require.NoError(t, err)
		assert.Equal(t, election.ListOpenElectionsResponse{
			OpenElections: []election.OpenElection{
				openElection3,
			},
		}, response)
	})
}
