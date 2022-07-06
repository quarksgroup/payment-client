package airtel

import "context"

// PaymentReq represents as single transaction
type PaymentReq struct {
	ID     string
	Ref    string
	Amount int64
	Phone  string
}

// PaymentResp reports the status of a requested transaction
type PaymentResp struct {
	Ref          string `json:"ref"`
	Status       bool   `json:"status"`
	Message      string `json:"message"`
	ResponseCode string `json:"response_code"`
}

// CollectionsService push funds to airtel payment API provider
type CollectionsService interface {
	// Push initializes a push of funds to an external wallet transaction.
	Push(context.Context, *PaymentReq) (*PaymentResp, *Response, error)
}
