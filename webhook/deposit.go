package webhook

// DepositData represents the data payload for deposit webhook events.
type DepositData struct {
	DirectID         string `json:"direct_id"`
	AccountID        string `json:"account_id"`
	AccountName      string `json:"account_name"`
	DepositID        string `json:"deposit_id"`
	ShortReferenceID string `json:"short_reference_id"`
	DepositCurrency  string `json:"deposit_currency"`
	DepositAmount    string `json:"deposit_amount"`
	DepositFee       string `json:"deposit_fee"`
	CreateTime       string `json:"create_time"`
	CompleteTime     string `json:"complete_time"`
	UpdateTime       string `json:"update_time"`
	DepositStatus    string `json:"deposit_status"`
	DepositReference string `json:"deposit_reference"`
}
