package service

import (
	"errors"
	"wallet/config"
	"wallet/model"
	"wallet/wallet"
)

type IService interface {
	GetWallets() ([]model.Wallet, error)
	GetWallet(user string) (*model.Wallet, error)
	CreateWallet(user string) (*model.Wallet, error)
	ChangeWalletAmount(user string, amount float64) (*model.Wallet, error)
	CheckWalletExists(user string) (*model.Wallet, bool)

}
type Service struct {
	wallet wallet.IWallet
	config *config.Config
}

func (s *Service) GetWallets() ([]model.Wallet, error) {
	wallets, err := s.wallet.GetWallets()
	if err != nil {
		return nil, err
	}
	v := make([]model.Wallet, 0, len(*wallets))
	for _, value := range *wallets {
		v = append(v, value)
	}

	return v, nil
}
func (s *Service) GetWallet(user string) (*model.Wallet, error) {
	wallet, err := s.wallet.GetWallet(user)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}
func (s *Service) CreateWallet(user string) (*model.Wallet, error) {
	wallet, exists := s.wallet.CheckWalletExists(user)
	if exists {
		return wallet, nil
	}
	wallet = s.wallet.CreateWallet(user, s.config.InitialBalance)
	return wallet, nil
}
func (s *Service) ChangeWalletAmount(user string, amount float64) (*model.Wallet, error) {
	existsWallet, exists := s.wallet.CheckWalletExists(user)
	if exists == false {
		return nil, errors.New("wallet Not Found")
	}
	exchange :=  WalletCanExchange(existsWallet, amount, s.config.MinimumBalance)
	if !exchange {

		return nil, errors.New("minimum amount exceeded")
	}

	wallet, err := s.wallet.ChangeWalletAmount(existsWallet, amount)
	if err != nil {
		return nil, err
	}
	return wallet, nil

}
func (s *Service) CheckWalletExists(user string) (*model.Wallet, bool) {
	wallet, exists := s.wallet.CheckWalletExists(user)
	return wallet, exists
}
func   WalletCanExchange(wallet *model.Wallet, amount float64, benchValue float64) bool {
	if wallet.Balance+amount < benchValue {
		return false
	}

	return true
}
func NewWalletService(wallet wallet.IWallet,config *config.Config) IService {
	return &Service{wallet: wallet,config: config}
}
