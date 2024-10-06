package models

import (
	"github.com/google/uuid"
)

type Category struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	Name     string    `json:"name"`
	ParentID uuid.UUID `json:"parent_id,omitempty"`
}
