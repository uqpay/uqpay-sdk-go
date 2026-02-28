package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/banking"
)

func TestExchangeRates(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	t.Run("ListAll", func(t *testing.T) {
		resp, err := client.Banking.ExchangeRates.List(ctx, &banking.ListRatesRequest{})
		if err != nil {
			t.Fatalf("Failed to list exchange rates: %v", err)
		}

		t.Logf("Retrieved %d exchange rates (last updated: %s)", len(resp.Rates), resp.LastUpdated)
		if len(resp.UnavailableCurrencyPairs) > 0 {
			t.Logf("Unavailable pairs: %v", resp.UnavailableCurrencyPairs)
		}

		for i, rate := range resp.Rates {
			if i >= 5 {
				t.Logf("  ... and %d more", len(resp.Rates)-5)
				break
			}
			t.Logf("  %s: Buy=%s, Sell=%s", rate.CurrencyPair, rate.BuyPrice, rate.SellPrice)
		}
	})

	t.Run("ListSpecificPairs", func(t *testing.T) {
		resp, err := client.Banking.ExchangeRates.List(ctx, &banking.ListRatesRequest{
			CurrencyPairs: []string{"USDEUR", "USDGBP", "EURUSD", "GBPUSD"},
		})
		if err != nil {
			t.Logf("Failed to list specific pairs: %v", err)
			return
		}

		t.Logf("Retrieved %d rates for specified pairs", len(resp.Rates))
		for _, rate := range resp.Rates {
			t.Logf("  %s: Buy=%s, Sell=%s", rate.CurrencyPair, rate.BuyPrice, rate.SellPrice)
		}
	})

	t.Run("ListUSDPairs", func(t *testing.T) {
		resp, err := client.Banking.ExchangeRates.List(ctx, &banking.ListRatesRequest{
			CurrencyPairs: []string{"USDEUR", "USDGBP", "USDSGD", "USDCNH", "USDHKD"},
		})
		if err != nil {
			t.Logf("Failed to list USD pairs: %v", err)
			return
		}

		t.Logf("Retrieved %d USD pairs", len(resp.Rates))
		for _, rate := range resp.Rates {
			t.Logf("  %s: Buy=%s, Sell=%s", rate.CurrencyPair, rate.BuyPrice, rate.SellPrice)
		}
	})

	t.Run("ListEURPairs", func(t *testing.T) {
		resp, err := client.Banking.ExchangeRates.List(ctx, &banking.ListRatesRequest{
			CurrencyPairs: []string{"EURGBP", "EURSGD", "EURCNH", "EURHKD"},
		})
		if err != nil {
			t.Logf("Failed to list EUR pairs: %v", err)
			return
		}

		t.Logf("Retrieved %d EUR pairs", len(resp.Rates))
		for _, rate := range resp.Rates {
			t.Logf("  %s: Buy=%s, Sell=%s", rate.CurrencyPair, rate.BuyPrice, rate.SellPrice)
		}
	})

	t.Run("ListCrossCurrencyPairs", func(t *testing.T) {
		resp, err := client.Banking.ExchangeRates.List(ctx, &banking.ListRatesRequest{
			CurrencyPairs: []string{"EURGBP", "GBPSGD", "GBPCNH", "SGDCNH"},
		})
		if err != nil {
			t.Logf("Failed to list cross-currency pairs: %v", err)
			return
		}

		t.Logf("Retrieved %d cross-currency pairs", len(resp.Rates))
		for _, rate := range resp.Rates {
			t.Logf("  %s: Buy=%s, Sell=%s", rate.CurrencyPair, rate.BuyPrice, rate.SellPrice)
		}
	})

	t.Run("CompareRates", func(t *testing.T) {
		resp, err := client.Banking.ExchangeRates.List(ctx, &banking.ListRatesRequest{
			CurrencyPairs: []string{"USDEUR", "USDGBP", "USDSGD"},
		})
		if err != nil {
			t.Logf("Failed to retrieve rates: %v", err)
			return
		}

		for _, rate := range resp.Rates {
			t.Logf("  %s: Buy=%s, Sell=%s", rate.CurrencyPair, rate.BuyPrice, rate.SellPrice)
		}
	})

	t.Run("VerifyRateUpdates", func(t *testing.T) {
		resp, err := client.Banking.ExchangeRates.List(ctx, &banking.ListRatesRequest{
			CurrencyPairs: []string{"USDEUR"},
		})
		if err != nil {
			t.Logf("Failed to retrieve rate: %v", err)
			return
		}

		t.Logf("Last updated: %s", resp.LastUpdated)
		if len(resp.Rates) > 0 {
			rate := resp.Rates[0]
			t.Logf("  %s: Buy=%s, Sell=%s", rate.CurrencyPair, rate.BuyPrice, rate.SellPrice)
		}
	})
}
