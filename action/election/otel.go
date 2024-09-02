package election

import (
	"go.opentelemetry.io/otel"
)

const instrumentationName = "github.com/pdt256/vote/action/election"

var (
	tracer = otel.Tracer(instrumentationName)
	meter  = otel.Meter(instrumentationName)
)
