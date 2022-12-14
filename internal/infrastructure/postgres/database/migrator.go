package database

import (
	"database/sql"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // must be imported this way (linter want's a comment)
)

type Migrator struct {
	postgresMigrator *migrate.Migrate
}

func NewMigrator(postgresDBConn *sql.DB, migrationsPath string) (*Migrator, error) {
	migrator := &Migrator{}
	if err := migrator.configure(postgresDBConn, migrationsPath); err != nil {
		return nil, errors.Wrap(err, "NewMigrator")
	}

	return migrator, nil
}

func (migrator *Migrator) Up() error {
	if err := migrator.postgresMigrator.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return errors.Wrap(err, "migrator.Up: failed to run migrations for Postgres DB")
		}
	}

	return nil
}

func (migrator *Migrator) WithLogger(logger migrate.Logger) *Migrator {
	migrator.postgresMigrator.Log = logger

	return migrator
}

func (migrator *Migrator) configure(postgresDBConn *sql.DB, migrationsPath string) error {
	config := &postgres.Config{MigrationsTable: "student_migrations"}

	driver, err := postgres.WithInstance(postgresDBConn, config)
	if err != nil {
		return errors.Wrap(err, "failed to create Postgres driver for migrator")
	}

	sourceURL := fmt.Sprintf("file://%s", migrationsPath)
	realMigrator, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		return errors.Wrap(err, "failed to create migrator instance")
	}

	migrator.postgresMigrator = realMigrator

	return nil
}
