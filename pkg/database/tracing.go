package database

import (
	"go.opentelemetry.io/otel/trace"
)

func NewPgxTracer(tp trace.TracerProvider) trace.Tracer {
	return tp.Tracer("pgx")
}
