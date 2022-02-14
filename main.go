package main

import (
	"fmt"
	"net/http"
	"wallet/config"
	"wallet/controller"
	"wallet/service"
	"wallet/wallet"
)

func main() {
	config := config.Get()
	wallet := wallet.NewWalletService()
	service := service.NewWalletService(wallet,config)
	myController := controller.NewController(service)
	http.HandleFunc("/", myController.WalletOps)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
