package banking

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// PayoutsClient handles payout operations
type PayoutsClient struct {
	client *common.APIClient
}

// PayoutConversion represents the conversion details in a payout
type PayoutConversion struct {
	CurrencyPair string `json:"currency_pair"` // e.g. "USDEUR"
	ClientRate   string `json:"client_rate"`   // exchange rate applied
}

// Payout represents a payout transaction
type Payout struct {
	PayoutID              string            `json:"payout_id"`                         // UUID, unique payout identifier
	ShortReferenceID      string            `json:"short_reference_id"`                // system-generated reference, e.g. "P220406-LLCVLRM"
	UniqueRequestID       string            `json:"unique_request_id,omitempty"`       // UUID, idempotency key from creation request
	PayoutCurrency        string            `json:"payout_currency"`                   // ISO 4217 currency code, beneficiary receives
	PayoutAmount          string            `json:"payout_amount"`                     // decimal, amount in payout_currency
	FeeAmount             string            `json:"fee_amount"`                        // decimal, fee charged for this payout
	FeePaidBy             string            `json:"fee_paid_by"`                       // OURS
	FeeCurrency           string            `json:"fee_currency"`                      // ISO 4217 currency code for the fee
	PayoutDate            string            `json:"payout_date"`                       // YYYY-MM-DD, date system submits payment
	PayoutMethod          string            `json:"payout_method"`                     // LOCAL or SWIFT
	PayoutReason          string            `json:"payout_reason"`                     // max 200 chars, transaction reason
	PayoutReference       string            `json:"payout_reference"`                  // max 100 chars, bank reference visible to beneficiary
	PayoutStatus          string            `json:"payout_status"`                     // READY_TO_SEND, PENDING, REJECTED, FAILED, or COMPLETED
	FailureReturnedAmount string            `json:"failure_returned_amount,omitempty"` // optional, amount returned on failure
	FailureReason         string            `json:"failure_reason,omitempty"`          // optional, explanation of failure
	QuoteID               string            `json:"quote_id,omitempty"`                // optional, UUID of pre-created quote
	PurposeCode           string            `json:"purpose_code,omitempty"`            // e.g. GOODS_PURCHASED, WAGES_SALARY
	Conversion            *PayoutConversion `json:"conversion,omitempty"`              // optional, present for cross-currency payouts
	CreateTime            string            `json:"create_time"`                       // ISO 8601 timestamp
	UpdateTime            string            `json:"update_time,omitempty"`             // ISO 8601 timestamp
	CompleteTime          *string           `json:"complete_time"`                     // nullable, ISO 8601 timestamp when completed
}

// PayoutInlineBeneficiary represents an inline beneficiary for payout creation
type PayoutInlineBeneficiary struct {
	EntityType     string          `json:"entity_type"`               // required, INDIVIDUAL or COMPANY
	FirstName      string          `json:"first_name,omitempty"`      // required if INDIVIDUAL, max 45 chars
	LastName       string          `json:"last_name,omitempty"`       // required if INDIVIDUAL, max 45 chars
	CompanyName    string          `json:"company_name,omitempty"`    // required if COMPANY, max 120 chars
	IDNumber       string          `json:"id_number,omitempty"`       // required for Mainland China residents with CNH/LOCAL
	Nickname       string          `json:"nickname,omitempty"`        // optional, max 120 chars
	Email          string          `json:"email,omitempty"`           // optional, valid email address
	PaymentMethod  string          `json:"payment_method"`            // required, LOCAL or SWIFT
	BankDetails    *BankDetails    `json:"bank_details"`              // required
	Address        *Address        `json:"address"`                   // required
	AdditionalInfo *AdditionalInfo `json:"additional_info,omitempty"` // optional
}

// PayoutDocumentation represents documentation attached to a payout
type PayoutDocumentation struct {
	File   string `json:"file,omitempty"`    // conditional, base64-encoded with MIME prefix (e.g. "data:image/jpeg;base64,...")
	FileID string `json:"file_id,omitempty"` // conditional, UUID of previously uploaded file
}

// CreatePayoutRequest represents a payout creation request
type CreatePayoutRequest struct {
	Currency        string                   `json:"currency"`                  // required, ISO 4217 currency code the payer sends
	Amount          string                   `json:"amount"`                    // required, decimal sending amount
	QuoteID         string                   `json:"quote_id,omitempty"`        // conditional, UUID required for cross-currency payouts
	PayoutCurrency  string                   `json:"payout_currency,omitempty"` // conditional, ISO 4217 must match quote's buy_currency
	PayoutAmount    string                   `json:"payout_amount,omitempty"`   // conditional, decimal must match quote's buy_amount
	PurposeCode     string                   `json:"purpose_code"`              // required, e.g. GOODS_PURCHASED, WAGES_SALARY, PERSONAL_REMITTANCE
	PayoutReference string                   `json:"payout_reference"`          // required, max 100 chars, appears in recipient's bank records
	FeePaidBy       string                   `json:"fee_paid_by"`               // required, OURS (SHARED currently unavailable)
	PayoutDate      string                   `json:"payout_date"`               // required, YYYY-MM-DD date system submits payment
	BeneficiaryID   string                   `json:"beneficiary_id,omitempty"`  // conditional, UUID of existing beneficiary; exclude if using inline beneficiary
	Beneficiary     *PayoutInlineBeneficiary `json:"beneficiary,omitempty"`     // conditional, inline beneficiary details; exclude if using beneficiary_id
	IsPayer         string                   `json:"is_payer,omitempty"`        // deprecated, Y or N
	PayerID         string                   `json:"payer_id,omitempty"`        // deprecated, UUID of the payer
	Documentation   []PayoutDocumentation    `json:"documentation,omitempty"`   // conditional, required for INR transfers with IFSC clearing system
}

// CreatePayoutResponse represents a payout creation response
type CreatePayoutResponse struct {
	PayoutID         string `json:"payout_id"`          // UUID, unique payout identifier
	ShortReferenceID string `json:"short_reference_id"` // system-generated reference, e.g. "P220406-LLCVLRM"
	PayoutStatus     string `json:"payout_status"`      // READY_TO_SEND, PENDING, REJECTED, FAILED, or COMPLETED
}

// ListPayoutsRequest represents a payout list request
type ListPayoutsRequest struct {
	PageSize     int    `json:"page_size"`     // required, 10-100 items per page
	PageNumber   int    `json:"page_number"`   // required, >= 1
	StartTime    string `json:"start_time"`    // optional, ISO 8601; defaults to last 30 days if omitted with end_time
	EndTime      string `json:"end_time"`      // optional, ISO 8601; defaults to last 30 days if omitted with start_time
	PayoutStatus string `json:"payout_status"` // optional, READY_TO_SEND | PENDING | REJECTED | FAILED | COMPLETED
}

// ListPayoutsResponse represents a payout list response
type ListPayoutsResponse struct {
	TotalPages int      `json:"total_pages"` // total number of available pages
	TotalItems int      `json:"total_items"` // total count of matching payouts
	Data       []Payout `json:"data"`        // array of payout objects
}

// PayoutPayerDetail represents the payer details in a payout detail response
type PayoutPayerDetail struct {
	PayerID             string `json:"payer_id"`                       // UUID, payer identifier
	EntityType          string `json:"entity_type"`                    // INDIVIDUAL or COMPANY
	Country             string `json:"country"`                        // ISO 3166-1 alpha-2 country code
	FirstName           string `json:"first_name,omitempty"`           // optional, present if INDIVIDUAL
	LastName            string `json:"last_name,omitempty"`            // optional, present if INDIVIDUAL
	CompanyName         string `json:"company_name,omitempty"`         // optional, present if COMPANY
	City                string `json:"city,omitempty"`                 // optional, payer city
	Address             string `json:"address,omitempty"`              // optional, payer street address
	State               string `json:"state,omitempty"`                // optional, payer state/province
	PostalCode          string `json:"postal_code,omitempty"`          // optional, payer postal/ZIP code
	DateBirth           string `json:"date_birth,omitempty"`           // optional, payer date of birth
	IdentificationType  string `json:"identification_type,omitempty"`  // optional, e.g. PASSPORT, NATIONAL_ID
	IdentificationValue string `json:"identification_value,omitempty"` // optional, identification document number
}

// PayoutDetailResponse represents a detailed payout response
type PayoutDetailResponse struct {
	Payout                                       // embedded payout fields
	AmountPayerPays           string             `json:"amount_payer_pays,omitempty"`           // decimal, total amount payer remits
	SourceCurrency            string             `json:"source_currency,omitempty"`             // ISO 4217 currency code, payer's currency
	SourceAmount              string             `json:"source_amount,omitempty"`               // decimal, amount paid by payer
	AmountBeneficiaryReceives string             `json:"amount_beneficiary_receives,omitempty"` // decimal, amount received by beneficiary
	Payer                     *PayoutPayerDetail `json:"payer,omitempty"`                       // payer entity details
	Beneficiary               *Beneficiary       `json:"beneficiary,omitempty"`                 // beneficiary entity details
}

// Create creates a new payout
// Optional RequestOptions can be provided to set custom headers like x-idempotency-key or x-on-behalf-of
func (c *PayoutsClient) Create(ctx context.Context, req *CreatePayoutRequest, opts ...*common.RequestOptions) (*CreatePayoutResponse, error) {
	var resp CreatePayoutResponse
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.PostWithOptions(ctx, "/v1/payouts", req, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to create payout: %w", err)
	}
	return &resp, nil
}

// List lists payouts with filters and pagination
// Optional RequestOptions can be provided to set custom headers like x-on-behalf-of
func (c *PayoutsClient) List(ctx context.Context, req *ListPayoutsRequest, opts ...*common.RequestOptions) (*ListPayoutsResponse, error) {
	var resp ListPayoutsResponse
	path := fmt.Sprintf("/v1/payouts?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)

	if req.StartTime != "" {
		path += fmt.Sprintf("&start_time=%s", req.StartTime)
	}
	if req.EndTime != "" {
		path += fmt.Sprintf("&end_time=%s", req.EndTime)
	}
	if req.PayoutStatus != "" {
		path += fmt.Sprintf("&payout_status=%s", req.PayoutStatus)
	}

	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.GetWithOptions(ctx, path, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to list payouts: %w", err)
	}
	return &resp, nil
}

// Get retrieves a specific payout by ID
// Optional RequestOptions can be provided to set custom headers like x-on-behalf-of
func (c *PayoutsClient) Get(ctx context.Context, payoutID string, opts ...*common.RequestOptions) (*PayoutDetailResponse, error) {
	var resp PayoutDetailResponse
	path := fmt.Sprintf("/v1/payouts/%s", payoutID)
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.GetWithOptions(ctx, path, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to get payout: %w", err)
	}
	return &resp, nil
}
