package model


type Wallet struct {
	User    string
	Balance float64
}
var Wallets map[string]Wallet