package database

import (
	"context"
	"fmt"

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

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.DatabaseHost, config.DatabasePort, config.DatabaseUser, config.DatabasePassword, config.DatabaseName, config.DatabaseSslMode)
	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	if config.DatabaseTracing {
		pgxConfig.ConnConfig.Tracer = NewTracer(d.Logger)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return err
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		d.Logger.Panicf("error connecting to database: %v", err)
	}

	migrator := NewMigrator(pool, config, d.Logger)

	err = migrator.Migrate()
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	d.Pool = pool
	d.Querier = db.New(pool)

	return nil
}
