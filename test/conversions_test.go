package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	uqpay "github.com/uqpay/uqpay-sdk-go"
	"github.com/uqpay/uqpay-sdk-go/banking"
)

// getAvailableConversionDate fetches available conversion dates and returns the first valid one
func getAvailableConversionDate(t *testing.T, client *uqpay.Client, ctx context.Context, from, to string) string {
	t.Helper()
	dates, err := client.Banking.Conversions.ListConversionDates(ctx, from, to)
	if err != nil {
		t.Fatalf("Failed to get conversion dates for %s->%s: %v", from, to, err)
	}
	for _, d := range dates {
		if d.Valid {
			return d.Date
		}
	}
	t.Skipf("No valid conversion dates available for %s->%s", from, to)
	return ""
}

func TestConversionCreateQuote(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	convDate := getAvailableConversionDate(t, client, ctx, "USD", "EUR")

	req := &banking.CreateQuoteRequest{
		SellCurrency:    "USD",
		SellAmount:      "100.00",
		BuyCurrency:     "EUR",
		ConversionDate:  convDate,
		TransactionType: "conversion",
	}

	t.Logf("Creating quote: %s -> %s, Amount: %s, Date: %s", req.SellCurrency, req.BuyCurrency, req.SellAmount, convDate)

	quote, err := client.Banking.Conversions.CreateQuote(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create quote: %v", err)
	}

	if quote.QuotePrice.QuoteID == "" {
		t.Error("Expected quote_id to be set")
	}
	if quote.QuotePrice.DirectRate == "" {
		t.Error("Expected direct_rate to be set")
	}

	t.Logf("Quote created: ID=%s, Rate=%s, Buy=%s %s",
		quote.QuotePrice.QuoteID, quote.QuotePrice.DirectRate, quote.BuyAmount, quote.BuyCurrency)
}

func TestConversionCreateQuoteWithSettlementDate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	convDate := getAvailableConversionDate(t, client, ctx, "USD", "EUR")

	req := &banking.CreateQuoteRequest{
		SellCurrency:    "USD",
		SellAmount:      "250.00",
		BuyCurrency:     "EUR",
		ConversionDate:  convDate,
		TransactionType: "conversion",
	}

	t.Logf("Creating quote: USD->EUR with specific date: %s", convDate)

	quote, err := client.Banking.Conversions.CreateQuote(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create quote: %v", err)
	}

	if quote.QuotePrice.QuoteID == "" {
		t.Error("Expected quote_id to be set")
	}

	t.Logf("Quote created: ID=%s, Rate=%s, Buy=%s EUR",
		quote.QuotePrice.QuoteID, quote.QuotePrice.DirectRate, quote.BuyAmount)
}

func TestConversionCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	convDate := getAvailableConversionDate(t, client, ctx, "USD", "EUR")

	// First create a quote
	quoteReq := &banking.CreateQuoteRequest{
		SellCurrency:    "USD",
		SellAmount:      "100.00",
		BuyCurrency:     "EUR",
		ConversionDate:  convDate,
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
		BuyCurrency:    "EUR",
		ConversionDate: convDate,
	}

	resp, err := client.Banking.Conversions.Create(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create conversion: %v", err)
	}

	if resp.ConversionID == "" {
		t.Error("Expected conversion_id to be set")
	}

	t.Logf("Conversion created: ID=%s, Ref=%s, Status=%s",
		resp.ConversionID, resp.ShortReferenceID, resp.Status)

	// Get the created conversion
	t.Run("GetConversion", func(t *testing.T) {
		conversion, err := client.Banking.Conversions.Get(ctx, resp.ConversionID)
		if err != nil {
			t.Fatalf("Failed to get conversion: %v", err)
		}

		if conversion.ConversionID != resp.ConversionID {
			t.Errorf("Expected conversion_id %s, got %s", resp.ConversionID, conversion.ConversionID)
		}

		t.Logf("Get OK: ID=%s, Status=%s, Rate=%s",
			conversion.ConversionID, conversion.ConversionStatus, conversion.ClientRate)
	})
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

	resp, err := client.Banking.Conversions.List(ctx, req)
	if err != nil {
		t.Fatalf("Failed to list conversions: %v", err)
	}

	t.Logf("Found %d conversions (total: %d, pages: %d)", len(resp.Data), resp.TotalItems, resp.TotalPages)

	for i, conv := range resp.Data {
		if i >= 3 {
			break
		}
		t.Logf("  %d: ID=%s, %s %s -> %s %s, Rate=%s, Status=%s",
			i+1, conv.ConversionID,
			conv.SellAmount, conv.SellCurrency,
			conv.BuyAmount, conv.BuyCurrency,
			conv.ClientRate, conv.ConversionStatus)
	}
}

func TestConversionListWithFilters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	req := &banking.ListConversionsRequest{
		PageSize:   10,
		PageNumber: 1,
		SellCurrency: "USD",
	}

	resp, err := client.Banking.Conversions.List(ctx, req)
	if err != nil {
		t.Fatalf("Failed to list conversions: %v", err)
	}

	t.Logf("Found %d conversions selling USD (total: %d)", len(resp.Data), resp.TotalItems)
}

func TestConversionListWithTimeRange(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -30)

	req := &banking.ListConversionsRequest{
		PageSize:   10,
		PageNumber: 1,
		StartTime:  startTime.UnixMilli(),
		EndTime:    endTime.UnixMilli(),
	}

	t.Logf("Listing conversions from %s to %s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))

	resp, err := client.Banking.Conversions.List(ctx, req)
	if err != nil {
		t.Fatalf("Failed to list conversions: %v", err)
	}

	t.Logf("Found %d conversions in the last 30 days (total: %d)", len(resp.Data), resp.TotalItems)
}

func TestConversionDates(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	dates, err := client.Banking.Conversions.ListConversionDates(ctx, "USD", "EUR")
	if err != nil {
		t.Fatalf("Failed to get conversion dates: %v", err)
	}

	if len(dates) == 0 {
		t.Error("Expected at least one conversion date")
	}

	t.Logf("Available conversion dates for USD->EUR: %d", len(dates))
	for i, d := range dates {
		t.Logf("  %d: %s (valid=%t)", i+1, d.Date, d.Valid)
	}

	// Verify date format (YYYY-MM-DD)
	for _, d := range dates {
		if len(d.Date) != 10 {
			t.Errorf("Expected date format YYYY-MM-DD, got %s", d.Date)
		}
	}
}

func TestConversionDatesMultipleCurrencyPairs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	pairs := []struct{ from, to string }{
		{"USD", "EUR"},
		{"USD", "GBP"},
		{"EUR", "USD"},
		{"GBP", "USD"},
	}

	for _, p := range pairs {
		t.Run(fmt.Sprintf("%s_%s", p.from, p.to), func(t *testing.T) {
			dates, err := client.Banking.Conversions.ListConversionDates(ctx, p.from, p.to)
			if err != nil {
				t.Logf("Failed: %v", err)
				return
			}
			validCount := 0
			for _, d := range dates {
				if d.Valid {
					validCount++
				}
			}
			t.Logf("%s->%s: %d dates (%d valid)", p.from, p.to, len(dates), validCount)
		})
	}
}

func TestConversionFullFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	// Step 1: Get conversion dates
	var convDate string
	t.Run("GetConversionDates", func(t *testing.T) {
		dates, err := client.Banking.Conversions.ListConversionDates(ctx, "USD", "EUR")
		if err != nil {
			t.Fatalf("Failed to get conversion dates: %v", err)
		}
		for _, d := range dates {
			if d.Valid {
				convDate = d.Date
				t.Logf("Using conversion date: %s", convDate)
				return
			}
		}
		t.Skip("No valid conversion dates available")
	})

	// Step 2: Create a quote
	var quoteID string
	t.Run("CreateQuote", func(t *testing.T) {
		req := &banking.CreateQuoteRequest{
			SellCurrency:    "USD",
			SellAmount:      "100.00",
			BuyCurrency:     "EUR",
			ConversionDate:  convDate,
			TransactionType: "conversion",
		}

		quote, err := client.Banking.Conversions.CreateQuote(ctx, req)
		if err != nil {
			t.Fatalf("Failed to create quote: %v", err)
		}

		quoteID = quote.QuotePrice.QuoteID
		t.Logf("Quote: ID=%s, Rate=%s, Buy=%s %s",
			quoteID, quote.QuotePrice.DirectRate, quote.BuyAmount, quote.BuyCurrency)
	})

	// Step 3: Create conversion using the quote
	var conversionID string
	t.Run("CreateConversion", func(t *testing.T) {
		req := &banking.CreateConversionRequest{
			QuoteID:        quoteID,
			SellCurrency:   "USD",
			SellAmount:     "100.00",
			BuyCurrency:    "EUR",
			ConversionDate: convDate,
		}

		conv, err := client.Banking.Conversions.Create(ctx, req)
		if err != nil {
			t.Fatalf("Failed to create conversion: %v", err)
		}

		conversionID = conv.ConversionID
		t.Logf("Conversion: ID=%s, Ref=%s, Status=%s",
			conv.ConversionID, conv.ShortReferenceID, conv.Status)
	})

	// Step 4: Get the conversion details
	t.Run("GetConversion", func(t *testing.T) {
		conversion, err := client.Banking.Conversions.Get(ctx, conversionID)
		if err != nil {
			t.Fatalf("Failed to get conversion: %v", err)
		}

		t.Logf("Get OK: ID=%s, Status=%s, %s %s -> %s %s",
			conversion.ConversionID, conversion.ConversionStatus,
			conversion.SellAmount, conversion.SellCurrency,
			conversion.BuyAmount, conversion.BuyCurrency)
	})

	// Step 5: Verify in list
	t.Run("ListConversions", func(t *testing.T) {
		resp, err := client.Banking.Conversions.List(ctx, &banking.ListConversionsRequest{
			PageSize: 10, PageNumber: 1,
		})
		if err != nil {
			t.Fatalf("Failed to list conversions: %v", err)
		}

		found := false
		for _, conv := range resp.Data {
			if conv.ConversionID == conversionID {
				found = true
				t.Logf("Found in list: Status=%s", conv.ConversionStatus)
				break
			}
		}
		if !found {
			t.Log("Note: Conversion not found in first page")
		}
	})
}

func TestConversionErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	t.Run("GetNonExistentConversion", func(t *testing.T) {
		_, err := client.Banking.Conversions.Get(ctx, "non-existent-id")
		if err == nil {
			t.Error("Expected error when getting non-existent conversion")
		}
		t.Logf("Got expected error: %v", err)
	})

	t.Run("CreateQuoteInvalidAmount", func(t *testing.T) {
		convDate := time.Now().Format("2006-01-02")
		req := &banking.CreateQuoteRequest{
			SellCurrency:    "USD",
			SellAmount:      "invalid",
			BuyCurrency:     "SGD",
			ConversionDate:  convDate,
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

	resp1, err := client.Banking.Conversions.List(ctx, &banking.ListConversionsRequest{
		PageSize: 5, PageNumber: 1,
	})
	if err != nil {
		t.Fatalf("Failed to list page 1: %v", err)
	}

	t.Logf("Page 1: %d items (Total: %d items, %d pages)", len(resp1.Data), resp1.TotalItems, resp1.TotalPages)

	if resp1.TotalPages <= 1 {
		t.Skip("Not enough data to test pagination")
	}

	resp2, err := client.Banking.Conversions.List(ctx, &banking.ListConversionsRequest{
		PageSize: 5, PageNumber: 2,
	})
	if err != nil {
		t.Fatalf("Failed to list page 2: %v", err)
	}

	t.Logf("Page 2: %d items", len(resp2.Data))

	if len(resp1.Data) > 0 && len(resp2.Data) > 0 {
		if resp1.Data[0].ConversionID == resp2.Data[0].ConversionID {
			t.Error("Expected different data on different pages")
		}
	}
}
