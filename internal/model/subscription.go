package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	Id          int       `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserId      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
}

type SubscriptionDTO struct {
	Id          int       `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserId      uuid.UUID `json:"user_id"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type CostRequest struct {
	StartDate   string
	EndDate     string
	UserId      uuid.UUID
	ServiceName string
}

type CostDTO struct {
	StartDate   time.Time
	EndDate     time.Time
	UserId      uuid.UUID
	ServiceName string
}

type SubRequest struct {
	ServiceName string    `json:"service_name" example:"Yandex Plus"`
	Price       int       `json:"price" example:"400"`
	UserId      uuid.UUID `json:"user_id" example:"UUID"`
	StartDate   string    `json:"start_date" example:"01-2025"`
	EndDate     string    `json:"end_date" example:"02-2025"`
}
