package listener

import (
	"go.opentelemetry.io/otel"
)

const instrumentationName = "github.com/inklabs/vote/listener"

var (
	tracer = otel.Tracer(instrumentationName)
	meter  = otel.Meter(instrumentationName)
)
