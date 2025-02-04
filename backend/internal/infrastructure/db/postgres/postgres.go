package postgres

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver

	"github.com/jmoiron/sqlx"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/infrastructure/config"
)

func NewPsqlDB(cfg *config.DBConfig) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DataBaseName,
		cfg.Password,
	)

	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return db, nil
}
