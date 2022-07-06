package airtel

import (
	"context"
	"fmt"
	"net/http"

	"github.com/quarksgroup/payment-client/payment/airtel"
)

type checkNumberService struct {
	client *wrapper
}

func (s *checkNumberService) Check(ctx context.Context, phone string) (*airtel.Number, *airtel.Response, error) {

	endpoint := fmt.Sprintf("standard/v1/users/%s", phone)

	header := http.Header{
		"X-country":  []string{s.client.Country},
		"X-Currency": []string{s.client.Currency},
	}

	out := new(checkResponse)

	res, err := s.client.do(ctx, "GET", endpoint, nil, out, header)

	if !out.Status.Success {
		return nil, nil, &airtel.Error{Code: http.StatusBadRequest, Message: out.Status.Message}
	}

	return convertNumberInfo(out), res, err
}

type checkResponse struct {
	Data struct {
		FirstName    string `json:"first_name"`
		Grade        string `json:"grade"`
		IsBarred     bool   `json:"is_barred"`
		IsPinSet     bool   `json:"is_pin_set"`
		LastName     string `json:"last_name"`
		Msisdn       string `json:"msisdn"`
		RegStatus    string `json:"reg_status"`
		Registration struct {
			Id     string `json:"id"`
			Status string `json:"status"`
		} `json:"registration"`
	} `json:"data"`
	Status struct {
		Code        string `json:"code"`
		Success     bool   `json:"success"`
		Result_code string `json:"result_code"`
		Message     string `json:"message"`
	} `json:"status"`
}

func convertNumberInfo(res *checkResponse) *airtel.Number {
	return &airtel.Number{
		Phone:     res.Data.Msisdn,
		FirstName: res.Data.FirstName,
		LastName:  res.Data.LastName,
		Status:    res.Status.Success,
		HasPin:    res.Data.IsPinSet,
	}
}
