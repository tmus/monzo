package monzo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Account represents a Monzo Account.
type Account struct {
	ID            string
	Closed        bool
	Created       string
	Description   string
	Type          AccountType
	Currency      Currency
	Country       string `json:"country_code"`
	AccountNumber string `json:"account_number"`
	SortCode      string `json:"sort_code"`

	// The monzo.Client is embedded here to enable a fluent API.
	client *Client
}

// AccountType is the way that Monzo identifies accounts internally.
type AccountType string

// PrepaidAccount is for accounts that were created before
// the UKRetailAccount existed and can no longer be opened.
const PrepaidAccount AccountType = "uk_prepaid"

// UKRetailAccount is a Current Account.
const UKRetailAccount AccountType = "uk_retail"

// UKRetailJointAccount is a Current Account shared by two
// Monzo users.
const UKRetailJointAccount AccountType = "uk_retail_joint"

// Balance returns the current balance for the Account that
// it is called on.
func (a Account) Balance() (Balance, error) {
	req, err := a.client.NewRequest(http.MethodGet, "balance", nil)
	if err != nil {
		return Balance{}, err
	}

	q := req.URL.Query()
	q.Add("account_id", a.ID)
	req.URL.RawQuery = q.Encode()

	resp, _ := a.client.Do(req)

	b := new(bytes.Buffer)
	b.ReadFrom(resp.Body)
	str := b.String()

	if resp.StatusCode != http.StatusOK {
		return Balance{}, fmt.Errorf("failed to fetch balance: %s", str)
	}

	var bal Balance
	if err := json.Unmarshal(b.Bytes(), &bal); err != nil {
		return Balance{}, err
	}

	return bal, nil
}
