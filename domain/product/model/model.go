package model

import (
	"time"

	"github.com/google/uuid"
)

// Product data structure.
type Product struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Category  string     `json:"category"`
	Qty       int        `json:"qty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
