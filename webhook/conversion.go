package webhook

// ConversionData represents the conversion information in conversion webhook events.
// This is returned in the data field for conversion.trade.* events.
type ConversionData struct {
	// AccountID is the unique identifier for the account
	AccountID string `json:"account_id"`

	// AccountName is the display name of the account
	AccountName string `json:"account_name"`

	// BuyAmount is the amount being bought as a string (e.g., "38.51")
	BuyAmount string `json:"buy_amount"`

	// BuyCurrency is the ISO 4217 currency code for the buy side (e.g., "SGD")
	BuyCurrency string `json:"buy_currency"`

	// ClientRate is the exchange rate applied to the conversion
	ClientRate string `json:"client_rate"`

	// ConversionID is the unique identifier for the conversion
	ConversionID string `json:"conversion_id"`

	// ConversionStatus is the current status of the conversion (e.g., "TRADE_SETTLED")
	ConversionStatus string `json:"conversion_status"`

	// ConversionWay indicates how the conversion was created (e.g., "API", "WEB")
	ConversionWay string `json:"conversion_way"`

	// CreateTime is the creation timestamp
	CreateTime string `json:"create_time"`

	// Creator is the name of the user who created the conversion
	Creator string `json:"creator"`

	// DirectID is the direct account identifier
	DirectID string `json:"direct_id"`

	// SellAmount is the amount being sold as a string (e.g., "100")
	SellAmount string `json:"sell_amount"`

	// SellCurrency is the ISO 4217 currency code for the sell side (e.g., "USD")
	SellCurrency string `json:"sell_currency"`

	// SettleTime is the settlement timestamp (if settled)
	SettleTime string `json:"settle_time,omitempty"`

	// ShortReferenceID is a short reference identifier for the conversion
	ShortReferenceID string `json:"short_reference_id"`
}

// Conversion status constants
const (
	// ConversionStatusTradeSettled indicates the trade has been settled
	ConversionStatusTradeSettled = "TRADE_SETTLED"

	// ConversionStatusAwaitingFunds indicates the conversion is awaiting funds
	ConversionStatusAwaitingFunds = "AWAITING_FUNDS"

	// ConversionStatusFundsArrived indicates the funds have arrived
	ConversionStatusFundsArrived = "FUNDS_ARRIVED"

	// ConversionStatusPending indicates the conversion is pending
	ConversionStatusPending = "PENDING"

	// ConversionStatusCompleted indicates the conversion is completed
	ConversionStatusCompleted = "COMPLETED"

	// ConversionStatusCanceled indicates the conversion was canceled
	ConversionStatusCanceled = "CANCELED"

	// ConversionStatusFailed indicates the conversion failed
	ConversionStatusFailed = "FAILED"
)

// Conversion way constants
const (
	// ConversionWayAPI indicates the conversion was created via API
	ConversionWayAPI = "API"

	// ConversionWayWeb indicates the conversion was created via web interface
	ConversionWayWeb = "WEB"
)
