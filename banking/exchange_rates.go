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
	CurrencyPair string `json:"currency_pair"` // e.g., "USDEUR"
	BuyPrice     string `json:"buy_price"`
	SellPrice    string `json:"sell_price"`
}

// ListRatesRequest represents a request to list exchange rates
type ListRatesRequest struct {
	CurrencyPairs []string `json:"currency_pairs,omitempty"` // optional: filter by specific currency pairs (e.g., ["USDEUR", "GBPUSD"])
}

// ListRatesResponse represents a response containing exchange rates
type ListRatesResponse struct {
	Rates                    []RateItem `json:"rates"`
	UnavailableCurrencyPairs []string   `json:"unavailable_currency_pairs"`
	LastUpdated              string     `json:"last_updated"`
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
