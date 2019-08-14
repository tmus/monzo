package monzo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// APIBase is the root of the Monzo API.
const APIBase string = "https://api.monzo.com"

// Client is the way to interact with the Monzo API.
type Client struct {
	Token string

	http.Client
}

// NewClient uses the passed token to create a new Monzo Client.
// No validation is done on this method so users should call
// the Ping method after creating a client to ensure that
// the new connection has been created successsfully.
func NewClient(token string) *Client {
	return &Client{
		Token: token,
	}
}

// Ping attempts to connect to the Monzo API using the given
// client.
func (c *Client) Ping() error {
	if c.Token == "" {
		return errors.New("error pinging Monzo API. Client token cannot be empty")
	}

	req, err := c.NewRequest(http.MethodGet, "ping/whoami", nil)
	if err != nil {
		return err
	}

	resp, _ := c.Do(req)

	if resp.StatusCode != http.StatusOK {
		b := new(bytes.Buffer)
		b.ReadFrom(resp.Body)
		str := b.String()
		return fmt.Errorf("error pinging Monzo API. JSON response: %s", str)
	}

	return nil
}

// NewRequest creates an *http.Request with some Monzo-specific
// sensible defaults, such as Authorization headers.
func (c *Client) NewRequest(method string, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, APIBase+"/"+endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.Token)
	return req, nil
}

func (c *Client) resourceRequest(resource string) (*http.Request, error) {
	req, err := c.NewRequest(http.MethodGet, resource, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// Account returns a single Account from the Monzo API. If the
// account is not found (does not exist), an error is returned.
func (c *Client) Account(id string) (Account, error) {
	accs, err := c.accounts()
	if err != nil {
		return Account{}, err
	}

	// Monzo doesn't have the capability to retrieve a single
	// account, so we need to get all the users accounts
	// and filter them down using the provided id.
	for _, acc := range accs {
		if acc.ID == id {
			return acc, nil
		}
	}

	return Account{}, fmt.Errorf("no account found with ID %s", id)
}

// Accounts returns a slice of Account structs, one for each of
// the Monzo accounts associated with the authentication.
func (c *Client) Accounts() ([]Account, error) {
	return c.accounts()
}

func (c *Client) accounts() ([]Account, error) {
	req, err := c.resourceRequest("accounts")
	if err != nil {
		return nil, err
	}

	resp, _ := c.Do(req)

	b := new(bytes.Buffer)
	b.ReadFrom(resp.Body)
	str := b.String()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch accounts: %s", str)
	}

	bytes := b.Bytes()
	var accounts []Account
	if err := unwrapJSON(bytes, "accounts", &accounts); err != nil {
		return nil, err
	}

	var accs []Account
	for _, acc := range accounts {
		// The API still returns the Monzo beta prepaid accounts.
		// These can't be actioned meaningfully, so they are
		// removed from the slice if they exist.
		if acc.Type != PrepaidAccount {
			// To provide a fluent API for the account, it needs
			// to know how to talk to Monzo. The monzo.Client
			// is embedded in the Account struct so that
			// calls can be passed to it.
			acc.client = c
			accs = append(accs, acc)
		}
	}

	return accs, nil
}

// Deposit creates a new Deposit struct. Monzo uses a 'dedupe_id'
// to ensure that the request is idempotent, so the deposit is
// not ran when it is created. To action the deposit, call
// the `Run` method on it.
func (a Account) Deposit(p Pot, amt int) (*Deposit, error) {
	endpoint := "/pots/" + p.ID + "/deposit"

	data := url.Values{}
	data.Add("source_account_id", a.ID)
	data.Add("amount", strconv.Itoa(amt))

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	data.Add("dedupe_id", strconv.FormatFloat(r.Float64(), 'f', 6, 64))

	req, err := a.client.NewRequest(http.MethodPut, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return &Deposit{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return &Deposit{req, &a.client.Client}, nil
}

// Withdraw creates a new Withdrawal struct. Monzo uses a 'dedupe_id'
// to ensure that the request is idempotent, so the withdrawal is
// not ran when it is created. To action the withdrawal, call
// the `Run` method on it.
func (a Account) Withdraw(p Pot, amt int) (*Withdrawal, error) {
	endpoint := "/pots/" + p.ID + "/withdraw"
	data := url.Values{}
	data.Add("destination_account_id", a.ID)
	data.Add("amount", strconv.Itoa(amt))

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	data.Add("dedupe_id", strconv.FormatFloat(r.Float64(), 'f', 6, 64))

	req, err := a.client.NewRequest(http.MethodPut, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return &Withdrawal{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return &Withdrawal{req, &a.client.Client}, nil
}

// unwrapJSON takes a JSON response and unmarshals the contents
// of the first item. Responses from Monzo are wrapped in a key
// pertaining to the resource, which needs removing before
// unmarshalling.
func unwrapJSON(data []byte, wrapper string, v interface{}) error {
	var objmap map[string]*json.RawMessage
	if err := json.Unmarshal(data, &objmap); err != nil {
		return err
	}

	if err := json.Unmarshal(*objmap[wrapper], v); err != nil {
		return err
	}

	return nil
}
