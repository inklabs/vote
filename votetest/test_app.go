package votetest

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/asynccommandstore"
	"github.com/inklabs/cqrs/cqrstest"
	"github.com/inklabs/cqrs/pkg/clock/provider/incrementingclock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inklabs/vote"
	"github.com/inklabs/vote/internal/authorization"
	"github.com/inklabs/vote/internal/electionrepository"
	"github.com/inklabs/vote/internal/electionrepository/inmemoryrepo"
	"github.com/inklabs/vote/internal/electionrepository/postgresrepo"
)

type testApp struct {
	t   *testing.T
	app cqrs.App

	EventDispatcher    cqrstest.RecordingEventDispatcher
	ElectionRepository electionrepository.Repository
	AsyncCommandStore  cqrs.AsyncCommandStore
	jwtSigningKey      []byte
	RegularUserID      string
	AdminUserID        string
}

func NewTestApp(t *testing.T) testApp {
	t.Helper()

	a := testApp{
		t:                 t,
		jwtSigningKey:     []byte("9742fed04ba648bcb476a13b9e3d87e3"),
		RegularUserID:     "de06e622-9169-4351-b14e-9109dfd9dee3",
		AdminUserID:       "e5bca084-bf48-4b31-8bd2-233cfd5b6c92",
		EventDispatcher:   cqrstest.NewRecordingEventDispatcher(),
		AsyncCommandStore: asynccommandstore.NewInMemory(),
	}

	if os.Getenv("PG_HOST") != "" {
		db := getTestDB(t)
		repository, err := postgresrepo.NewFromDB(db)
		require.NoError(t, err)

		ctx := cqrstest.TimeoutContext(t)
		require.NoError(t, repository.InitDB(ctx))

		truncateTables(t, db)

		a.ElectionRepository = repository
	} else {
		a.ElectionRepository = inmemoryrepo.New()
	}

	a.app = vote.NewApp(
		vote.WithEventDispatcher(a.EventDispatcher),
		vote.WithAuthorization(authorization.NewJWTAuthorization(a.jwtSigningKey)),
		vote.WithClock(incrementingclock.NewFromZero()),
		vote.WithAsyncCommandStore(a.AsyncCommandStore),
		vote.WithElectionRepository(a.ElectionRepository),
	)

	return a
}

func (a *testApp) ExecuteCommand(ctx context.Context, command cqrs.Command) (cqrs.CommandResponse, error) {
	return a.app.CommandBus().Execute(ctx, command)
}

func (a *testApp) EnqueueCommand(ctx context.Context, asyncCommand cqrs.AsyncCommand) (cqrs.AsyncCommandResponse, error) {
	return a.app.AsyncCommandBus().Enqueue(ctx, asyncCommand)
}

func (a *testApp) ExecuteQuery(ctx context.Context, query cqrs.Query) (cqrs.QueryResponse, error) {
	return a.app.QueryBus().Execute(ctx, query)
}

func (a *testApp) GetAuthenticatedUserContext() context.Context {
	return context.WithValue(cqrstest.TimeoutContext(a.t), "authorization", a.getUserToken())
}

func (a *testApp) getUserToken() string {
	return a.getSignedBearerToken(authorization.JWTClaims{
		Email:   "john.user@example.com",
		UserID:  a.RegularUserID,
		IsAdmin: false,
	})
}

func (a *testApp) getAdminToken() string {
	return a.getSignedBearerToken(authorization.JWTClaims{
		Email:   "john.admin@example.com",
		UserID:  a.AdminUserID,
		IsAdmin: true,
	})
}

func (a *testApp) getSignedBearerToken(claims authorization.JWTClaims) string {
	signedToken, err := authorization.NewSignedToken(claims, a.jwtSigningKey)
	require.NoError(a.t, err)
	return "Bearer " + signedToken
}

func getTestDB(t *testing.T) *sql.DB {
	config, err := postgresrepo.NewConfigFromEnvironment()
	require.NoError(t, err)

	db, err := postgresrepo.NewDB(config)
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, db.Close())
	})

	return db
}

func truncateTables(t *testing.T, db *sql.DB) {
	ctx := cqrstest.TimeoutContext(t)
	sqlStatements := []string{
		"TRUNCATE TABLE vote_ranked_proposal CASCADE",
		"TRUNCATE TABLE vote CASCADE",
		"TRUNCATE TABLE proposal CASCADE",
		"TRUNCATE TABLE election CASCADE",
	}

	for _, sqlStatement := range sqlStatements {
		_, err := db.ExecContext(ctx, sqlStatement)
		require.NoError(t, err, sqlStatement)
	}
}
