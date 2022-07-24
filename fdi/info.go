package fdi

import (
	"time"
)

// Info represents information about a transaction
type Info struct {
	ID        string
	Amount    float64
	Cost      float64
	Status    string
	Type      string
	CreatedAt time.Time
}
