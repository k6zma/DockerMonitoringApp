package dto

import "time"

type GetContainerStatusResponse struct {
	ID                 int64     `json:"id"`
	IPAddress          string    `json:"ip_address"`
	PingTime           float64   `json:"ping_time"`
	LastSuccessfulPing time.Time `json:"last_successful_ping"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type DeleteContainerStatusResponse struct {
	Message string `json:"message"`
}
