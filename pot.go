package monzo

import (
	"bytes"
	"fmt"
	"net/http"
)

// Pot represents a Monzo pot.
type Pot struct {
	ID       string
	Name     string
	Balance  int
	Currency string
	Created  string
	Updated  string
	Deleted  bool
}

// AllPots retrieves all the users pots from the Monzo API,
// even the ones that have been deleted.
func (c *Client) AllPots() ([]Pot, error) {
	return c.pots()
}

// Pots returns a slice of Pots that belong to the user.
// Only pots that haven't been deleted are returned.
func (c *Client) Pots() ([]Pot, error) {
	pots, err := c.pots()
	if err != nil {
		return pots, err
	}

	var filtered []Pot
	for _, p := range pots {
		if !p.Deleted {
			filtered = append(filtered, p)
		}
	}

	return filtered, nil
}

// Pot returns a single Pot from the Monzo API.
func (c *Client) Pot(id string) (Pot, error) {
	pots, err := c.pots()
	if err != nil {
		return Pot{}, err
	}

	// Monzo doesn't have the capability to retrieve a single
	// pot, so we need to get all the user's pots and filter
	// them down using the provided pot id.
	for _, p := range pots {
		if p.ID == id {
			return p, nil
		}
	}

	return Pot{}, fmt.Errorf("no pot found with ID %s", id)
}

func (c *Client) pots() ([]Pot, error) {
	req, err := c.resourceRequest("pots")
	if err != nil {
		return nil, err
	}

	resp, _ := c.Do(req)

	b := new(bytes.Buffer)
	b.ReadFrom(resp.Body)
	str := b.String()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch pots: %s", str)
	}

	bytes := b.Bytes()
	var pots []Pot
	if err := unwrapJSON(bytes, "pots", &pots); err != nil {
		return nil, err
	}

	return pots, nil
}
