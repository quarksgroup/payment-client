package fdi

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

// InfoService returns information about a passed in item
type InfoService interface {
	// TransactionInfo takes a transaction's id and returns info about it.
	TransactionInfo(context.Context, string) (*Info, *Response, error)
}
