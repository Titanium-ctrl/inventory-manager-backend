package models

import "github.com/google/uuid"

type SKUAttributes struct {
	SkuID          uuid.UUID `json:"sku_id"`
	AttributeID    uuid.UUID `json:"attribute_id"`
	AttributeValue string    `json:"attr_value"`
	UserID         uuid.UUID `json:"user_id"`
}
