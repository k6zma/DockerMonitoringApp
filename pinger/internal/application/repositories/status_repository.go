package repositories

import "context"

type StatusRepository interface {
	UpdateStatus(ctx context.Context, ip string, pingTime float64) error
	CreateStatus(ctx context.Context, ip string, pingTime float64) error
}
