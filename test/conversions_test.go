package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/uqpay/uqpay-sdk-go/banking"
)

func TestConversionCreateQuote(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	// Create a quote
	today := time.Now().Format("2006-01-02")
	req := &banking.CreateQuoteRequest{
		SellCurrency:    "USD",
		SellAmount:      "100.00",
		BuyCurrency:     "SGD",
		ConversionDate:  today,
		TransactionType: "conversion",
	}

	t.Logf("Creating quote: %s -> %s, Amount: %s", req.SellCurrency, req.BuyCurrency, req.SellAmount)

	quote, err := client.Banking.Conversions.CreateQuote(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create quote: %v", err)
	}

	if quote.QuotePrice.QuoteID == "" {
		t.Error("Expected quote_id to be set")
	}
	if quote.SellCurrency != req.SellCurrency {
		t.Errorf("Expected sell_currency %s, got %s", req.SellCurrency, quote.SellCurrency)
	}
	if quote.BuyCurrency != req.BuyCurrency {
		t.Errorf("Expected buy_currency %s, got %s", req.BuyCurrency, quote.BuyCurrency)
	}
	if quote.QuotePrice.DirectRate == "" {
		t.Error("Expected direct_rate to be set")
	}
	if quote.BuyAmount == "" {
		t.Error("Expected buy_amount to be set")
	}

	t.Logf("Quote created successfully:")
	t.Logf("  Quote ID: %s", quote.QuotePrice.QuoteID)
	t.Logf("  Sell: %s %s", quote.SellAmount, quote.SellCurrency)
	t.Logf("  Buy: %s %s", quote.BuyAmount, quote.BuyCurrency)
	t.Logf("  Rate: %s", quote.QuotePrice.DirectRate)
	t.Logf("  Currency Pair: %s", quote.QuotePrice.CurrencyPair)
}

func TestConversionCreateQuoteWithSettlementDate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	// Get conversion dates first
	dates, err := client.Banking.Conversions.ListConversionDates(ctx, "USD", "SGD")
	if err != nil {
		t.Fatalf("Failed to get conversion dates: %v", err)
	}

	if len(dates) == 0 {
		t.Skip("No conversion dates available")
	}

	// Use the first available date
	conversionDate := dates[0].Date

	// Create a quote with specific conversion date
	req := &banking.CreateQuoteRequest{
		SellCurrency:    "USD",
		SellAmount:      "100.00",
		BuyCurrency:     "SGD",
		ConversionDate:  conversionDate,
		TransactionType: "conversion",
	}

	t.Logf("Creating quote with conversion date: %s", conversionDate)

	quote, err := client.Banking.Conversions.CreateQuote(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create quote: %v", err)
	}

	if quote.QuotePrice.QuoteID == "" {
		t.Error("Expected quote_id to be set")
	}

	t.Logf("Quote created with conversion date:")
	t.Logf("  Quote ID: %s", quote.QuotePrice.QuoteID)
	t.Logf("  Rate: %s", quote.QuotePrice.DirectRate)
}

func TestConversionCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	// First create a quote to get a valid quote_id
	today := time.Now().Format("2006-01-02")
	quoteReq := &banking.CreateQuoteRequest{
		SellCurrency:    "USD",
		SellAmount:      "100.00",
		BuyCurrency:     "SGD",
		ConversionDate:  today,
		TransactionType: "conversion",
	}

	quote, err := client.Banking.Conversions.CreateQuote(ctx, quoteReq)
	if err != nil {
		t.Fatalf("Failed to create quote: %v", err)
	}

	t.Logf("Quote created: %s", quote.QuotePrice.QuoteID)

	// Create a conversion using the quote
	req := &banking.CreateConversionRequest{
		QuoteID:        quote.QuotePrice.QuoteID,
		SellCurrency:   "USD",
		SellAmount:     "100.00",
		BuyCurrency:    "SGD",
		ConversionDate: today,
	}

	t.Logf("Creating conversion: %s -> %s, Amount: %s, Date: %s", req.SellCurrency, req.BuyCurrency, req.SellAmount, req.ConversionDate)

	resp, err := client.Banking.Conversions.Create(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create conversion: %v", err)
	}

	if resp.ConversionID == "" {
		t.Error("Expected conversion_id to be set")
	}
	if resp.ShortReferenceID == "" {
		t.Error("Expected short_reference_id to be set")
	}

	t.Logf("Conversion created successfully:")
	t.Logf("  Conversion ID: %s", resp.ConversionID)
	t.Logf("  Short Reference ID: %s", resp.ShortReferenceID)
	t.Logf("  Sell: %s %s", resp.SellAmount, resp.SellCurrency)
	t.Logf("  Buy: %s %s", resp.BuyAmount, resp.BuyCurrency)
	t.Logf("  Status: %s", resp.Status)

	// Get the created conversion
	t.Run("GetConversion", func(t *testing.T) {
		t.Logf("Getting conversion: %s", resp.ConversionID)

		conversion, err := client.Banking.Conversions.Get(ctx, resp.ConversionID)
		if err != nil {
			t.Fatalf("Failed to get conversion: %v", err)
		}

		if conversion.ConversionID != resp.ConversionID {
			t.Errorf("Expected conversion_id %s, got %s", resp.ConversionID, conversion.ConversionID)
		}

		t.Logf("Retrieved conversion:")
		t.Logf("  ID: %s", conversion.ConversionID)
		t.Logf("  Short Reference ID: %s", conversion.ShortReferenceID)
		t.Logf("  Status: %s", conversion.ConversionStatus)
	})
}

func TestConversionCreateWithQuote(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	// First, create a quote
	today := time.Now().Format("2006-01-02")
	quoteReq := &banking.CreateQuoteRequest{
		SellCurrency:    "USD",
		SellAmount:      "100.00",
		BuyCurrency:     "SGD",
		ConversionDate:  today,
		TransactionType: "conversion",
	}

	t.Logf("Creating quote for conversion...")

	quote, err := client.Banking.Conversions.CreateQuote(ctx, quoteReq)
	if err != nil {
		t.Fatalf("Failed to create quote: %v", err)
	}

	t.Logf("Quote created: %s with rate %s", quote.QuotePrice.QuoteID, quote.QuotePrice.DirectRate)

	// Now create conversion using the quote
	convReq := &banking.CreateConversionRequest{
		QuoteID:        quote.QuotePrice.QuoteID,
		SellCurrency:   "USD",
		SellAmount:     "100.00",
		BuyCurrency:    "SGD",
		ConversionDate: today,
	}

	t.Logf("Creating conversion with quote ID: %s", quote.QuotePrice.QuoteID)

	conv, err := client.Banking.Conversions.Create(ctx, convReq)
	if err != nil {
		t.Fatalf("Failed to create conversion: %v", err)
	}

	t.Logf("Conversion created successfully:")
	t.Logf("  Conversion ID: %s", conv.ConversionID)
	t.Logf("  Short Reference ID: %s", conv.ShortReferenceID)
	t.Logf("  Sell: %s %s", conv.SellAmount, conv.SellCurrency)
	t.Logf("  Buy: %s %s", conv.BuyAmount, conv.BuyCurrency)
	t.Logf("  Status: %s", conv.Status)
}

func TestConversionList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	req := &banking.ListConversionsRequest{
		PageSize:   10,
		PageNumber: 1,
	}

	t.Logf("Listing conversions...")

	resp, err := client.Banking.Conversions.List(ctx, req)
	if err != nil {
		t.Fatalf("Failed to list conversions: %v", err)
	}

	if resp.TotalPages < 0 {
		t.Error("Expected total_pages to be >= 0")
	}
	if resp.TotalItems < 0 {
		t.Error("Expected total_items to be >= 0")
	}

	t.Logf("Listed conversions:")
	t.Logf("  Total Pages: %d", resp.TotalPages)
	t.Logf("  Total Items: %d", resp.TotalItems)
	t.Logf("  Current Page Items: %d", len(resp.Data))

	for i, conv := range resp.Data {
		t.Logf("  Conversion %d:", i+1)
		t.Logf("    ID: %s", conv.ConversionID)
		t.Logf("    Short Reference ID: %s", conv.ShortReferenceID)
		t.Logf("    From: %s %s", conv.AmountFrom, conv.CurrencyFrom)
		t.Logf("    To: %s %s", conv.AmountTo, conv.CurrencyTo)
		t.Logf("    Rate: %s", conv.Rate)
		t.Logf("    Status: %s", conv.ConversionStatus)
		t.Logf("    Create Time: %s", conv.CreateTime)
	}
}

func TestConversionListWithFilters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	// Test with status filter
	req := &banking.ListConversionsRequest{
		PageSize:         10,
		PageNumber:       1,
		ConversionStatus: "COMPLETED",
	}

	t.Logf("Listing conversions with status filter: %s", req.ConversionStatus)

	resp, err := client.Banking.Conversions.List(ctx, req)
	if err != nil {
		t.Fatalf("Failed to list conversions: %v", err)
	}

	t.Logf("Found %d completed conversions", resp.TotalItems)

	// Verify all returned conversions have the correct status
	for _, conv := range resp.Data {
		if conv.ConversionStatus != req.ConversionStatus {
			t.Errorf("Expected conversion status %s, got %s", req.ConversionStatus, conv.ConversionStatus)
		}
	}
}

func TestConversionListWithCurrencyFilters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	// Test with currency filters
	req := &banking.ListConversionsRequest{
		PageSize:     10,
		PageNumber:   1,
		CurrencyFrom: "USD",
		CurrencyTo:   "EUR",
	}

	t.Logf("Listing conversions: %s -> %s", req.CurrencyFrom, req.CurrencyTo)

	resp, err := client.Banking.Conversions.List(ctx, req)
	if err != nil {
		t.Fatalf("Failed to list conversions: %v", err)
	}

	t.Logf("Found %d conversions for %s -> %s", resp.TotalItems, req.CurrencyFrom, req.CurrencyTo)

	// Verify all returned conversions have the correct currencies
	for _, conv := range resp.Data {
		if conv.CurrencyFrom != req.CurrencyFrom {
			t.Errorf("Expected currency_from %s, got %s", req.CurrencyFrom, conv.CurrencyFrom)
		}
		if conv.CurrencyTo != req.CurrencyTo {
			t.Errorf("Expected currency_to %s, got %s", req.CurrencyTo, conv.CurrencyTo)
		}
	}
}

func TestConversionListWithTimeRange(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	// Get conversions from last 30 days
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -30)

	req := &banking.ListConversionsRequest{
		PageSize:   10,
		PageNumber: 1,
		StartTime:  startTime.Format(time.RFC3339),
		EndTime:    endTime.Format(time.RFC3339),
	}

	t.Logf("Listing conversions from %s to %s", req.StartTime, req.EndTime)

	resp, err := client.Banking.Conversions.List(ctx, req)
	if err != nil {
		t.Fatalf("Failed to list conversions: %v", err)
	}

	t.Logf("Found %d conversions in the last 30 days", resp.TotalItems)
}

func TestConversionDates(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	currencyFrom := "USD"
	currencyTo := "EUR"

	t.Logf("Getting conversion dates for %s -> %s", currencyFrom, currencyTo)

	dates, err := client.Banking.Conversions.ListConversionDates(ctx, currencyFrom, currencyTo)
	if err != nil {
		t.Fatalf("Failed to get conversion dates: %v", err)
	}

	if len(dates) == 0 {
		t.Error("Expected at least one conversion date")
	}

	t.Logf("Available conversion dates: %d", len(dates))

	for i, date := range dates {
		t.Logf("  Date %d:", i+1)
		t.Logf("    Date: %s", date.Date)
		t.Logf("    First Cutoff: %s", date.FirstCutoff)
		t.Logf("    Second Cutoff: %s", date.SecondCutoff)
		t.Logf("    Optimized: %t", date.OptimizedDate)
	}

	// Verify date format (YYYY-MM-DD)
	for _, date := range dates {
		if len(date.Date) != 10 {
			t.Errorf("Expected date format YYYY-MM-DD, got %s", date.Date)
		}
		if date.FirstCutoff == "" {
			t.Error("Expected first_cutoff to be set")
		}
		if date.SecondCutoff == "" {
			t.Error("Expected second_cutoff to be set")
		}
	}
}

func TestConversionFullFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	// Step 1: Get conversion dates
	t.Run("GetConversionDates", func(t *testing.T) {
		dates, err := client.Banking.Conversions.ListConversionDates(ctx, "USD", "EUR")
		if err != nil {
			t.Fatalf("Failed to get conversion dates: %v", err)
		}

		t.Logf("Available conversion dates: %d", len(dates))
		if len(dates) > 0 {
			t.Logf("  First available date: %s", dates[0].Date)
		}
	})

	// Step 2: Create a quote
	var quoteID string
	var quoteRate string

	t.Run("CreateQuote", func(t *testing.T) {
		today := time.Now().Format("2006-01-02")
		req := &banking.CreateQuoteRequest{
			SellCurrency:    "USD",
			SellAmount:      "100.00",
			BuyCurrency:     "EUR",
			ConversionDate:  today,
			TransactionType: "conversion",
		}

		quote, err := client.Banking.Conversions.CreateQuote(ctx, req)
		if err != nil {
			t.Fatalf("Failed to create quote: %v", err)
		}

		quoteID = quote.QuotePrice.QuoteID
		quoteRate = quote.QuotePrice.DirectRate

		t.Logf("Quote created:")
		t.Logf("  ID: %s", quote.QuotePrice.QuoteID)
		t.Logf("  Rate: %s", quote.QuotePrice.DirectRate)
		t.Logf("  Buy Amount: %s %s", quote.BuyAmount, quote.BuyCurrency)
	})

	// Step 3: Create conversion using the quote
	var conversionID string

	t.Run("CreateConversion", func(t *testing.T) {
		today := time.Now().Format("2006-01-02")
		req := &banking.CreateConversionRequest{
			QuoteID:        quoteID,
			SellCurrency:   "USD",
			SellAmount:     "100.00",
			BuyCurrency:    "EUR",
			ConversionDate: today,
		}

		conv, err := client.Banking.Conversions.Create(ctx, req)
		if err != nil {
			t.Fatalf("Failed to create conversion: %v", err)
		}

		conversionID = conv.ConversionID

		t.Logf("Conversion created:")
		t.Logf("  ID: %s", conv.ConversionID)
		t.Logf("  Short Reference ID: %s", conv.ShortReferenceID)
		t.Logf("  Sell: %s %s", conv.SellAmount, conv.SellCurrency)
		t.Logf("  Buy: %s %s", conv.BuyAmount, conv.BuyCurrency)
	})

	// Step 4: Get the conversion details
	t.Run("GetConversion", func(t *testing.T) {
		conversion, err := client.Banking.Conversions.Get(ctx, conversionID)
		if err != nil {
			t.Fatalf("Failed to get conversion: %v", err)
		}

		t.Logf("Retrieved conversion:")
		t.Logf("  ID: %s", conversion.ConversionID)
		t.Logf("  Status: %s", conversion.ConversionStatus)
		t.Logf("  Rate: %s", conversion.Rate)
		t.Logf("  From: %s %s", conversion.AmountFrom, conversion.CurrencyFrom)
		t.Logf("  To: %s %s", conversion.AmountTo, conversion.CurrencyTo)

		if conversion.Rate != quoteRate {
			t.Logf("  Note: Rate changed from quote (%s) to final (%s)", quoteRate, conversion.Rate)
		}
	})

	// Step 5: List conversions and verify our conversion is in the list
	t.Run("ListConversions", func(t *testing.T) {
		req := &banking.ListConversionsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Banking.Conversions.List(ctx, req)
		if err != nil {
			t.Fatalf("Failed to list conversions: %v", err)
		}

		found := false
		for _, conv := range resp.Data {
			if conv.ConversionID == conversionID {
				found = true
				t.Logf("Found our conversion in the list:")
				t.Logf("  Status: %s", conv.ConversionStatus)
				break
			}
		}

		if !found {
			t.Logf("Note: Conversion not found in first page of results")
		} else {
			t.Log("Successfully verified conversion in list")
		}
	})
}

func TestConversionErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	// Test getting non-existent conversion
	t.Run("GetNonExistentConversion", func(t *testing.T) {
		_, err := client.Banking.Conversions.Get(ctx, "non-existent-id")
		if err == nil {
			t.Error("Expected error when getting non-existent conversion")
		}
		t.Logf("Got expected error: %v", err)
	})

	// Test creating conversion with invalid currency
	t.Run("CreateConversionInvalidCurrency", func(t *testing.T) {
		today := time.Now().Format("2006-01-02")
		req := &banking.CreateConversionRequest{
			QuoteID:        "invalid-quote-id",
			SellCurrency:   "INVALID",
			SellAmount:     "100.00",
			BuyCurrency:    "EUR",
			ConversionDate: today,
		}

		_, err := client.Banking.Conversions.Create(ctx, req)
		if err == nil {
			t.Error("Expected error when creating conversion with invalid currency")
		}
		t.Logf("Got expected error: %v", err)
	})

	// Test creating quote with invalid amount
	t.Run("CreateQuoteInvalidAmount", func(t *testing.T) {
		today := time.Now().Format("2006-01-02")
		req := &banking.CreateQuoteRequest{
			SellCurrency:    "USD",
			SellAmount:      "invalid",
			BuyCurrency:     "SGD",
			ConversionDate:  today,
			TransactionType: "conversion",
		}

		_, err := client.Banking.Conversions.CreateQuote(ctx, req)
		if err == nil {
			t.Error("Expected error when creating quote with invalid amount")
		}
		t.Logf("Got expected error: %v", err)
	})
}

func TestConversionPagination(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	// Get first page
	req1 := &banking.ListConversionsRequest{
		PageSize:   5,
		PageNumber: 1,
	}

	resp1, err := client.Banking.Conversions.List(ctx, req1)
	if err != nil {
		t.Fatalf("Failed to list conversions page 1: %v", err)
	}

	t.Logf("Page 1: %d items (Total: %d items, %d pages)",
		len(resp1.Data), resp1.TotalItems, resp1.TotalPages)

	if resp1.TotalPages <= 1 {
		t.Skip("Not enough data to test pagination")
	}

	// Get second page
	req2 := &banking.ListConversionsRequest{
		PageSize:   5,
		PageNumber: 2,
	}

	resp2, err := client.Banking.Conversions.List(ctx, req2)
	if err != nil {
		t.Fatalf("Failed to list conversions page 2: %v", err)
	}

	t.Logf("Page 2: %d items", len(resp2.Data))

	// Verify pages have different data
	if len(resp1.Data) > 0 && len(resp2.Data) > 0 {
		if resp1.Data[0].ConversionID == resp2.Data[0].ConversionID {
			t.Error("Expected different data on different pages")
		}
	}
}

func TestConversionDatesMultipleCurrencyPairs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	currencyPairs := []struct {
		from string
		to   string
	}{
		{"USD", "EUR"},
		{"USD", "GBP"},
		{"EUR", "USD"},
		{"GBP", "USD"},
	}

	for _, pair := range currencyPairs {
		t.Run(fmt.Sprintf("%s_%s", pair.from, pair.to), func(t *testing.T) {
			dates, err := client.Banking.Conversions.ListConversionDates(ctx, pair.from, pair.to)
			if err != nil {
				t.Logf("Failed to get conversion dates for %s -> %s: %v", pair.from, pair.to, err)
				return
			}

			t.Logf("Available dates for %s -> %s: %d", pair.from, pair.to, len(dates))

			if len(dates) > 0 {
				// Find optimized date
				for _, date := range dates {
					if date.OptimizedDate {
						t.Logf("  Optimized date: %s", date.Date)
						break
					}
				}
			}
		})
	}
}
