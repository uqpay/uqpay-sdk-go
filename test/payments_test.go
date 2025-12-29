package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/payment"
)

// ============================================================================
// Payment Intents Tests
// ============================================================================

func TestPaymentIntents(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	var createdIntentID string

	t.Run("Create", func(t *testing.T) {
		req := &payment.CreatePaymentIntentRequest{
			Amount:          "100.00",
			Currency:        "USD",
			MerchantOrderID: "test-order-001",
			Description:     "Test payment intent",
			ReturnURL:       "https://example.com/return",
			Metadata: map[string]string{
				"test": "true",
			},
		}

		resp, err := client.Payment.PaymentIntents.Create(ctx, req)
		if err != nil {
			t.Logf("Create payment intent returned error: %v", err)
			return
		}

		createdIntentID = resp.PaymentIntentID
		t.Logf("Payment intent created successfully")
		t.Logf("   ID: %s", resp.PaymentIntentID)
		t.Logf("   Amount: %s %s", resp.Amount, resp.Currency)
		t.Logf("   Status: %s", resp.IntentStatus)
		t.Logf("   Merchant Order ID: %s", resp.MerchantOrderID)
		t.Logf("   Description: %s", resp.Description)
		t.Logf("   Created: %s", resp.CreateTime)
	})

	t.Run("Get", func(t *testing.T) {
		// First list to get an ID if we don't have one
		if createdIntentID == "" {
			listReq := &payment.ListPaymentIntentsRequest{
				PageSize:   1,
				PageNumber: 1,
			}
			listResp, err := client.Payment.PaymentIntents.List(ctx, listReq)
			if err != nil {
				t.Logf("List payment intents failed, skipping Get test: %v", err)
				return
			}
			if len(listResp.Data) == 0 {
				t.Log("No payment intents available, skipping Get test")
				return
			}
			createdIntentID = listResp.Data[0].PaymentIntentID
		}

		resp, err := client.Payment.PaymentIntents.Get(ctx, createdIntentID)
		if err != nil {
			t.Logf("Get payment intent returned error: %v", err)
			return
		}

		t.Logf("Payment intent retrieved successfully")
		t.Logf("   ID: %s", resp.PaymentIntentID)
		t.Logf("   Amount: %s %s", resp.Amount, resp.Currency)
		t.Logf("   Status: %s", resp.IntentStatus)
		t.Logf("   Merchant Order ID: %s", resp.MerchantOrderID)
	})

	t.Run("List", func(t *testing.T) {
		req := &payment.ListPaymentIntentsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Payment.PaymentIntents.List(ctx, req)
		if err != nil {
			t.Logf("List payment intents returned error: %v", err)
			return
		}

		t.Logf("Found %d payment intents (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			for i, intent := range resp.Data {
				if i >= 3 {
					t.Logf("   ... and %d more", len(resp.Data)-3)
					break
				}
				t.Logf("   Intent %d: ID=%s, Amount=%s %s, Status=%s",
					i+1, intent.PaymentIntentID, intent.Amount, intent.Currency, intent.IntentStatus)
			}
		} else {
			t.Log("No payment intents found")
		}
	})

	t.Run("ListWithFilters", func(t *testing.T) {
		statuses := []string{"requires_payment_method", "requires_confirmation", "succeeded", "canceled"}

		for _, status := range statuses {
			req := &payment.ListPaymentIntentsRequest{
				PageSize:   5,
				PageNumber: 1,
				Status:     status,
			}

			resp, err := client.Payment.PaymentIntents.List(ctx, req)
			if err != nil {
				t.Logf("%s intents: error - %v", status, err)
				continue
			}

			t.Logf("%s intents: %d found", status, resp.TotalItems)
		}
	})

	t.Run("Update", func(t *testing.T) {
		if createdIntentID == "" {
			t.Log("No payment intent ID available, skipping Update test")
			return
		}

		req := &payment.UpdatePaymentIntentRequest{
			Description: "Updated test payment intent",
			Metadata: map[string]string{
				"updated": "true",
			},
		}

		resp, err := client.Payment.PaymentIntents.Update(ctx, createdIntentID, req)
		if err != nil {
			t.Logf("Update payment intent returned error: %v", err)
			return
		}

		t.Logf("Payment intent updated successfully")
		t.Logf("   ID: %s", resp.PaymentIntentID)
		t.Logf("   Description: %s", resp.Description)
		t.Logf("   Updated: %s", resp.UpdateTime)
	})

	t.Run("Cancel", func(t *testing.T) {
		if createdIntentID == "" {
			t.Log("No payment intent ID available, skipping Cancel test")
			return
		}

		req := &payment.CancelPaymentIntentRequest{
			CancellationReason: "requested_by_customer",
		}

		resp, err := client.Payment.PaymentIntents.Cancel(ctx, createdIntentID, req)
		if err != nil {
			t.Logf("Cancel payment intent returned error: %v", err)
			return
		}

		t.Logf("Payment intent canceled successfully")
		t.Logf("   ID: %s", resp.PaymentIntentID)
		t.Logf("   Status: %s", resp.IntentStatus)
	})
}

// ============================================================================
// Payment Attempts Tests
// ============================================================================

func TestPaymentAttempts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		req := &payment.ListPaymentAttemptsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Payment.PaymentAttempts.List(ctx, req)
		if err != nil {
			t.Logf("List payment attempts returned error: %v", err)
			return
		}

		t.Logf("Found %d payment attempts (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			attempt := resp.Data[0]
			t.Logf("First attempt:")
			t.Logf("   ID: %s", attempt.AttemptID)
			t.Logf("   Amount: %s %s", attempt.Amount, attempt.Currency)
			t.Logf("   Captured: %s, Refunded: %s", attempt.CapturedAmount, attempt.RefundedAmount)
			t.Logf("   Status: %s", attempt.AttemptStatus)
			if attempt.FailureCode != "" {
				t.Logf("   Failure: %s", attempt.FailureCode)
			}
			t.Logf("   Created: %s", attempt.CreateTime)
		} else {
			t.Log("No payment attempts found")
		}
	})

	t.Run("Get", func(t *testing.T) {
		// First list to get a valid ID
		listReq := &payment.ListPaymentAttemptsRequest{
			PageSize:   1,
			PageNumber: 1,
		}

		listResp, err := client.Payment.PaymentAttempts.List(ctx, listReq)
		if err != nil {
			t.Logf("List payment attempts failed, skipping Get test: %v", err)
			return
		}

		if len(listResp.Data) == 0 {
			t.Log("No payment attempts available, skipping Get test")
			return
		}

		attemptID := listResp.Data[0].AttemptID

		resp, err := client.Payment.PaymentAttempts.Get(ctx, attemptID)
		if err != nil {
			t.Logf("Get payment attempt returned error: %v", err)
			return
		}

		t.Logf("Payment attempt retrieved successfully")
		t.Logf("   ID: %s", resp.AttemptID)
		t.Logf("   Amount: %s %s", resp.Amount, resp.Currency)
		t.Logf("   Status: %s", resp.AttemptStatus)
	})

	t.Run("ListWithFilters", func(t *testing.T) {
		// Test filtering by status
		statuses := []string{"pending", "succeeded", "failed"}

		for _, status := range statuses {
			req := &payment.ListPaymentAttemptsRequest{
				PageSize:   5,
				PageNumber: 1,
				Status:     status,
			}

			resp, err := client.Payment.PaymentAttempts.List(ctx, req)
			if err != nil {
				t.Logf("%s attempts: error - %v", status, err)
				continue
			}

			t.Logf("%s attempts: %d found", status, resp.TotalItems)
		}
	})
}

// ============================================================================
// Payment Balances Tests
// ============================================================================

func TestPaymentBalances(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		resp, err := client.Payment.Balances.List(ctx)
		if err != nil {
			t.Logf("List payment balances returned error: %v", err)
			return
		}

		t.Logf("Found %d payment balances", len(resp.Data))

		if len(resp.Data) > 0 {
			for i, balance := range resp.Data {
				t.Logf("Balance %d: %s", i+1, balance.Currency)
				t.Logf("   Available: %s", balance.AvailableBalance)
				t.Logf("   Payable: %s", balance.PayableBalance)
				t.Logf("   Pending: %s", balance.PendingBalance)
				t.Logf("   Reserved: %s", balance.ReservedBalance)
			}
		} else {
			t.Log("No payment balances found")
		}
	})

	t.Run("Get", func(t *testing.T) {
		currency := "USD"

		resp, err := client.Payment.Balances.Get(ctx, currency)
		if err != nil {
			t.Logf("Get payment balance for %s returned error: %v", currency, err)
			return
		}

		t.Logf("Balance retrieved for %s", currency)
		t.Logf("   Balance ID: %s", resp.BalanceID)
		t.Logf("   Available: %s", resp.AvailableBalance)
		t.Logf("   Payable: %s", resp.PayableBalance)
		t.Logf("   Pending: %s", resp.PendingBalance)
		t.Logf("   Reserved: %s", resp.ReservedBalance)
	})

	t.Run("GetMultipleCurrencies", func(t *testing.T) {
		currencies := []string{"USD", "EUR", "GBP", "AUD", "SGD"}

		for _, currency := range currencies {
			resp, err := client.Payment.Balances.Get(ctx, currency)
			if err != nil {
				t.Logf("%s: error - %v", currency, err)
				continue
			}

			t.Logf("%s: Available=%s, Payable=%s",
				currency, resp.AvailableBalance, resp.PayableBalance)
		}
	})
}

// ============================================================================
// Payment Payouts Tests
// ============================================================================

func TestPaymentPayouts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		req := &payment.ListPayoutsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Payment.Payouts.List(ctx, req)
		if err != nil {
			t.Logf("List payouts returned error: %v", err)
			return
		}

		t.Logf("Found %d payouts (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			payout := resp.Data[0]
			t.Logf("First payout:")
			t.Logf("   ID: %s", payout.PayoutID)
			t.Logf("   Amount: %s %s", payout.PayoutAmount, payout.PayoutCurrency)
			t.Logf("   Status: %s", payout.PayoutStatus)
			t.Logf("   Note: %s", payout.InternalNote)
			t.Logf("   Descriptor: %s", payout.StatementDescriptor)
			t.Logf("   Created: %s", payout.CreateTime)
		} else {
			t.Log("No payouts found")
		}
	})

	t.Run("Get", func(t *testing.T) {
		// First list to get a valid ID
		listReq := &payment.ListPayoutsRequest{
			PageSize:   1,
			PageNumber: 1,
		}

		listResp, err := client.Payment.Payouts.List(ctx, listReq)
		if err != nil {
			t.Logf("List payouts failed, skipping Get test: %v", err)
			return
		}

		if len(listResp.Data) == 0 {
			t.Log("No payouts available, skipping Get test")
			return
		}

		payoutID := listResp.Data[0].PayoutID

		resp, err := client.Payment.Payouts.Get(ctx, payoutID)
		if err != nil {
			t.Logf("Get payout returned error: %v", err)
			return
		}

		t.Logf("Payout retrieved successfully")
		t.Logf("   ID: %s", resp.PayoutID)
		t.Logf("   Amount: %s %s", resp.PayoutAmount, resp.PayoutCurrency)
		t.Logf("   Status: %s", resp.PayoutStatus)
		t.Logf("   Note: %s", resp.InternalNote)
	})

	t.Run("ListWithFilters", func(t *testing.T) {
		statuses := []string{"pending", "processing", "completed", "failed"}

		for _, status := range statuses {
			req := &payment.ListPayoutsRequest{
				PageSize:   5,
				PageNumber: 1,
				Status:     status,
			}

			resp, err := client.Payment.Payouts.List(ctx, req)
			if err != nil {
				t.Logf("%s payouts: error - %v", status, err)
				continue
			}

			t.Logf("%s payouts: %d found", status, resp.TotalItems)
		}
	})
}

// ============================================================================
// Payment Refunds Tests
// ============================================================================

func TestPaymentRefunds(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		req := &payment.ListRefundsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Payment.Refunds.List(ctx, req)
		if err != nil {
			t.Logf("List refunds returned error: %v", err)
			return
		}

		t.Logf("Found %d refunds (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			refund := resp.Data[0]
			t.Logf("First refund:")
			t.Logf("   ID: %s", refund.PaymentRefundID)
			t.Logf("   Payment Attempt ID: %s", refund.PaymentAttemptID)
			t.Logf("   Amount: %s %s", refund.Amount, refund.Currency)
			t.Logf("   Status: %s", refund.RefundStatus)
			t.Logf("   Reason: %s", refund.Reason)
			t.Logf("   Created: %s", refund.CreateTime)
		} else {
			t.Log("No refunds found")
		}
	})

	t.Run("Get", func(t *testing.T) {
		// First list to get a valid ID
		listReq := &payment.ListRefundsRequest{
			PageSize:   1,
			PageNumber: 1,
		}

		listResp, err := client.Payment.Refunds.List(ctx, listReq)
		if err != nil {
			t.Logf("List refunds failed, skipping Get test: %v", err)
			return
		}

		if len(listResp.Data) == 0 {
			t.Log("No refunds available, skipping Get test")
			return
		}

		refundID := listResp.Data[0].PaymentRefundID

		resp, err := client.Payment.Refunds.Get(ctx, refundID)
		if err != nil {
			t.Logf("Get refund returned error: %v", err)
			return
		}

		t.Logf("Refund retrieved successfully")
		t.Logf("   ID: %s", resp.PaymentRefundID)
		t.Logf("   Payment Attempt ID: %s", resp.PaymentAttemptID)
		t.Logf("   Amount: %s %s", resp.Amount, resp.Currency)
		t.Logf("   Status: %s", resp.RefundStatus)
	})

	t.Run("ListWithFilters", func(t *testing.T) {
		statuses := []string{"pending", "succeeded", "failed"}

		for _, status := range statuses {
			req := &payment.ListRefundsRequest{
				PageSize:   5,
				PageNumber: 1,
				Status:     status,
			}

			resp, err := client.Payment.Refunds.List(ctx, req)
			if err != nil {
				t.Logf("%s refunds: error - %v", status, err)
				continue
			}

			t.Logf("%s refunds: %d found", status, resp.TotalItems)
		}
	})
}

// ============================================================================
// Payment Reports Tests
// ============================================================================

func TestPaymentReports(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	t.Run("ListSettlements", func(t *testing.T) {
		req := &payment.ListSettlementsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Payment.Reports.ListSettlements(ctx, req)
		if err != nil {
			t.Logf("List settlements returned error: %v", err)
			return
		}

		t.Logf("Found %d settlements (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			settlement := resp.Data[0]
			t.Logf("First settlement:")
			t.Logf("   ID: %s", settlement.SettlementID)
			t.Logf("   Payment Intent: %s", settlement.PaymentIntentID)
			t.Logf("   Transaction Amount: %s %s", settlement.TransactionAmount, settlement.TransactionCurrency)
			t.Logf("   Settlement Amount: %s %s", settlement.SettlementAmount, settlement.SettlementCurrency)
			t.Logf("   Net Settlement: %s", settlement.NetSettlementAmount)
			t.Logf("   Total Fee: %s", settlement.TotalFeeAmount)
			t.Logf("   Status: %s", settlement.SettlementStatus)
			t.Logf("   Settlement Date: %s", settlement.SettlementDate)
		} else {
			t.Log("No settlements found")
		}
	})

	t.Run("ListSettlementsWithDateRange", func(t *testing.T) {
		// Test with date range (last 30 days)
		req := &payment.ListSettlementsRequest{
			SettledStartTime: "2024-01-01T00:00:00Z",
			SettledEndTime:   "2024-12-31T23:59:59Z",
			PageSize:         10,
			PageNumber:       1,
		}

		resp, err := client.Payment.Reports.ListSettlements(ctx, req)
		if err != nil {
			t.Logf("List settlements with date range returned error: %v", err)
			return
		}

		t.Logf("Found %d settlements in date range (total: %d)", len(resp.Data), resp.TotalItems)

		if len(resp.Data) > 0 {
			for i, settlement := range resp.Data {
				if i >= 3 {
					t.Logf("   ... and %d more", len(resp.Data)-3)
					break
				}
				t.Logf("Settlement %d: ID=%s, Amount=%s %s, Status=%s, Date=%s",
					i+1, settlement.SettlementID, settlement.SettlementAmount, settlement.SettlementCurrency, settlement.SettlementStatus, settlement.SettlementDate)
			}
		}
	})
}
