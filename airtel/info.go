package airtel

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
