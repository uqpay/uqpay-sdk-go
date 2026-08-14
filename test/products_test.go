package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/v2/issuing"
)

func TestProducts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	t.Run("ListProducts", func(t *testing.T) {
		req := &issuing.ListProductsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		t.Log("📋 Listing card products...")

		resp, err := client.Issuing.Products.List(ctx, req)
		if err != nil {
			t.Fatalf("❌ Failed to list products: %v", err)
		}

		t.Logf("✅ Found %d products (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("   Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			t.Log("\n📦 Product Details:")
			for i, product := range resp.Data {
				t.Logf("   [%d] Product ID: %s", i+1, product.ProductID)
				t.Logf("       Card BIN: %s", product.CardBin)
				t.Logf("       Card Scheme: %s", product.CardScheme)
				t.Logf("       Card Currency: %s", product.CardCurrency)
				t.Logf("       Mode Type: %s", product.ModeType)
				t.Logf("       Card Forms: %v", product.CardForm)
				t.Logf("       Status: %s", product.ProductStatus)
				t.Logf("       Max Card Quota: %d", product.MaxCardQuota)
				t.Logf("       No PIN Payment: %v", product.NoPinPaymentAmount)
				t.Logf("       Created: %s", product.CreateTime)
				t.Logf("")
			}
		}
	})
}
