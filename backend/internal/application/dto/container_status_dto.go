package dto

import "time"

type ContainerStatusDTO struct {
	ID                 int64
	Name               string
	IPAddress          string
	Status             string
	PingTime           float64
	LastSuccessfulPing time.Time
	UpdatedAt          time.Time
	CreatedAt          time.Time
}

type ContainerStatusFilter struct {
	IPAddress    *string
	ID           *int64
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
