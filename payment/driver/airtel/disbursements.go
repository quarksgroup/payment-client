package airtel

import (
	"context"
	"net/http"

	"github.com/quarksgroup/payment-client/payment/airtel"
)

type disbursementService struct {
	client *wrapper
}

//Pull...
func (d *disbursementService) Pull(ctx context.Context, req *airtel.PaymentReq) (*airtel.PaymentResp, *airtel.Response, error) {

	endpoint := "standard/v1/disbursements/"

	in := new(pullRequest)
	in.Reference = req.Ref
	in.Pin = d.client.EncryptedPin
	in.Payee.Msisdn = req.Phone
	in.Transaction.Amount = req.Amount
	in.Transaction.Id = req.ID

	header := http.Header{
		"X-country":  []string{d.client.Country},
		"X-Currency": []string{d.client.Currency},
	}

	out := new(pullResponse)

	res, err := d.client.do(ctx, "POST", endpoint, in, out, header)

	return convertPull(out), res, err
}

//pullRequest represent push collection request payload
type pullRequest struct {
	Reference string `json:"reference"`
	Pin       string `json:"pin"`
	Payee     struct {
		Msisdn string `json:"msisdn"`
	} `json:"payee"`
	Transaction struct {
		Id     string `json:"id"`
		Amount int64  `json:"amount"`
	} `json:"transaction"`
}

//pullResponse represents push collection response body
type pullResponse struct {
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

func convertPull(res *pullResponse) *airtel.PaymentResp {
	return &airtel.PaymentResp{
		Ref:          res.Data.Transaction.Id,
		Status:       res.Status.Success,
		ResponseCode: res.Status.ResponseCode,
		Message:      res.Status.Message,
	}
}
