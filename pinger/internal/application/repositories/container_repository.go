package repositories

import "context"

type ContainerRepository interface {
	GetIPs(ctx context.Context) ([]string, error)
}
