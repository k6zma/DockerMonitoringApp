package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	dockerClient "github.com/docker/docker/client"

	"github.com/k6zma/DockerMonitoringApp/pinger/internal/application/repositories"
	"github.com/k6zma/DockerMonitoringApp/pinger/internal/domain"
	"github.com/k6zma/DockerMonitoringApp/pinger/internal/infrastructure/config"
	"github.com/k6zma/DockerMonitoringApp/pinger/pkg/utils"
)

type DockerContainerRepo struct {
	client *dockerClient.Client
	logger utils.LoggerInterface
}

func NewDockerContainerRepo(
	cfg *config.Config,
	logger utils.LoggerInterface,
) (repositories.ContainerRepository, error) {
	client, err := dockerClient.NewClientWithOpts(
		dockerClient.WithHost("unix://"+cfg.Docker.SocketPath),
		dockerClient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("docker client init failed: %w", err)
	}

	return &DockerContainerRepo{client: client, logger: logger}, nil
}

func (r *DockerContainerRepo) GetContainers(ctx context.Context) ([]domain.ContainerInfo, error) {
	r.logger.Debug("Getting containers list")
	containers, err := r.client.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		r.logger.Errorf("Container list failed: %v", err)
		return nil, fmt.Errorf("container list failed: %w", err)
	}

	for i := range containers {
		r.logger.Debugf("Found container: %s with IP: %s", containers[i].Names[0], containers[i].NetworkSettings.Networks)
	}

	containerList := make([]domain.ContainerInfo, 0, len(containers))
	for i := range containers {
		var ip string
		for _, n := range containers[i].NetworkSettings.Networks {
			ip = n.IPAddress
			break
		}

		containerList = append(containerList, domain.ContainerInfo{
			IP:     ip,
			Name:   containers[i].Names[0],
			Status: containers[i].State,
		})
	}

	return containerList, nil
}
