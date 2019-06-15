package monzo

// Currency is a country code for a specific currency.
type Currency string

// Currencies that can be passed to the Monzo API.
//
// More will probably exist but haven't been tested.
const (
	CurrencyGBP Currency = "GBP"
	CurrencyUSD Currency = "USD"
)
