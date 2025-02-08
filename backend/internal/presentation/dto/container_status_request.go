package dto

import "time"

type CreateContainerStatusRequest struct {
	IPAddress          string    `json:"ip_address"           validate:"required,ip"`
	Name               string    `json:"name"`
	Status             string    `json:"status"               validate:"required,oneof=created restarting running removing paused exited dead"`
	PingTime           float64   `json:"ping_time"`
	LastSuccessfulPing time.Time `json:"last_successful_ping" validate:"required"`
}

type UpdateContainerStatusRequest struct {
	Name               string    `json:"name"`
	Status             string    `json:"status"               validate:"omitempty,oneof=created restarting running removing paused exited dead"`
	PingTime           float64   `json:"ping_time"`
	LastSuccessfulPing time.Time `json:"last_successful_ping,omitempty"`
}
