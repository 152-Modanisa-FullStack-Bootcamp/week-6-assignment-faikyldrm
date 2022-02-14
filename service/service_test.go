package service_test

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"wallet/config"
	"wallet/mock"
	"wallet/service"

	"wallet/model"
)

var testConfig = &config.Config{InitialBalance: 0, MinimumBalance: -100}

func TestGetWallets(t *testing.T) {
	cases := []struct {
		mockResultExpected *map[string]model.Wallet
		resultExpected     []model.Wallet
		error              error
		caseName           string
	}{
		{&map[string]model.Wallet{
			"faik": {
				User:    "faik",
				Balance: 0,
			},
		},
			[]model.Wallet{
				{
					User:    "faik",
					Balance: 0,
				},
			},
			nil,
			"Only Result",
		}, {
			&map[string]model.Wallet{
				"faik": {
					User:    "faik",
					Balance: 0,
				},
				"Veli": {
					User:    "Veli",
					Balance: 0,
				},
			},
			[]model.Wallet{
				{
					User:    "Veli",
					Balance: 0,
				},
				{
					User:    "faik",
					Balance: 0,
				},

			}, nil,
			"Multiple Result",
		}, {
			nil,
			nil,
			errors.New("test error"),
			"Nil Case",
		},
	}
	wallet := mock.NewMockIWallet(gomock.NewController(t))
	service := service.NewWalletService(wallet, testConfig)
	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			defer modelNiller()
			wallet.EXPECT().GetWallets().Return(c.mockResultExpected, c.error).Times(1)
			result, err := service.GetWallets()
			assert.Equal(t, c.error, err)
			assert.ElementsMatch(t, c.resultExpected, result)
		})
	}
}

func TestGetWallet(t *testing.T) {
	cases := []struct {
		param              string
		mockResultExpected *model.Wallet
		error              error
		caseName           string
	}{
		{"faik",
			&model.Wallet{User: "faik", Balance: 0},
			nil,
			"first",
		},
		{"faik",
			&model.Wallet{},
			nil,
			"empty",
		},
		{"faik",
			nil,
			errors.New("no mather error"),
			"error case",
		},
		{"veli",
			nil,
			errors.New("data layer error"),
			"error case",
		},
	}

	wallet := mock.NewMockIWallet(gomock.NewController(t))
	service := service.NewWalletService(wallet, testConfig)
	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			wallet.EXPECT().GetWallet(c.param).Return(c.mockResultExpected, c.error).Times(1)
			result, err := service.GetWallet(c.param)
			assert.Equal(t, c.error, err)
			assert.Equal(t, c.mockResultExpected, result)
		})

	}

}

func TestCreateWallet(t *testing.T) {
	cases := []struct {
		param                   string
		mockResultExpected      *model.Wallet
		mockCheckResponseBool   bool
		mockCheckResponseWallet *model.Wallet
		error                   error
		caseName                string
	}{
		{"faik",
			&model.Wallet{User: "faik", Balance: testConfig.InitialBalance},
			false,
			nil,
			nil,
			"happy path",
		},
		{"faik",
			&model.Wallet{User: "faik", Balance: testConfig.InitialBalance},
			true,
			&model.Wallet{User: "faik", Balance: testConfig.InitialBalance},
			nil,
			"exists wallet",
		},
	}
	wallet := mock.NewMockIWallet(gomock.NewController(t))
	service := service.NewWalletService(wallet, testConfig)
	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			defer modelNiller()
			wallet.EXPECT().CheckWalletExists(c.param).Return(c.mockCheckResponseWallet, c.mockCheckResponseBool).Times(1)
			if !c.mockCheckResponseBool {
				wallet.EXPECT().CreateWallet(c.param, testConfig.InitialBalance).Return(c.mockResultExpected).Times(1)
			}
			result, err := service.CreateWallet(c.param)
			assert.Equal(t, c.error, err)
			assert.Equal(t, c.mockResultExpected, result)
		})
	}
}

func TestCheckWalletExists(t *testing.T) {
	cases := []struct {
		param            string
		mockResultWallet *model.Wallet
		mockResultBool   bool
		error            error
		caseName         string
	}{
		{"faik",
			&model.Wallet{User: "faik", Balance: 0},
			true,
			nil,
			"happy path",
		},
		{"faik",
			&model.Wallet{},
			false,
			nil,
			"happy path",
		},
	}
	wallet := mock.NewMockIWallet(gomock.NewController(t))
	service := service.NewWalletService(wallet, testConfig)
	for _, c := range cases {
		t.Run(c.param, func(t *testing.T) {
			wallet.EXPECT().CheckWalletExists(c.param).Return(c.mockResultWallet, c.mockResultBool).Times(1)
			existsWallet, exists := service.CheckWalletExists(c.param)
			assert.Equal(t, c.mockResultWallet, existsWallet)
			assert.Equal(t, c.mockResultBool, exists)
		})
	}
}
func TestWalletCanChange(t *testing.T) {
	cases := []struct {
		wallet   *model.Wallet
		amount   float64
		expected bool
		caseName string
	}{
		{
			&model.Wallet{User: "faik", Balance: testConfig.InitialBalance},
			100,
			true,
			"happy path",
		},
		{
			&model.Wallet{User: "faik", Balance: testConfig.InitialBalance},
			-500,
			false,
			"happy path",
		},
	}

	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			result := service.WalletCanExchange(c.wallet, c.amount, testConfig.MinimumBalance)
			assert.Equal(t, c.expected, result)
		})
	}

}
func TestChangeWalletAmount(t *testing.T) {
	cases := []struct {
		param                   string
		amount                  float64
		mockCheckResponseBool   bool
		mockCheckResponseWallet *model.Wallet
		changeWalletResult      *model.Wallet
		mockCanWalletChange     bool
		error                   error
		caseName                string
	}{
		{
			"faik",
			100,
			true,
			&model.Wallet{User: "faik", Balance: testConfig.InitialBalance},
			&model.Wallet{User: "faik", Balance: testConfig.InitialBalance + 100},
			true,
			nil,
			"happy path",
		},
		{
			"faik",
			100,
			false,
			&model.Wallet{User: "faik", Balance: testConfig.InitialBalance},
			nil,
			true,
			errors.New("wallet Not Found"),
			"Wallet Not Found",
		},
		{
			"faik",
			-1000,
			true,
			&model.Wallet{User: "faik", Balance: testConfig.InitialBalance},
			nil,
			false,
			errors.New("minimum amount exceeded"),
			"Wallet Not Found",
		},
		{
			"faik",
			100,
			true,
			&model.Wallet{User: "faik", Balance: testConfig.InitialBalance},
			nil,
			true,
			errors.New("any error"),
			"Wallet Not Found",
		},
	}
	wallet := mock.NewMockIWallet(gomock.NewController(t))
	service := service.NewWalletService(wallet, testConfig)
	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			wallet.EXPECT().CheckWalletExists(c.param).Return(c.mockCheckResponseWallet, c.mockCheckResponseBool).Times(1)
			if c.mockCheckResponseBool && c.mockCanWalletChange {
				wallet.EXPECT().ChangeWalletAmount(c.mockCheckResponseWallet, c.amount).Return(c.changeWalletResult, c.error).Times(1)
			}
			result, err := service.ChangeWalletAmount(c.param, c.amount)
			fmt.Println(result)
			assert.Equal(t, c.changeWalletResult, result)
			assert.Equal(t, c.error, err)
		})
	}

}

func modelNiller() {
	model.Wallets = nil
}
