package coinbase

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
)

const (
	API_ENDPOINT = "https://coinbase.com/api/v1"
)

// Coinbase represents a coinbase client wrapper
type Coinbase struct {
	restClient *rClient
}

// New return a new Coinbase
func New(apiKey, apiSecret string) *Coinbase {
	return &Coinbase{
		restClient: newRestClient(apiKey, apiSecret),
	}
}

// Balance returns the user's accounts
func (c *Coinbase) GetAccounts() (accounts Accounts, err error) {
	resp, err := c.restClient.Do("GET", "accounts", "")
	if err = resp.HandleErr(err, []int{200}); err != nil {
		return
	}
	err = json.Unmarshal(resp.Body, &accounts)
	return
}

// GetPrimaryAccountBalance returns balance for the primary account
func (c *Coinbase) GetPrimaryAccountBalance() (balance float64, err error) {
	accounts, err := c.GetAccounts()
	if err != nil {
		return
	}
	// search primamry account.
	// it must be the fisrt one but...©
	var account Account
	for _, account = range accounts.Accounts {
		if account.Primary {
			break
		}
	}

	// double check in case there is no primary account on first page
	// it should never happen but...©
	if !account.Primary {
		return balance, errors.New("No primary account found.")
	}

	balance = account.Balance.Amount.asFloat
	return
}

// SendMoney call send_money for tansaction transaction
func (c *Coinbase) SendMoney(transaction *SmTransaction) (response SmMoneyResponse, err error) {
	payload, err := json.Marshal(transaction)
	if err != nil {
		return
	}
	r, err := c.restClient.Do("POST", "transactions/send_money", string(payload))
	log.Println(string(r.Body), err)
	if err = r.HandleErr(err, []int{200}); err != nil {
		return
	}
	err = json.Unmarshal(r.Body, &response)
	return
}

// GetTransactionDetails return transaction with ID id
func (c *Coinbase) GetTransactionDetails(id string) (transaction Transaction, err error) {
	r, err := c.restClient.Do("GET", "transactions/"+url.QueryEscape(id), "")
	if err = r.HandleErr(err, []int{200}); err != nil {
		return
	}
	var rr struct {
		Transaction Transaction
	}
	err = json.Unmarshal(r.Body, &rr)
	if err != nil {
		return
	}
	transaction = rr.Transaction
	return
}
