package usecases

import (
	"fmt"
	"time"

	"github.com/k6zma/DockerMonitoringApp/backend/internal/application/dto"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/application/repositories"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/domain"
	"github.com/k6zma/DockerMonitoringApp/backend/pkg/utils"
)

type ContainerStatusUseCase struct {
	repo repositories.ContainerStatusRepository
}

func NewContainerStatusUseCase(repo repositories.ContainerStatusRepository) *ContainerStatusUseCase {
	return &ContainerStatusUseCase{
		repo: repo,
	}
}

func (uc *ContainerStatusUseCase) FindContainerStatuses(filter *dto.ContainerStatusFilter) ([]*dto.ContainerStatusDTO, error) {
	utils.LoggerInstance.Debugf("USECASES: finding container statuses with filter: %+v", filter)

	statuses, err := uc.repo.Find(filter)
	if err != nil {
		utils.LoggerInstance.Errorf("USECASES: failed to fetch container statuses: %v", err)
		return nil, fmt.Errorf("failed to fetch container statuses: %w", err)
	}

	var dtos = make([]*dto.ContainerStatusDTO, 0, len(statuses))
	for _, status := range statuses {
		dtos = append(dtos, mapDomainToDTO(status))
	}

	utils.LoggerInstance.Debugf("USECASES: found %d container statuses", len(dtos))

	return dtos, nil
}

func (uc *ContainerStatusUseCase) CreateContainerStatus(statusDTO *dto.ContainerStatusDTO) (*dto.ContainerStatusDTO, error) {
	utils.LoggerInstance.Debugf("USECASES: creating container status: %+v", statusDTO)

	newStatus := &domain.ContainerStatus{
		IPAddress:          statusDTO.IPAddress,
		PingTime:           statusDTO.PingTime,
		LastSuccessfulPing: statusDTO.LastSuccessfulPing,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	err := uc.repo.Create(newStatus)
	if err != nil {
		utils.LoggerInstance.Errorf("USECASES: failed to create container status: %v", err)
		return nil, fmt.Errorf("failed to create container status: %w", err)
	}

	utils.LoggerInstance.Debugf("Created container status record")

	return mapDomainToDTO(newStatus), nil
}

func (uc *ContainerStatusUseCase) UpdateContainerStatus(ip string, statusDTO *dto.ContainerStatusDTO) error {
	utils.LoggerInstance.Debugf("USECASES: updating container status for IP: %s with data: %+v", ip, statusDTO)

	existing, err := uc.repo.Find(&dto.ContainerStatusFilter{IPAddress: &ip})
	if err != nil {
		utils.LoggerInstance.Errorf("USECASES: error fetching container status for IP %s: %v", ip, err)
		return fmt.Errorf("error fetching container status: %w", err)
	}

	if len(existing) == 0 {
		utils.LoggerInstance.Errorf("USECASES: error fetching container status with IP %s not found", ip)
		return fmt.Errorf("container status with IP %s not found", ip)
	}

	status := existing[0]

	if statusDTO.PingTime != 0 {
		status.PingTime = statusDTO.PingTime
	}
	if !statusDTO.LastSuccessfulPing.IsZero() {
		status.LastSuccessfulPing = statusDTO.LastSuccessfulPing
	}

	status.UpdatedAt = time.Now()

	err = uc.repo.Update(status)
	if err != nil {
		utils.LoggerInstance.Errorf("USECASES: failed to update container status for IP %s: %v", ip, err)
		return fmt.Errorf("failed to update container status: %w", err)
	}

	utils.LoggerInstance.Debugf("Successfully updated container status for IP: %s", ip)

	return nil
}

func (uc *ContainerStatusUseCase) DeleteContainerStatusByIP(ip string) error {
	utils.LoggerInstance.Debugf("USECASES: deleting container status for IP: %s", ip)

	existing, err := uc.repo.Find(&dto.ContainerStatusFilter{IPAddress: &ip})
	if err != nil {
		utils.LoggerInstance.Errorf("USECASES: error checking container status for IP %s: %v", ip, err)
		return fmt.Errorf("error checking container status: %w", err)
	}

	if len(existing) == 0 {
		utils.LoggerInstance.Warnf("USECASES: attempted to delete non-existent container status for IP: %s", ip)
		return fmt.Errorf("container status with IP %s not found", ip)
	}

	err = uc.repo.DeleteByIP(ip)
	if err != nil {
		utils.LoggerInstance.Errorf("USECASES: failed to delete container status for IP %s: %v", ip, err)
		return fmt.Errorf("failed to delete container status: %w", err)
	}

	utils.LoggerInstance.Debugf("USECASES: successfully deleted container status for IP: %s", ip)

	return nil
}

func mapDomainToDTO(status *domain.ContainerStatus) *dto.ContainerStatusDTO {
	return &dto.ContainerStatusDTO{
		ID:                 status.ID,
		IPAddress:          status.IPAddress,
		PingTime:           status.PingTime,
		LastSuccessfulPing: status.LastSuccessfulPing,
		UpdatedAt:          status.UpdatedAt,
		CreatedAt:          status.CreatedAt,
	}
}
