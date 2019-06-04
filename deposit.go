package monzo

import (
	"bytes"
	"fmt"
	"net/http"
)

// Deposit represents a deposit that is ready to be made against
// Monzo. It is read-only. There are no side effects if the
// deposit is not ran.
type Deposit struct {
	Request *http.Request
	Client  *http.Client
}

// Run executes the deposit against the Monzo API. An error is
// only returned if the deposit fails to run. If the deposit
// has already ran against the account, it is not ran again
// and an error is not returned.
func (d Deposit) Run() error {
	resp, _ := d.Client.Do(d.Request)

	b := new(bytes.Buffer)
	b.ReadFrom(resp.Body)
	str := b.String()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch accounts: %s", str)
	}

	return nil
}
