package database

import (
	"context"
	"fmt"

	"github.com/exaring/otelpgx"
	"github.com/jeancarloshp/calorieai/internal/domain"

	"github.com/jeancarloshp/calorieai/pkg/database/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Logger  domain.Logger
	Querier db.Querier
	Pool    *pgxpool.Pool
}

func New(logger domain.Logger) *Database {
	return &Database{
		Logger: logger,
	}
}

func (d *Database) NewConnection(config *domain.Config) error {
	ctx := context.Background()

	pgxConfig, err := pgxpool.ParseConfig(config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	if config.TracingEnabled && config.DatabaseTracing {
		pgxConfig.ConnConfig.Tracer = otelpgx.NewTracer()
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return err
	}

	if config.TracingEnabled && config.DatabaseTracing {
		if err := otelpgx.RecordStats(pool); err != nil {
			return fmt.Errorf("unable to record database stats: %w", err)
		}
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		d.Logger.Panicf("error connecting to database: %v", err)
	}

	migrator, err := NewMigrator(pool, config.DatabaseMigrationURL, d.Logger)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	err = migrator.Migrate()
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	d.Pool = pool
	d.Querier = db.New(pool)

	return nil
}
