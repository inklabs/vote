package authorization

import (
	"go.opentelemetry.io/otel"
)

const instrumentationName = "github.com/inklabs/vote/internal/authorization/delay-auth"

var (
	tracer = otel.Tracer(instrumentationName)
	meter  = otel.Meter(instrumentationName)
)
