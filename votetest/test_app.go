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
}

func NewTestApp(t *testing.T) testApp {
	t.Helper()

	a := testApp{
		t:                  t,
		EventDispatcher:    cqrstest.NewRecordingEventDispatcher(),
		ElectionRepository: electionrepository.NewInMemory(),
	}

	a.app = vote.NewApp(
		vote.WithEventDispatcher(a.EventDispatcher),
		vote.WithAuthorization(cqrstest.NewPassThruAuth()),
		vote.WithClock(incrementingclock.NewFromZero()),
		vote.WithAsyncCommandStore(asynccommandstore.NewInMemory()),
		vote.WithElectionRepository(a.ElectionRepository),
	)

	return a
}

func (a *testApp) ExecuteCommand(command cqrs.Command) (*cqrs.CommandResponse, error) {
	ctx := cqrstest.TimeoutContext(a.t)
	return a.app.CommandBus().Execute(ctx, command)
}
