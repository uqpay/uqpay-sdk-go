package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/issuing"
)

func TestCards(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		req := &issuing.ListCardsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Issuing.Cards.List(ctx, req)
		if err != nil {
			t.Logf("List cards returned: %v", err)
			return
		}

		t.Logf("✅ Found %d cards (total: %d)", len(resp.Data), resp.TotalItems)

		if len(resp.Data) > 0 {
			card := resp.Data[0]
			t.Logf("First card: ID=%s, Status=%s, Balance=%s %s",
				card.CardID, card.CardStatus, card.AvailableBalance, card.CardCurrency)
		}
	})

	t.Run("Recharge", func(t *testing.T) {
		t.Skip("Skipping - requires valid card ID and balance")

		cardID := "test-card-id"
		req := &issuing.CardOrderRequest{
			Amount: "100.00",
		}

		order, err := client.Issuing.Cards.Recharge(ctx, cardID, req)
		if err != nil {
			t.Fatalf("Failed to recharge card: %v", err)
		}

		t.Logf("✅ Recharge order created: ID=%s, Status=%s, Amount=%s",
			order.CardOrderID, order.OrderStatus, order.Amount)
	})

	t.Run("Withdraw", func(t *testing.T) {
		t.Skip("Skipping - requires valid card ID with sufficient balance")

		cardID := "test-card-id"
		req := &issuing.CardOrderRequest{
			Amount: "50.00",
		}

		order, err := client.Issuing.Cards.Withdraw(ctx, cardID, req)
		if err != nil {
			t.Fatalf("Failed to withdraw from card: %v", err)
		}

		t.Logf("✅ Withdraw order created: ID=%s, Status=%s, Amount=%s",
			order.CardOrderID, order.OrderStatus, order.Amount)
	})
}
