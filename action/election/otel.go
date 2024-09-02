package election

import (
	"go.opentelemetry.io/otel"
)

const instrumentationName = "github.com/inklabs/vote/action/election"

var (
	tracer = otel.Tracer(instrumentationName)
	meter  = otel.Meter(instrumentationName)
)
