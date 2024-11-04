package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	CompanyID uuid.UUID `json:"company_id"`
	Role      string    `json:"role"`
	UpdatedAt time.Time `json:"updated_at"`
}
