package simulator

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/v2/common"
)

type IssuingClient struct{ client *common.APIClient }

type AuthorizationRequest struct {
	CardID               string  `json:"card_id"`
	TransactionAmount    float64 `json:"transaction_amount"`
	TransactionCurrency  string  `json:"transaction_currency"`
	MerchantName         string  `json:"merchant_name"`
	MerchantCategoryCode string  `json:"merchant_category_code"`
}

type AuthorizationResponse struct {
	CardID               string                 `json:"card_id"`
	CardNumber           string                 `json:"card_number"`
	CardholderID         string                 `json:"cardholder_id"`
	TransactionID        string                 `json:"transaction_id"`
	TransactionType      string                 `json:"transaction_type"`
	CardAvailableBalance float64                `json:"card_available_balance"`
	AuthorizationCode    string                 `json:"authorization_code"`
	BillingAmount        float64                `json:"billing_amount"`
	BillingCurrency      string                 `json:"billing_currency"`
	TransactionAmount    float64                `json:"transaction_amount"`
	TransactionCurrency  string                 `json:"transaction_currency"`
	TransactionTime      string                 `json:"transaction_time"`
	PostedTime           string                 `json:"posted_time"`
	MerchantData         map[string]interface{} `json:"merchant_data"`
	FailureReason        string                 `json:"failure_reason"`
	TransactionStatus    string                 `json:"transaction_status"`
}

type ReversalRequest struct {
	TransactionID string `json:"transaction_id"`
}

func (c *IssuingClient) Authorize(ctx context.Context, req *AuthorizationRequest, opts ...*common.RequestOptions) (*AuthorizationResponse, error) {
	var resp AuthorizationResponse
	if err := c.client.PostWithOptions(ctx, "/v1/simulation/issuing/authorization", req, &resp, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to simulate authorization: %w", err)
	}
	return &resp, nil
}

func (c *IssuingClient) Reverse(ctx context.Context, req *ReversalRequest, opts ...*common.RequestOptions) (*AuthorizationResponse, error) {
	var resp AuthorizationResponse
	if err := c.client.PostWithOptions(ctx, "/v1/simulation/issuing/reversal", req, &resp, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to simulate reversal: %w", err)
	}
	return &resp, nil
}
