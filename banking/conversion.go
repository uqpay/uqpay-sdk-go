package banking

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// ConversionClient handles conversion operations
type ConversionClient struct {
	client *common.APIClient
}

// Conversion represents a currency conversion
type Conversion struct {
	ConversionID     string `json:"conversion_id"`
	ShortReferenceID string `json:"short_reference_id"`
	CurrencyFrom     string `json:"currency_from"`
	CurrencyTo       string `json:"currency_to"`
	AmountFrom       string `json:"amount_from"`
	AmountTo         string `json:"amount_to"`
	Rate             string `json:"rate"`
	ConversionStatus string `json:"conversion_status"` // COMPLETED, PENDING, FAILED
	CreateTime       string `json:"create_time"`
	CompletedTime    string `json:"completed_time,omitempty"`
	SettlementDate   string `json:"settlement_date,omitempty"`
}

// CreateConversionRequest represents a conversion creation request
type CreateConversionRequest struct {
	QuoteID        string `json:"quote_id"`              // required, UUID from quote
	SellCurrency   string `json:"sell_currency"`         // required, ISO 4217 currency code
	SellAmount     string `json:"sell_amount,omitempty"` // provide either sell_amount or buy_amount
	BuyCurrency    string `json:"buy_currency"`          // required, ISO 4217 currency code
	BuyAmount      string `json:"buy_amount,omitempty"`  // provide either sell_amount or buy_amount
	ConversionDate string `json:"conversion_date"`       // required, format: YYYY-MM-DD (only current date supported)
}

// CreateConversionResponse represents a conversion creation response
type CreateConversionResponse struct {
	ConversionID     string `json:"conversion_id"`
	ShortReferenceID string `json:"short_reference_id"`
	SellCurrency     string `json:"sell_currency"`
	SellAmount       string `json:"sell_amount"`
	BuyCurrency      string `json:"buy_currency"`
	BuyAmount        string `json:"buy_amount"`
	CreatedDate      string `json:"created_date"`
	CurrencyPair     string `json:"currency_pair"`
	Reference        string `json:"reference"`
	Status           string `json:"status"`
}

// ListConversionsRequest represents a conversion list request
type ListConversionsRequest struct {
	PageSize         int    `json:"page_size"`         // required, 10-100
	PageNumber       int    `json:"page_number"`       // required, >=1
	StartTime        string `json:"start_time"`        // optional, ISO8601
	EndTime          string `json:"end_time"`          // optional, ISO8601
	ConversionStatus string `json:"conversion_status"` // optional: COMPLETED, PENDING, FAILED
	CurrencyFrom     string `json:"currency_from"`     // optional
	CurrencyTo       string `json:"currency_to"`       // optional
}

// ListConversionsResponse represents a conversion list response
type ListConversionsResponse struct {
	TotalPages int          `json:"total_pages"`
	TotalItems int          `json:"total_items"`
	Data       []Conversion `json:"data"`
}

// CreateQuoteRequest represents a quote creation request
type CreateQuoteRequest struct {
	SellCurrency    string `json:"sell_currency"`    // required, ISO 4217 currency code
	SellAmount      string `json:"sell_amount"`      // amount to sell
	BuyCurrency     string `json:"buy_currency"`     // required, ISO 4217 currency code
	BuyAmount       string `json:"buy_amount"`       // amount to buy
	ConversionDate  string `json:"conversion_date"`  // required, format: YYYY-MM-DD
	TransactionType string `json:"transaction_type"` // required, e.g., "conversion"
}

// QuoteValidity represents the validity period of a quote
type QuoteValidity struct {
	ValidFrom int64 `json:"valid_from"` // Unix timestamp in milliseconds
	ValidTo   int64 `json:"valid_to"`   // Unix timestamp in milliseconds
}

// QuotePrice represents the price details of a quote
type QuotePrice struct {
	CurrencyPair string        `json:"currency_pair"`
	DirectRate   string        `json:"direct_rate"`
	InverseRate  string        `json:"inverse_rate"`
	QuoteID      string        `json:"quote_id"`
	Validity     QuoteValidity `json:"validity"`
}

// CreateQuoteResponse represents a quote creation response
type CreateQuoteResponse struct {
	SellCurrency string     `json:"sell_currency"`
	SellAmount   string     `json:"sell_amount"`
	BuyCurrency  string     `json:"buy_currency"`
	BuyAmount    string     `json:"buy_amount"`
	QuotePrice   QuotePrice `json:"quote_price"`
}

// ConversionDate represents available conversion dates for a currency pair
type ConversionDate struct {
	Date          string `json:"date"`           // format: YYYY-MM-DD
	FirstCutoff   string `json:"first_cutoff"`   // ISO8601 timestamp
	SecondCutoff  string `json:"second_cutoff"`  // ISO8601 timestamp
	OptimizedDate bool   `json:"optimized_date"` // whether this is the optimal conversion date
}

// List lists conversions
func (c *ConversionClient) List(ctx context.Context, req *ListConversionsRequest) (*ListConversionsResponse, error) {
	var resp ListConversionsResponse
	path := fmt.Sprintf("/v1/conversion?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)

	if req.StartTime != "" {
		path += fmt.Sprintf("&start_time=%s", req.StartTime)
	}
	if req.EndTime != "" {
		path += fmt.Sprintf("&end_time=%s", req.EndTime)
	}
	if req.ConversionStatus != "" {
		path += fmt.Sprintf("&conversion_status=%s", req.ConversionStatus)
	}
	if req.CurrencyFrom != "" {
		path += fmt.Sprintf("&currency_from=%s", req.CurrencyFrom)
	}
	if req.CurrencyTo != "" {
		path += fmt.Sprintf("&currency_to=%s", req.CurrencyTo)
	}

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list conversions: %w", err)
	}
	return &resp, nil
}

// Create creates a new conversion
func (c *ConversionClient) Create(ctx context.Context, req *CreateConversionRequest) (*CreateConversionResponse, error) {
	var resp CreateConversionResponse
	if err := c.client.Post(ctx, "/v1/conversion", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create conversion: %w", err)
	}
	return &resp, nil
}

// Get retrieves a specific conversion
func (c *ConversionClient) Get(ctx context.Context, conversionID string) (*Conversion, error) {
	var resp Conversion
	path := fmt.Sprintf("/v1/conversion/%s", conversionID)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get conversion: %w", err)
	}
	return &resp, nil
}

// ListConversionDates retrieves available conversion dates for a currency pair
func (c *ConversionClient) ListConversionDates(ctx context.Context, currencyFrom, currencyTo string) ([]ConversionDate, error) {
	var resp []ConversionDate
	path := fmt.Sprintf("/v1/conversion/conversion_dates?currency_from=%s&currency_to=%s", currencyFrom, currencyTo)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list conversion dates: %w", err)
	}
	return resp, nil
}

// CreateQuote creates a new conversion quote
func (c *ConversionClient) CreateQuote(ctx context.Context, req *CreateQuoteRequest) (*CreateQuoteResponse, error) {
	var resp CreateQuoteResponse
	if err := c.client.Post(ctx, "/v1/conversion/quote", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create quote: %w", err)
	}
	return &resp, nil
}
