package dto

import "time"

type ContainerStatusDTO struct {
	ContainerID        string
	Name               string
	IPAddress          string
	Status             string
	PingTime           float64
	LastSuccessfulPing time.Time
	UpdatedAt          time.Time
	CreatedAt          time.Time
}

type ContainerStatusFilter struct {
	ContainerID  *string
	IPAddress    *string
	Name         *string
	Status       *string
	PingTimeMin  *float64
	PingTimeMax  *float64
	CreatedAtGte *time.Time
	CreatedAtLte *time.Time
	UpdatedAtGte *time.Time
	UpdatedAtLte *time.Time
	Limit        *int
}
