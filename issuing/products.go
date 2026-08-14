package issuing

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/v2/common"
)

// ProductsClient handles product operations
type ProductsClient struct {
	client *common.APIClient
}

// NoPinPaymentLimit represents a no-pin payment limit
type NoPinPaymentLimit struct {
	Amount   string `json:"amount"` // API returns string
	Currency string `json:"currency"`
}

// CardProduct represents a card product
type CardProduct struct {
	ProductID          string              `json:"product_id"`
	ModeType           string              `json:"mode_type"`
	CardBin            string              `json:"card_bin"`
	CardForm           []string            `json:"card_form"`
	MaxCardQuota       int                 `json:"max_card_quota"`
	CardScheme         string              `json:"card_scheme"`
	CardCurrency       []string            `json:"card_currency"` // API returns array
	ProductStatus      string              `json:"product_status"`
	NoPinPaymentAmount []NoPinPaymentLimit `json:"no_pin_payment_amount"` // Array of payment limits
	CreateTime         string              `json:"create_time"`
	UpdateTime         string              `json:"update_time"`
}

// ListProductsRequest represents a product list request
type ListProductsRequest struct {
	PageSize   int `json:"page_size"`
	PageNumber int `json:"page_number"`
}

// ListProductsResponse represents a product list response
type ListProductsResponse struct {
	TotalPages int           `json:"total_pages"`
	TotalItems int           `json:"total_items"`
	Data       []CardProduct `json:"data"`
}

// List lists card products
func (c *ProductsClient) List(ctx context.Context, req *ListProductsRequest, opts ...*common.RequestOptions) (*ListProductsResponse, error) {
	var resp ListProductsResponse
	path := fmt.Sprintf("/v1/issuing/products?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)
	if err := c.client.GetWithOptions(ctx, path, &resp, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return &resp, nil
}
