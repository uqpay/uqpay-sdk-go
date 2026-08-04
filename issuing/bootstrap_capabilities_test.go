package issuing

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestBootstrapCapabilitiesMatchPublishedContract(t *testing.T) {
	client, requests := newRequestOptionsTestClient(t)
	ctx := context.Background()

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		query  map[string]string
		invoke func() error
	}{
		{
			name: "elevate limit", method: http.MethodPost,
			path: "/v1/issuing/cards/card_123/elevate_limit", body: `"limit_amount":1000`,
			invoke: func() error {
				_, err := client.Cards.ElevateLimit(ctx, "card_123", &ElevateLimitRequest{LimitAmount: 1000})
				return err
			},
		},
		{
			name: "enroll network protection", method: http.MethodPost,
			path: "/v1/issuing/cards/card_123/risk", body: `"action_code":"41"`,
			invoke: func() error {
				_, err := client.Cards.EnrollNetworkProtection(ctx, "card_123", &EnrollNetworkProtectionRequest{
					RiskControl: "network_protection", ActionCode: "41",
				})
				return err
			},
		},
		{
			name: "remove network protection", method: http.MethodDelete,
			path: "/v1/issuing/cards/card_123/risk", body: `"risk_control":"network_protection"`,
			invoke: func() error {
				_, err := client.Cards.RemoveNetworkProtection(ctx, "card_123", &RemoveNetworkProtectionRequest{
					RiskControl: "network_protection",
				})
				return err
			},
		},
		{
			name: "manage pin", method: http.MethodPost,
			path: "/v1/issuing/cards/manage/pin", body: `"pin":"1234"`,
			invoke: func() error {
				_, err := client.Cards.ManagePIN(ctx, &ManageCardPINRequest{
					CardID: "card_123", Type: "SET", PIN: "1234",
				})
				return err
			},
		},
		{
			name: "list arts", method: http.MethodGet,
			path: "/v1/issuing/cards/arts", query: map[string]string{"card_product_id": "product_123"},
			invoke: func() error {
				_, err := client.Cards.ListArts(ctx, &ListCardArtsRequest{CardProductID: "product_123"})
				return err
			},
		},
		{
			name: "set default art", method: http.MethodPost,
			path: "/v1/issuing/cards/arts/default", body: `"card_art_id":"art_123"`,
			invoke: func() error {
				_, err := client.Cards.SetDefaultArt(ctx, &SetDefaultCardArtRequest{CardArtID: "art_123"})
				return err
			},
		},
		{
			name: "list merchant brands", method: http.MethodGet,
			path:  "/v1/issuing/merchant_brands",
			query: map[string]string{"display_name": "Grab", "page_number": "1", "page_size": "10"},
			invoke: func() error {
				_, err := client.MerchantBrands.List(ctx, &ListMerchantBrandsRequest{
					DisplayName: "Grab", PageNumber: 1, PageSize: 10,
				})
				return err
			},
		},
		{
			name: "claim unsolicited refund", method: http.MethodPost,
			path: "/v1/issuing/transactions/unsolicited_refund/release",
			body: `"related_transaction_id":"txn_123"`,
			invoke: func() error {
				_, err := client.Transactions.ClaimUnsolicitedRefund(ctx, &ClaimUnsolicitedRefundRequest{
					RelatedTransactionID: "txn_123",
				})
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.invoke(); err != nil {
				t.Fatalf("operation returned an error: %v", err)
			}
			got := <-requests
			if got.method != tt.method {
				t.Errorf("method = %q, want %q", got.method, tt.method)
			}
			if got.path != tt.path {
				t.Errorf("path = %q, want %q", got.path, tt.path)
			}
			if tt.body != "" && !strings.Contains(got.body, tt.body) {
				t.Errorf("body = %q, want it to contain %q", got.body, tt.body)
			}
			query, err := url.ParseQuery(got.query)
			if err != nil {
				t.Fatalf("invalid query %q: %v", got.query, err)
			}
			for key, want := range tt.query {
				if gotValue := query.Get(key); gotValue != want {
					t.Errorf("query %s = %q, want %q", key, gotValue, want)
				}
			}
		})
	}
}
