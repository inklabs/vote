package vote

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/asynccommandbus"
	"github.com/inklabs/cqrs/asynccommandstore"
	"github.com/inklabs/cqrs/broker/inmemory"
	"github.com/inklabs/cqrs/broker/nats"
	"github.com/inklabs/cqrs/commandbus"
	"github.com/inklabs/cqrs/cqrstest"
	"github.com/inklabs/cqrs/eventdispatcher"
	"github.com/inklabs/cqrs/pkg/clock"
	"github.com/inklabs/cqrs/pkg/clock/provider/systemclock"
	"github.com/inklabs/cqrs/querybus"
	natsClient "github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	noopM "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/trace"
	noopT "go.opentelemetry.io/otel/trace/noop"

	"github.com/inklabs/vote/action/election"
	"github.com/inklabs/vote/event"
	"github.com/inklabs/vote/internal/authorization"
	"github.com/inklabs/vote/internal/electionrepository"
	"github.com/inklabs/vote/internal/electionrepository/inmemoryrepo"
	"github.com/inklabs/vote/internal/electionrepository/postgresrepo"
	"github.com/inklabs/vote/listener"
)

//go:generate go run github.com/inklabs/cqrs/cmd/domaingenerator -module github.com/inklabs/vote
//go:generate go run github.com/inklabs/cqrs/cmd/httpgenerator
//go:generate go run github.com/inklabs/cqrs/cmd/grpcgenerator
//go:generate go run github.com/inklabs/cqrs/cmd/sdkgenerator -js ./web/src/plugins/jsSDK_gen.js
//go:generate go run github.com/inklabs/cqrs/cmd/cligenerator

//go:embed domain.gob
var DomainBytes []byte

var Version = "dev-build"

type app struct {
	commandBus             cqrs.CommandBus
	asyncCommandBus        cqrs.AsyncCommandBus
	queryBus               cqrs.QueryBus
	eventDispatcher        cqrs.EventDispatcher
	asyncCommandStore      cqrs.AsyncCommandStore
	authorization          cqrs.Authorization
	clock                  clock.Clock
	useSyncLocalCommandBus bool
	ctxShutdowns           []func(ctx context.Context) error
	meterProvider          metric.MeterProvider
	tracerProvider         trace.TracerProvider

	electionRepository electionrepository.Repository
}

type Option func(a *app)

func WithAsyncCommandStore(store cqrs.AsyncCommandStore) Option {
	return func(a *app) {
		a.asyncCommandStore = store
	}
}

func WithEventDispatcher(eventDispatcher cqrs.EventDispatcher) Option {
	return func(a *app) {
		a.eventDispatcher = eventDispatcher
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

func WithCtxShutdown(shutdowns ...func(ctx context.Context) error) Option {
	return func(a *app) {
		a.ctxShutdowns = append(a.ctxShutdowns, shutdowns...)
	}
}

func WithTelemetry(meterProvider metric.MeterProvider, tracerProvider trace.TracerProvider) Option {
	return func(a *app) {
		a.meterProvider = meterProvider
		a.tracerProvider = tracerProvider
	}
}

func NewApp(opts ...Option) *app {
	a := &app{
		clock:              systemclock.New(),
		authorization:      cqrstest.NewPassThruAuth(),
		asyncCommandStore:  asynccommandstore.NewInMemory(),
		electionRepository: inmemoryrepo.New(),
		meterProvider:      otel.GetMeterProvider(),
		tracerProvider:     otel.GetTracerProvider(),
	}

	a.eventDispatcher = eventdispatcher.NewConcurrentLocal(
		a.meterProvider,
		a.tracerProvider,
		log.Default(),
		a.GetEventListeners(),
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
		a.meterProvider,
		a.tracerProvider,
	)

	a.asyncCommandBus = asynccommandbus.NewConcurrentLocal(
		commandHandlerRegistry,
		a.eventDispatcher,
		a.asyncCommandStore,
		a.clock,
		a.authorization,
		a.meterProvider,
		a.tracerProvider,
	)

	a.queryBus = querybus.NewLocal(
		queryHandlerRegistry,
		a.authorization,
		a.meterProvider,
		a.tracerProvider,
	)

	return a
}

func NewProdApp() *app {
	//resource := NewResource()
	//
	//conn := NewOLTPConn()
	//tracerProvider := NewOTLPTracerProvider(resource, conn)
	//meterProvider := NewOLTPMeterProvider(resource, conn)
	tracerProvider := noopT.NewTracerProvider()
	meterProvider := noopM.NewMeterProvider()

	eventDispatcher := newDistributedEventDispatcher(meterProvider, tracerProvider)

	repository := inmemoryrepo.New()
	asyncCommandStore := asynccommandstore.NewInMemory()

	//config := getPostgresConfig()
	//repository, err := postgresrepo.NewFromConfig(config)
	//if err != nil {
	//	log.Fatalf("error loading repository: %s", err)
	//}
	//ctx := context.Background()
	//err = repository.InitDB(ctx)
	//if err != nil {
	//	log.Fatalf("error initializing repository: %s", err)
	//}

	//postgresConfig := asynccommandstore.PostgresConfig{
	//	Host:       config.Host,
	//	Port:       config.Port,
	//	User:       config.User,
	//	Password:   config.Password,
	//	DBName:     config.DBName,
	//	SearchPath: config.SearchPath,
	//}
	//
	//asyncCommandStore, err := asynccommandstore.NewPostgresFromConfig(postgresConfig, GetAsyncCommands())
	//if err != nil {
	//	log.Fatalf("error getting async command store: %s", err)
	//}
	//err = asyncCommandStore.InitDB(ctx)
	//if err != nil {
	//	log.Fatalf("error initializing async command store: %s", err)
	//}

	return NewApp(
		WithAuthorization(authorization.NewDelayAuth()),
		WithAsyncCommandStore(asyncCommandStore),
		WithEventDispatcher(eventDispatcher),
		WithElectionRepository(repository),
		WithTelemetry(meterProvider, tracerProvider),
		//WithCtxShutdown(
		//	tracerProvider.Shutdown,
		//	meterProvider.Shutdown,
		//),
	)
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

func (a *app) MeterProvider() metric.MeterProvider {
	return a.meterProvider
}

func (a *app) TracerProvider() trace.TracerProvider {
	return a.tracerProvider
}

func (a *app) Stop() {
	a.asyncCommandBus.Stop()
	a.eventDispatcher.Stop()
	_ = a.asyncCommandStore.Close()
	ctx := context.Background()
	for _, shutdown := range a.ctxShutdowns {
		_ = shutdown(ctx)
	}
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
		election.NewListOpenElectionsHandler(a.electionRepository),
		election.NewListProposalsHandler(a.electionRepository),
		election.NewGetElectionHandler(a.electionRepository),
		election.NewGetProposalDetailsHandler(a.electionRepository),
		election.NewGetElectionResultsHandler(a.electionRepository),
	}
}

func (a *app) GetEventListeners() []cqrs.EventListener {
	return []cqrs.EventListener{
		listener.NewElectionWinnerVoterNotification(),
		listener.NewElectionWinnerMediaNotification(),
	}
}

func newDistributedEventDispatcher(meterProvider metric.MeterProvider, tracerProvider trace.TracerProvider) cqrs.EventDispatcher {
	eventRegistry := cqrs.NewEventRegistry()
	event.BindEvents(eventRegistry)

	eventSerializer := cqrs.NewEventPayloadSerializer(eventRegistry)

	logger := log.Default()

	const queueName = "vote-events"
	//publisher := GetRabbitMQBroker(logger)
	//publisher := GetNatsBroker(meterProvider, tracerProvider, logger)
	publisher := inmemory.NewBroker(meterProvider, tracerProvider, logger)

	eventDispatcher, err := eventdispatcher.NewDistributedEventDispatcher(
		queueName,
		publisher,
		eventSerializer,
		meterProvider,
		tracerProvider,
		logger,
	)
	if err != nil {
		panic(fmt.Errorf("failed to create rabbitmq dispatcher: %w", err))
	}

	return eventDispatcher
}

func GetNatsBroker(meterProvider metric.MeterProvider, tracerProvider trace.TracerProvider, logger *log.Logger) cqrs.Broker {
	return nats.NewBroker(
		natsClient.DefaultURL,
		meterProvider,
		tracerProvider,
		logger,
	)
}

//func GetRabbitMQBroker(logger *log.Logger) cqrs.Broker {
//	return rabbitmq.NewBroker(
//		"amqp://guest:guest@localhost:5672/",
//		logger,
//	)
//}

func getPostgresConfig() postgresrepo.Config {
	return postgresrepo.Config{
		Host:       "localhost",
		Port:       "5432",
		User:       "admin",
		Password:   "root",
		DBName:     "vote_demo",
		SearchPath: "",
	}
}
