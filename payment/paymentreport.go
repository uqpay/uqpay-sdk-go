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
	ID             string `json:"id"`
	Amount         string `json:"amount,omitempty"`
	Currency       string `json:"currency,omitempty"`
	Status         string `json:"status,omitempty"`
	SettledAt      string `json:"settled_at,omitempty"`
	TransactionFee string `json:"transaction_fee,omitempty"`
	NetAmount      string `json:"net_amount,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
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
