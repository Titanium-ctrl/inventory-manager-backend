package models

import (
	"time"

	"github.com/google/uuid"
)

type Barcode struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	SkuID        uuid.UUID `json:"sku_id"`
	BarcodeName  string    `json:"barcode_name"`
	BarcodeValue string    `json:"barcode_value"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
