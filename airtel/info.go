package airtel

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/quarksgroup/payment-client/client"
)

type Kind string

const (
	Cashin  Kind = "cashin"
	Cashout Kind = "cashout"
)

// Number represents information about a phone number details
type Number struct {
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Status    bool   `json:"status"`
	HasPin    bool   `json:"has_pin"`
}

//TxInfo respresent information about transaction
type TxInfo struct {
	Ref    string
	Status string
	Type   Kind
}

//Abrivated transaction status
const (
	ts  = "TS"  //Transaction Success
	tf  = "TF"  //Transaction Failed
	ta  = "TA"  //Transaction Ambiguous
	tip = "TIP" //Transaction in Progress
)

//ConvertStatus convert transaction status to common status value
func ConvertStatus(status string) string {

	switch strings.ToUpper(status) {
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

//NumberInfo this is responsible for quering phone number information if is registered
func (c *Client) NumberInfo(ctx context.Context, phone string) (*Number, *client.Response, error) {

	endpoint := fmt.Sprintf("standard/v1/users/%s", phone)

	headers := http.Header{
		"X-country":  []string{c.Country},
		"X-Currency": []string{c.Currency},
	}

	out := new(checkResponse)

	res, err := c.do(ctx, "GET", endpoint, nil, out, headers, true)

	if !out.Status.Success {
		return nil, nil, &Error{Code: http.StatusBadRequest, Message: out.Status.Message}
	}

	return convertNumberInfo(out), res, err
}

//PullInfo responsible for quering the transaction information for the transaction made using
//collection api of airtel payment
func (c *Client) PullInfo(ctx context.Context, ref string) (*TxInfo, *client.Response, error) {
	endpoint := fmt.Sprintf("standard/v1/payments/%s", ref)

	headers := http.Header{
		"X-country":  []string{c.Country},
		"X-Currency": []string{c.Currency},
	}

	out := new(txInfo)

	res, err := c.do(ctx, "GET", endpoint, nil, out, headers, true)

	if !out.Status.Success {
		return nil, nil, &Error{Code: http.StatusBadRequest, Message: out.Status.Message}
	}

	return convertTxInfo(out, Cashin), res, err
}

//PushInfo responsible for quering the distrubuted information for the transaction made using
//distribution or cashout api of airtel payment
func (c *Client) PushInfo(ctx context.Context, ref string) (*TxInfo, *client.Response, error) {
	endpoint := fmt.Sprintf("standard/v1/disbursements/%s", ref)

	headers := http.Header{
		"X-country":  []string{c.Country},
		"X-Currency": []string{c.Currency},
	}

	out := new(txInfo)

	res, err := c.do(ctx, "GET", endpoint, nil, out, headers, true)

	if !out.Status.Success {
		return nil, nil, &Error{Code: http.StatusBadRequest, Message: out.Status.Message}
	}

	return convertTxInfo(out, Cashout), res, err
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

func convertTxInfo(res *txInfo, kind Kind) *TxInfo {
	return &TxInfo{
		Ref:    res.Data.Transaction.Id,
		Status: ConvertStatus(res.Data.Transaction.Status),
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

func convertNumberInfo(res *checkResponse) *Number {
	return &Number{
		Phone:     res.Data.Msisdn,
		FirstName: res.Data.FirstName,
		LastName:  res.Data.LastName,
		Status:    res.Status.Success,
		HasPin:    res.Data.IsPinSet,
	}
}
