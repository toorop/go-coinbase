package coinbase

import (
	"encoding/json"
	"strconv"
	"time"
)

const TIME_FORMAT = "2006-01-02T15:04:05-07:00"

type jTime struct {
	time.Time
}

func (jt *jTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.Parse(TIME_FORMAT, s)
	if err != nil {
		return err
	}
	jt.Time = t
	return nil
}

func (jt jTime) MarshalJSON() ([]byte, error) {
	return json.Marshal((*time.Time)(&jt.Time).Format(TIME_FORMAT))
}

// Amount represents amount ...
// For some strange reason coinbase return a string for amount
// this type allows us to have the float value directly in the struct
type amount struct {
	asStr   string
	asFloat float64
}

func (a *amount) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	a.asFloat = t
	a.asStr = s
	return nil
}

// Accout represents an account
type Account struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Balance struct {
		Amount   amount `json:"amount"`
		Currency string `json:"currency"`
	} `json:"balance"`
	NativeBalance struct {
		Amount   amount `json:"amount"`
		Currency string `json:"currency"`
	} `json:"native_balance"`
	CreatedAt string `json:"created_at"`
	Primary   bool   `json:"primary"`
	Active    bool   `json:"active"`
}

// Accounts  represents a getAccounts response
type Accounts struct {
	Accounts    []Account `json:"accounts"`
	TotalCount  int       `json:"total_count"`
	NumPages    int       `json:"num_pages"`
	CurrentPage int       `json:"current_page"`
}

// smTransaction represents the transactions parameter for send_money call
type SmTransaction struct {
	// An email address or a bitcoin address
	To string `json:"to"`
	// A string amount that will be converted to BTC, such as ‘1’ or ‘1.234567’.
	// Also must be >= ‘0.01’ or it will shown an error.
	Amount            string `json:"amount"`
	AmountString      string `json:"amount_string,omitempty"`
	AmountCurrencyIso string `json:"amount_currency_iso,omitempty"`
	// Notes
	Notes string `json:"notes,omitempty"`
	// User fee
	UserFee string `json:"user_fee,omitempty"`
	// Referrer ID
	ReferrerId string `json:"referrer_id,omitempty"`
	// idem
	Idem string `json:"idem,omitempty"`
	// instant Buy
	InstantBuy bool `json:"instant_buy,omitempty"`
	// order_id
	OrderId string `json:"order_ids,omitempty"`
}

// Tansaction represents a transaction
type Transaction struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	Hsh       string `json:"hsh"`
	Notes     string `json:"notes"`
	Idem      string `json:"idem"`
	Amount    struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Request bool   `json:"request"`
	Status  string `json:"status"`
	Sender  struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"sender"`
	Recipient struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"recipient"`
	RecipientAddress string `json:"recipient_address"`
}

// smMoneyResponse represents a response to a send_money_call
type SmMoneyResponse struct {
	Success     bool        `json:"success"`
	Error       []string    `json:"error"`
	Transaction Transaction `json:"transaction"`
}
