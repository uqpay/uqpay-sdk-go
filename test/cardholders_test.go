package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/v2/issuing"
)

func TestCardholders(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		req := &issuing.ListCardholdersRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Issuing.Cardholders.List(ctx, req)
		if err != nil {
			t.Logf("List cardholders returned: %v", err)
			return
		}

		t.Logf("✅ Found %d cardholders (total: %d)", len(resp.Data), resp.TotalItems)

		if len(resp.Data) > 0 {
			cardholder := resp.Data[0]
			t.Logf("First cardholder: ID=%s, Name=%s %s, Email=%s",
				cardholder.CardholderID, cardholder.FirstName, cardholder.LastName, cardholder.Email)
		}
	})
}
