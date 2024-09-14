package vote

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v4"
	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/asynccommandbus"
	"github.com/inklabs/cqrs/asynccommandstore"
	"github.com/inklabs/cqrs/commandbus"
	"github.com/inklabs/cqrs/cqrstest"
	"github.com/inklabs/cqrs/eventdispatcher"
	"github.com/inklabs/cqrs/eventdispatcher/distributed"
	"github.com/inklabs/cqrs/eventdispatcher/distributed/provider/nats"
	"github.com/inklabs/cqrs/eventdispatcher/distributed/provider/rabbitmq"
	"github.com/inklabs/cqrs/pkg/clock"
	"github.com/inklabs/cqrs/pkg/clock/provider/systemclock"
	"github.com/inklabs/cqrs/querybus"
	natsClient "github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	metricNoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/trace"
	traceNoop "go.opentelemetry.io/otel/trace/noop"

	"github.com/inklabs/vote/action/election"
	"github.com/inklabs/vote/event"
	"github.com/inklabs/vote/internal/authorization"
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
	ctxShutdowns           []func(ctx context.Context) error

	electionRepository electionrepository.Repository
	meterProvider      metric.MeterProvider
	tracerProvider     trace.TracerProvider
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

func WithTelemetry(meterProvider metric.MeterProvider, tracerProvider trace.TracerProvider) Option {
	return func(a *app) {
		a.meterProvider = meterProvider
		a.tracerProvider = tracerProvider
	}
}

func WithCtxShutdown(shutdowns ...func(ctx context.Context) error) Option {
	return func(a *app) {
		a.ctxShutdowns = append(a.ctxShutdowns, shutdowns...)
	}
}

func NewApp(opts ...Option) *app {
	a := &app{
		clock:              systemclock.New(),
		authorization:      cqrstest.NewPassThruAuth(),
		asyncCommandStore:  asynccommandstore.NewInMemory(),
		electionRepository: electionrepository.NewInMemory(),
		meterProvider:      metricNoop.NewMeterProvider(),
		tracerProvider:     traceNoop.NewTracerProvider(),
	}

	a.eventDispatcher = eventdispatcher.NewConcurrentLocal(
		log.Default(),
		a.GetEventListeners(),
		a.meterProvider,
		a.tracerProvider,
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
	)

	return a
}

func NewProdApp() *app {
	resource := NewResource()

	tracerProvider := NewJaegerTracerProvider(resource)
	meterProvider := NewOLTPMeterProvider(resource)
	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)

	asyncCommandStore := asynccommandstore.NewBadger(
		badger.DefaultOptions("./.badger.db").
			WithLogger(nil),
		GetAsyncCommands(),
	)

	eventDispatcher := newDistributedEventDispatcher(meterProvider, tracerProvider)

	return NewApp(
		WithAuthorization(authorization.NewDelayAuth()),
		WithAsyncCommandStore(asyncCommandStore),
		WithTelemetry(meterProvider, tracerProvider),
		WithEventDispatcher(eventDispatcher),
		WithCtxShutdown(
			tracerProvider.Shutdown,
			meterProvider.Shutdown,
		),
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

func (a *app) GetMeterProvider() metric.MeterProvider {
	return a.meterProvider
}

func (a *app) GetTracerProvider() trace.TracerProvider {
	return a.tracerProvider
}

func newDistributedEventDispatcher(
	meterProvider metric.MeterProvider,
	tracerProvider trace.TracerProvider,
) cqrs.EventDispatcher {
	eventRegistry := cqrs.NewEventRegistry()
	event.BindEvents(eventRegistry)

	eventSerializer := cqrs.NewEventPayloadSerializer(eventRegistry)

	logger := log.Default()

	const queueName = "vote-events"
	//publisher := GetRabbitMQBroker(logger, meterProvider, tracerProvider)
	publisher := GetNatsBroker(logger, meterProvider, tracerProvider)

	eventDispatcher, err := distributed.NewEventDispatcher(
		queueName,
		publisher,
		eventSerializer,
		logger,
		meterProvider,
		tracerProvider,
	)
	if err != nil {
		panic(fmt.Errorf("failed to create rabbitmq dispatcher: %w", err))
	}

	return eventDispatcher
}

func GetNatsBroker(
	logger *log.Logger,
	meterProvider metric.MeterProvider,
	tracerProvider trace.TracerProvider,
) distributed.Broker {
	return nats.NewBroker(
		natsClient.DefaultURL,
		logger,
		meterProvider,
		tracerProvider,
	)
}

func GetRabbitMQBroker(
	logger *log.Logger,
	meterProvider metric.MeterProvider,
	tracerProvider trace.TracerProvider,
) distributed.Broker {
	return rabbitmq.NewBroker(
		"amqp://guest:guest@localhost:5672/",
		logger,
		meterProvider,
		tracerProvider,
	)
}
