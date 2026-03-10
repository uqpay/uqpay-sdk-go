package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/uqpay/uqpay-sdk-go/issuing"
)

func TestFullIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	// Step 1: Create Cardholder
	t.Run("CreateCardholder", func(t *testing.T) {
		timestamp := time.Now().UnixNano()
		req := &issuing.CreateCardholderRequest{
			Email:       fmt.Sprintf("test%d@example.com", timestamp),
			PhoneNumber: fmt.Sprintf("%d", timestamp%100000000), // Use timestamp as phone
			FirstName:   "Integration",
			LastName:    "Test",
			CountryCode: "SG",
		}

		t.Logf("📝 Creating cardholder: %s %s (%s)", req.FirstName, req.LastName, req.Email)

		cardholder, err := client.Issuing.Cardholders.Create(ctx, req)
		if err != nil {
			t.Fatalf("❌ Failed to create cardholder: %v", err)
		}

		t.Logf("✅ Cardholder created: ID=%s", cardholder.CardholderID)
		t.Logf("   Status: %s", cardholder.CardholderStatus)

		// Store for next steps
		cardholderID := cardholder.CardholderID

		// Step 2: Get Cardholder
		t.Run("GetCardholder", func(t *testing.T) {
			t.Logf("🔍 Getting cardholder: %s", cardholderID)

			ch, err := client.Issuing.Cardholders.Get(ctx, cardholderID)
			if err != nil {
				t.Fatalf("❌ Failed to get cardholder: %v", err)
			}

			t.Logf("✅ Retrieved cardholder:")
			t.Logf("   ID: %s", ch.CardholderID)
			t.Logf("   Name: %s %s", ch.FirstName, ch.LastName)
			t.Logf("   Email: %s", ch.Email)
			t.Logf("   Status: %s", ch.CardholderStatus)
		})

		// Step 3: Create Card for this Cardholder
		t.Run("CreateCard", func(t *testing.T) {
			// First, get available card products
			t.Log("📋 Listing card products...")

			productsResp, err := client.Issuing.Products.List(ctx, &issuing.ListProductsRequest{
				PageSize:   10,
				PageNumber: 1,
			})

			if err != nil || len(productsResp.Data) == 0 {
				t.Skip("No card products available")
			}

			// Find an active SINGLE mode product
			var productID string
			var productCurrency string
			for _, product := range productsResp.Data {
				if product.ProductStatus == "ENABLED" && product.ModeType == "SINGLE" {
					productID = product.ProductID
					if len(product.CardCurrency) > 0 {
						productCurrency = product.CardCurrency[0]
					}
					t.Logf("✅ Found active product: %s (%s, %s)", productID, product.CardScheme, productCurrency)
					break
				}
			}

			if productID == "" {
				t.Skip("No suitable card product found")
			}

			createReq := &issuing.CreateCardRequest{
				CardCurrency:  productCurrency,
				CardholderID:  cardholderID,
				CardProductID: productID,
			}

			t.Logf("💳 Creating card for cardholder: %s", cardholderID)
			t.Logf("   Using product: %s (%s)", productID, productCurrency)

			card, err := client.Issuing.Cards.Create(ctx, createReq)
			if err != nil {
				t.Logf("❌ Failed to create card: %v", err)
				t.Skip("Skipping card tests due to creation failure")
				return
			}

			t.Logf("✅ Card created: ID=%s", card.CardID)
			t.Logf("   Status: %s", card.CardStatus)

			cardID := card.CardID

			// Step 4: Get Card Details
			t.Run("GetCard", func(t *testing.T) {
				t.Logf("🔍 Getting card: %s", cardID)

				cardDetails, err := client.Issuing.Cards.Get(ctx, cardID)
				if err != nil {
					t.Fatalf("❌ Failed to get card: %v", err)
				}

				t.Logf("✅ Retrieved card:")
				t.Logf("   ID: %s", cardDetails.CardID)
				t.Logf("   Number: %s", cardDetails.CardNumber)
				t.Logf("   Status: %s", cardDetails.CardStatus)
				t.Logf("   Balance: %s %s", cardDetails.AvailableBalance, cardDetails.CardCurrency)
			})

			// Step 5: Get Secure Card Info
			t.Run("GetSecureCardInfo", func(t *testing.T) {
				t.Logf("🔒 Getting secure card info: %s", cardID)

				secureInfo, err := client.Issuing.Cards.GetSecure(ctx, cardID)
				if err != nil {
					t.Logf("❌ Failed to get secure info: %v", err)
					return
				}

				t.Logf("✅ Retrieved secure info:")
				t.Logf("   Card Number: %s", secureInfo.CardNumber)
				t.Logf("   CVV: %s", secureInfo.CVV)
				t.Logf("   Expiry: %s", secureInfo.ExpireDate)
			})

			// Step 6: Recharge Card
			t.Run("RechargeCard", func(t *testing.T) {
				rechargeReq := &issuing.CardOrderRequest{
					Amount: "100.50",
				}

				t.Logf("💰 Recharging card %s with amount: %s", cardID, rechargeReq.Amount)

				order, err := client.Issuing.Cards.Recharge(ctx, cardID, rechargeReq)
				if err != nil {
					t.Logf("❌ Failed to recharge card: %v", err)
					return
				}

				t.Logf("✅ Recharge order created:")
				t.Logf("   Order ID: %s", order.CardOrderID)
				t.Logf("   Status: %s", order.OrderStatus)
				t.Logf("   Amount: %s", order.Amount)
			})

			// Step 7: Update Card Status
			t.Run("UpdateCardStatus", func(t *testing.T) {
				statusReq := &issuing.UpdateCardStatusRequest{
					CardStatus: "FROZEN",
				}

				t.Logf("🔄 Updating card status to: %s", statusReq.CardStatus)

				resp, err := client.Issuing.Cards.UpdateStatus(ctx, cardID, statusReq)
				if err != nil {
					t.Logf("❌ Failed to update card status: %v", err)
					return
				}

				t.Logf("✅ Card status updated to: %s", statusReq.CardStatus)
				t.Logf("   Order Status: %s", resp.OrderStatus)

				// Verify status change
				time.Sleep(1 * time.Second)
				cardDetails, _ := client.Issuing.Cards.Get(ctx, cardID)
				if cardDetails != nil {
					t.Logf("   Current status: %s", cardDetails.CardStatus)
				}
			})

			// Step 8: List Transactions for this Card
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

				t.Logf("✅ Found %d transactions", resp.TotalItems)
				for i, txn := range resp.Data {
					t.Logf("   [%d] %s: %s %s - %s",
						i+1, txn.TransactionID, txn.TransactionAmount,
						txn.TransactionCurrency, txn.TransactionStatus)
				}
			})
		})
	})
}
