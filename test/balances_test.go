package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/banking"
)

func TestBalances(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	t.Run("Get", func(t *testing.T) {
		resp, err := client.Banking.Balances.Get(ctx, "USD")
		if err != nil {
			t.Logf("Get balance returned error: %v", err)
			return
		}

		t.Logf("Balance: ID=%s, Currency=%s, Available=%s, Status=%s",
			resp.BalanceID, resp.Currency, resp.AvailableBalance, resp.BalanceStatus)
		t.Logf("  Prepaid=%s, Margin=%s, Frozen=%s", resp.PrepaidBalance, resp.MarginBalance, resp.FrozenBalance)
		t.Logf("  Created=%s, LastTrade=%s", resp.CreateTime, resp.LastTradeTime)
	})

	t.Run("List", func(t *testing.T) {
		resp, err := client.Banking.Balances.List(ctx, &banking.ListBalancesRequest{
			PageSize: 10, PageNumber: 1,
		})
		if err != nil {
			t.Logf("List balances returned error: %v", err)
			return
		}

		t.Logf("Found %d balances (total: %d, pages: %d)", len(resp.Data), resp.TotalItems, resp.TotalPages)
		for i, b := range resp.Data {
			if i >= 3 {
				break
			}
			t.Logf("  %d: %s - Available=%s, Status=%s", i+1, b.Currency, b.AvailableBalance, b.BalanceStatus)
		}
	})

	t.Run("ListTransactions", func(t *testing.T) {
		resp, err := client.Banking.Balances.ListTransactions(ctx, &banking.ListBalanceTransactionsRequest{
			PageSize: 10, PageNumber: 1,
		})
		if err != nil {
			t.Logf("List transactions returned error: %v", err)
			return
		}

		t.Logf("Found %d transactions (total: %d)", len(resp.Data), resp.TotalItems)
		if len(resp.Data) > 0 {
			txn := resp.Data[0]
			t.Logf("  First: ID=%s, Type=%s, Amount=%s %s, Status=%s, CreditDebit=%s",
				txn.TransactionID, txn.TransactionType, txn.Amount, txn.Currency,
				txn.TransactionStatus, txn.CreditDebitType)
			t.Logf("    Ref=%s, Way=%s, Complete=%s", txn.ReferenceID, txn.TransactionWay, txn.CompleteTime)
		}
	})

	t.Run("ListTransactionsWithFilters", func(t *testing.T) {
		resp, err := client.Banking.Balances.ListTransactions(ctx, &banking.ListBalanceTransactionsRequest{
			PageSize:        10,
			PageNumber:      1,
			Currency:        "USD",
			TransactionType: "CONVERSION",
		})
		if err != nil {
			t.Logf("List transactions with filters returned error: %v", err)
			return
		}

		t.Logf("Found %d USD CONVERSION transactions (total: %d)", len(resp.Data), resp.TotalItems)
	})

	t.Run("ListTransactionsByType", func(t *testing.T) {
		types := []string{"DEPOSIT", "PAYOUT", "TRANSFER", "CONVERSION", "FEE"}
		for _, txnType := range types {
			resp, err := client.Banking.Balances.ListTransactions(ctx, &banking.ListBalanceTransactionsRequest{
				PageSize: 10, PageNumber: 1, TransactionType: txnType,
			})
			if err != nil {
				t.Logf("  %s: error - %v", txnType, err)
				continue
			}
			t.Logf("  %s: %d found", txnType, resp.TotalItems)
		}
	})
}
