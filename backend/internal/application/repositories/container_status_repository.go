package repositories

import (
	"github.com/k6zma/DockerMonitoringApp/backend/internal/application/dto"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/domain"
)

type ContainerStatusRepository interface {
	Find(filter *dto.ContainerStatusFilter) ([]*domain.ContainerStatus, error)
	Create(status *domain.ContainerStatus) error
	Update(status *domain.ContainerStatus) error
	DeleteByIP(ipAddress string) error
}
