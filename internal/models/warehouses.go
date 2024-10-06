package models

import (
	"time"

	"github.com/google/uuid"
)

type Warehouse struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"user_id"`
	Name    string    `json:"name"`
	Address struct {
		AddressLine1 string `json:"address_line_1"`
		AddressLine2 string `json:"address_line_2,omitempty"`
		TownCity     string `json:"town_city"`
		StateCounty  string `json:"state_county"`
		PostZipCode  string `json:"post_zip_code"`
		Country      string `json:"country"`
	} `json:"address"`
	Latitude  float32   `json:"latitude"`
	Longitude float32   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WarehouseDatabase struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Name         string    `json:"name"`
	AddressLine1 string    `json:"address_line_1"`
	AddressLine2 string    `json:"address_line_2,omitempty"`
	TownCity     string    `json:"town_city"`
	StateCounty  string    `json:"state_county"`
	PostZipCode  string    `json:"post_zip_code"`
	Country      string    `json:"country"`
	Latitude     float32   `json:"latitude"`
	Longitude    float32   `json:"longitude"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
