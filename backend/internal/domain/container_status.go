package domain

import "time"

type ContainerStatus struct {
	ID                 int64     `db:"id"`
	IPAddress          string    `db:"ip_address"`
	PingTime           float64   `db:"ping_time"`
	LastSuccessfulPing time.Time `db:"last_successful_ping"`
	UpdatedAt          time.Time `db:"updated_at"`
	CreatedAt          time.Time `db:"created_at"`
}
