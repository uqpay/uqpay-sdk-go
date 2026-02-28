package payment

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// PaymentBalancesClient handles payment balance operations
type PaymentBalancesClient struct {
	client *common.APIClient
}

// ============================================================================
// Request Structures
// ============================================================================

// ListBalancesRequest represents a balances list request
type ListBalancesRequest struct {
	PageSize   int `json:"page_size"`   // Number of items per page
	PageNumber int `json:"page_number"` // Page number (1-based)
}

// ============================================================================
// Response Structures
// ============================================================================

// Balance represents a currency account balance
type Balance struct {
	BalanceID        string `json:"balance_id,omitempty"`
	Currency         string `json:"currency"`
	AvailableBalance string `json:"available_balance,omitempty"`
	PayableBalance   string `json:"payable_balance,omitempty"`
	PendingBalance   string `json:"pending_balance,omitempty"`
	ReservedBalance  string `json:"reserved_balance,omitempty"`
	MarginBalance    string `json:"margin_balance,omitempty"`
	FrozenBalance    string `json:"frozen_balance,omitempty"`
}

// ListBalancesResponse represents a paginated list of currency balances
type ListBalancesResponse struct {
	TotalPages int       `json:"total_pages"`
	TotalItems int       `json:"total_items"`
	Data       []Balance `json:"data"`
}

// ============================================================================
// API Methods
// ============================================================================

// Get retrieves the balance for a specific currency
func (c *PaymentBalancesClient) Get(ctx context.Context, currency string) (*Balance, error) {
	var resp Balance
	path := fmt.Sprintf("/v2/payment/balances/%s", currency)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return &resp, nil
}

// List returns a paginated list of currency account balances
func (c *PaymentBalancesClient) List(ctx context.Context, req *ListBalancesRequest) (*ListBalancesResponse, error) {
	var resp ListBalancesResponse

	path := "/v2/payment/balances"
	separator := "?"

	if req.PageSize > 0 {
		path += fmt.Sprintf("%spage_size=%d", separator, req.PageSize)
		separator = "&"
	}
	if req.PageNumber > 0 {
		path += fmt.Sprintf("%spage_number=%d", separator, req.PageNumber)
	}

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list balances: %w", err)
	}
	return &resp, nil
}
