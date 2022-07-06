package airtel

import (
	"context"
	"net/http"

	"github.com/iradukunda1/payment-staging/payment/airtel"
)

type collectionService struct {
	client *wrapper
}

//Push...
func (s *collectionService) Push(ctx context.Context, req *airtel.PaymentReq) (*airtel.PaymentResp, *airtel.Response, error) {

	endpoint := "merchant/v1/payments/"

	in := new(pushRequest)
	in.Reference = req.Ref
	in.Subscriber.Country = s.client.Country
	in.Subscriber.Currency = s.client.Currency
	in.Subscriber.Msisdn = req.Phone
	in.Transaction.Amount = req.Amount
	in.Transaction.Id = req.ID

	header := http.Header{
		"X-country":  []string{s.client.Country},
		"X-Currency": []string{s.client.Currency},
	}

	out := new(pushResponse)

	res, err := s.client.do(ctx, "POST", endpoint, in, out, header)

	// if !out.Status.Success {
	// 	return nil, nil, &airtel.Error{Code: http.StatusBadRequest, Message: out.Status.Message}
	// }

	return convertPush(out), res, err
}

//pushRequest represent push collection request payload
type pushRequest struct {
	Reference  string `json:"reference"`
	Subscriber struct {
		Country  string `json:"country"`
		Currency string `json:"currency"`
		Msisdn   string `json:"msisdn"`
	} `json:"subscriber"`
	Transaction struct {
		Id     string `json:"id"`
		Amount int64  `json:"amount"`
	} `json:"transaction"`
}

//pushResponse represents push collection response body
type pushResponse struct {
	Data struct {
		Transaction struct {
			Id     string `json:"id"`
			Status string `json:"status"`
		} `json:"transaction"`
	} `json:"data"`
	Status struct {
		ResponseCode string `json:"response_code"`
		Code         string `json:"code"`
		Success      bool   `json:"success"`
		Result_code  string `json:"result_code"`
		Message      string `json:"message"`
	} `json:"status"`
}

func convertPush(res *pushResponse) *airtel.PaymentResp {
	return &airtel.PaymentResp{
		Ref:          res.Data.Transaction.Id,
		Status:       res.Status.Success,
		ResponseCode: res.Status.ResponseCode,
		Message:      res.Status.Message,
	}
}
