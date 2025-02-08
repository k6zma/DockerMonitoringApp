package repositories

import "context"

type StatusRepository interface {
	UpdateStatus(ctx context.Context, ip string, pingTime int64, name string, status string) error
	CreateStatus(ctx context.Context, ip string, pingTime int64, name string, status string) error
}
