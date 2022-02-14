package wallet

import (
	"errors"
	"wallet/config"
	"wallet/model"
)

var testConfig = &config.Config{InitialBalance: 100, MinimumBalance: -100}

type IWallet interface {
	GetWallets() (*map[string]model.Wallet, error)
	GetWallet(user string) (*model.Wallet, error)
	CreateWallet(user string, initialAmount float64) *model.Wallet
	ChangeWalletAmount(opWallet *model.Wallet, amount float64) (*model.Wallet, error)
	CheckWalletExists(user string) (*model.Wallet, bool)
}
type Wallet struct {
}

func (w *Wallet) CreateWallet(user string, initialAmount float64) *model.Wallet {
	newWallet := model.Wallet{User: user, Balance: initialAmount}
	if len(model.Wallets) == 0 {
		model.Wallets = make(map[string]model.Wallet, 1)
	}
	model.Wallets[user] = newWallet
	return &newWallet
}

func (w *Wallet) ChangeWalletAmount(opWallet *model.Wallet, amount float64) (*model.Wallet, error) {
	if opWallet == nil {
		return nil, errors.New("wallet cannot be Nil")
	}
	value:= model.Wallets[opWallet.User]
	value.Balance+=amount
	model.Wallets[opWallet.User]=value
	return &value, nil
}

func (w *Wallet) GetWallets() (*map[string]model.Wallet, error) {
	if len(model.Wallets) > 0 {
		return &model.Wallets, nil
	} else {
		empty := make(map[string]model.Wallet, 1)
		return &empty, nil
	}

}

func (w *Wallet) GetWallet(user string) (*model.Wallet, error) {
	wallet, exists := model.Wallets[user]
	if exists {
		return &wallet, nil
	} else {
		return nil, errors.New("user not exists")
	}

}

func (w *Wallet) CheckWalletExists(user string) (*model.Wallet, bool) {
	wallet, exists := model.Wallets[user]
	if exists {
		return &wallet, true
	} else {
		return nil, false
	}
}
func NewWalletService() IWallet {
	return &Wallet{}
}
