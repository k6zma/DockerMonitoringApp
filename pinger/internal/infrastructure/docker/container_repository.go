package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	dockerClient "github.com/docker/docker/client"

	"github.com/k6zma/DockerMonitoringApp/pinger/internal/application/repositories"
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

func (r *DockerContainerRepo) GetIPs(ctx context.Context) ([]string, error) {
	containers, err := r.client.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("container list failed: %w", err)
	}

	ips := make([]string, 0, len(containers))
	for i := range containers {
		for _, n := range containers[i].NetworkSettings.Networks {
			ips = append(ips, n.IPAddress)
			break
		}
	}

	return ips, nil
}
