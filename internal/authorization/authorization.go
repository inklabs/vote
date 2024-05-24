package authorization

import (
	"context"

	"github.com/inklabs/cqrs"
)

type Context interface {
	Context() context.Context
	Email() string
	UserID() string
	IsAdmin() bool
}

type CommandVerifier interface {
	VerifyAuthorization(ctx Context, command cqrs.Command) error
}

type AsyncCommandVerifier interface {
	VerifyAuthorization(ctx Context, command cqrs.AsyncCommand) error
}

type QueryVerifier interface {
	VerifyAuthorization(ctx Context, query cqrs.Query) error
}
