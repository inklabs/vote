package vote

import (
	_ "embed"
	"log"

	"github.com/dgraph-io/badger/v4"
	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/asynccommandbus"
	"github.com/inklabs/cqrs/asynccommandstore"
	"github.com/inklabs/cqrs/commandbus"
	"github.com/inklabs/cqrs/cqrstest"
	"github.com/inklabs/cqrs/eventdispatcher"
	"github.com/inklabs/cqrs/pkg/clock"
	"github.com/inklabs/cqrs/pkg/clock/provider/systemclock"
	"github.com/inklabs/cqrs/querybus"

	"github.com/inklabs/vote/action/election"
	"github.com/inklabs/vote/internal/electionrepository"
	"github.com/inklabs/vote/listener"
)

//go:generate go run github.com/inklabs/cqrs/cmd/domaingenerator -module github.com/inklabs/vote
//go:generate go run github.com/inklabs/cqrs/cmd/httpgenerator
//go:generate go run github.com/inklabs/cqrs/cmd/grpcgenerator
//go:generate go run github.com/inklabs/cqrs/cmd/sdkgenerator
//go:generate go run github.com/inklabs/cqrs/cmd/cligenerator

//go:embed domain.gob
var DomainBytes []byte

type app struct {
	commandBus             cqrs.CommandBus
	asyncCommandBus        cqrs.AsyncCommandBus
	queryBus               cqrs.QueryBus
	eventDispatcher        cqrs.EventDispatcher
	asyncCommandStore      cqrs.AsyncCommandStore
	authorization          cqrs.Authorization
	clock                  clock.Clock
	useSyncLocalCommandBus bool

	electionRepository electionrepository.Repository
}

type Option func(a *app)

func WithAsyncCommandStore(store cqrs.AsyncCommandStore) Option {
	return func(a *app) {
		a.asyncCommandStore = store
	}
}

func WithEventDispatcher(dispatcher cqrs.EventDispatcher) Option {
	return func(a *app) {
		a.eventDispatcher = dispatcher
	}
}

func WithAuthorization(authorization cqrs.Authorization) Option {
	return func(a *app) {
		a.authorization = authorization
	}
}

func WithClock(clock clock.Clock) Option {
	return func(a *app) {
		a.clock = clock
	}
}

func WithElectionRepository(repository electionrepository.Repository) Option {
	return func(a *app) {
		a.electionRepository = repository
	}
}

func WithSyncLocalAsyncCommandBus() Option {
	return func(a *app) {
		a.useSyncLocalCommandBus = true
	}
}

func NewProdApp() *app {
	return NewApp(
		WithAsyncCommandStore(
			asynccommandstore.NewBadgerAsyncCommandStore(
				badger.DefaultOptions("./.badger.db").
					WithLogger(nil),
				GetAsyncCommands(),
			),
		),
	)
}

func NewApp(opts ...Option) *app {
	a := &app{
		clock:             systemclock.New(),
		authorization:     cqrstest.NewPassThruAuth(),
		asyncCommandStore: asynccommandstore.NewInMemory(),
	}

	a.eventDispatcher = eventdispatcher.NewConcurrentLocal(
		log.Default(),
		a.getDomainEventListeners(),
	)

	for _, opt := range opts {
		opt(a)
	}

	commandHandlerRegistry := cqrs.NewCommandHandlerRegistry(
		a.getCommandHandlers(),
		a.getAsyncCommandHandlers(),
	)
	queryHandlerRegistry := cqrs.NewQueryHandlerRegistry(
		a.getQueryHandlers(),
	)

	a.commandBus = commandbus.NewLocal(
		commandHandlerRegistry,
		a.eventDispatcher,
		a.authorization,
	)

	if a.useSyncLocalCommandBus {
		a.asyncCommandBus = asynccommandbus.NewSyncLocal(
			commandHandlerRegistry,
			a.eventDispatcher,
			a.asyncCommandStore,
			a.clock,
			a.authorization,
		)
	} else {
		a.asyncCommandBus = asynccommandbus.NewConcurrentLocal(
			commandHandlerRegistry,
			a.eventDispatcher,
			a.asyncCommandStore,
			a.clock,
			a.authorization,
		)
	}

	a.queryBus = querybus.NewLocal(
		queryHandlerRegistry,
		a.authorization,
	)

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

func (a *app) Stop() {
	a.asyncCommandBus.Stop()
	a.eventDispatcher.Stop()
	_ = a.asyncCommandStore.Close()
}

func (a *app) getCommandHandlers() []cqrs.CommandHandler {
	return []cqrs.CommandHandler{
		election.NewCommenceElectionHandler(a.electionRepository, a.clock),
		election.NewMakeProposalHandler(a.electionRepository, a.clock),
		election.NewCastVoteHandler(a.electionRepository, a.clock),
	}
}

func (a *app) getAsyncCommandHandlers() []cqrs.AsyncCommandHandler {
	return []cqrs.AsyncCommandHandler{
		election.NewCloseElectionByOwnerHandler(a.electionRepository, a.clock),
	}
}

func (a *app) getQueryHandlers() []cqrs.QueryHandler {
	return []cqrs.QueryHandler{
		election.NewListOpenElectionsHandler(),
		election.NewListProposalsHandler(),
		election.NewGetProposalDetailsHandler(),
		election.NewGetElectionResultsHandler(a.electionRepository),
	}
}

func (a *app) getDomainEventListeners() []cqrs.EventListener {
	return []cqrs.EventListener{
		listener.NewElectionWinnerVoterNotification(),
		listener.NewElectionWinnerMediaNotification(),
	}
}
