package votetest

import (
	"context"
	"testing"

	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/asynccommandstore"
	"github.com/inklabs/cqrs/cqrstest"
	"github.com/inklabs/cqrs/pkg/clock/provider/incrementingclock"
	"github.com/stretchr/testify/require"

	"github.com/inklabs/vote"
	"github.com/inklabs/vote/internal/authorization"
	"github.com/inklabs/vote/internal/electionrepository"
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
		t:                  t,
		jwtSigningKey:      []byte("9742fed04ba648bcb476a13b9e3d87e3"),
		RegularUserID:      "de06e622-9169-4351-b14e-9109dfd9dee3",
		AdminUserID:        "e5bca084-bf48-4b31-8bd2-233cfd5b6c92",
		EventDispatcher:    cqrstest.NewRecordingEventDispatcher(),
		AsyncCommandStore:  asynccommandstore.NewInMemory(),
		ElectionRepository: electionrepository.NewInMemory(),
	}

	a.app = vote.NewApp(
		vote.WithEventDispatcher(a.EventDispatcher),
		vote.WithAuthorization(authorization.NewJWTAuthorization(a.jwtSigningKey)),
		vote.WithClock(incrementingclock.NewFromZero()),
		vote.WithAsyncCommandStore(a.AsyncCommandStore),
		vote.WithElectionRepository(a.ElectionRepository),
		vote.WithSyncLocalAsyncCommandBus(),
	)

	return a
}

func (a *testApp) ExecuteCommand(ctx context.Context, command cqrs.Command) (*cqrs.CommandResponse, error) {
	return a.app.CommandBus().Execute(ctx, command)
}

// EnqueueCommand executes the async command synchronously via vote.WithSyncLocalAsyncCommandBus
func (a *testApp) EnqueueCommand(ctx context.Context, asyncCommand cqrs.AsyncCommand) (*cqrs.AsyncCommandResponse, error) {
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
