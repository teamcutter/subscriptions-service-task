package model

import "github.com/google/uuid"

type Subscription struct {
	ID int64 `json:"-" db:"id"`
	ServiceName string `json:"service_name" db:"service_name"`
	Price int `json:"price" db:"price"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	StartDate string `json:"start_date" db:"start_date"`
	EndDate string `json:"end_date,omitempty" db:"end_date"`
	CreatedAt string `json:"-" db:"created_at"`
}