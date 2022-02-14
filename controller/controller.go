package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"wallet/service"
)

type IController interface {
	WalletOps(w http.ResponseWriter, r *http.Request)
}
type Controller struct {
	service service.IService
}

func (c *Controller) WalletOps(w http.ResponseWriter, r *http.Request) {
	type postBodyRequest struct {
		Amount float64 `json:"amount"`
	}
	userName:=strings.Replace(r.RequestURI,"/","",1)
	if r.Method == "GET" {
		if r.RequestURI != "/" {
			wallet, err := c.service.GetWallet(userName)
			if err != nil {
				errorHandle(err, w)
				return
			}
			formatResponse(wallet, w)
		} else {
			wallets, err := c.service.GetWallets()
			fmt.Println(wallets)
			if err != nil {
				errorHandle(err, w)
				return
			}
			formatResponse(wallets, w)
		}

	} else if r.Method == "PUT" {

		wallet, err := c.service.CreateWallet(userName)
		if err != nil {
			errorHandle(err, w)
			return
		}
		formatResponse(wallet, w)

	} else if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errorHandle(err, w)
			return
		}

		requestBody := &postBodyRequest{}
		parseErr := json.Unmarshal(body, requestBody)

		if parseErr != nil {
			errorHandle(parseErr, w)
			return
		}

		changedWallet, err := c.service.ChangeWalletAmount(userName, requestBody.Amount)
		if err != nil {
			errorHandle(err, w)
			return
		}
		formatResponse(changedWallet, w)
	}

}
func formatResponse(response interface{}, w http.ResponseWriter) {
	json, err := json.Marshal(&response)
	if err != nil {
		errorHandle(err, w)
		return
	} else {
		w.Header().Add("content-type", "application/json")
		w.Write(json)
	}

}
func errorHandle(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
	return
}
func NewController(service service.IService) IController {
	return &Controller{service: service}
}
