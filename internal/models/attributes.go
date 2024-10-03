package models

import "github.com/google/uuid"

type Attribute struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
	CreatedAt string    `json:"created_at"`
}
