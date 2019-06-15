package monzo

import (
	"encoding/json"
	"fmt"
)

// Item represents a single line on a Receipt.
type Item struct {
	description string
	amount      int
	currency    Currency
	quantity    int
	unit        string

	// subItems can only be nested one level deep.
	//
	// Monzo does not complain if a nested subItem is passed,
	// but this should probably be checked at library-level.
	// TODO: Prevent subItems being nested.
	subItems []*Item
}

// Receipt represents a summary of spending against a transaction.
type Receipt struct {
	TransactionID string
	ExternalID    string
	Total         int
	Currency      Currency
	Items         []*Item
}

// MakeReceiptItem creates an Item representing a single line
// of a receipt.
func MakeReceiptItem(desc string, amount int, currency Currency) *Item {
	return &Item{
		description: desc,
		amount:      amount,
		currency:    currency,
		unit:        "unit",
		quantity:    1,
	}
}

// AddSubItem adds a sub item to an existing Item. Monzo only
// supports one level of nesting Items.
func (i *Item) AddSubItem(subitem *Item) {
	i.subItems = append(i.subItems, subitem)
}

// Unit adds a unit measurement to an Item.
// For example, "kgs" or "piece".
func (i *Item) Unit(unit string) {
	i.unit = unit
}

// Quantity adds a quantity to an Item. The default quantity
// when creating an item is one.
func (i *Item) Quantity(quant int) {
	i.quantity = quant
}

// MakeReceipt creates a Receipt to store against a Transaction.
func MakeReceipt(externalID string) *Receipt {
	return &Receipt{
		ExternalID: externalID,
	}
}

// AddItem adds a number of Items to a Receipt.
func (r *Receipt) AddItem(is ...*Item) {
	r.Items = append(r.Items, is...)
}

// SetTransaction determines which Transaction the Receipt
// should be saved against.
func (r *Receipt) SetTransaction(transaction string) {
	r.TransactionID = transaction
}

// MarshalJSON converts the Receipt into a format usable by
// the Monzo API.
func (r *Receipt) MarshalJSON() ([]byte, error) {
	if len(r.Items) == 0 {
		return nil, fmt.Errorf("a receipt must contain at least one item")
	}

	r.calculateTotal()
	r.determineCurrency()

	return json.Marshal(struct {
		TransactionID string   `json:"transaction_id"`
		ExternalID    string   `json:"external_id"`
		Total         int      `json:"total"`
		Currency      Currency `json:"currency"`
		Items         []*Item  `json:"items"`
	}{
		TransactionID: r.TransactionID,
		ExternalID:    r.ExternalID,
		Total:         r.Total,
		Currency:      r.Currency,
		Items:         r.Items,
	})
}

// calculateTotal saves the total cost of all receipt Items
// in the Receipt's Total field.
func (r *Receipt) calculateTotal() {
	r.Total = 0

	for _, item := range r.Items {
		r.Total = r.Total + item.amount
	}
}

// determineCurrency sets the currency for the receipt. A receipt
// can only be for a single currency, so just grab the first
// item's currency and use that.
func (r *Receipt) determineCurrency() {
	r.Currency = r.Items[0].currency
}

// MarshalJSON converts the Item into a json object that can
// be read by the Monzo API.
func (i *Item) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Description string   `json:"description"`
		Quantity    int      `json:"quantity"`
		Unit        string   `json:"unit"`
		Amount      int      `json:"amount"`
		Currency    Currency `json:"currency"`
		SubItems    []*Item  `json:"sub_items"`
	}{
		Description: i.description,
		Quantity:    i.quantity,
		Unit:        i.unit,
		Amount:      i.amount,
		Currency:    i.currency,
		SubItems:    i.subItems,
	})
}
