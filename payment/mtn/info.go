package mtn

import (
	"context"
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

// InfoService returns information about a transaction
type InfoService interface {
	// Info takes a transaction's id and returns info about it.
	Info(context.Context, string) (*Info, *Response, error)
}
