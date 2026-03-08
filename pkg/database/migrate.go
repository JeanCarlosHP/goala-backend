package database

import (
	"context"
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jeancarloshp/calorieai/internal/domain"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Migrator struct {
	Conn         *pgxpool.Pool
	logger       domain.Logger
	migrator     *migrate.Migrate
	sourceDriver source.Driver
}

func NewMigrator(conn *pgxpool.Pool, migrationURL string, logger domain.Logger) (*Migrator, error) {
	var src source.Driver
	src, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("create iofs source driver: %w", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", src, migrationURL)
	if err != nil {
		return nil, fmt.Errorf("create migrate instance: %w", err)
	}

	return &Migrator{
		Conn:         conn,
		logger:       logger,
		migrator:     migrator,
		sourceDriver: src,
	}, nil
}

func (mg *Migrator) Migrate() (err error) {
	defer func() {
		if srcErr, dbErr := mg.migrator.Close(); srcErr != nil || dbErr != nil {
			if srcErr != nil {
				mg.logger.Warn("failed to close migration source", "error", srcErr)
			}
			if dbErr != nil {
				mg.logger.Warn("failed to close migration database connection", "error", dbErr)
			}
		}
	}()

	v, d, err := mg.migrator.Version()
	if err != nil {
		return fmt.Errorf("get current migration version: %w", err)
	}

	latestVersion, err := mg.latestSourceVersion()
	if err != nil {
		return fmt.Errorf("get latest migration version: %w", err)
	}

	if v == latestVersion && !d {
		mg.logger.Info("database is up to date")
		return nil
	}

	mg.logger.Info("running migrations", "current_version", v, "latest_version", latestVersion)

	err = mg.migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		if strings.Contains(err.Error(), "Dirty database version") {
			mg.logger.Warn("database is dirty")

			err = mg.verifyDatabaseVersion()
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

func (mg *Migrator) latestSourceVersion() (uint, error) {
	var v uint
	var err error
	v, err = mg.sourceDriver.First()
	if err != nil {
		return 0, err
	}

	for {
		var nextVersion uint
		nextVersion, err = mg.sourceDriver.Next(v)
		if err == os.ErrNotExist {
			break
		} else if pathErr, ok := err.(*os.PathError); ok && pathErr.Err == os.ErrNotExist {
			break
		} else if err != nil {
			return 0, err
		}
		v = nextVersion
	}

	return v, nil
}

func (mg *Migrator) verifyDatabaseVersion() error {
	version, isDirty, err := mg.migrator.Version()
	if err != nil {
		return fmt.Errorf("get current migration version: %w", err)
	}

	latestVersionOnFile, err := mg.latestSourceVersion()
	if err != nil {
		return fmt.Errorf("check latest migration version on file: %w", err)
	}

	switch {
	case version == latestVersionOnFile && isDirty:
		err := mg.forceDatabaseVersion(version)
		if err != nil {
			return err
		}

	case version < latestVersionOnFile:
		err = mg.forceDatabaseVersion(latestVersionOnFile - 1)
		if err != nil {
			return fmt.Errorf("force set migration to latest version on file: %w", err)
		}

	case version == 1:
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

func (mg *Migrator) forceDatabaseVersion(version uint) error {
	mg.logger.Info("force setting migration version", "version", version)
	err := mg.migrator.Force(int(version))
	if err != nil {
		return fmt.Errorf("force set migration version: %w", err)
	}

	return nil
}
