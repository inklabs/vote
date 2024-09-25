package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/eventdispatcher"

	"github.com/inklabs/vote"
	"github.com/inklabs/vote/event"
)

func main() {
	fmt.Println("Vote - Subscriber Daemon")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app := vote.NewProdApp()

	eventRegistry := cqrs.NewEventRegistry()
	event.BindEvents(eventRegistry)

	eventSerializer := cqrs.NewEventPayloadSerializer(eventRegistry)
	logger := log.Default()

	const queueName = "vote-events"
	//consumer := vote.GetRabbitMQBroker(logger, app.GetMeterProvider(), app.GetTracerProvider())
	consumer := vote.GetNatsBroker(logger)

	subscriber, err := eventdispatcher.NewDistributedEventSubscriber(
		queueName,
		consumer,
		eventSerializer,
		logger,
		app.GetEventListeners(),
	)
	if err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()

	fmt.Println("Shutting down Subscriber Daemon")
	subscriber.Stop()
	app.Stop()
}
