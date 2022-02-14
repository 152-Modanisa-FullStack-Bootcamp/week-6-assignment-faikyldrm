package wallet_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"wallet/config"
	"wallet/model"
	"wallet/wallet"
)

var testConfig = &config.Config{InitialBalance: 100, MinimumBalance: -100}

func TestGetWalletsEmpty(t *testing.T) {
	defer modelNiller()
	var cases = []struct {
		expected *map[string]model.Wallet
		caseName string
	}{
		{
			expected: &map[string]model.Wallet{},
			caseName: "emptyCase",
		},
	}
	walletService := wallet.NewWalletService()
	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			result, _ := walletService.GetWallets()
			assert.Equal(t, c.expected, result)
		})
	}

}
func TestGetWallets(t *testing.T) {
	defer modelNiller()
	expectedWallet := make(map[string]model.Wallet, 1)
	expectedWallet["faik"] = model.Wallet{
		User:    "faik",
		Balance: testConfig.InitialBalance,
	}
	walletService := wallet.NewWalletService()
	walletService.CreateWallet("faik",testConfig.InitialBalance)
	wallets, _ := walletService.GetWallets()
	assert.Equal(t, &expectedWallet, wallets)

}


//happy path
func TestCreateWallet(t *testing.T) {
	defer modelNiller()
	var cases = []struct {
		param    string
		expected *model.Wallet
		caseName string
	}{{
		"faik", &model.Wallet{
			User:    "faik",
			Balance: testConfig.InitialBalance,
		},
		"First",
	},
		{
			"test", &model.Wallet{User: "test", Balance: testConfig.InitialBalance}, "heyyo",
		},
	}
	walletService := wallet.NewWalletService()

	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			result := walletService.CreateWallet(c.param,testConfig.InitialBalance)
			assert.Equal(t, c.expected, result)
		})

	}

}
func TestChangeWalletAmount(t *testing.T) {
	defer modelNiller()
	model.Wallets=make(map[string]model.Wallet)
	walletService := wallet.NewWalletService()
	expectedWallet := model.Wallet{
		User:    "faik",
		Balance: testConfig.InitialBalance,
	}
	model.Wallets["faik"]=expectedWallet
	wallet, err := walletService.ChangeWalletAmount(&expectedWallet, 100)
	assert.Nil(t, err)
	assert.Equal(t, float64(200), wallet.Balance)
}
func TestChangeWalletAmountNil(t *testing.T) {
	defer modelNiller()
	walletService := wallet.NewWalletService()
	wallet, err := walletService.ChangeWalletAmount(nil, 100)
	assert.Nil(t, wallet)
	assert.Error(t, err)
}
func TestGetWalletEmpty(t *testing.T) {
	defer modelNiller()
	walletService := wallet.NewWalletService()
	wallet, err := walletService.GetWallet("faik")
	assert.Equal(t, errors.New("user not exists"),err)
	assert.Nil(t, wallet)

}
func TestGetWalletWithWallet(t *testing.T)  {
	defer modelNiller()
	model.Wallets=make(map[string]model.Wallet)
	testWallet:=model.Wallet{
		User:    "faik",
		Balance: testConfig.InitialBalance,
	}
	model.Wallets["faik"] = testWallet
	walletService := wallet.NewWalletService()
	wallet, err := walletService.GetWallet("faik")
	assert.Nil(t, err)
	assert.Equal(t, &testWallet,wallet)
}
func TestCheckNoWalletExists(t *testing.T)  {
	defer modelNiller()
	walletService := wallet.NewWalletService()
	_, err := walletService.CheckWalletExists("faik")
	assert.Equal(t, err,false)
}
func TestCheckWalletExists(t *testing.T) {
	defer modelNiller()
	model.Wallets=make(map[string]model.Wallet)
	testWallet:=model.Wallet{
		User:    "faik",
		Balance: testConfig.InitialBalance,
	}
	model.Wallets["faik"]=testWallet
	walletService := wallet.NewWalletService()
	wallet, err := walletService.CheckWalletExists("faik")
	assert.Equal(t, true,err)
	assert.Equal(t, &testWallet,wallet)
}
func modelNiller() {
	model.Wallets = nil
}