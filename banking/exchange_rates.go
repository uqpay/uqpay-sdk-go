package banking

import (
	"context"
	"fmt"
	"strings"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// ExchangeRatesClient handles exchange rate operations
type ExchangeRatesClient struct {
	client *common.APIClient
}

// RateItem represents an exchange rate for a currency pair
type RateItem struct {
	CurrencyPair string `json:"currency_pair"`        // Required. 6-letter uppercase currency pair code, e.g., "USDEUR"
	BuyPrice     string `json:"buy_price,omitempty"`  // Required. Buy price, rounded to 4 decimal places
	SellPrice    string `json:"sell_price,omitempty"` // Required. Sell price, rounded to 4 decimal places
}

// ListRatesRequest represents a request to list exchange rates
type ListRatesRequest struct {
	CurrencyPairs []string `json:"currency_pairs,omitempty"` // Optional. Up to 100 comma-separated 6-letter uppercase pairs, e.g., ["USDEUR", "GBPUSD"]. If omitted, all available pairs returned
}

// ListRatesResponse represents a response containing exchange rates
type ListRatesResponse struct {
	Rates                    []RateItem `json:"rates"`                      // Required. Array of rate items for each currency pair
	UnavailableCurrencyPairs []string   `json:"unavailable_currency_pairs"` // Optional. List of unsupported 6-letter currency pair codes from the request
	LastUpdated              string     `json:"last_updated"`               // Optional. ISO 8601 datetime of the last rate update
}

// listRatesDataWrapper wraps the API response
type listRatesDataWrapper struct {
	Data ListRatesResponse `json:"data"`
}

// List retrieves current exchange rates
// Optionally filter by specific currency pairs
func (c *ExchangeRatesClient) List(ctx context.Context, req *ListRatesRequest) (*ListRatesResponse, error) {
	var wrapper listRatesDataWrapper
	path := "/v1/exchange/rates"

	// Add currency_pairs query parameter if specified
	if req != nil && len(req.CurrencyPairs) > 0 {
		pairs := strings.Join(req.CurrencyPairs, ",")
		path += fmt.Sprintf("?currency_pairs=%s", pairs)
	}

	if err := c.client.Get(ctx, path, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to list exchange rates: %w", err)
	}
	return &wrapper.Data, nil
}
