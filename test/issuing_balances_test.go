package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/issuing"
)

func TestIssuingBalances(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	t.Run("Retrieve", func(t *testing.T) {
		req := &issuing.RetrieveBalanceRequest{
			Currency: "USD",
		}

		resp, err := client.Issuing.Balances.Retrieve(ctx, req)
		if err != nil {
			t.Logf("❌ Retrieve issuing balance returned error: %v", err)
			return
		}

		t.Logf("✅ Issuing balance retrieved for %s", req.Currency)
		t.Logf("💰 Balance ID: %s", resp.BalanceID)
		t.Logf("   Currency: %s", resp.Currency)
		t.Logf("   Available: %.2f", resp.AvailableBalance)
		t.Logf("   Margin: %.2f", resp.MarginBalance)
		t.Logf("   Frozen: %.2f", resp.FrozenBalance)
		t.Logf("   Status: %s", resp.BalanceStatus)
		t.Logf("   Created: %s", resp.CreateTime)
		t.Logf("   Last Trade: %s", resp.LastTradeTime)
	})

	t.Run("List", func(t *testing.T) {
		req := &issuing.ListBalancesRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Issuing.Balances.List(ctx, req)
		if err != nil {
			t.Logf("❌ List issuing balances returned error: %v", err)
			return
		}

		t.Logf("✅ Found %d issuing balances (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("📊 Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			for i, balance := range resp.Data {
				t.Logf("💰 Balance %d: %s - Available: %.2f, Margin: %.2f, Status: %s",
					i+1, balance.Currency, balance.AvailableBalance, balance.MarginBalance, balance.BalanceStatus)
			}
		} else {
			t.Logf("ℹ️  No issuing balances found")
		}
	})

	t.Run("ListTransactions", func(t *testing.T) {
		req := &issuing.ListBalanceTransactionsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Issuing.Balances.ListTransactions(ctx, req)
		if err != nil {
			t.Logf("❌ List issuing balance transactions returned error: %v", err)
			return
		}

		t.Logf("✅ Found %d issuing balance transactions (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("📊 Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			txn := resp.Data[0]
			t.Logf("🔍 First transaction:")
			t.Logf("   ID: %s", txn.TransactionID)
			t.Logf("   Short ID: %s", txn.ShortTransactionID)
			t.Logf("   Account ID: %s", txn.AccountID)
			t.Logf("   Balance ID: %s", txn.BalanceID)
			t.Logf("   Type: %s", txn.TransactionType)
			t.Logf("   Amount: %.2f %s", txn.Amount, txn.Currency)
			t.Logf("   Status: %s", txn.TransactionStatus)
			t.Logf("   Ending Balance: %.2f", txn.EndingBalance)
			t.Logf("   Description: %s", txn.Description)
			t.Logf("   Created: %s", txn.CreateTime)
			t.Logf("   Completed: %s", txn.CompleteTime)
		} else {
			t.Logf("ℹ️  No issuing balance transactions found")
		}
	})

	t.Run("ListTransactionsWithTimeFilter", func(t *testing.T) {
		req := &issuing.ListBalanceTransactionsRequest{
			PageSize:   10,
			PageNumber: 1,
			StartTime:  "2024-01-01T00:00:00+00:00",
			EndTime:    "2024-12-31T23:59:59+00:00",
		}

		resp, err := client.Issuing.Balances.ListTransactions(ctx, req)
		if err != nil {
			t.Logf("❌ List issuing balance transactions with time filter returned error: %v", err)
			return
		}

		t.Logf("✅ Found %d transactions in date range (total: %d)", len(resp.Data), resp.TotalItems)

		if len(resp.Data) > 0 {
			for i, txn := range resp.Data {
				if i >= 5 {
					t.Logf("   ... and %d more", len(resp.Data)-5)
					break
				}
				t.Logf("💰 Transaction %d: %s %.2f %s, Status: %s",
					i+1, txn.TransactionType, txn.Amount, txn.Currency, txn.TransactionStatus)
			}
		}
	})
}
