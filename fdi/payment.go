package fdi

// Provider is the supported payment channel
type Provider string

// Payment represents as single transaction
type Payment struct {
	ID       string
	Amount   float64
	Wallet   string
	Provider Provider
}

// Status reports the status of a requested transaction
type Status struct {
	Ref   string
	GRef  string
	State string
}
