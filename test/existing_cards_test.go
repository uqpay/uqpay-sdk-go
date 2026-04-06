package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/issuing"
)

// TestExistingCards tests operations on existing cards fetched dynamically from the API.
func TestExistingCards(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	// ── List cards ─────────────────────────────────────────────────────────────
	t.Run("List", func(t *testing.T) {
		resp, err := client.Issuing.Cards.List(ctx, &issuing.ListCardsRequest{
			PageSize:   10,
			PageNumber: 1,
		})
		if err != nil {
			t.Fatalf("List cards failed: %v", err)
		}

		t.Logf("✅ Found %d cards (total: %d, pages: %d)", len(resp.Data), resp.TotalItems, resp.TotalPages)
		for i, card := range resp.Data {
			t.Logf("  [%d] %s | %s | %s %s | %s",
				i+1, card.CardID, card.CardStatus, card.AvailableBalance, card.CardCurrency, card.CardNumber)
		}
	})

	// Get the first card for subsequent tests
	listResp, err := client.Issuing.Cards.List(ctx, &issuing.ListCardsRequest{
		PageSize:   1,
		PageNumber: 1,
	})
	if err != nil || len(listResp.Data) == 0 {
		t.Skip("No existing cards found, skipping per-card tests")
	}

	card := listResp.Data[0]
	cardID := card.CardID
	t.Logf("Using card: %s (status=%s, balance=%s %s)", cardID, card.CardStatus, card.AvailableBalance, card.CardCurrency)

	// ── Get card details ────────────────────────────────────────────────────────
	t.Run("Get", func(t *testing.T) {
		resp, err := client.Issuing.Cards.Get(ctx, cardID)
		if err != nil {
			t.Fatalf("Get card failed: %v", err)
		}

		if resp.CardID != cardID {
			t.Errorf("CardID mismatch: got %s, want %s", resp.CardID, cardID)
		}
		if resp.CardStatus == "" {
			t.Error("CardStatus should not be empty")
		}

		t.Logf("✅ Card: %s | BIN=%s | Scheme=%s | Status=%s | Limit=%s | Balance=%s %s",
			resp.CardID, resp.CardBIN, resp.CardScheme, resp.CardStatus,
			resp.CardLimit, resp.AvailableBalance, resp.CardCurrency)
		t.Logf("   Cardholder: %s (%s %s)", resp.Cardholder.CardholderID, resp.Cardholder.FirstName, resp.Cardholder.LastName)
	})

	// ── Get secure card info ───────────────────────────────────────────────────
	t.Run("GetSecure", func(t *testing.T) {
		resp, err := client.Issuing.Cards.GetSecure(ctx, cardID)
		if err != nil {
			t.Logf("GetSecure returned: %v", err)
			return
		}

		if resp.CardNumber == "" {
			t.Error("CardNumber should not be empty")
		}

		t.Logf("✅ Secure info: PAN=%s | CVV=%s | Expiry=%s", resp.CardNumber, resp.CVV, resp.ExpireDate)
	})

	// ── Create PAN token ───────────────────────────────────────────────────────
	t.Run("CreatePANToken", func(t *testing.T) {
		resp, err := client.Issuing.Cards.CreatePANToken(ctx, cardID)
		if err != nil {
			t.Logf("CreatePANToken returned: %v", err)
			return
		}

		if resp.Token == "" {
			t.Error("Token should not be empty")
		}

		t.Logf("✅ PAN token: %s (expires in %ds, at %s)", resp.Token, resp.ExpiresIn, resp.ExpiresAt)
	})

	// ── Freeze / Unfreeze ──────────────────────────────────────────────────────
	t.Run("FreezeUnfreeze", func(t *testing.T) {
		if card.CardStatus != "ACTIVE" {
			t.Skipf("Card status is %s, skipping freeze/unfreeze", card.CardStatus)
		}

		freezeResp, err := client.Issuing.Cards.UpdateStatus(ctx, cardID, &issuing.UpdateCardStatusRequest{
			CardStatus: "FROZEN",
		})
		if err != nil {
			t.Logf("Freeze returned: %v", err)
			return
		}
		t.Logf("✅ Frozen: OrderStatus=%s", freezeResp.OrderStatus)

		// Note: unfreeze may fail immediately if freeze order is still PROCESSING
		unfreezeResp, err := client.Issuing.Cards.UpdateStatus(ctx, cardID, &issuing.UpdateCardStatusRequest{
			CardStatus: "ACTIVE",
		})
		if err != nil {
			t.Logf("Unfreeze returned (may be pending): %v", err)
		} else {
			t.Logf("✅ Unfrozen: OrderStatus=%s", unfreezeResp.OrderStatus)
		}
	})

	// ── Recharge ───────────────────────────────────────────────────────────────
	t.Run("Recharge", func(t *testing.T) {
		if card.CardStatus != "ACTIVE" {
			t.Skipf("Card status is %s, skipping recharge", card.CardStatus)
		}

		order, err := client.Issuing.Cards.Recharge(ctx, cardID, &issuing.CardOrderRequest{
			Amount: 1.00,
		})
		if err != nil {
			t.Logf("Recharge returned: %v", err)
			return
		}

		if order.CardOrderID == "" {
			t.Error("CardOrderID should not be empty")
		}

		t.Logf("✅ Recharge order: %s | Status=%s | Amount=%.2f", order.CardOrderID, order.OrderStatus, order.Amount)
	})

	// ── List transactions ──────────────────────────────────────────────────────
	t.Run("ListTransactions", func(t *testing.T) {
		resp, err := client.Issuing.Transactions.List(ctx, &issuing.ListTransactionsRequest{
			PageSize:   10,
			PageNumber: 1,
			CardID:     cardID,
		})
		if err != nil {
			t.Logf("List transactions returned: %v", err)
			return
		}

		t.Logf("✅ Found %d transactions (total: %d)", len(resp.Data), resp.TotalItems)
		for i, txn := range resp.Data {
			t.Logf("  [%d] %s | %s | %s %s | %s",
				i+1, txn.TransactionID, txn.TransactionType,
				txn.TransactionAmount, txn.TransactionCurrency, txn.TransactionStatus)
		}
	})

	// ── Get transaction detail ─────────────────────────────────────────────────
	t.Run("GetTransaction", func(t *testing.T) {
		txnList, err := client.Issuing.Transactions.List(ctx, &issuing.ListTransactionsRequest{
			PageSize:   1,
			PageNumber: 1,
		})
		if err != nil || len(txnList.Data) == 0 {
			t.Skip("No transactions found, skipping GetTransaction")
		}

		txnID := txnList.Data[0].TransactionID
		txn, err := client.Issuing.Transactions.Get(ctx, txnID)
		if err != nil {
			t.Fatalf("Get transaction failed: %v", err)
		}

		if txn.TransactionID != txnID {
			t.Errorf("TransactionID mismatch: got %s, want %s", txn.TransactionID, txnID)
		}

		t.Logf("✅ Transaction: %s | %s | %s %s | %s | %s",
			txn.TransactionID, txn.TransactionType,
			txn.TransactionAmount, txn.TransactionCurrency,
			txn.TransactionStatus, txn.TransactionTime)

		if txn.MerchantData != nil {
			t.Logf("   Merchant: %s | CategoryCode=%s | %s, %s", txn.MerchantData.Name, txn.MerchantData.CategoryCode, txn.MerchantData.City, txn.MerchantData.Country)
		}
		if txn.BillingAmount != "" {
			t.Logf("   Billing: %s %s", txn.BillingAmount, txn.BillingCurrency)
		}
	})
}
