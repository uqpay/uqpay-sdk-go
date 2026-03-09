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
	CurrencyPair string `json:"currency_pair"`
	ClientRate   string `json:"client_rate"`
}

// Payout represents a payout transaction
type Payout struct {
	PayoutID              string            `json:"payout_id"`
	ShortReferenceID      string            `json:"short_reference_id"`
	UniqueRequestID       string            `json:"unique_request_id,omitempty"`
	PayoutCurrency        string            `json:"payout_currency"`
	PayoutAmount          string            `json:"payout_amount"`
	FeeAmount             string            `json:"fee_amount"`
	FeePaidBy             string            `json:"fee_paid_by"`
	FeeCurrency           string            `json:"fee_currency"`
	PayoutDate            string            `json:"payout_date"`
	PayoutMethod          string            `json:"payout_method"`
	PayoutReason          string            `json:"payout_reason"`
	PayoutReference       string            `json:"payout_reference"`
	PayoutStatus          string            `json:"payout_status"` // READY_TO_SEND, PENDING, REJECTED, FAILED, COMPLETED
	FailureReturnedAmount string            `json:"failure_returned_amount,omitempty"`
	FailureReason         string            `json:"failure_reason,omitempty"`
	QuoteID               string            `json:"quote_id,omitempty"`
	PurposeCode           string            `json:"purpose_code,omitempty"`
	Conversion            *PayoutConversion `json:"conversion,omitempty"`
	CreateTime            string            `json:"create_time"`
	UpdateTime            string            `json:"update_time,omitempty"`
	CompleteTime          *string           `json:"complete_time"` // nullable
}

// PayoutInlineBeneficiary represents an inline beneficiary for payout creation
type PayoutInlineBeneficiary struct {
	EntityType     string          `json:"entity_type"`               // required: INDIVIDUAL or COMPANY
	FirstName      string          `json:"first_name,omitempty"`      // required if INDIVIDUAL
	LastName       string          `json:"last_name,omitempty"`       // required if INDIVIDUAL
	CompanyName    string          `json:"company_name,omitempty"`    // required if COMPANY
	IDNumber       string          `json:"id_number,omitempty"`       // required when account currency = COP
	Nickname       string          `json:"nickname,omitempty"`        // optional
	Email          string          `json:"email,omitempty"`           // optional
	PaymentMethod  string          `json:"payment_method"`            // required: LOCAL or SWIFT
	BankDetails    *BankDetails    `json:"bank_details"`              // required (uses BankDetails from beneficiaries.go)
	Address        *Address        `json:"address"`                   // required (uses Address from beneficiaries.go)
	AdditionalInfo *AdditionalInfo `json:"additional_info,omitempty"` // optional
}

// PayoutDocumentation represents documentation attached to a payout
type PayoutDocumentation struct {
	File   string `json:"file,omitempty"`    // base64-encoded file
	FileID string `json:"file_id,omitempty"` // file ID from upload API
}

// CreatePayoutRequest represents a payout creation request
type CreatePayoutRequest struct {
	Currency        string                   `json:"currency"`                  // required, ISO 4217
	Amount          string                   `json:"amount"`                    // required
	QuoteID         string                   `json:"quote_id,omitempty"`        // optional, UUID
	PayoutCurrency  string                   `json:"payout_currency,omitempty"` // conditional, required when quote_id specified
	PayoutAmount    string                   `json:"payout_amount,omitempty"`   // conditional, required when quote_id specified
	PurposeCode     string                   `json:"purpose_code"`              // required
	PayoutReference string                   `json:"payout_reference"`          // required, max 100 chars
	FeePaidBy       string                   `json:"fee_paid_by"`               // required, "OURS"
	PayoutDate      string                   `json:"payout_date"`               // required, YYYY-MM-DD
	BeneficiaryID   string                   `json:"beneficiary_id,omitempty"`  // conditional, either this or beneficiary
	Beneficiary     *PayoutInlineBeneficiary `json:"beneficiary,omitempty"`     // conditional, inline beneficiary
	IsPayer         string                   `json:"is_payer,omitempty"`        // deprecated, "Y" or "N"
	PayerID         string                   `json:"payer_id,omitempty"`        // deprecated, UUID of the payer
	Documentation   []PayoutDocumentation    `json:"documentation,omitempty"`   // optional
}

// CreatePayoutResponse represents a payout creation response
type CreatePayoutResponse struct {
	PayoutID         string `json:"payout_id"`
	ShortReferenceID string `json:"short_reference_id"`
	PayoutStatus     string `json:"payout_status"`
}

// ListPayoutsRequest represents a payout list request
type ListPayoutsRequest struct {
	PageSize     int    `json:"page_size"`     // required, 10-100
	PageNumber   int    `json:"page_number"`   // required, >=1
	StartTime    string `json:"start_time"`    // optional, ISO8601
	EndTime      string `json:"end_time"`      // optional, ISO8601
	PayoutStatus string `json:"payout_status"` // optional: READY_TO_SEND, PENDING, REJECTED, FAILED, COMPLETED
}

// ListPayoutsResponse represents a payout list response
type ListPayoutsResponse struct {
	TotalPages int      `json:"total_pages"`
	TotalItems int      `json:"total_items"`
	Data       []Payout `json:"data"`
}

// PayoutPayerDetail represents the payer details in a payout detail response
type PayoutPayerDetail struct {
	PayerID             string `json:"payer_id"`
	EntityType          string `json:"entity_type"`
	Country             string `json:"country"`
	FirstName           string `json:"first_name,omitempty"`
	LastName            string `json:"last_name,omitempty"`
	CompanyName         string `json:"company_name,omitempty"`
	City                string `json:"city,omitempty"`
	Address             string `json:"address,omitempty"`
	State               string `json:"state,omitempty"`
	PostalCode          string `json:"postal_code,omitempty"`
	DateBirth           string `json:"date_birth,omitempty"`
	IdentificationType  string `json:"identification_type,omitempty"`
	IdentificationValue string `json:"identification_value,omitempty"`
}

// PayoutDetailResponse represents a detailed payout response
type PayoutDetailResponse struct {
	Payout
	AmountPayerPays           string             `json:"amount_payer_pays,omitempty"`
	SourceCurrency            string             `json:"source_currency,omitempty"`
	SourceAmount              string             `json:"source_amount,omitempty"`
	AmountBeneficiaryReceives string             `json:"amount_beneficiary_receives,omitempty"`
	Payer                     *PayoutPayerDetail `json:"payer,omitempty"`
	Beneficiary               *Beneficiary       `json:"beneficiary,omitempty"`
}

// Create creates a new payout
func (c *PayoutsClient) Create(ctx context.Context, req *CreatePayoutRequest) (*CreatePayoutResponse, error) {
	var resp CreatePayoutResponse
	if err := c.client.Post(ctx, "/v1/payouts", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create payout: %w", err)
	}
	return &resp, nil
}

// List lists payouts with filters and pagination
func (c *PayoutsClient) List(ctx context.Context, req *ListPayoutsRequest) (*ListPayoutsResponse, error) {
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

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list payouts: %w", err)
	}
	return &resp, nil
}

// Get retrieves a specific payout by ID
func (c *PayoutsClient) Get(ctx context.Context, payoutID string) (*PayoutDetailResponse, error) {
	var resp PayoutDetailResponse
	path := fmt.Sprintf("/v1/payouts/%s", payoutID)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get payout: %w", err)
	}
	return &resp, nil
}
