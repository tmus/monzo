package monzo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Transaction represents a single item on the Monzo feed.
type Transaction struct {
	ID       string
	Amount   int
	Currency Currency

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

// Transactions gets a number of transactions for an account.
func (a Account) Transactions(limit int) ([]Transaction, error) {
	params := make(map[string]string)
	params["limit"] = strconv.Itoa(limit)
	return a.getTransactions(params)
}

// Transaction returns a single transaction for an account.
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

// Note stores a string against the Transaction.
func (t Transaction) Note(note string) error {
	data := make(map[string]string)
	data["notes"] = note
	return t.AddMetadata(data)
}

// AddMetadata saves Metadata against a Transaction.
//
// Currently this is not visible in the Monzo App.
func (t Transaction) AddMetadata(meta map[string]string) error {
	endpoint := "/transactions/" + t.ID

	data := url.Values{}

	for key, value := range meta {
		data.Add("metadata["+key+"]", value)
	}

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
		return fmt.Errorf("failed to add metadata: %s", str)
	}

	return nil
}

// AddReceipt saves the given Receipt against the Transaction.
func (t Transaction) AddReceipt(r *Receipt) error {
	r.SetTransaction(t.ID)

	data, err := json.Marshal(r)
	if err != nil {
		return err
	}

	req, err := t.client.NewRequest(
		http.MethodPut,
		"/transaction-receipts",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, _ := t.client.Do(req)

	b := new(bytes.Buffer)
	b.ReadFrom(resp.Body)
	str := b.String()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add receipt: %s", str)
	}

	return nil
}
