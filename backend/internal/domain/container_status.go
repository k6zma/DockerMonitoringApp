package domain

import "time"

type ContainerStatus struct {
	ContainerID        string    `db:"container_id"`
	Name               string    `db:"name"`
	IPAddress          string    `db:"ip_address"`
	Status             string    `db:"status"`
	PingTime           float64   `db:"ping_time"`
	LastSuccessfulPing time.Time `db:"last_successful_ping"`
	UpdatedAt          time.Time `db:"updated_at"`
	CreatedAt          time.Time `db:"created_at"`
}
