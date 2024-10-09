package models

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	SkuID      uuid.UUID `json:"sku_id"`
	LocationID uuid.UUID `json:"location_id"`
	UserID     uuid.UUID `json:"user_id"`
	Quantity   int       `json:"quantity"`
	UpdatedAt  time.Time `json:"updated_at"`
}
