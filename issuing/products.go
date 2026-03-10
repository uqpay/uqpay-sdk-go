package issuing

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// ProductsClient handles product operations
type ProductsClient struct {
	client *common.APIClient
}

// NoPinPaymentLimit represents a no-pin payment limit
type NoPinPaymentLimit struct {
	Amount   string `json:"amount"`   // Required. Maximum no-PIN transaction amount (returned as string)
	Currency string `json:"currency"` // Required. ISO 4217 currency code
}

// CardProduct represents a card product
type CardProduct struct {
	ProductID          string              `json:"product_id"`            // Required. Unique product identifier (UUID)
	ModeType           string              `json:"mode_type"`             // Required. SINGLE or SHARE
	CardBin            string              `json:"card_bin"`              // Required. Card number prefix (BIN)
	CardForm           []string            `json:"card_form"`             // Required. Supported forms: VIR (virtual) or PHY (physical)
	MaxCardQuota       int                 `json:"max_card_quota"`        // Optional. Maximum number of cards issuable under account
	CardScheme         string              `json:"card_scheme"`           // Required. Card scheme, e.g. VISA
	CardCurrency       []string            `json:"card_currency"`         // Required. Supported ISO 4217 currencies for card issuance
	ProductStatus      string              `json:"product_status"`        // Required. ENABLED or DISABLED
	NoPinPaymentAmount []NoPinPaymentLimit `json:"no_pin_payment_amount"` // Required. No-PIN transaction limits per currency
	CreateTime         string              `json:"create_time"`           // Required. ISO timestamp of creation
	UpdateTime         string              `json:"update_time"`           // Required. ISO timestamp of last update
}

// ListProductsRequest represents a product list request
type ListProductsRequest struct {
	PageSize   int `json:"page_size"`   // Required. Max items per page, min: 10, max: 100, default: 10
	PageNumber int `json:"page_number"` // Required. Page number to retrieve, min: 1, default: 1
}

// ListProductsResponse represents a product list response
type ListProductsResponse struct {
	TotalPages int           `json:"total_pages"` // Total number of available pages
	TotalItems int           `json:"total_items"` // Total count of available items
	Data       []CardProduct `json:"data"`        // Array of card product objects
}

// List lists card products
func (c *ProductsClient) List(ctx context.Context, req *ListProductsRequest) (*ListProductsResponse, error) {
	var resp ListProductsResponse
	path := fmt.Sprintf("/v1/issuing/products?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return &resp, nil
}
