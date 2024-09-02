package authorization

import (
	"context"
	"time"

	"github.com/inklabs/cqrs"
)

type delayAuth struct{}

func NewDelayAuth() *delayAuth {
	return &delayAuth{}
}

func (a *delayAuth) VerifyCommand(ctx context.Context, _ cqrs.CommandHandler, _ cqrs.Command) error {
	_, span := tracer.Start(ctx, "auth.verify-command")
	defer span.End()

	time.Sleep(1 * time.Millisecond)
	return nil
}

func (a *delayAuth) VerifyAsyncCommand(ctx context.Context, _ cqrs.AsyncCommandHandler, _ cqrs.AsyncCommand) error {
	_, span := tracer.Start(ctx, "auth.verify-async-command")
	defer span.End()

	time.Sleep(1 * time.Millisecond)
	return nil
}

func (a *delayAuth) VerifyQuery(ctx context.Context, _ cqrs.QueryHandler, _ cqrs.Query) error {
	_, span := tracer.Start(ctx, "auth.verify-query")
	defer span.End()

	time.Sleep(1 * time.Millisecond)
	return nil
}

func (a *delayAuth) VerifyRequest(ctx context.Context) error {
	_, span := tracer.Start(ctx, "auth.verify-request")
	defer span.End()

	time.Sleep(1 * time.Millisecond)
	return nil
}
