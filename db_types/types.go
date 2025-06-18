package db_types

import (
	"time"

	errors "github.com/snakehunterr/hacs_app/db_types/errors"
)

type Client struct {
	ID         int64     `json:"client_id"`
	Name       string    `json:"client_name"`
	IsAdmin    bool      `json:"is_admin"`
	LastEdited time.Time `json:"last_edited"`
}

type Expense struct {
	ID         int64     `json:"expense_id"`
	Date       time.Time `json:"expense_date"`
	Amount     float64   `json:"expense_amount"`
	LastEdited time.Time `json:"last_edited"`
}

type Room struct {
	ID          int64     `json:"room_id"`
	ClientID    int64     `json:"client_id"`
	PeopleCount uint8     `json:"room_people_count"`
	Area        float64   `json:"room_area"`
	LastEdited  time.Time `json:"last_edited"`
}

type Payment struct {
	ID         int64     `json:"payment_id"`
	ClientID   int64     `json:"client_id"`
	RoomID     int64     `json:"room_id"`
	Date       time.Time `json:"payment_date"`
	Amount     float64   `json:"payment_amount"`
	LastEdited time.Time `json:"last_edited"`
}

type APIResponse struct {
	Error   *errors.APIError `json:"error,omitempty"`
	Message string           `json:"message,omitempty"`
}
