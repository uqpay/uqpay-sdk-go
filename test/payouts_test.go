package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/banking"
)

func TestPayouts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		resp, err := client.Banking.Payouts.List(ctx, &banking.ListPayoutsRequest{
			PageSize: 10, PageNumber: 1,
		})
		if err != nil {
			t.Logf("List payouts returned error: %v", err)
			return
		}

		t.Logf("Found %d payouts (total: %d, pages: %d)", len(resp.Data), resp.TotalItems, resp.TotalPages)
		if len(resp.Data) > 0 {
			p := resp.Data[0]
			t.Logf("  First: ID=%s, Amount=%s %s, Status=%s, Fee=%s %s",
				p.PayoutID, p.PayoutAmount, p.PayoutCurrency,
				p.PayoutStatus, p.FeeAmount, p.FeeCurrency)
			t.Logf("    Reason=%s, PurposeCode=%s, Date=%s",
				p.PayoutReason, p.PurposeCode, p.PayoutDate)
		}
	})

	t.Run("ListWithFilters", func(t *testing.T) {
		resp, err := client.Banking.Payouts.List(ctx, &banking.ListPayoutsRequest{
			PageSize: 10, PageNumber: 1, PayoutStatus: "COMPLETED", Currency: "USD",
		})
		if err != nil {
			t.Logf("List with filters returned error: %v", err)
			return
		}
		t.Logf("Found %d completed USD payouts (total: %d)", len(resp.Data), resp.TotalItems)
	})

	t.Run("ListByStatus", func(t *testing.T) {
		for _, status := range []string{"PENDING", "READY_TO_SEND", "COMPLETED", "FAILED"} {
			resp, err := client.Banking.Payouts.List(ctx, &banking.ListPayoutsRequest{
				PageSize: 10, PageNumber: 1, PayoutStatus: status,
			})
			if err != nil {
				t.Logf("  %s: error - %v", status, err)
				continue
			}
			t.Logf("  %s: %d found", status, resp.TotalItems)
		}
	})

	t.Run("Create", func(t *testing.T) {
		t.Skip("Skipping payout creation to avoid transaction costs")
	})

	t.Run("CreateWithInlineBeneficiary", func(t *testing.T) {
		t.Skip("Skipping inline beneficiary payout creation to avoid transaction costs")
	})

	t.Run("Get", func(t *testing.T) {
		listResp, err := client.Banking.Payouts.List(ctx, &banking.ListPayoutsRequest{
			PageSize: 10, PageNumber: 1,
		})
		if err != nil {
			t.Logf("Failed to list payouts: %v", err)
			return
		}
		if len(listResp.Data) == 0 {
			t.Skip("No payouts available to test Get")
		}

		id := listResp.Data[0].PayoutID
		resp, err := client.Banking.Payouts.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get payout: %v", err)
		}

		t.Logf("Get OK: ID=%s, Amount=%s %s, Status=%s",
			resp.PayoutID, resp.PayoutAmount, resp.PayoutCurrency, resp.PayoutStatus)
		t.Logf("  Fee=%s %s, Reason=%s, PurposeCode=%s",
			resp.FeeAmount, resp.FeeCurrency, resp.PayoutReason, resp.PurposeCode)
		t.Logf("  Created=%s, Ref=%s", resp.CreateTime, resp.ShortReferenceID)
		if resp.CompleteTime != nil && *resp.CompleteTime != "" {
			t.Logf("  Completed: %s", *resp.CompleteTime)
		}
		if resp.FailureReason != "" {
			t.Logf("  FailureReason: %s", resp.FailureReason)
		}
	})
}
