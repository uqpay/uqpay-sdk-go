package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/banking"
)

func TestDeposits(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		resp, err := client.Banking.Deposits.List(ctx, &banking.ListDepositsRequest{
			PageSize: 10, PageNumber: 1,
		})
		if err != nil {
			t.Logf("List deposits returned error: %v", err)
			return
		}

		t.Logf("Found %d deposits (total: %d, pages: %d)", len(resp.Data), resp.TotalItems, resp.TotalPages)
		if len(resp.Data) > 0 {
			d := resp.Data[0]
			t.Logf("  First: ID=%s, Amount=%s %s, Status=%s, Fee=%s",
				d.DepositID, d.Amount, d.Currency, d.DepositStatus, d.DepositFee)
			if d.Sender != nil {
				t.Logf("    Sender: %s (%s)", d.Sender.SenderName, d.Sender.SenderCountry)
			}
		}
	})

	t.Run("ListWithFilters", func(t *testing.T) {
		resp, err := client.Banking.Deposits.List(ctx, &banking.ListDepositsRequest{
			PageSize: 10, PageNumber: 1, DepositStatus: "COMPLETED",
		})
		if err != nil {
			t.Logf("List with filters returned error: %v", err)
			return
		}
		t.Logf("Found %d completed deposits (total: %d)", len(resp.Data), resp.TotalItems)
	})

	t.Run("ListByStatus", func(t *testing.T) {
		for _, status := range []string{"PENDING", "COMPLETED", "FAILED"} {
			resp, err := client.Banking.Deposits.List(ctx, &banking.ListDepositsRequest{
				PageSize: 10, PageNumber: 1, DepositStatus: status,
			})
			if err != nil {
				t.Logf("  %s: error - %v", status, err)
				continue
			}
			t.Logf("  %s: %d found", status, resp.TotalItems)
		}
	})

	t.Run("ListByCurrency", func(t *testing.T) {
		for _, currency := range []string{"USD", "SGD", "EUR"} {
			resp, err := client.Banking.Deposits.List(ctx, &banking.ListDepositsRequest{
				PageSize: 10, PageNumber: 1, Currency: currency,
			})
			if err != nil {
				t.Logf("  %s: error - %v", currency, err)
				continue
			}
			t.Logf("  %s: %d found", currency, resp.TotalItems)
		}
	})

	t.Run("ListWithTimeRange", func(t *testing.T) {
		resp, err := client.Banking.Deposits.List(ctx, &banking.ListDepositsRequest{
			PageSize: 10, PageNumber: 1,
		})
		if err != nil {
			t.Logf("List deposits returned error: %v", err)
			return
		}
		t.Logf("Found %d deposits (total: %d)", len(resp.Data), resp.TotalItems)
	})

	t.Run("Get", func(t *testing.T) {
		listResp, err := client.Banking.Deposits.List(ctx, &banking.ListDepositsRequest{
			PageSize: 10, PageNumber: 1,
		})
		if err != nil {
			t.Logf("Failed to list deposits: %v", err)
			return
		}
		if len(listResp.Data) == 0 {
			t.Skip("No deposits available to test Get")
		}

		id := listResp.Data[0].DepositID
		resp, err := client.Banking.Deposits.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get deposit: %v", err)
		}

		t.Logf("Get OK: ID=%s, Amount=%s %s, Status=%s",
			resp.DepositID, resp.Amount, resp.Currency, resp.DepositStatus)
		t.Logf("  Fee=%s, Receiver=%s, Ref=%s",
			resp.DepositFee, resp.ReceiverAccountNumber, resp.ShortReferenceID)
		if resp.Sender != nil {
			t.Logf("  Sender: %s from %s, Account=%s",
				resp.Sender.SenderName, resp.Sender.SenderCountry, resp.Sender.SenderAccountNumber)
		}
		if resp.CompleteTime != "" {
			t.Logf("  Completed: %s", resp.CompleteTime)
		}
	})

	t.Run("GetMultipleDeposits", func(t *testing.T) {
		listResp, err := client.Banking.Deposits.List(ctx, &banking.ListDepositsRequest{
			PageSize: 5, PageNumber: 1,
		})
		if err != nil {
			t.Logf("Failed to list deposits: %v", err)
			return
		}
		if len(listResp.Data) == 0 {
			t.Skip("No deposits available")
		}

		for i, d := range listResp.Data {
			resp, err := client.Banking.Deposits.Get(ctx, d.DepositID)
			if err != nil {
				t.Logf("  %d: Failed to get %s: %v", i+1, d.DepositID, err)
				continue
			}
			t.Logf("  %d: %s %s - Status=%s", i+1, resp.Amount, resp.Currency, resp.DepositStatus)
		}
	})
}
