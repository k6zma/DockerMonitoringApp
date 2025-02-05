package dto

import "time"

type CreateContainerStatusRequest struct {
	IPAddress          string    `json:"ip_address" validate:"required,ip"`
	PingTime           float64   `json:"ping_time" validate:"required,min=1"`
	LastSuccessfulPing time.Time `json:"last_successful_ping" validate:"required"`
}

type UpdateContainerStatusRequest struct {
	PingTime           float64   `json:"ping_time,omitempty"`
	LastSuccessfulPing time.Time `json:"last_successful_ping,omitempty"`
}
