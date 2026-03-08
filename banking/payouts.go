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

// Payout represents a payout transaction
type Payout struct {
	PayoutID              string  `json:"payout_id"`
	ShortReferenceID      string  `json:"short_reference_id"`
	UniqueRequestID       string  `json:"unique_request_id,omitempty"`
	PayoutCurrency        string  `json:"payout_currency"`
	PayoutAmount          string  `json:"payout_amount"`
	FeeAmount             string  `json:"fee_amount"`
	FeePaidBy             string  `json:"fee_paid_by"`
	FeeCurrency           string  `json:"fee_currency"`
	PayoutDate            string  `json:"payout_date"`
	PayoutMethod          string  `json:"payout_method"`
	PayoutReason          string  `json:"payout_reason"`
	PayoutReference       string  `json:"payout_reference"`
	PayoutStatus          string  `json:"payout_status"` // PENDING, READY_TO_SEND, COMPLETED, FAILED, CANCELLED
	FailureReturnedAmount string  `json:"failure_returned_amount,omitempty"`
	FailureReason         string  `json:"failure_reason,omitempty"`
	QuoteID               string  `json:"quote_id,omitempty"`
	PurposeCode           string  `json:"purpose_code,omitempty"`
	CreateTime            string  `json:"create_time"`
	UpdateTime            string  `json:"update_time,omitempty"`
	CompleteTime          *string `json:"complete_time"` // nullable
}

// PayoutInlineBeneficiary represents an inline beneficiary for payout creation
type PayoutInlineBeneficiary struct {
	EntityType     string          `json:"entity_type"`                // required: INDIVIDUAL or COMPANY
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
	IsPayer         string                   `json:"is_payer,omitempty"`        // optional, "Y" or "N"
	Documentation   []PayoutDocumentation    `json:"documentation,omitempty"`   // optional
}

// CreatePayoutResponse represents a payout creation response
type CreatePayoutResponse struct {
	PayoutID         string `json:"payout_id"`
	ShortReferenceID string `json:"short_reference_id"`
	Status           string `json:"status"`
	CreateTime       string `json:"create_time"`
}

// ListPayoutsRequest represents a payout list request
type ListPayoutsRequest struct {
	PageSize      int    `json:"page_size"`      // required, 10-100
	PageNumber    int    `json:"page_number"`    // required, >=1
	StartTime     string `json:"start_time"`     // optional, ISO8601
	EndTime       string `json:"end_time"`       // optional, ISO8601
	PayoutStatus  string `json:"payout_status"`  // optional: PENDING, PROCESSING, COMPLETED, FAILED, CANCELLED, ALL
	Currency      string `json:"currency"`       // optional, filter by currency
	BeneficiaryID string `json:"beneficiary_id"` // optional, filter by beneficiary
}

// ListPayoutsResponse represents a payout list response
type ListPayoutsResponse struct {
	TotalPages int      `json:"total_pages"`
	TotalItems int      `json:"total_items"`
	Data       []Payout `json:"data"`
}

// PayoutDetailResponse represents a detailed payout response
type PayoutDetailResponse struct {
	Payout
	// Additional fields that might be present in detail view
	TransactionDetails *TransactionDetails `json:"transaction_details,omitempty"`
}

// TransactionDetails represents additional transaction information
type TransactionDetails struct {
	TransactionID     string `json:"transaction_id,omitempty"`
	ProcessingTime    string `json:"processing_time,omitempty"`
	SettlementTime    string `json:"settlement_time,omitempty"`
	ExchangeRate      string `json:"exchange_rate,omitempty"`
	ProcessorResponse string `json:"processor_response,omitempty"`
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
	if req.Currency != "" {
		path += fmt.Sprintf("&currency=%s", req.Currency)
	}
	if req.BeneficiaryID != "" {
		path += fmt.Sprintf("&beneficiary_id=%s", req.BeneficiaryID)
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
