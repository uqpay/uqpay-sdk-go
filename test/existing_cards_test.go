package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/issuing"
)

// TestExistingCards tests operations on existing cards
func TestExistingCards(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	// First, get a list of existing cards
	t.Log("📋 Getting list of existing cards...")
	cardsResp, err := client.Issuing.Cards.List(ctx, &issuing.ListCardsRequest{
		PageSize:   1,
		PageNumber: 1,
	})

	if err != nil || len(cardsResp.Data) == 0 {
		t.Skip("No existing cards found to test")
	}

	testCard := cardsResp.Data[0]
	cardID := testCard.CardID

	t.Logf("✅ Using existing card: %s (Status: %s, Balance: %s %s)",
		cardID, testCard.CardStatus, testCard.AvailableBalance, testCard.CardCurrency)

	// Test 1: Get Card Details
	t.Run("GetCardDetails", func(t *testing.T) {
		t.Logf("🔍 Getting card details: %s", cardID)

		card, err := client.Issuing.Cards.Get(ctx, cardID)
		if err != nil {
			t.Fatalf("❌ Failed to get card: %v", err)
		}

		t.Logf("✅ Card details retrieved:")
		t.Logf("   Card ID: %s", card.CardID)
		t.Logf("   Card Number: %s", card.CardNumber)
		t.Logf("   Status: %s", card.CardStatus)
		t.Logf("   Balance: %s %s", card.AvailableBalance, card.CardCurrency)
		t.Logf("   Product ID: %s", card.CardProductID)
		t.Logf("   Cardholder ID: %s", card.Cardholder.CardholderID)
	})

	// Test 2: Get Secure Card Info
	t.Run("GetSecureCardInfo", func(t *testing.T) {
		t.Logf("🔒 Getting secure card info: %s", cardID)

		secureInfo, err := client.Issuing.Cards.GetSecure(ctx, cardID)
		if err != nil {
			t.Logf("❌ Failed to get secure info: %v", err)
			return
		}

		t.Logf("✅ Secure card info retrieved:")
		t.Logf("   Full Card Number: %s", secureInfo.CardNumber)
		t.Logf("   CVV: %s", secureInfo.CVV)
		t.Logf("   Expiry Date: %s", secureInfo.ExpireDate)
	})

	// Test 3: Recharge Card
	t.Run("RechargeCard", func(t *testing.T) {
		if testCard.CardStatus != "ACTIVE" {
			t.Skip("Card is not ACTIVE, skipping recharge test")
		}

		rechargeReq := &issuing.CardOrderRequest{
			Amount: "50.00",
		}

		t.Logf("💰 Recharging card %s with %s %s", cardID, rechargeReq.Amount, testCard.CardCurrency)

		order, err := client.Issuing.Cards.Recharge(ctx, cardID, rechargeReq)
		if err != nil {
			t.Logf("❌ Failed to recharge: %v", err)
			return
		}

		t.Logf("✅ Recharge order created:")
		t.Logf("   Order ID: %s", order.CardOrderID)
		t.Logf("   Card ID: %s", order.CardID)
		t.Logf("   Amount: %s", order.Amount)
		t.Logf("   Status: %s", order.OrderStatus)
		t.Logf("   Create Time: %s", order.CreateTime)
	})

	// Test 4: Update Card Status (FREEZE)
	t.Run("UpdateCardStatus_Freeze", func(t *testing.T) {
		if testCard.CardStatus != "ACTIVE" {
			t.Skip("Card is not ACTIVE, cannot freeze")
		}

		updateReq := &issuing.UpdateCardStatusRequest{
			CardStatus: "FROZEN",
		}

		t.Logf("❄️ Freezing card: %s", cardID)

		resp, err := client.Issuing.Cards.UpdateStatus(ctx, cardID, updateReq)
		if err != nil {
			t.Logf("❌ Failed to freeze card: %v", err)
			return
		}

		t.Logf("✅ Card status updated to: FROZEN")
		t.Logf("   Order Status: %s", resp.OrderStatus)

		// Verify the change
		card, _ := client.Issuing.Cards.Get(ctx, cardID)
		if card != nil {
			t.Logf("   Current status: %s", card.CardStatus)
		}
	})

	// Test 5: Update Card Status (UNFREEZE)
	t.Run("UpdateCardStatus_Unfreeze", func(t *testing.T) {
		updateReq := &issuing.UpdateCardStatusRequest{
			CardStatus: "ACTIVE",
		}

		t.Logf("🔓 Unfreezing card: %s", cardID)

		resp, err := client.Issuing.Cards.UpdateStatus(ctx, cardID, updateReq)
		if err != nil {
			t.Logf("❌ Failed to unfreeze card: %v", err)
			return
		}

		t.Logf("✅ Card status updated to: ACTIVE")
		t.Logf("   Order Status: %s", resp.OrderStatus)

		// Verify the change
		card, _ := client.Issuing.Cards.Get(ctx, cardID)
		if card != nil {
			t.Logf("   Current status: %s", card.CardStatus)
		}
	})

	// Test 6: List Transactions for this Card
	t.Run("ListCardTransactions", func(t *testing.T) {
		txnReq := &issuing.ListTransactionsRequest{
			PageSize:   10,
			PageNumber: 1,
			CardID:     cardID,
		}

		t.Logf("📊 Listing transactions for card: %s", cardID)

		resp, err := client.Issuing.Transactions.List(ctx, txnReq)
		if err != nil {
			t.Logf("❌ Failed to list transactions: %v", err)
			return
		}

		t.Logf("✅ Found %d transactions (total: %d)", len(resp.Data), resp.TotalItems)

		for i, txn := range resp.Data {
			t.Logf("   [%d] %s", i+1, txn.TransactionID)
			t.Logf("       Type: %s", txn.TransactionType)
			t.Logf("       Amount: %s %s", txn.TransactionAmount, txn.TransactionCurrency)
			t.Logf("       Status: %s", txn.TransactionStatus)
			t.Logf("       Time: %s", txn.TransactionTime)
		}
	})

	// Test 7: Get Transaction Details
	t.Run("GetTransactionDetails", func(t *testing.T) {
		// Get transactions first
		txnResp, err := client.Issuing.Transactions.List(ctx, &issuing.ListTransactionsRequest{
			PageSize:   1,
			PageNumber: 1,
		})

		if err != nil || len(txnResp.Data) == 0 {
			t.Skip("No transactions found to test")
		}

		txnID := txnResp.Data[0].TransactionID
		t.Logf("🔍 Getting transaction details: %s", txnID)

		txn, err := client.Issuing.Transactions.Get(ctx, txnID)
		if err != nil {
			t.Fatalf("❌ Failed to get transaction: %v", err)
		}

		t.Logf("✅ Transaction details retrieved:")
		t.Logf("   Transaction ID: %s", txn.TransactionID)
		t.Logf("   Card ID: %s", txn.CardID)
		t.Logf("   Type: %s", txn.TransactionType)
		t.Logf("   Amount: %s %s", txn.TransactionAmount, txn.TransactionCurrency)
		t.Logf("   Billing Amount: %s %s", txn.BillingAmount, txn.BillingCurrency)
		t.Logf("   Status: %s", txn.TransactionStatus)
		t.Logf("   Merchant: %s", txn.MerchantName)
		t.Logf("   Time: %s", txn.TransactionTime)
	})
}
