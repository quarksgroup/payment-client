package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/quarksgroup/payment-client/airtel"
)

//NumberInfo this is responsible for quering phone number information if is registered
func (c *Client) NumberInfo(ctx context.Context, phone string) (*airtel.Number, *airtel.Response, error) {

	if err := c.renewToken(ctx); err != nil {
		return nil, nil, err
	}

	endpoint := fmt.Sprintf("standard/v1/users/%s", phone)

	header := http.Header{
		"X-country":  []string{c.Client.Country},
		"X-Currency": []string{c.Client.Currency},
	}

	out := new(checkResponse)

	res, err := c.do(ctx, "GET", endpoint, nil, out, header)

	if !out.Status.Success {
		return nil, nil, &airtel.Error{Code: http.StatusBadRequest, Message: out.Status.Message}
	}

	return convertNumberInfo(out), res, err
}

//PullInfo responsible for quering the transaction information for the transaction made using
//collection api of airtel payment
func (c *Client) PullInfo(ctx context.Context, ref string) (*airtel.TxInfo, *airtel.Response, error) {

	if err := c.renewToken(ctx); err != nil {
		return nil, nil, err
	}

	endpoint := fmt.Sprintf("standard/v1/payments/%s", ref)

	header := http.Header{
		"X-country":  []string{c.Client.Country},
		"X-Currency": []string{c.Client.Currency},
	}

	out := new(txInfo)

	res, err := c.do(ctx, "GET", endpoint, nil, out, header)

	if !out.Status.Success {
		return nil, nil, &airtel.Error{Code: http.StatusBadRequest, Message: out.Status.Message}
	}

	return convertTxInfo(out, airtel.Cashin), res, err
}

//PushInfo responsible for quering the distrubuted information for the transaction made using
//distribution or cashout api of airtel payment
func (c *Client) PushInfo(ctx context.Context, ref string) (*airtel.TxInfo, *airtel.Response, error) {

	if err := c.renewToken(ctx); err != nil {
		return nil, nil, err
	}

	endpoint := fmt.Sprintf("standard/v1/disbursements/%s", ref)

	header := http.Header{
		"X-country":  []string{c.Client.Country},
		"X-Currency": []string{c.Client.Currency},
	}

	out := new(txInfo)

	res, err := c.do(ctx, "GET", endpoint, nil, out, header)

	if !out.Status.Success {
		return nil, nil, &airtel.Error{Code: http.StatusBadRequest, Message: out.Status.Message}
	}

	return convertTxInfo(out, airtel.Cashout), res, err
}

type txInfo struct {
	Data struct {
		Transaction struct {
			AirtelMoneyId string `json:"airtel_money_id"`
			Id            string `json:"id"`
			Message       string `json:"message"`
			Status        string `json:"status"`
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

func convertTxInfo(res *txInfo, kind airtel.Kind) *airtel.TxInfo {
	return &airtel.TxInfo{
		Ref:    res.Data.Transaction.Id,
		Status: convertTxStatus(res.Data.Transaction.Status),
		Type:   kind,
	}
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

const (
	ts  = "TS"  //Transaction Success
	tf  = "TF"  //Transaction Failed
	ta  = "TA"  //Transaction Ambiguous
	tip = "TIP" //Transaction in Progress
)

//convertTxStatus convert transaction status to common status value
func convertTxStatus(status string) string {
	switch status {
	case ts:
		return "successful"
	case tf:
		return "failed"
	case ta:
		return "failed"
	case tip:
		return "pending"
	default:
		return "failed"
	}
}
