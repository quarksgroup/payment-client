package airtel

// Payment represents as single transaction
type Payment struct {
	ID     string
	Ref    string
	Amount int64
	Phone  string
}

// Status reports the status of a requested transaction
type Status struct {
	Ref          string `json:"ref"`
	Status       bool   `json:"status"`
	Message      string `json:"message"`
	ResponseCode string `json:"response_code"`
}
