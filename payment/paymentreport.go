package payment

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// PaymentReportsClient handles payment reports operations
type PaymentReportsClient struct {
	client *common.APIClient
}

// ============================================================================
// Request Structures
// ============================================================================

// ListSettlementsRequest represents a settlements list request
type ListSettlementsRequest struct {
	PaymentIntentID   string `json:"payment_intent_id"`   // Optional. Filter by payment intent ID, e.g. "PI1945730395043532800"
	SettlementBatchID string `json:"settlement_batch_id"` // Optional. Filter by settlement batch ID, e.g. "SB1947180993781698560"
	SettledStartTime  string `json:"settled_start_time"`  // Optional. Inclusive start date filter, format: YYYY-MM-DD (UTC+8). Must be used together with settled_end_time; max range is one month
	SettledEndTime    string `json:"settled_end_time"`    // Optional. Inclusive end date filter, format: YYYY-MM-DD (UTC+8). Must be used together with settled_start_time; max range is one month
	PageSize          int    `json:"page_size"`           // Required. Number of items per page, range: 1-100
	PageNumber        int    `json:"page_number"`         // Required. Page number for pagination, minimum: 1
}

// ============================================================================
// Response Structures
// ============================================================================

// Settlement represents a settlement record
type Settlement struct {
	SettlementID          string `json:"settlement_id"`                     // Unique settlement identifier, UUID format
	AccountID             string `json:"account_id,omitempty"`              // Optional. Associated account identifier, UUID format
	AccountName           string `json:"account_name,omitempty"`            // Optional. Account display name
	SourceType            string `json:"source_type,omitempty"`             // Optional. Transaction source classification, e.g. "PAYMENT"
	TransactionType       string `json:"transaction_type,omitempty"`        // Optional. Transaction classification, e.g. "PAYMENT"
	MerchantOrderID       string `json:"merchant_order_id,omitempty"`       // Optional. Merchant system order reference
	PaymentIntentID       string `json:"payment_intent_id,omitempty"`       // Optional. Associated payment intent ID, e.g. "PI..."
	PaymentMethod         string `json:"payment_method,omitempty"`          // Optional. Payment method used, e.g. card, alipaycn, alipayhk, unionpay, wechatpay
	TransactionCreateDate string `json:"transaction_create_date,omitempty"` // Optional. Transaction creation timestamp, ISO 8601 format (UTC)
	TransactionAmount     string `json:"transaction_amount,omitempty"`      // Optional. Original transaction amount, numeric string
	TransactionCurrency   string `json:"transaction_currency,omitempty"`    // Optional. ISO 4217 currency code of the transaction
	TransactionDate       string `json:"transaction_date,omitempty"`        // Optional. Transaction completion timestamp, ISO 8601 format (UTC+8)
	SettlementAmount      string `json:"settlement_amount,omitempty"`       // Optional. Gross settlement amount before fees, numeric string
	SettlementCurrency    string `json:"settlement_currency,omitempty"`     // Optional. ISO 4217 currency code of the settlement
	NetSettlementAmount   string `json:"net_settlement_amount,omitempty"`   // Optional. Settlement amount after all fees deducted, numeric string
	ExchangeRate          string `json:"exchange_rate,omitempty"`           // Optional. FX rate applied between transaction and settlement currencies
	FeeCurrency           string `json:"fee_currency,omitempty"`            // Optional. ISO 4217 currency code for fee amounts
	InterchangeFee        string `json:"interchange_fee,omitempty"`         // Optional. Card network interchange fee, numeric string
	SchemeFee             string `json:"scheme_fee,omitempty"`              // Optional. Payment scheme fee, numeric string
	TranscationFee        string `json:"transcation_fee,omitempty"`         // Optional. Transaction processing fee, numeric string. Note: API field has typo ("transcation")
	ReturnFee             string `json:"return_fee,omitempty"`              // Optional. Refund-related fee, numeric string
	TotalFeeAmount        string `json:"total_fee_amount,omitempty"`        // Optional. Aggregate of all fees, numeric string
	SettlementStatus      string `json:"settlement_status,omitempty"`       // Optional. Settlement state, e.g. "SUCCESS"
	SettlementBatchID     string `json:"settlement_batch_id,omitempty"`     // Optional. Batch grouping identifier, e.g. "SB..."
	SettlementCreateDate  string `json:"settlement_create_date,omitempty"`  // Optional. Settlement record creation timestamp, ISO 8601 format (UTC+8)
	SettlementDate        string `json:"settlement_date,omitempty"`         // Optional. Settlement completion timestamp, ISO 8601 format (UTC+8)
}

// ListSettlementsResponse represents a paginated list of settlements
type ListSettlementsResponse struct {
	TotalPages int          `json:"total_pages"` // Total number of pages available
	TotalItems int          `json:"total_items"` // Total number of settlement records
	Data       []Settlement `json:"data"`        // List of settlement records for current page
}

// ============================================================================
// API Methods
// ============================================================================

// ListSettlements returns a paginated list of settlements with optional date filters
// Note: When both date params are specified, max interval is one month
func (c *PaymentReportsClient) ListSettlements(ctx context.Context, req *ListSettlementsRequest) (*ListSettlementsResponse, error) {
	var resp ListSettlementsResponse

	path := "/v2/payment/settlements"
	separator := "?"

	if req.PaymentIntentID != "" {
		path += fmt.Sprintf("%spayment_intent_id=%s", separator, req.PaymentIntentID)
		separator = "&"
	}
	if req.SettlementBatchID != "" {
		path += fmt.Sprintf("%ssettlement_batch_id=%s", separator, req.SettlementBatchID)
		separator = "&"
	}
	if req.SettledStartTime != "" {
		path += fmt.Sprintf("%ssettled_start_time=%s", separator, req.SettledStartTime)
		separator = "&"
	}
	if req.SettledEndTime != "" {
		path += fmt.Sprintf("%ssettled_end_time=%s", separator, req.SettledEndTime)
		separator = "&"
	}
	if req.PageSize > 0 {
		path += fmt.Sprintf("%spage_size=%d", separator, req.PageSize)
		separator = "&"
	}
	if req.PageNumber > 0 {
		path += fmt.Sprintf("%spage_number=%d", separator, req.PageNumber)
	}

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list settlements: %w", err)
	}
	return &resp, nil
}
