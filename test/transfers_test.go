package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/banking"
)

func TestTransfers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		resp, err := client.Banking.Transfers.List(ctx, &banking.ListTransfersRequest{
			PageSize: 10, PageNumber: 1,
		})
		if err != nil {
			t.Logf("List transfers returned error: %v", err)
			return
		}

		t.Logf("Found %d transfers (total: %d, pages: %d)", len(resp.Data), resp.TotalItems, resp.TotalPages)
		if len(resp.Data) > 0 {
			tr := resp.Data[0]
			t.Logf("  First: ID=%s, Amount=%s %s, Status=%s",
				tr.TransferID, tr.TransferAmount, tr.TransferCurrency, tr.TransferStatus)
			t.Logf("    From=%s, To=%s, Reason=%s",
				tr.SourceAccountName, tr.DestinationAccountName, tr.TransferReason)
		}
	})

	t.Run("ListWithFilters", func(t *testing.T) {
		resp, err := client.Banking.Transfers.List(ctx, &banking.ListTransfersRequest{
			PageSize: 10, PageNumber: 1, TransferStatus: "completed", Currency: "USD",
		})
		if err != nil {
			t.Logf("List with filters returned error: %v", err)
			return
		}
		t.Logf("Found %d completed USD transfers (total: %d)", len(resp.Data), resp.TotalItems)
	})

	t.Run("Create", func(t *testing.T) {
		t.Skip("Skipping create test - requires valid source and target account IDs")
	})

	t.Run("Get", func(t *testing.T) {
		listResp, err := client.Banking.Transfers.List(ctx, &banking.ListTransfersRequest{
			PageSize: 10, PageNumber: 1,
		})
		if err != nil {
			t.Logf("Failed to list transfers: %v", err)
			return
		}
		if len(listResp.Data) == 0 {
			t.Skip("No transfers available to test Get")
		}

		id := listResp.Data[0].TransferID
		resp, err := client.Banking.Transfers.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get transfer: %v", err)
		}

		t.Logf("Get OK: ID=%s, Amount=%s %s, Status=%s",
			resp.TransferID, resp.TransferAmount, resp.TransferCurrency, resp.TransferStatus)
		t.Logf("  From=%s, To=%s, Ref=%s",
			resp.SourceAccountName, resp.DestinationAccountName, resp.ShortReferenceID)
		if resp.CompleteTime != "" {
			t.Logf("  Completed: %s", resp.CompleteTime)
		}
	})
}
