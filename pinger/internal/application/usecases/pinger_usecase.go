package usecases

import (
	"context"
	"fmt"
	"strings"
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

	uc.logger.Debug("Tickers started")
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
	uc.logger.Debug("Checking containers")

	containers, err := uc.containerRepo.GetContainers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get container info: %w", err)
	}

	containerInfos := make([]string, 0, len(containers))
	for _, container := range containers {
		containerInfos = append(containerInfos, fmt.Sprintf("%s (%s) [%s]", container.Name, container.IP, container.Status))
	}
	uc.logger.Debugf("Discovered %d containers: %s", len(containers), strings.Join(containerInfos, ", "))

	uc.logger.Debug("Pinging containers")
	for _, container := range containers {
		result, err := uc.ping(container)
		if err != nil {
			uc.logger.Warnf("Ping failed for container %s (%s) [%s]: %v", container.Name, container.IP, container.Status, err)
			continue
		}

		uc.logger.Debugf("Ping successful for container %s (%s) [%s]: %.2f ms", container.Name, container.IP, container.Status, result.PingTime)

		result.Name = container.Name
		result.Status = container.Status

		if err := uc.updateStatus(ctx, result); err != nil {
			uc.logger.Errorf("Failed to update status for container %s (%s) [%s]: %v", container.Name, container.IP, container.Status, err)
		}
	}

	return nil
}

func (uc *PingerUsecase) ping(container domain.ContainerInfo) (*domain.PingResult, error) {
	uc.logger.Debugf("Pinging container %s (%s) [%s]", container.Name, container.IP, container.Status)

	pinger, err := probing.NewPinger(container.IP)
	if err != nil {
		return nil, fmt.Errorf("ping init failed: %w", err)
	}

	pinger.Count = 3
	pinger.Timeout = 5 * time.Second
	pinger.SetPrivileged(true)

	if err := pinger.Run(); err != nil {
		return nil, fmt.Errorf("ping execution failed: %w", err)
	}

	stats := pinger.Statistics()
	uc.logger.Debugf("Ping stats for container %s (%s) [%s]: %+v", container.Name, container.IP, container.Status, stats)

	var pingTime int64 = -1
	if stats.PacketsRecv > 0 {
		pingTime = stats.AvgRtt.Microseconds()
	}

	uc.logger.Debugf("Ping time for container %s (%s) [%s]: %.2f ms", container.Name, container.IP, container.Status, pingTime)

	return &domain.PingResult{
		IP:       container.IP,
		Name:     container.Name,
		Status:   container.Status,
		Success:  stats.PacketsRecv > 0,
		PingTime: pingTime,
	}, nil
}

func (uc *PingerUsecase) updateStatus(ctx context.Context, result *domain.PingResult) error {
	if err := uc.statusRepo.UpdateStatus(ctx, result.IP, result.PingTime, result.Name, result.Status); err != nil {
		uc.logger.Warnf("Update failed for container %s (%s) [%s], trying to create: %v",
			result.Name, result.IP, result.Status, err)

		if err := uc.statusRepo.CreateStatus(ctx, result.IP, result.PingTime, result.Name, result.Status); err != nil {
			return fmt.Errorf("create status failed for container %s (%s) [%s]: %w",
				result.Name, result.IP, result.Status, err)
		}
	}

	return nil
}
