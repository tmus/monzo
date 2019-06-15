package monzo

// Webhook is an endpoint that Monzo will send events to.
type Webhook struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	URL       string `json:"url"`
}
