package database

import (
	"context"
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	pgxv5driver "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jeancarloshp/calorieai/internal/domain"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Migrator struct {
	Conn   *pgxpool.Pool
	config *domain.Config
	logger domain.Logger
}

func NewMigrator(conn *pgxpool.Pool, config *domain.Config, logger domain.Logger) *Migrator {
	return &Migrator{
		Conn:   conn,
		config: config,
		logger: logger,
	}
}

func (mg *Migrator) Migrate() (err error) {
	var src source.Driver
	src, err = iofs.New(migrations, "migrations")
	if err != nil {
		return
	}

	var dst database.Driver
	db := stdlib.OpenDBFromPool(mg.Conn)
	dst, err = pgxv5driver.WithInstance(db, &pgxv5driver.Config{
		MultiStatementMaxSize: pgxv5driver.DefaultMultiStatementMaxSize,
	})
	if err != nil {
		db.Close()
		return fmt.Errorf("create pgx driver: %w", err)
	}
	var m *migrate.Migrate
	m, err = migrate.NewWithInstance("embed", src, "pgx5", dst)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer m.Close()

	v, d, err := dst.Version()
	if err != nil {
		return fmt.Errorf("get current migration version: %w", err)
	}

	l, err := mg.latestSourceVersion(src)
	if err != nil {
		return fmt.Errorf("get latest migration version: %w", err)
	}

	if v == l && !d {
		mg.logger.Info("database is up to date")
		return nil
	}

	mg.logger.Info("running migrations", "current_version", v, "latest_version", l)

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		if strings.Contains(err.Error(), "Dirty database version") {
			mg.logger.Info("database is dirty")
			err := mg.verifyDatabaseVersion(m, src, dst)
			if err != nil {
				return fmt.Errorf("verify database version: %w", err)
			}
		} else {
			return fmt.Errorf("run migration %w", err)
		}
	}

	mg.logger.Info("migrations ran successfully")

	return nil
}

func (mg *Migrator) latestSourceVersion(sourceDriver source.Driver) (int, error) {
	var v uint
	var err error
	v, err = sourceDriver.First()
	if err != nil {
		return 0, err
	}

	for {
		var nextVersion uint
		nextVersion, err = sourceDriver.Next(v)
		if err == os.ErrNotExist {
			break
		} else if pathErr, ok := err.(*os.PathError); ok && pathErr.Err == os.ErrNotExist {
			break
		} else if err != nil {
			return 0, err
		}
		v = nextVersion
	}

	return int(v), nil
}

func (mg *Migrator) verifyDatabaseVersion(migrator *migrate.Migrate, sourceDriver source.Driver, destinationDriver database.Driver) error {
	versionInDB, d, err := destinationDriver.Version()
	if err != nil {
		return fmt.Errorf("get current migration version in DB: %w", err)
	}

	latestVersionOnFile, err := mg.latestSourceVersion(sourceDriver)
	if err != nil {
		return fmt.Errorf("check latest migration version on file: %w", err)
	}

	switch {
	case versionInDB == latestVersionOnFile && d:
		err := mg.forceDatabaseVersion(migrator, versionInDB)
		if err != nil {
			return err
		}

	case versionInDB < latestVersionOnFile:
		err = mg.forceDatabaseVersion(migrator, latestVersionOnFile-1)
		if err != nil {
			return fmt.Errorf("force set migration to latest version on file: %w", err)
		}

	case versionInDB == 1:
		ct, err := mg.Conn.Exec(context.Background(), "DROP TABLE IF EXISTS schema_migrations")
		if err != nil {
			return fmt.Errorf("drop migration: %w", err)
		}
		if ct.RowsAffected() > 0 {
			mg.logger.Info("table schema_migrations dropped")
		}
	}

	return mg.Migrate()
}

func (mg *Migrator) forceDatabaseVersion(migrator *migrate.Migrate, version int) error {
	mg.logger.Info("force setting migration version", "version", version)
	err := migrator.Force(version)
	if err != nil {
		return fmt.Errorf("force set migration version: %w", err)
	}

	return nil
}
