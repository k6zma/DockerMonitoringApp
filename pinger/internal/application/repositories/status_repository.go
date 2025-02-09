package repositories

import (
	"context"

	"github.com/k6zma/DockerMonitoringApp/pinger/internal/domain"
)

type StatusRepository interface {
	UpdateStatus(ctx context.Context, containerID string, pingTime int64, name, status string, success bool) error
	CreateStatus(ctx context.Context, containerID, ip string, pingTime int64, name, status string) error
	DeleteStatus(ctx context.Context, containerID string) error
	GetStatuses(ctx context.Context) ([]domain.PingResult, error)
}
