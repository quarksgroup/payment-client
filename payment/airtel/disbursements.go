package airtel

import "context"

// DisbursementService pulls a disbursement funds from airtel payment API provider
type DisbursementService interface {
	// Pull initializes a Pull of funds to an external wallet transaction.
	Pull(context.Context, *PaymentReq) (*PaymentResp, *Response, error)
}
