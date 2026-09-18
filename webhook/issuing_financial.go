package webhook

// NetworkProtectionFeeData is an issuing webhook payload. Decimal amounts are strings.
type NetworkProtectionFeeData struct {
	CardID              string `json:"card_id,omitempty"`
	AccountID           string `json:"account_id,omitempty"`
	BillingPeriod       string `json:"billing_period,omitempty"`
	CardCount           int    `json:"card_count,omitempty"`
	ActionCode          string `json:"action_code,omitempty"`
	TransactionID       string `json:"transaction_id,omitempty"`
	TransactionType     string `json:"transaction_type,omitempty"`
	Remark              string `json:"remark,omitempty"`
	TransactionAmount   string `json:"transaction_amount,omitempty"`
	TransactionCurrency string `json:"transaction_currency,omitempty"`
	TransactionTime     string `json:"transaction_time,omitempty"`
	PostedTime          string `json:"posted_time,omitempty"`
	TransactionStatus   string `json:"transaction_status,omitempty"`
	BalanceAmount       string `json:"balance_amount,omitempty"`
}

// IssuingTransferStatusChangedData is an issuing webhook payload. Decimal amounts are strings.
type IssuingTransferStatusChangedData struct {
	TransferID           string `json:"transfer_id,omitempty"`
	Status               string `json:"status,omitempty"`
	PreviousStatus       string `json:"previous_status,omitempty"`
	SourceAccountID      string `json:"source_account_id,omitempty"`
	DestinationAccountID string `json:"destination_account_id,omitempty"`
	Currency             string `json:"currency,omitempty"`
	Amount               string `json:"amount,omitempty"`
	Remark               string `json:"remark,omitempty"`
	SucceededAt          string `json:"succeeded_at,omitempty"`
	CreateTime           string `json:"create_time,omitempty"`
}
