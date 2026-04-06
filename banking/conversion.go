package banking

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

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
	AccountName      string `json:"account_name,omitempty"`
	Creator          string `json:"creator,omitempty"`
	SellCurrency     string `json:"sell_currency"`
	BuyCurrency      string `json:"buy_currency"`
	SellAmount       string `json:"sell_amount"`
	BuyAmount        string `json:"buy_amount"`
	ClientRate       string `json:"client_rate"`
	ConversionStatus string `json:"conversion_status"` // FUNDS_ARRIVED, TRADE_SETTLED, PENDING, etc.
	CreateTime       string `json:"create_time"`
	SettleTime       string `json:"settle_time,omitempty"`
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
	StartTime        int64  `json:"start_time"`        // optional, Unix timestamp in milliseconds
	EndTime          int64  `json:"end_time"`          // optional, Unix timestamp in milliseconds
	ConversionStatus string `json:"conversion_status"` // optional: FUNDS_ARRIVED, TRADE_SETTLED, etc.
	SellCurrency     string `json:"sell_currency"`     // optional
	BuyCurrency      string `json:"buy_currency"`      // optional
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
	Date  string `json:"date"`  // format: YYYY-MM-DD
	Valid bool   `json:"valid"` // whether this date is available for conversion
}

// List lists conversions
func (c *ConversionClient) List(ctx context.Context, req *ListConversionsRequest) (*ListConversionsResponse, error) {
	var resp ListConversionsResponse
	params := url.Values{}
	params.Set("page_size", strconv.Itoa(req.PageSize))
	params.Set("page_number", strconv.Itoa(req.PageNumber))
	if req.StartTime != 0 {
		params.Set("start_time", strconv.FormatInt(req.StartTime, 10))
	}
	if req.EndTime != 0 {
		params.Set("end_time", strconv.FormatInt(req.EndTime, 10))
	}
	if req.ConversionStatus != "" {
		params.Set("conversion_status", req.ConversionStatus)
	}
	if req.SellCurrency != "" {
		params.Set("sell_currency", req.SellCurrency)
	}
	if req.BuyCurrency != "" {
		params.Set("buy_currency", req.BuyCurrency)
	}
	path := "/v1/conversion?" + params.Encode()

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
