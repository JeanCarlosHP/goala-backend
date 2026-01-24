package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jeancarloshp/calorieai/internal/domain"
)

type Tracer struct {
	logger domain.Logger
}

func NewTracer(logger domain.Logger) *Tracer {
	return &Tracer{
		logger: logger,
	}
}

func (tracer *Tracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	query := data.SQL
	for i, arg := range data.Args {
		query = strings.Replace(query, fmt.Sprintf("$%d", i+1), fmt.Sprintf("'%v'", arg), 1)
	}

	tracer.logger.Debug("Executing query", "query", query)

	return ctx
}

func (tracer *Tracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}
