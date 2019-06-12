package monzo

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	ID       string
	Amount   int
	Currency AccountCurrency

	Description string
	Notes       string

	Created       time.Time
	DeclineReason string
	IsLoad        bool
	Settled       time.Time
	Category      string
	// Merchant Merchant
	Metadata interface{}

	client *Client
}

// Transactions gets a number of transactions for an account from the Monzo API.
// A bit useless atm because Monzo's API returns transactions in ascending order.
func (a Account) Transactions(limit int) ([]Transaction, error) {
	params := make(map[string]string)
	params["limit"] = strconv.Itoa(limit)
	return a.getTransactions(params)
}

func (a Account) Transaction(id string) (Transaction, error) {
	req, err := a.client.resourceRequest("transactions/" + id)
	if err != nil {
		return Transaction{}, err
	}

	resp, _ := a.client.Do(req)

	b := new(bytes.Buffer)
	b.ReadFrom(resp.Body)
	str := b.String()

	if resp.StatusCode != http.StatusOK {
		return Transaction{}, fmt.Errorf("failed to fetch transaction: %s", str)
	}

	bytes := b.Bytes()
	var transaction Transaction
	if err := unwrapJSON(bytes, "transaction", &transaction); err != nil {
		return Transaction{}, err
	}

	transaction.client = a.client

	return transaction, nil
}

// TransactionsSince returns the transactions that have occured since a given Time.
func (a Account) TransactionsSince(ts time.Time, limit int) ([]Transaction, error) {
	params := make(map[string]string)
	params["limit"] = strconv.Itoa(limit)
	params["since"] = ts.Format(time.RFC3339)

	return a.getTransactions(params)
}

// TransactionsBefore returns the transactions that occured before a given Time.
func (a Account) TransactionsBefore(ts time.Time, limit int) ([]Transaction, error) {
	params := make(map[string]string)
	params["limit"] = strconv.Itoa(limit)
	params["before"] = ts.Format(time.RFC3339)

	return a.getTransactions(params)
}

// TransactionsBetween returns the transactions that happened between two Times.
func (a Account) TransactionsBetween(since time.Time, before time.Time) ([]Transaction, error) {
	params := make(map[string]string)
	params["since"] = since.Format(time.RFC3339)
	params["before"] = before.Format(time.RFC3339)

	return a.getTransactions(params)
}

func (a Account) getTransactions(params map[string]string) ([]Transaction, error) {
	req, err := a.client.resourceRequest("transactions")
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("account_id", a.ID)

	for param, value := range params {
		q.Add(param, value)
	}

	req.URL.RawQuery = q.Encode()

	resp, _ := a.client.Do(req)

	b := new(bytes.Buffer)
	b.ReadFrom(resp.Body)
	str := b.String()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch transactions: %s", str)
	}

	bytes := b.Bytes()
	var transactions []Transaction
	if err := unwrapJSON(bytes, "transactions", &transactions); err != nil {
		return nil, err
	}

	for _, tx := range transactions {
		tx.client = a.client
	}

	return transactions, nil
}

func (t Transaction) Note(note string) error {
	endpoint := "/transactions/" + t.ID

	data := url.Values{}
	data.Add("metadata[notes]", note)

	req, err := t.client.NewRequest(http.MethodPatch, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := t.client.Do(req)

	b := new(bytes.Buffer)
	b.ReadFrom(resp.Body)
	str := b.String()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create note: %s", str)
	}

	return nil
}
