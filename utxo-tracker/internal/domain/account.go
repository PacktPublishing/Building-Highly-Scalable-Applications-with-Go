package domain

// AccountName is unique per user
type AccountName string

// Username is unique per application instance
type UserName string

// This represents the xPub, yPub, zPub etc
type XPubEtc string

// The type of Bitcoin account
type AccountType string

const (
	ATTapRoot      AccountType = "taproot"       // starts with bc1p
	ATNativeSegwit AccountType = "native-segwit" // starts with bc1
	ATSegwit       AccountType = "segwit"        // starts with 3
	ATLegacy       AccountType = "legacy"        // starts with 1
	ATNotDefined   AccountType = ""
)

// Address is a Bitcoin address
type Address string

// Account Represents a Bitcoin Account
type Account struct {
	ID   AccountName
	User UserName
	XPub XPubEtc
	Type AccountType
}

// AccountAddresses is the addresses that was found for a specific XPubEtc
type AccountAddresses struct {
	Account AccountName
	Items   []Address
}
