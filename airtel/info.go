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
	Ts  = "TS"  //Transaction Success
	Tf  = "TF"  //Transaction Failed
	Ta  = "TA"  //Transaction Ambiguous
	Tip = "TIP" //Transaction in Progress
)

//ConvertStatus convert transaction status to common status value
func ConvertStatus(status string) string {

	switch strings.ToUpper(status) {
	case Ts:
		return "successful"
	case Tf:
		return "failed"
	case Ta:
		return "failed"
	case Tip:
		return "pending"
	default:
		return "failed"
	}
}
