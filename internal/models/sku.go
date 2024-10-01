package models

import (
	"time"

	"github.com/google/uuid"
)

type SKU struct {
	ID        uuid.UUID `json:"id"`
	ProductID uint      `json:"product_id"`
	SKU       string    `json:"sku"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
