package webhook

// PayoutConversion represents the conversion details in a payout webhook.
type PayoutConversion struct {
	ClientRate   string `json:"client_rate"`
	CurrencyPair string `json:"currency_pair"`
}

// PayoutData represents the data payload for payout webhook events.
type PayoutData struct {
	AccountID                  string            `json:"account_id"`
	Amount                     string            `json:"amount"`
	AuthorisationStepsRequired string            `json:"authorisation_steps_required"`
	BeneficiaryID              string            `json:"beneficiary_id"`
	Conversion                 *PayoutConversion `json:"conversion,omitempty"`
	Currency                   string            `json:"currency"`
	DirectID                   string            `json:"direct_id"`
	FailureReason              string            `json:"failure_reason"`
	FailureReturnedAmount      string            `json:"failure_returned_amount"`
	FeeAmount                  string            `json:"fee_amount"`
	FeeCurrency                string            `json:"fee_currency"`
	FeePaidBy                  string            `json:"fee_paid_by"`
	PaymentDate                string            `json:"payment_date"`
	PaymentType                string            `json:"payment_type"`
	PayoutAmount               string            `json:"payout_amount"`
	PayoutCurrency             string            `json:"payout_currency"`
	PayoutID                   string            `json:"payout_id"`
	PayoutWay                  string            `json:"payout_way"`
	QuoteID                    string            `json:"quote_id"`
	Reason                     string            `json:"reason"`
	Reference                  string            `json:"reference"`
	ShortReference             string            `json:"short_reference"`
	ShortReferenceID           string            `json:"short_reference_id"`
	Status                     string            `json:"status"`
	UniqueRequestID            string            `json:"unique_request_id"`
}
