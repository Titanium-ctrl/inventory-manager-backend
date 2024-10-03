package models

import (
	"time"

	"github.com/google/uuid"
)

type Attribute struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
