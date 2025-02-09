package usecases

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	probing "github.com/prometheus-community/pro-bing"

	"github.com/k6zma/DockerMonitoringApp/pinger/internal/application/repositories"
	"github.com/k6zma/DockerMonitoringApp/pinger/internal/domain"
	"github.com/k6zma/DockerMonitoringApp/pinger/pkg/utils"
)

type PingerUsecase struct {
	containerRepo repositories.ContainerRepository
	statusRepo    repositories.StatusRepository
	interval      time.Duration
	logger        utils.LoggerInterface
}

func NewPingerUsecase(
	cr repositories.ContainerRepository,
	sr repositories.StatusRepository,
	inter time.Duration,
	logger utils.LoggerInterface,
) *PingerUsecase {
	return &PingerUsecase{
		containerRepo: cr,
		statusRepo:    sr,
		interval:      inter,
		logger:        logger,
	}
}

func (uc *PingerUsecase) Run(ctx context.Context) error {
	uc.logger.Infof("Starting monitoring with interval %v", uc.interval)

	uc.logger.Debugf("Ticker interval: %v", uc.interval)
	ticker := time.NewTicker(uc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			uc.logger.Info("Shutting down pinger service")
			return nil
		case <-ticker.C:
			if err := uc.checkContainers(ctx); err != nil {
				uc.logger.Errorf("Monitoring cycle failed: %v", err)
			}
		}
	}
}

func (uc *PingerUsecase) checkContainers(ctx context.Context) error {
	containers, err := uc.containerRepo.GetContainers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get container info: %w", err)
	}

	activeContainerIDs := make(map[string]bool)
	containerInfos := make([]string, 0, len(containers))
	for _, container := range containers {
		activeContainerIDs[container.ContainerID] = true
		containerInfos = append(containerInfos,
			fmt.Sprintf("%s (ID: %s, IP: %s) [%s]", container.Name, container.ContainerID, container.IP, container.Status))
	}
	uc.logger.Debugf("Discovered %d containers: %s", len(containers), strings.Join(containerInfos, ", "))

	uc.logger.Debug("Pinging containers")
	var wg sync.WaitGroup

	for _, container := range containers {
		wg.Add(1)
		go func(container domain.ContainerInfo) {
			defer wg.Done()

			var result *domain.PingResult
			if container.IP == "" {
				uc.logger.Warnf("No IP for container %s (ID: %s), updating status as %s", container.Name, container.ContainerID, container.Status)
				result = &domain.PingResult{
					ContainerID: container.ContainerID,
					IP:          "",
					Name:        container.Name,
					Status:      container.Status,
					Success:     false,
					PingTime:    0,
					LastPing:    time.Now().Format(time.RFC3339),
				}
			} else {
				res, err := uc.ping(container)
				if err != nil {
					uc.logger.Warnf("Ping failed for container %s (ID: %s, IP: %s) [%s]: %v",
						container.Name, container.ContainerID, container.IP, container.Status, err)
					res = &domain.PingResult{
						ContainerID: container.ContainerID,
						IP:          container.IP,
						Name:        container.Name,
						Status:      container.Status,
						Success:     false,
						PingTime:    0,
						LastPing:    time.Now().Format(time.RFC3339),
					}
				}
				result = res
			}

			if err := uc.updateStatus(ctx, result); err != nil {
				uc.logger.Errorf("Failed to update status for container %s (ID: %s, IP: %s) [%s]: %v",
					container.Name, container.ContainerID, container.IP, container.Status, err)
			}
		}(container)
	}
	wg.Wait()

	if err := uc.cleanupStatuses(ctx, activeContainerIDs); err != nil {
		uc.logger.Errorf("Cleanup statuses failed: %v", err)
		return fmt.Errorf("cleanup statuses failed: %w", err)
	}

	return nil
}

func (uc *PingerUsecase) ping(container domain.ContainerInfo) (*domain.PingResult, error) {
	uc.logger.Debugf("Pinging container %s (ID: %s, IP: %s) [%s]",
		container.Name, container.ContainerID, container.IP, container.Status)

	pinger, err := probing.NewPinger(container.IP)
	if err != nil {
		uc.logger.Errorf("Ping init failed for container %s (ID: %s, IP: %s) [%s]: %v",
			container.Name, container.ContainerID, container.IP, container.Status, err)
		return nil, fmt.Errorf("ping init failed: %w", err)
	}

	pinger.Count = 100
	pinger.Timeout = 2 * time.Second
	pinger.SetPrivileged(true)

	if err := pinger.Run(); err != nil {
		uc.logger.Errorf("Ping execution failed for container %s (ID: %s, IP: %s) [%s]: %v",
			container.Name, container.ContainerID, container.IP, container.Status, err)
		return nil, fmt.Errorf("ping execution failed: %w", err)
	}

	stats := pinger.Statistics()
	uc.logger.Debugf("Ping stats for container %s (ID: %s, IP: %s) [%s]: %+v",
		container.Name, container.ContainerID, container.IP, container.Status, stats)

	var pingTime int64 = -1
	if stats.PacketsRecv > 0 {
		pingTime = stats.AvgRtt.Microseconds()
	}

	uc.logger.Debugf("Ping time for container %s (ID: %s, IP: %s) [%s]: %.2f ms",
		container.Name, container.ContainerID, container.IP, container.Status, pingTime)

	return &domain.PingResult{
		ContainerID: container.ContainerID,
		IP:          container.IP,
		Name:        container.Name,
		Status:      container.Status,
		Success:     stats.PacketsRecv > 0,
		PingTime:    pingTime,
	}, nil
}

func (uc *PingerUsecase) updateStatus(ctx context.Context, result *domain.PingResult) error {
	if err := uc.statusRepo.UpdateStatus(ctx, result.ContainerID, result.PingTime, result.Name, result.Status, result.Success); err != nil {
		uc.logger.Warnf("Update failed for container %s (ID: %s, IP: %s) [%s], trying to create: %v",
			result.Name, result.ContainerID, result.IP, result.Status, err)

		if err := uc.statusRepo.CreateStatus(ctx, result.ContainerID, result.IP, result.PingTime, result.Name, result.Status); err != nil {
			uc.logger.Errorf("Create status failed for container %s (ID: %s, IP: %s) [%s]: %v",
				result.Name, result.ContainerID, result.IP, result.Status, err)
			return fmt.Errorf("create status failed for container %s (ID: %s, IP: %s) [%s]: %w",
				result.Name, result.ContainerID, result.IP, result.Status, err)
		}
	}

	return nil
}

func (uc *PingerUsecase) cleanupStatuses(ctx context.Context, activeContainerIDs map[string]bool) error {
	uc.logger.Debug("Cleaning up statuses")
	statuses, err := uc.statusRepo.GetStatuses(ctx)
	if err != nil {
		uc.logger.Errorf("Failed to get statuses: %v", err)
		return fmt.Errorf("failed to get statuses: %w", err)
	}

	for _, status := range statuses {
		uc.logger.Debugf("Status: %+v", status)
	}

	for _, status := range statuses {
		if status.ContainerID == "" {
			uc.logger.Debugf("Skipping deletion for container %s with empty container_id", status.Name)
			continue
		}

		if !activeContainerIDs[status.ContainerID] {
			uc.logger.Debugf("Container with container_id %s not found among active containers. Deleting its record.", status.ContainerID)
			if err := uc.statusRepo.DeleteStatus(ctx, status.ContainerID); err != nil {
				uc.logger.Errorf("Failed to delete status for container_id %s: %v", status.ContainerID, err)
				return fmt.Errorf("failed to delete status for container_id %s: %w", status.ContainerID, err)
			} else {
				uc.logger.Debugf("Successfully deleted status for container %s with container_id %s", status.Name, status.ContainerID)
			}
		}
	}

	return nil
}
