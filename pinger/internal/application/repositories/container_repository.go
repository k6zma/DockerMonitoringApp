package repositories

import (
	"context"

	"github.com/k6zma/DockerMonitoringApp/pinger/internal/domain"
)

type ContainerRepository interface {
	GetContainers(ctx context.Context) ([]domain.ContainerInfo, error)
}
