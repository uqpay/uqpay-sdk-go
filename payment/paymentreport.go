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
	SettledStartTime string `json:"settled_start_time"` // Start of settlement date range (ISO8601)
	SettledEndTime   string `json:"settled_end_time"`   // End of settlement date range (ISO8601)
	PageSize         int    `json:"page_size"`          // Number of items per page
	PageNumber       int    `json:"page_number"`        // Page number (1-based)
}

// ============================================================================
// Response Structures
// ============================================================================

// Settlement represents a settlement record
type Settlement struct {
	SettlementID          string `json:"settlement_id"`
	AccountID             string `json:"account_id,omitempty"`
	AccountName           string `json:"account_name,omitempty"`
	SourceType            string `json:"source_type,omitempty"`
	TransactionType       string `json:"transaction_type,omitempty"`
	MerchantOrderID       string `json:"merchant_order_id,omitempty"`
	PaymentIntentID       string `json:"payment_intent_id,omitempty"`
	PaymentMethod         string `json:"payment_method,omitempty"`
	TransactionCreateDate string `json:"transaction_create_date,omitempty"`
	TransactionAmount     string `json:"transaction_amount,omitempty"`
	TransactionCurrency   string `json:"transaction_currency,omitempty"`
	TransactionDate       string `json:"transaction_date,omitempty"`
	SettlementAmount      string `json:"settlement_amount,omitempty"`
	SettlementCurrency    string `json:"settlement_currency,omitempty"`
	NetSettlementAmount   string `json:"net_settlement_amount,omitempty"`
	ExchangeRate          string `json:"exchange_rate,omitempty"`
	FeeCurrency           string `json:"fee_currency,omitempty"`
	InterchangeFee        string `json:"interchange_fee,omitempty"`
	SchemeFee             string `json:"scheme_fee,omitempty"`
	TranscationFee        string `json:"transcation_fee,omitempty"` // Note: API uses "transcation" (typo)
	ReturnFee             string `json:"return_fee,omitempty"`
	TotalFeeAmount        string `json:"total_fee_amount,omitempty"`
	SettlementStatus      string `json:"settlement_status,omitempty"`
	SettlementBatchID     string `json:"settlement_batch_id,omitempty"`
	SettlementCreateDate  string `json:"settlement_create_date,omitempty"`
	SettlementDate        string `json:"settlement_date,omitempty"`
}

// ListSettlementsResponse represents a paginated list of settlements
type ListSettlementsResponse struct {
	TotalPages int          `json:"total_pages"`
	TotalItems int          `json:"total_items"`
	Data       []Settlement `json:"data"`
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
