package main

import (
	"github.com/Toorop/go-coinbase"
	"log"
)

const (
	APIKEY    = ""
	APISECRET = ""
)

func main() {
	cb := coinbase.New(APIKEY, APISECRET)

	// GetAccounts
	accounts, err := cb.GetAccounts()
	log.Println(err, accounts)

	// send 0.001 BTC to 1HgpsmxV52eAjDcoNpVGpYEhGfgN7mM1JB
	toSend := &coinbase.SmTransaction{
		Amount:  "0.001",
		To:      "1HgpsmxV52eAjDcoNpVGpYEhGfgN7mM1JB",
		UserFee: "0.0002",
	}
	r, err := cb.SendMoney(toSend)
	log.Println(err, r)
	log.Println(r.Transaction.Id)

	// get transaction details
	transaction, err := cb.GetTransactionDetails(r.Transaction.Id)
	log.Println(err, transaction)
	log.Println(transaction.Hsh)
}
