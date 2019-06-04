package monzo

import (
	"bytes"
	"fmt"
	"net/http"
)

// Withdrawal represents moving money out of a pot and back
// into a Monzo account. It is read-only. There are no side
// effects if the withdrawal is not ran.
type Withdrawal struct {
	Request *http.Request
	Client  *http.Client
}

// Run executes the withdrawal against the Monzo API. An error is
// only returned if the withdrawal fails to run. If the action
// has already ran against the account, it is not ran again
// and an error is not returned.
func (d Withdrawal) Run() error {
	resp, _ := d.Client.Do(d.Request)

	b := new(bytes.Buffer)
	b.ReadFrom(resp.Body)
	str := b.String()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to withdraw: %s", str)
	}

	return nil
}
