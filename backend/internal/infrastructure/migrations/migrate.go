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

func ApplyMigrations(db *sqlx.DB, migrationsPath string) error {
	utils.LoggerInstance.Debugf("MIGRATIONS: starting ApplyMigrations with path: %s", migrationsPath)
	driver, err := pgx.WithInstance(db.DB, &pgx.Config{})
	if err != nil {
		utils.LoggerInstance.Errorf("MIGRATIONS: failed to create migration driver: %v", err)
		return fmt.Errorf("failed to create migration driver: %w", err)
	}
	utils.LoggerInstance.Debug("MIGRATIONS: migration driver created successfully")

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		utils.LoggerInstance.Errorf("MIGRATIONS: failed to create migrate instance: %v", err)
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	utils.LoggerInstance.Debug("MIGRATIONS: migrate instance created successfully")

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		utils.LoggerInstance.Errorf("MIGRATIONS: failed to apply migrations: %v", err)
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	utils.LoggerInstance.Debug("MIGRATIONS: applied migrations successfully or no change found")

	return nil
}

func RollbackMigrations(db *sqlx.DB, migrationsPath string) error {
	utils.LoggerInstance.Debugf("MIGRATIONS: starting RollbackMigrations with path: %s", migrationsPath)
	driver, err := pgx.WithInstance(db.DB, &pgx.Config{})
	if err != nil {
		utils.LoggerInstance.Errorf("MIGRATIONS: failed to create migration driver: %v", err)
		return fmt.Errorf("failed to create migration driver: %w", err)
	}
	utils.LoggerInstance.Debug("MIGRATIONS: migration driver created successfully")

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		utils.LoggerInstance.Errorf("MIGRATIONS: failed to create migrate instance: %v", err)
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	utils.LoggerInstance.Debug("MIGRATIONS: migrate instance created successfully")

	if err = m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		utils.LoggerInstance.Errorf("MIGRATIONS: failed to rollback migration: %v", err)
		return fmt.Errorf("failed to rollback migration: %w", err)
	}
	utils.LoggerInstance.Debug("MIGRATIONS: rolled back migration successfully or no change found")

	return nil
}

func DropMigrations(db *sqlx.DB, migrationsPath string) error {
	utils.LoggerInstance.Debugf("MIGRATIONS: starting DropMigrations with path: %s", migrationsPath)
	driver, err := pgx.WithInstance(db.DB, &pgx.Config{})
	if err != nil {
		utils.LoggerInstance.Errorf("MIGRATIONS: failed to create migration driver: %v", err)
		return fmt.Errorf("failed to create migration driver: %w", err)
	}
	utils.LoggerInstance.Debug("MIGRATIONS: migration driver created successfully")

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		utils.LoggerInstance.Errorf("MIGRATIONS: failed to create migrate instance: %v", err)
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	utils.LoggerInstance.Debug("MIGRATIONS: migrate instance created successfully")

	if err = m.Drop(); err != nil {
		utils.LoggerInstance.Errorf("MIGRATIONS: failed to drop all migrations: %v", err)
		return fmt.Errorf("failed to drop all migrations: %w", err)
	}
	utils.LoggerInstance.Debug("MIGRATIONS: dropped all migrations successfully")

	return nil
}
