package controller_test

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"wallet/controller"
	"wallet/mock"
	"wallet/model"
)

//happy Path
func TestWalletOpsGetWallets(t *testing.T) {
	cases := []struct {
		expectedResult []model.Wallet
		expectedError  error
		caseName       string
	}{
		{[]model.Wallet{},
			nil,
			"happy path",
		},
		{[]model.Wallet{
			{User: "faik", Balance: 100},
		},
			nil,
			"single Item",
		},
		{[]model.Wallet{
			{User: "faik", Balance: 100},
			{User: "veli", Balance: 200},
			{User: "yeni", Balance: 300},
		},
			nil,
			"multiple Item",
		},
		{[]model.Wallet{},
			errors.New("no mather"),
			"multiple Item",
		},
	}
	service := mock.NewMockIService(gomock.NewController(t))
	controller := controller.NewController(service)

	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			service.EXPECT().GetWallets().Return(c.expectedResult, c.expectedError).Times(1)
			controller.WalletOps(w, r)
			result := &[]model.Wallet{}
			err := json.Unmarshal(w.Body.Bytes(), result)
			if err != nil {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			}
			assert.Equal(t, &c.expectedResult, result)
		})
	}

}
func TestWalletOpsGetWallet(t *testing.T) {
	cases := []struct {
		param          string
		expectedResult *model.Wallet
		expectedError  error
		caseName       string
	}{
		{
			param: "faik",
			expectedResult: &model.Wallet{
				User:    "faik",
				Balance: 100,
			},

			caseName: "happy Path",
		},
		{
			param:          "faik",
			expectedResult: &model.Wallet{},

			caseName: "happy Path",
		},
		{
			param:          "faik",
			expectedResult: &model.Wallet{},
			expectedError:  errors.New("no mather"),
			caseName:       "happy Path",
		},
	}
	service := mock.NewMockIService(gomock.NewController(t))
	controller := controller.NewController(service)
	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {

			r := httptest.NewRequest(http.MethodGet, "/"+c.param, nil)
			w := httptest.NewRecorder()
			service.EXPECT().GetWallet(c.param).Return(c.expectedResult, c.expectedError).Times(1)
			controller.WalletOps(w, r)
			result := &model.Wallet{}
			err := json.Unmarshal(w.Body.Bytes(), result)
			if err != nil {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			}
			assert.Equal(t, c.expectedResult, result)
		})
	}
}

func TestWalletOpsCreateWallet(t *testing.T) {
	cases := []struct {
		param          string
		expectedResult *model.Wallet
		expectedError  error
		caseName       string
	}{
		{
			"faik",
			&model.Wallet{
				User:    "faik",
				Balance: 100,
			},
			nil,
			"happy Path",
		},
		{
			"faik",
			&model.Wallet{
				User:    "faik",
				Balance: 100,
			},
			errors.New("boo"),
			"happy Path",
		},
	}

	service := mock.NewMockIService(gomock.NewController(t))
	controller := controller.NewController(service)
	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPut, "/"+c.param, nil)
			w := httptest.NewRecorder()
			service.EXPECT().CreateWallet(c.param).Return(c.expectedResult, c.expectedError).Times(1)
			controller.WalletOps(w, r)
			result := &model.Wallet{}
			err := json.Unmarshal(w.Body.Bytes(), result)
			if err != nil {
				assert.Equal(t, w.Code, http.StatusInternalServerError)
			} else {
				assert.Equal(t, "application/json", w.Header().Get("content-type"))
				assert.Equal(t, c.expectedResult, result)
			}

		})
	}

}

func TestWalletOpsChangeAmount(t *testing.T) {
	type bodyReq struct {
		Amount float64 `json:"amount"`
	}
	cases := []struct {
		param          string
		bodyParam      bodyReq
		expectedResult *model.Wallet
		expectedError  error
		caseName       string
	}{
		{
			param:          "faik",
			bodyParam:      bodyReq{Amount: 100},
			expectedResult: &model.Wallet{User: "faik", Balance: 200},
			expectedError:  nil,
			caseName:       "happy path",
		},
		{
			param:          "faik",
			bodyParam:      bodyReq{Amount: 100},
			expectedResult: &model.Wallet{User: "faik", Balance: 200},
			expectedError:  errors.New("some Error"),
			caseName:       "happy path",
		},
	}
	service := mock.NewMockIService(gomock.NewController(t))
	controller := controller.NewController(service)

	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			marshal, _ := json.Marshal(c.bodyParam)
			r := httptest.NewRequest(http.MethodPost, "/"+c.param, strings.NewReader(string(marshal)))
			w := httptest.NewRecorder()
			service.EXPECT().ChangeWalletAmount(c.param, c.bodyParam.Amount).Return(c.expectedResult, c.expectedError).Times(1)
			controller.WalletOps(w, r)
			result := &model.Wallet{}
			err := json.Unmarshal(w.Body.Bytes(), result)
			if err != nil {
				assert.Equal(t, w.Code, http.StatusInternalServerError)
			} else {
				assert.Equal(t, "application/json", w.Header().Get("content-type"))
				assert.Equal(t, c.expectedResult, result)
			}

		})
	}

}
