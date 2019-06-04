package monzo

import "fmt"

// Balance returns the data from the Monzo API's /balance endpoint.
// All values are returned in pence.
type Balance struct {
	Balance     int
	Total       int `json:"total_balance"`
	WithSavings int `json:"balance_including_flexible_savings"`
	Currency    string
}

func (b Balance) String() string {
	return fmt.Sprintf(
		"%s %.2f (Total: %s %.2f)",
		b.Currency,
		float64(b.Balance/100),
		b.Currency,
		float64(b.Total/100),
	)
}
