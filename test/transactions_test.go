package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/v2/issuing"
)

func TestTransactions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		req := &issuing.ListTransactionsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Issuing.Transactions.List(ctx, req)
		if err != nil {
			t.Logf("List transactions returned: %v", err)
			return
		}

		t.Logf("✅ Found %d transactions (total: %d)", len(resp.Data), resp.TotalItems)

		if len(resp.Data) > 0 {
			txn := resp.Data[0]
			t.Logf("First transaction: ID=%s, Amount=%s %s, Status=%s",
				txn.TransactionID, txn.TransactionAmount, txn.TransactionCurrency, txn.TransactionStatus)
		}
	})
}
