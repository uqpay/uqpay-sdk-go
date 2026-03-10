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
	ConversionID     string `json:"conversion_id"`          // UUID, unique conversion identifier
	ShortReferenceID string `json:"short_reference_id"`     // system-generated reference, e.g. P220406-LLCVLRM
	AccountName      string `json:"account_name,omitempty"` // optional, customer account name
	Creator          string `json:"creator,omitempty"`      // optional, entity that initiated the transaction
	SellCurrency     string `json:"sell_currency"`          // ISO 4217 currency code being sold
	BuyCurrency      string `json:"buy_currency"`           // ISO 4217 currency code being purchased
	SellAmount       string `json:"sell_amount"`            // amount of sell currency
	BuyAmount        string `json:"buy_amount"`             // amount of buy currency
	ClientRate       string `json:"client_rate"`            // exchange rate applied to the transaction
	ConversionStatus string `json:"conversion_status"`      // PROCESSING, AWAITING_FUNDS, TRADE_SETTLED, or FUNDS_ARRIVED
	CreateTime       string `json:"create_time"`            // ISO 8601 timestamp of transaction initiation
	SettleTime       string `json:"settle_time,omitempty"`  // optional, ISO 8601 timestamp of settlement completion
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
	ConversionID     string `json:"conversion_id"`      // UUID, unique conversion identifier
	ShortReferenceID string `json:"short_reference_id"` // system-generated reference, e.g. P220406-LLCVLRM
	SellCurrency     string `json:"sell_currency"`      // ISO 4217 currency code sold
	SellAmount       string `json:"sell_amount"`        // final amount sold
	BuyCurrency      string `json:"buy_currency"`       // ISO 4217 currency code purchased
	BuyAmount        string `json:"buy_amount"`         // final amount purchased
	CreatedDate      string `json:"created_date"`       // ISO 8601 transaction timestamp
	CurrencyPair     string `json:"currency_pair"`      // pair notation, e.g. USDSGD
	Reference        string `json:"reference"`          // system reference, e.g. XC240822-TXCRW2EI
	Status           string `json:"status"`             // PROCESSING, AWAITING_FUNDS, TRADE_SETTLED, or FUNDS_ARRIVED
}

// ListConversionsRequest represents a conversion list request
type ListConversionsRequest struct {
	PageSize         int    `json:"page_size"`         // required, 1-100
	PageNumber       int    `json:"page_number"`       // required, >= 1
	StartTime        int64  `json:"start_time"`        // optional, Unix timestamp in milliseconds
	EndTime          int64  `json:"end_time"`          // optional, Unix timestamp in milliseconds
	ConversionStatus string `json:"conversion_status"` // optional, PROCESSING, AWAITING_FUNDS, TRADE_SETTLED, or FUNDS_ARRIVED
	SellCurrency     string `json:"sell_currency"`     // optional, ISO 4217 currency code
	BuyCurrency      string `json:"buy_currency"`      // optional, ISO 4217 currency code
}

// ListConversionsResponse represents a conversion list response
type ListConversionsResponse struct {
	TotalPages int          `json:"total_pages"` // total number of result pages
	TotalItems int          `json:"total_items"` // total number of items
	Data       []Conversion `json:"data"`        // list of conversion records
}

// CreateQuoteRequest represents a quote creation request
type CreateQuoteRequest struct {
	SellCurrency    string `json:"sell_currency"`    // required, ISO 4217 currency code being sold
	SellAmount      string `json:"sell_amount"`      // optional, provide either sell_amount or buy_amount, not both
	BuyCurrency     string `json:"buy_currency"`     // required, ISO 4217 currency code being purchased
	BuyAmount       string `json:"buy_amount"`       // optional, provide either sell_amount or buy_amount, not both
	ConversionDate  string `json:"conversion_date"`  // required, format: YYYY-MM-DD, must be a valid business day
	TransactionType string `json:"transaction_type"` // optional, "conversion" (default) or "payout"
}

// QuoteValidity represents the validity period of a quote
type QuoteValidity struct {
	ValidFrom int64 `json:"valid_from"` // Unix timestamp in milliseconds, when quote becomes active
	ValidTo   int64 `json:"valid_to"`   // Unix timestamp in milliseconds, expires ~75 seconds after creation
}

// QuotePrice represents the price details of a quote
type QuotePrice struct {
	CurrencyPair string        `json:"currency_pair"` // pair notation, e.g. USDSGD
	DirectRate   string        `json:"direct_rate"`   // exchange rate from sell to buy currency
	InverseRate  string        `json:"inverse_rate"`  // reciprocal of direct rate
	QuoteID      string        `json:"quote_id"`      // UUID, use this when creating a conversion
	Validity     QuoteValidity `json:"validity"`      // time window during which the quote is valid
}

// CreateQuoteResponse represents a quote creation response
type CreateQuoteResponse struct {
	SellCurrency string     `json:"sell_currency"` // ISO 4217 currency code being sold
	SellAmount   string     `json:"sell_amount"`   // calculated or provided sell amount
	BuyCurrency  string     `json:"buy_currency"`  // ISO 4217 currency code being purchased
	BuyAmount    string     `json:"buy_amount"`    // calculated or provided buy amount
	QuotePrice   QuotePrice `json:"quote_price"`   // exchange rate details and quote identifier
}

// ConversionDate represents available conversion dates for a currency pair
type ConversionDate struct {
	Date  string `json:"date"`  // format: YYYY-MM-DD
	Valid bool   `json:"valid"` // whether this date is available for conversion
}

// ListConversionDatesResponse represents a response for listing available conversion dates
type ListConversionDatesResponse struct {
	Data []ConversionDate `json:"data"` // list of available conversion dates
}

// List lists conversions
// Optional RequestOptions can be provided to set custom headers like x-on-behalf-of
func (c *ConversionClient) List(ctx context.Context, req *ListConversionsRequest, opts ...*common.RequestOptions) (*ListConversionsResponse, error) {
	var resp ListConversionsResponse
	path := fmt.Sprintf("/v1/conversion?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)

	if req.StartTime != 0 {
		path += fmt.Sprintf("&start_time=%d", req.StartTime)
	}
	if req.EndTime != 0 {
		path += fmt.Sprintf("&end_time=%d", req.EndTime)
	}
	if req.ConversionStatus != "" {
		path += fmt.Sprintf("&conversion_status=%s", req.ConversionStatus)
	}
	if req.SellCurrency != "" {
		path += fmt.Sprintf("&sell_currency=%s", req.SellCurrency)
	}
	if req.BuyCurrency != "" {
		path += fmt.Sprintf("&buy_currency=%s", req.BuyCurrency)
	}

	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.GetWithOptions(ctx, path, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to list conversions: %w", err)
	}
	return &resp, nil
}

// Create creates a new conversion
// Optional RequestOptions can be provided to set custom headers like x-idempotency-key or x-on-behalf-of
func (c *ConversionClient) Create(ctx context.Context, req *CreateConversionRequest, opts ...*common.RequestOptions) (*CreateConversionResponse, error) {
	var resp CreateConversionResponse
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.PostWithOptions(ctx, "/v1/conversion", req, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to create conversion: %w", err)
	}
	return &resp, nil
}

// Get retrieves a specific conversion
// Optional RequestOptions can be provided to set custom headers like x-on-behalf-of
func (c *ConversionClient) Get(ctx context.Context, conversionID string, opts ...*common.RequestOptions) (*Conversion, error) {
	var resp Conversion
	path := fmt.Sprintf("/v1/conversion/%s", conversionID)
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.GetWithOptions(ctx, path, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to get conversion: %w", err)
	}
	return &resp, nil
}

// ListConversionDates retrieves available conversion dates for a currency pair
// Optional RequestOptions can be provided to set custom headers like x-on-behalf-of
func (c *ConversionClient) ListConversionDates(ctx context.Context, currencyFrom, currencyTo string, opts ...*common.RequestOptions) (*ListConversionDatesResponse, error) {
	var dates []ConversionDate
	path := fmt.Sprintf("/v1/conversion/conversion_dates?currency_from=%s&currency_to=%s", currencyFrom, currencyTo)
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.GetWithOptions(ctx, path, &dates, opt); err != nil {
		return nil, fmt.Errorf("failed to list conversion dates: %w", err)
	}
	return &ListConversionDatesResponse{Data: dates}, nil
}

// CreateQuote creates a new conversion quote
// Optional RequestOptions can be provided to set custom headers like x-idempotency-key or x-on-behalf-of
func (c *ConversionClient) CreateQuote(ctx context.Context, req *CreateQuoteRequest, opts ...*common.RequestOptions) (*CreateQuoteResponse, error) {
	var resp CreateQuoteResponse
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.PostWithOptions(ctx, "/v1/conversion/quote", req, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to create quote: %w", err)
	}
	return &resp, nil
}
