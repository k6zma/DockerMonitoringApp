package migrations

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file" // import for side effects
	"github.com/jmoiron/sqlx"

	"github.com/k6zma/DockerMonitoringApp/backend/pkg/utils"
)

type Migrations struct {
	db             *sqlx.DB
	migrationsPath string
	logger         utils.LoggerInterface
}

func NewMigrate(db *sqlx.DB, migrationsPath string, logger utils.LoggerInterface) *Migrations {
	return &Migrations{
		db:             db,
		migrationsPath: migrationsPath,
		logger:         logger,
	}
}

func (m *Migrations) ApplyMigrations() error {
	m.logger.Debugf("MIGRATIONS: starting ApplyMigrations with path: %s", m.migrationsPath)

	driver, err := pgx.WithInstance(m.db.DB, &pgx.Config{})
	if err != nil {
		m.logger.Errorf("MIGRATIONS: failed to create migration driver: %v", err)
		return fmt.Errorf("failed to create migration driver: %w", err)
	}
	m.logger.Debug("MIGRATIONS: migration driver created successfully")

	mig, err := migrate.NewWithDatabaseInstance(
		"file://"+m.migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		m.logger.Errorf("MIGRATIONS: failed to create migrate instance: %v", err)
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	m.logger.Debug("MIGRATIONS: migrate instance created successfully")

	if err = mig.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		m.logger.Errorf("MIGRATIONS: failed to apply migrations: %v", err)
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	m.logger.Debug("MIGRATIONS: applied migrations successfully or no change found")

	return nil
}

func (m *Migrations) RollbackMigrations() error {
	m.logger.Debugf("MIGRATIONS: starting RollbackMigrations with path: %s", m.migrationsPath)

	driver, err := pgx.WithInstance(m.db.DB, &pgx.Config{})
	if err != nil {
		m.logger.Errorf("MIGRATIONS: failed to create migration driver: %v", err)
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m.logger.Debug("MIGRATIONS: migration driver created successfully")

	mig, err := migrate.NewWithDatabaseInstance(
		"file://"+m.migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		m.logger.Errorf("MIGRATIONS: failed to create migrate instance: %v", err)
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	m.logger.Debug("MIGRATIONS: migrate instance created successfully")

	if err = mig.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		m.logger.Errorf("MIGRATIONS: failed to rollback migration: %v", err)
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	m.logger.Debug("MIGRATIONS: rolled back migration successfully or no change found")

	return nil
}

func (m *Migrations) DropMigrations() error {
	m.logger.Debugf("MIGRATIONS: starting DropMigrations with path: %s", m.migrationsPath)

	driver, err := pgx.WithInstance(m.db.DB, &pgx.Config{})
	if err != nil {
		m.logger.Errorf("MIGRATIONS: failed to create migration driver: %v", err)
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m.logger.Debug("MIGRATIONS: migration driver created successfully")

	mig, err := migrate.NewWithDatabaseInstance(
		"file://"+m.migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		m.logger.Errorf("MIGRATIONS: failed to create migrate instance: %v", err)
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	m.logger.Debug("MIGRATIONS: migrate instance created successfully")

	if err = mig.Drop(); err != nil {
		m.logger.Errorf("MIGRATIONS: failed to drop all migrations: %v", err)
		return fmt.Errorf("failed to drop all migrations: %w", err)
	}

	m.logger.Debug("MIGRATIONS: dropped all migrations successfully")

	return nil
}
