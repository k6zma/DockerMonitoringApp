package repositories

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/k6zma/DockerMonitoringApp/backend/internal/application/dto"
	appRepo "github.com/k6zma/DockerMonitoringApp/backend/internal/application/repositories"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/domain"
	"github.com/k6zma/DockerMonitoringApp/backend/pkg/utils"
)

type ContainerStatusRepositoryImpl struct {
	db     *sqlx.DB
	logger utils.LoggerInterface
}

func NewContainerStatusRepositoryImpl(db *sqlx.DB, logger utils.LoggerInterface) appRepo.ContainerStatusRepository {
	return &ContainerStatusRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *ContainerStatusRepositoryImpl) Find(filter *dto.ContainerStatusFilter) ([]*domain.ContainerStatus, error) {
	r.logger.Debugf("REPOSITORIES: executing Find with filter: %+v", *filter)

	query := `
		SELECT id, ip_address, ping_time, last_successful_ping, created_at, updated_at
		FROM container_status
	`

	var conditions []string
	var args []interface{}
	argCounter := 1

	if filter.IPAddress != nil {
		conditions = append(conditions, fmt.Sprintf("ip_address = $%d", argCounter))
		args = append(args, filter.IPAddress)
		argCounter++
	}

	if filter.ID != nil {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argCounter))
		args = append(args, filter.ID)
		argCounter++
	}

	if filter.PingTimeMin != nil {
		conditions = append(conditions, fmt.Sprintf("ping_time >= $%d", argCounter))
		args = append(args, *filter.PingTimeMin)
		argCounter++
	}

	if filter.PingTimeMax != nil {
		conditions = append(conditions, fmt.Sprintf("ping_time <= $%d", argCounter))
		args = append(args, *filter.PingTimeMax)
		argCounter++
	}

	if filter.CreatedAtGte != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argCounter))
		args = append(args, *filter.CreatedAtGte)
		argCounter++
	}

	if filter.CreatedAtLte != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argCounter))
		args = append(args, *filter.CreatedAtLte)
		argCounter++
	}

	if filter.UpdatedAtGte != nil {
		conditions = append(conditions, fmt.Sprintf("updated_at >= $%d", argCounter))
		args = append(args, *filter.UpdatedAtGte)
		argCounter++
	}

	if filter.UpdatedAtLte != nil {
		conditions = append(conditions, fmt.Sprintf("updated_at <= $%d", argCounter))
		args = append(args, *filter.UpdatedAtLte)
		argCounter++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	if filter.Limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", argCounter)
		args = append(args, *filter.Limit)
	}

	r.logger.Debugf("REPOSITORIES: final Query: %s, Args: %+v", query, args)

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		r.logger.Errorf("REPOSITORIES: failed to execute query: %v\n", err)
		return nil, fmt.Errorf("database query error: %w", err)
	}

	var results []*domain.ContainerStatus
	for rows.Next() {
		var status domain.ContainerStatus
		var pingTime float64

		err := rows.Scan(&status.ID, &status.IPAddress, &pingTime, &status.LastSuccessfulPing, &status.CreatedAt, &status.UpdatedAt)
		if err != nil {
			r.logger.Errorf("REPOSITORIES: failed to scan row: %v\n", err)
			return nil, fmt.Errorf("database scan error: %w", err)
		}

		status.PingTime = pingTime
		results = append(results, &status)
	}

	r.logger.Debugf("REPOSITORIES: query executed successfully, found %d records", len(results))

	return results, nil
}

func (r *ContainerStatusRepositoryImpl) Create(status *domain.ContainerStatus) error {
	r.logger.Debugf("REPOSITORIES: creating container status record: %+v", status)

	query := `
		INSERT INTO container_status (ip_address, ping_time, last_successful_ping, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.db.QueryRowx(
		query,
		status.IPAddress,
		status.PingTime,
		status.LastSuccessfulPing,
		status.CreatedAt,
		status.UpdatedAt,
	).Scan(&status.ID)
	if err != nil {
		r.logger.Errorf("REPOSITORIES: failed to create container status: %v", err)
		return fmt.Errorf("failed to create container status: %w", err)
	}

	r.logger.Debugf("REPOSITORIES: container status created with ID: %d", status.ID)

	return nil
}

func (r *ContainerStatusRepositoryImpl) Update(status *domain.ContainerStatus) error {
	r.logger.Debugf("REPOSITORIES: updating container status record for IP: %s", status.IPAddress)

	query := `
		UPDATE container_status
		SET ping_time = $1, last_successful_ping = $2, updated_at = $3
		WHERE ip_address = $4
	`

	_, err := r.db.Exec(
		query,
		status.PingTime,
		status.LastSuccessfulPing,
		status.UpdatedAt,
		status.IPAddress,
	)
	if err != nil {
		r.logger.Errorf("REPOSITORIES: failed to update container status for IP %s: %v", status.IPAddress, err)
		return fmt.Errorf("failed to update container status: %w", err)
	}

	r.logger.Debugf("REPOSITORIES: container status for IP %s updated successfully", status.IPAddress)

	return nil
}

func (r *ContainerStatusRepositoryImpl) DeleteByIP(ipAddress string) error {
	r.logger.Debugf("REPOSITORIES: deleting container status record for IP: %s", ipAddress)

	query := `
		DELETE FROM container_status
		WHERE ip_address = $1
	`

	_, err := r.db.Exec(query, ipAddress)
	if err != nil {
		r.logger.Errorf("REPOSITORIES: failed to delete container status for IP %s: %v", ipAddress, err)
		return fmt.Errorf("failed to delete container status: %w", err)
	}

	r.logger.Debugf("REPOSITORIES: container status for IP %s deleted successfully", ipAddress)

	return nil
}
