package repositories

import (
	"context"

	"github.com/k6zma/DockerMonitoringApp/pinger/internal/domain"
)

type StatusRepository interface {
	UpdateStatus(ctx context.Context, ip string, pingTime int64, name string, status string) error
	CreateStatus(ctx context.Context, ip string, pingTime int64, name string, status string) error
	DeleteStatus(ctx context.Context, ip string) error
	GetStatuses(ctx context.Context) ([]domain.PingResult, error)
}
