package monzo

// Balance returns the data from the Monzo API. All values are
// returned in pence.
type Balance struct {
	Balance     int
	Total       int `json:"total_balance"`
	WithSavings int `json:"balance_including_flexible_savings"`
	Currency    Currency
}
