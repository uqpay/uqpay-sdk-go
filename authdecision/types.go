// Package authdecision handles PGP-encrypted real-time card authorization
// decisions sent by UQPAY to a customer's HTTP endpoint.
package authdecision

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

var decimalPattern = regexp.MustCompile(`^-?(0|[1-9][0-9]*)(\.[0-9]+)?([eE][+-]?[0-9]+)?$`)

// KeyPair contains ASCII-armored RSA keys for the authorization decision flow.
type KeyPair struct {
	PublicKey  string
	PrivateKey string
}

// Config contains the PGP keys used to process authorization decisions.
// PrivateKey and UQPayPublicKey accept either ASCII-armored keys or paths ending
// in .asc, .pgp, or .gpg.
type Config struct {
	PrivateKey     string
	UQPayPublicKey string
	Passphrase     string
}

// Transaction is the decrypted authorization decision request from UQPAY.
// Monetary values are strings to preserve decimal precision. The decoder also
// accepts JSON numbers for compatibility with older payload examples.
type Transaction struct {
	TransactionID                   string `json:"transaction_id"`
	TransactionType                 int    `json:"transaction_type"`
	CardID                          string `json:"card_id"`
	ProcessingCode                  string `json:"processing_code"`
	BillingAmount                   string `json:"billing_amount"`
	TransactionAmount               string `json:"transaction_amount"`
	AuthAmount                      string `json:"auth_amount"`
	DateOfTransaction               string `json:"date_of_transaction"`
	BillingCurrencyCode             string `json:"billing_currency_code"`
	TransactionCurrencyCode         string `json:"transaction_currency_code"`
	AuthCurrencyCode                string `json:"auth_currency_code"`
	CardBalance                     string `json:"card_balance"`
	MerchantID                      string `json:"merchant_id"`
	MerchantName                    string `json:"merchant_name"`
	MerchantCategoryCode            string `json:"merchant_category_code"`
	MerchantCity                    string `json:"merchant_city"`
	MerchantCountry                 string `json:"merchant_country"`
	TerminalID                      string `json:"terminal_id"`
	PosEntryMode                    string `json:"pos_entry_mode"`
	PosConditionCode                string `json:"pos_condition_code"`
	PosEnv                          string `json:"pos_env"`
	ECI                             string `json:"eci"`
	PinEntryCapability              string `json:"pin_entry_capability"`
	RetrievalReferenceNumber        string `json:"retrieval_reference_number"`
	SystemTraceAuditNumber          string `json:"system_trace_audit_number"`
	AcquiringInstitutionCountryCode string `json:"acquiring_institution_country_code"`
	AcquiringInstitutionID          string `json:"acquiring_institution_id"`
	WalletType                      string `json:"wallet_type"`
}

// UnmarshalJSON accepts monetary values encoded as either JSON strings or JSON
// numbers while preserving their original decimal representation.
func (t *Transaction) UnmarshalJSON(data []byte) error {
	type transactionWire struct {
		TransactionID                   string       `json:"transaction_id"`
		TransactionType                 int          `json:"transaction_type"`
		CardID                          string       `json:"card_id"`
		ProcessingCode                  string       `json:"processing_code"`
		BillingAmount                   decimalValue `json:"billing_amount"`
		TransactionAmount               decimalValue `json:"transaction_amount"`
		AuthAmount                      decimalValue `json:"auth_amount"`
		DateOfTransaction               string       `json:"date_of_transaction"`
		BillingCurrencyCode             string       `json:"billing_currency_code"`
		TransactionCurrencyCode         string       `json:"transaction_currency_code"`
		AuthCurrencyCode                string       `json:"auth_currency_code"`
		CardBalance                     decimalValue `json:"card_balance"`
		MerchantID                      string       `json:"merchant_id"`
		MerchantName                    string       `json:"merchant_name"`
		MerchantCategoryCode            string       `json:"merchant_category_code"`
		MerchantCity                    string       `json:"merchant_city"`
		MerchantCountry                 string       `json:"merchant_country"`
		TerminalID                      string       `json:"terminal_id"`
		PosEntryMode                    string       `json:"pos_entry_mode"`
		PosConditionCode                string       `json:"pos_condition_code"`
		PosEnv                          string       `json:"pos_env"`
		ECI                             string       `json:"eci"`
		PinEntryCapability              string       `json:"pin_entry_capability"`
		RetrievalReferenceNumber        string       `json:"retrieval_reference_number"`
		SystemTraceAuditNumber          string       `json:"system_trace_audit_number"`
		AcquiringInstitutionCountryCode string       `json:"acquiring_institution_country_code"`
		AcquiringInstitutionID          string       `json:"acquiring_institution_id"`
		WalletType                      string       `json:"wallet_type"`
	}

	var wire transactionWire
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&wire); err != nil {
		return err
	}

	*t = Transaction{
		TransactionID:                   wire.TransactionID,
		TransactionType:                 wire.TransactionType,
		CardID:                          wire.CardID,
		ProcessingCode:                  wire.ProcessingCode,
		BillingAmount:                   string(wire.BillingAmount),
		TransactionAmount:               string(wire.TransactionAmount),
		AuthAmount:                      string(wire.AuthAmount),
		DateOfTransaction:               wire.DateOfTransaction,
		BillingCurrencyCode:             wire.BillingCurrencyCode,
		TransactionCurrencyCode:         wire.TransactionCurrencyCode,
		AuthCurrencyCode:                wire.AuthCurrencyCode,
		CardBalance:                     string(wire.CardBalance),
		MerchantID:                      wire.MerchantID,
		MerchantName:                    wire.MerchantName,
		MerchantCategoryCode:            wire.MerchantCategoryCode,
		MerchantCity:                    wire.MerchantCity,
		MerchantCountry:                 wire.MerchantCountry,
		TerminalID:                      wire.TerminalID,
		PosEntryMode:                    wire.PosEntryMode,
		PosConditionCode:                wire.PosConditionCode,
		PosEnv:                          wire.PosEnv,
		ECI:                             wire.ECI,
		PinEntryCapability:              wire.PinEntryCapability,
		RetrievalReferenceNumber:        wire.RetrievalReferenceNumber,
		SystemTraceAuditNumber:          wire.SystemTraceAuditNumber,
		AcquiringInstitutionCountryCode: wire.AcquiringInstitutionCountryCode,
		AcquiringInstitutionID:          wire.AcquiringInstitutionID,
		WalletType:                      wire.WalletType,
	}
	return nil
}

type decimalValue string

func (v *decimalValue) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		return fmt.Errorf("amount must be a string or number")
	}
	if len(data) > 0 && data[0] == '"' {
		var value string
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		if !decimalPattern.MatchString(value) {
			return fmt.Errorf("invalid decimal %q", value)
		}
		*v = decimalValue(value)
		return nil
	}

	var number json.Number
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&number); err != nil {
		return fmt.Errorf("amount must be a string or number: %w", err)
	}
	*v = decimalValue(number.String())
	return nil
}

// Result is the customer's approve or decline decision.
type Result struct {
	ResponseCode       string
	PartnerReferenceID string
}

// DecisionFunc receives a decrypted transaction and returns an authorization
// response. Implementations should honor ctx and finish before the configured
// UQPAY timeout window.
type DecisionFunc func(ctx context.Context, transaction Transaction) (Result, error)

// HandlerOptions configures the net/http authorization decision handler.
type HandlerOptions struct {
	Decide DecisionFunc
	// DecisionTimeout should be shorter than the 1-5 second timeout configured
	// with UQPAY. A zero value relies on the request context without adding a
	// separate SDK deadline.
	DecisionTimeout time.Duration
	// MaxBodyBytes limits the encrypted request body. Zero uses 1 MiB.
	MaxBodyBytes int64
	// OnError is called before the handler aborts the HTTP response. Aborting
	// allows UQPAY's configured timeout action to decide the transaction.
	OnError func(error)
}
