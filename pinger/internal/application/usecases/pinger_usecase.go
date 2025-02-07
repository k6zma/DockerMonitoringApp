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
	ips, err := uc.containerRepo.GetIPs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get container IPs: %w", err)
	}
	uc.logger.Debugf("Discovered %d containers: %s", len(ips), strings.Join(ips, ", "))

	uc.logger.Debug("Pinging containers")
	for _, ip := range ips {
		result, err := uc.ping(ip)

		uc.logger.Debugf("Ping result for %s: %+v", ip, result)
		if err != nil {
			uc.logger.Warnf("Ping failed for %s: %v", ip, err)
			continue
		}

		uc.logger.Debugf("Ping result for %s: %.2f ms", ip, result.PingTime)
		if err := uc.updateStatus(ctx, ip, result); err != nil {
			uc.logger.Errorf("Failed to update status for %s: %v", ip, err)
		}
	}

	return nil
}

func (uc *PingerUsecase) ping(ip string) (*domain.PingResult, error) {
	uc.logger.Debugf("Pinging %s", ip)

	pinger, err := probing.NewPinger(ip)
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
	uc.logger.Debugf("Ping stats for %s: %+v", ip, stats)

	var pingTime float64
	if stats.PacketsRecv > 0 {
		pingTime = float64(stats.AvgRtt.Milliseconds())
	}

	uc.logger.Debugf("Ping time for %s: %.2f ms", ip, pingTime)

	return &domain.PingResult{
		IP:       ip,
		Success:  stats.PacketsRecv > 0,
		PingTime: pingTime,
	}, nil
}

func (uc *PingerUsecase) updateStatus(
	ctx context.Context,
	ip string,
	result *domain.PingResult,
) error {
	if err := uc.statusRepo.UpdateStatus(ctx, ip, result.PingTime); err != nil {
		uc.logger.Warnf("Update failed for %s, trying to create: %v", ip, err)

		if err := uc.statusRepo.CreateStatus(ctx, ip, result.PingTime); err != nil {
			return fmt.Errorf("create status failed: %w", err)
		}
	}

	return nil
}
