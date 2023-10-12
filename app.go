package vote

import (
	"log"

	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/asynccommandbus"
	"github.com/inklabs/cqrs/asynccommandstore"
	"github.com/inklabs/cqrs/commandbus"
	"github.com/inklabs/cqrs/cqrstest"
	"github.com/inklabs/cqrs/eventdispatcher"
	"github.com/inklabs/cqrs/pkg/clock/provider/systemclock"
	"github.com/inklabs/cqrs/querybus"

	"github.com/inklabs/vote/action/election"
)

//go:generate go run github.com/inklabs/cqrs/cmd/domaingenerator -module github.com/inklabs/vote

type app struct {
	commandBus      cqrs.CommandBus
	asyncCommandBus cqrs.AsyncCommandBus
	queryBus        cqrs.QueryBus
}

func NewApp() *app {
	domainEventListeners := getDomainEventListeners()
	commandHandlerRegistry := cqrs.NewCommandHandlerRegistry(
		getCommandHandlers(),
		getAsyncCommandHandlers(),
	)
	queryHandlerRegistry := cqrs.NewQueryHandlerRegistry(getQueryHandlers())
	asyncCommandStore := asynccommandstore.NewInMemory()
	logger := log.Default()
	clock := systemclock.New()
	eventDispatcher := eventdispatcher.NewConcurrentLocal(logger, domainEventListeners)
	authorization := cqrstest.NewPassThruAuth()

	commandBus := commandbus.NewLocal(
		commandHandlerRegistry,
		eventDispatcher,
		authorization,
	)

	asyncCommandBus := asynccommandbus.NewConcurrentLocal(
		commandHandlerRegistry,
		eventDispatcher,
		asyncCommandStore,
		clock,
		authorization,
	)

	queryBus := querybus.NewLocal(
		queryHandlerRegistry,
		authorization,
	)

	a := &app{
		commandBus:      commandBus,
		asyncCommandBus: asyncCommandBus,
		queryBus:        queryBus,
	}

	return a
}

func (a *app) CommandBus() cqrs.CommandBus {
	return a.commandBus
}

func (a *app) AsyncCommandBus() cqrs.AsyncCommandBus {
	return a.asyncCommandBus
}

func (a *app) QueryBus() cqrs.QueryBus {
	return a.queryBus
}

func getCommandHandlers() []cqrs.CommandHandler {
	return []cqrs.CommandHandler{
		election.NewCommenceElectionHandler(),
	}
}

func getAsyncCommandHandlers() []cqrs.AsyncCommandHandler {
	return []cqrs.AsyncCommandHandler{
		election.NewCloseElectionByOwnerHandler(),
	}
}

func getQueryHandlers() []cqrs.QueryHandler {
	return []cqrs.QueryHandler{}
}

func getDomainEventListeners() []cqrs.EventListener {
	return []cqrs.EventListener{}
}
