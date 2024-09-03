package authorization

import (
	"go.opentelemetry.io/otel"
)

const (
	instrumentationName = "github.com/inklabs/vote/internal/authorization/delay-auth"

	UserIDKey = "user.id"
)

var (
	tracer = otel.Tracer(instrumentationName)
	meter  = otel.Meter(instrumentationName)
)
