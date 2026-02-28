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

// PayoutBeneficiary represents a payout beneficiary
type PayoutBeneficiary struct {
	BeneficiaryID   string                `json:"beneficiary_id,omitempty"`
	BeneficiaryName string                `json:"beneficiary_name,omitempty"`
	BankDetails     *PayoutBankDetails    `json:"bank_details,omitempty"`
	WalletDetails   *WalletDetails        `json:"wallet_details,omitempty"`
	ContactDetails  *PayoutContactDetails `json:"contact_details,omitempty"`
}

// PayoutBankDetails represents bank account information for payouts
type PayoutBankDetails struct {
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	BankCode      string `json:"bank_code"`
	BankName      string `json:"bank_name,omitempty"`
	BranchCode    string `json:"branch_code,omitempty"`
	RoutingNumber string `json:"routing_number,omitempty"`
	SwiftCode     string `json:"swift_code,omitempty"`
	IBAN          string `json:"iban,omitempty"`
	AccountType   string `json:"account_type,omitempty"` // SAVINGS, CHECKING, BUSINESS
}

// WalletDetails represents mobile wallet information
type WalletDetails struct {
	WalletProvider string `json:"wallet_provider"` // MPESA, AIRTEL_MONEY, MTN_MONEY, etc.
	WalletNumber   string `json:"wallet_number"`
	WalletName     string `json:"wallet_name,omitempty"`
}

// PayoutContactDetails represents beneficiary contact information for payouts
type PayoutContactDetails struct {
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Address     string `json:"address,omitempty"`
	City        string `json:"city,omitempty"`
	Country     string `json:"country,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
}

// CreatePayoutRequest represents a payout creation request
type CreatePayoutRequest struct {
	// Option 1: Reference existing beneficiary
	BeneficiaryID string `json:"beneficiary_id,omitempty"`

	// Option 2: Provide full beneficiary details inline
	Beneficiary *PayoutBeneficiary `json:"beneficiary,omitempty"`

	// Required fields
	Currency      string `json:"currency"`       // required, e.g., "USD", "KES", "UGX"
	Amount        string `json:"amount"`         // required, decimal string
	PayoutPurpose string `json:"payout_purpose"` // required, e.g., "salary", "vendor_payment", "refund"

	// Optional fields
	Description string `json:"description,omitempty"`
	Reference   string `json:"reference,omitempty"` // Client's internal reference
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
