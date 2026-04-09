package authdecision

import "context"

// KeyPair holds an armored PGP public and private key pair.
type KeyPair struct {
	PublicKey  string
	PrivateKey string
}

// Config holds the PGP configuration for authorization decision processing.
// PrivateKey and UQPayPublicKey accept armored PGP key strings or file paths
// ending in .asc, .pgp, or .gpg.
type Config struct {
	PrivateKey     string
	UQPayPublicKey string
	Passphrase     string
}

// Transaction represents an incoming card authorization transaction.
type Transaction struct {
	TransactionID                  string  `json:"transaction_id"`
	TransactionType                int     `json:"transaction_type"`
	CardID                         string  `json:"card_id"`
	ProcessingCode                 string  `json:"processing_code"`
	BillingAmount                  float64 `json:"billing_amount"`
	TransactionAmount              float64 `json:"transaction_amount"`
	AuthAmount                     float64 `json:"auth_amount"`
	CardBalance                    float64 `json:"card_balance"`
	DateOfTransaction              string  `json:"date_of_transaction"`
	BillingCurrencyCode            string  `json:"billing_currency_code"`
	TransactionCurrencyCode        string  `json:"transaction_currency_code"`
	AuthCurrencyCode               string  `json:"auth_currency_code"`
	MerchantID                     string  `json:"merchant_id"`
	MerchantName                   string  `json:"merchant_name"`
	MerchantCategoryCode           string  `json:"merchant_category_code"`
	MerchantCity                   string  `json:"merchant_city"`
	MerchantCountry                string  `json:"merchant_country"`
	TerminalID                     string  `json:"terminal_id"`
	PosEntryMode                   string  `json:"pos_entry_mode"`
	PosConditionCode               string  `json:"pos_condition_code"`
	PinEntryCapability             string  `json:"pin_entry_capability"`
	RetrievalReferenceNumber       string  `json:"retrieval_reference_number"`
	SystemTraceAuditNumber         string  `json:"system_trace_audit_number"`
	AcquiringInstitutionCountryCode string  `json:"acquiring_institution_country_code"`
	AcquiringInstitutionID         string  `json:"acquiring_institution_id"`
	WalletType                     string  `json:"wallet_type"`
}

// Result represents the authorization decision response.
type Result struct {
	ResponseCode       string `json:"response_code"`
	PartnerReferenceID string `json:"partner_reference_id"`
}

// HandlerOptions configures the authorization decision HTTP handler.
type HandlerOptions struct {
	// Decide is called for each incoming transaction to make an authorization decision.
	Decide func(ctx context.Context, tx Transaction) (Result, error)
	// OnError is called when any error occurs during request processing.
	OnError func(err error)
}
