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

func NewContainerStatusRepositoryImpl(
	db *sqlx.DB,
	logger utils.LoggerInterface,
) appRepo.ContainerStatusRepository {
	return &ContainerStatusRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *ContainerStatusRepositoryImpl) Find(
	filter *dto.ContainerStatusFilter,
) ([]*domain.ContainerStatus, error) {
	r.logger.Debugf("REPOSITORIES: executing Find with filter: %+v", *filter)

	query := `
		SELECT container_id, ip_address, name, status, ping_time, last_successful_ping, created_at, updated_at
		FROM container_status
	`

	var conditions []string
	var args []interface{}
	argCounter := 1

	if filter.IPAddress != nil {
		conditions = append(conditions, fmt.Sprintf("ip_address = $%d", argCounter))
		args = append(args, *filter.IPAddress)
		argCounter++
	}

	if filter.ContainerID != nil {
		conditions = append(conditions, fmt.Sprintf("container_id = $%d", argCounter))
		args = append(args, *filter.ContainerID)
		argCounter++
	}

	if filter.Name != nil {
		conditions = append(conditions, fmt.Sprintf("name = $%d", argCounter))
		args = append(args, *filter.Name)
		argCounter++
	}

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCounter))
		args = append(args, *filter.Status)
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

		err := rows.Scan(
			&status.ContainerID,
			&status.IPAddress,
			&status.Name,
			&status.Status,
			&pingTime,
			&status.LastSuccessfulPing,
			&status.CreatedAt,
			&status.UpdatedAt,
		)
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
		INSERT INTO container_status (container_id, ip_address, name, status, ping_time, last_successful_ping, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING container_id
	`

	err := r.db.QueryRowx(query,
		status.ContainerID,
		status.IPAddress,
		status.Name,
		status.Status,
		status.PingTime,
		status.LastSuccessfulPing,
		status.CreatedAt,
		status.UpdatedAt,
	).Scan(&status.ContainerID)
	if err != nil {
		r.logger.Errorf("REPOSITORIES: failed to create container status: %v", err)
		return fmt.Errorf("failed to create container status: %w", err)
	}

	r.logger.Debugf("REPOSITORIES: container status created with ID: %d", status.ContainerID)

	return nil
}

func (r *ContainerStatusRepositoryImpl) Update(status *domain.ContainerStatus) error {
	r.logger.Debugf("REPOSITORIES: updating container status record for ID: %s, IP: %s", status.ContainerID, status.IPAddress)

	query := `
		UPDATE container_status
		SET name = $1, status = $2, ping_time = $3, last_successful_ping = $4, updated_at = $5, ip_address = $6
		WHERE container_id = $7
	`

	_, err := r.db.Exec(query,
		status.Name,
		status.Status,
		status.PingTime,
		status.LastSuccessfulPing,
		status.UpdatedAt,
		status.IPAddress,
		status.ContainerID,
	)
	if err != nil {
		r.logger.Errorf(
			"REPOSITORIES: failed to update container status for ID %s, IP %s: %v",
			status.ContainerID,
			status.IPAddress,
			err,
		)
		return fmt.Errorf("failed to update container status: %w", err)
	}

	r.logger.Debugf(
		"REPOSITORIES: container status for ID %s, IP %s updated successfully",
		status.ContainerID,
		status.IPAddress,
	)

	return nil
}

func (r *ContainerStatusRepositoryImpl) DeleteByContainerID(containerID string) error {
	r.logger.Debugf("REPOSITORIES: deleting container status record for container id: %s", containerID)

	query := `
		DELETE FROM container_status
		WHERE container_id = $1
	`

	_, err := r.db.Exec(query, containerID)
	if err != nil {
		r.logger.Errorf(
			"REPOSITORIES: failed to delete container status for container id %s: %v",
			containerID,
			err,
		)
		return fmt.Errorf("failed to delete container status: %w", err)
	}

	r.logger.Debugf("REPOSITORIES: container status for container id %s deleted successfully", containerID)

	return nil
}
