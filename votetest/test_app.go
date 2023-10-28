package votetest

import (
	"testing"

	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/asynccommandstore"
	"github.com/inklabs/cqrs/cqrstest"
	"github.com/inklabs/cqrs/pkg/clock/provider/incrementingclock"

	"github.com/inklabs/vote"
	"github.com/inklabs/vote/internal/electionrepository"
)

type testApp struct {
	t   *testing.T
	app cqrs.App

	EventDispatcher    cqrstest.RecordingEventDispatcher
	ElectionRepository electionrepository.Repository
	AsyncCommandStore  cqrs.AsyncCommandStore
}

func NewTestApp(t *testing.T) testApp {
	t.Helper()

	a := testApp{
		t:                  t,
		EventDispatcher:    cqrstest.NewRecordingEventDispatcher(),
		AsyncCommandStore:  asynccommandstore.NewInMemory(),
		ElectionRepository: electionrepository.NewInMemory(),
	}

	a.app = vote.NewApp(
		vote.WithEventDispatcher(a.EventDispatcher),
		vote.WithAuthorization(cqrstest.NewPassThruAuth()),
		vote.WithClock(incrementingclock.NewFromZero()),
		vote.WithAsyncCommandStore(a.AsyncCommandStore),
		vote.WithElectionRepository(a.ElectionRepository),
		vote.WithSyncLocalAsyncCommandBus(),
	)

	return a
}

func (a *testApp) ExecuteCommand(command cqrs.Command) (*cqrs.CommandResponse, error) {
	ctx := cqrstest.TimeoutContext(a.t)
	return a.app.CommandBus().Execute(ctx, command)
}

// EnqueueCommand executes the async command synchronously via vote.WithSyncLocalAsyncCommandBus
func (a *testApp) EnqueueCommand(asyncCommand cqrs.AsyncCommand) (*cqrs.AsyncCommandResponse, error) {
	ctx := cqrstest.TimeoutContext(a.t)
	return a.app.AsyncCommandBus().Enqueue(ctx, asyncCommand)
}

func (a *testApp) ExecuteQuery(query cqrs.Query) (cqrs.QueryResponse, error) {
	ctx := cqrstest.TimeoutContext(a.t)
	return a.app.QueryBus().Execute(ctx, query)
}
