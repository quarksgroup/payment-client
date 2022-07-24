package airtel

import "strings"

type Kind string

const (
	Cashin  Kind = "cashin"
	Cashout Kind = "cashout"
)

// Number represents information about a phone number details
type Number struct {
	Phone     string
	FirstName string
	LastName  string
	Status    bool
	HasPin    bool
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
