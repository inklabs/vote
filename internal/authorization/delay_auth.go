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

func (a *delayAuth) VerifyCommand(_ context.Context, _ cqrs.CommandHandler, _ cqrs.Command) error {
	time.Sleep(2 * time.Millisecond)
	return nil
}

func (a *delayAuth) VerifyAsyncCommand(_ context.Context, _ cqrs.AsyncCommandHandler, _ cqrs.AsyncCommand) error {
	time.Sleep(2 * time.Millisecond)
	return nil
}

func (a *delayAuth) VerifyQuery(_ context.Context, _ cqrs.QueryHandler, _ cqrs.Query) error {
	time.Sleep(2 * time.Millisecond)
	return nil
}

func (a *delayAuth) VerifyRequest(_ context.Context) error {
	time.Sleep(2 * time.Millisecond)
	return nil
}
