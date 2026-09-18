package issuing

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/v3/common"
	"github.com/uqpay/uqpay-sdk-go/v3/configuration"
)

func TestFrozenIssuingResponses(t *testing.T) {
	raw, err := os.ReadFile("testdata/issuing-responses.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Name, Operation, Path string
		Body                  json.RawMessage
	}
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatal(err)
	}
	for _, fixture := range cases {
		t.Run(fixture.Operation+"/"+fixture.Name, func(t *testing.T) {
			requests := make(chan *http.Request, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests <- r
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(fixture.Body)
			}))
			defer server.Close()
			api := common.NewAPIClient(&configuration.Configuration{Environment: &configuration.Environment{BaseURL: server.URL}, HTTPClient: server.Client()}, &staticTokenProvider{token: "offline"})
			client := NewClient(api)
			ctx := context.Background()
			var actual interface{}
			var card *RetrieveCardResponse
			switch fixture.Operation {
			case "cards.list":
				var page *ListCardsResponse
				page, err = client.Cards.List(ctx, &ListCardsRequest{PageSize: 10, PageNumber: 1})
				actual = page
				if err == nil {
					if len(page.Data) != 1 {
						t.Fatal(page)
					}
					card = &page.Data[0]
				}
			case "cards.get":
				card, err = client.Cards.Get(ctx, "card-1")
				actual = card
			case "cardholders.list":
				actual, err = client.Cardholders.List(ctx, &ListCardholdersRequest{PageSize: 10, PageNumber: 1})
			case "cardholders.get":
				actual, err = client.Cardholders.Get(ctx, "holder-1")
			case "products.list":
				actual, err = client.Products.List(ctx, &ListProductsRequest{PageSize: 10, PageNumber: 1})
			case "cards.status":
				actual, err = client.Cards.UpdateStatus(ctx, "card-1", &UpdateCardStatusRequest{CardStatus: "FROZEN"})
			default:
				t.Fatal(fixture.Operation)
			}
			if err != nil {
				t.Fatal(err)
			}
			req := <-requests
			method := "GET"
			if fixture.Operation == "cards.status" {
				method = "POST"
			}
			if req.Method != method || req.URL.Path != fixture.Path {
				t.Fatal(req.Method, req.URL)
			}
			var want interface{}
			if err := json.Unmarshal(fixture.Body, &want); err != nil {
				t.Fatal(err)
			}
			if card != nil {
				fields := want.(map[string]interface{})
				if fixture.Operation == "cards.list" {
					fields = fields["data"].([]interface{})[0].(map[string]interface{})
				}
				// FlexibleString and FlexibleStringMap expose normalized public values.
				if fixture.Operation == "cards.get" {
					fields["card_limit"] = "1.25"
				}
				var metadata map[string]string
				if value, ok := fields["metadata"]; ok && value != nil {
					bytes, _ := json.Marshal(value)
					if text, ok := value.(string); ok {
						bytes = []byte(text)
					}
					if len(bytes) > 0 {
						if err := json.Unmarshal(bytes, &metadata); err != nil {
							t.Fatal(err)
						}
					}
				}
				if !reflect.DeepEqual(map[string]string(card.Metadata), metadata) {
					t.Fatalf("metadata: %#v want %#v", card.Metadata, metadata)
				}
				delete(fields, "metadata") // Checked above, including nil versus an allocated empty map.
				risk, exists := fields["risk_controls"]
				if !exists || risk == nil {
					if card.RiskControls != nil {
						t.Fatal("invented risk controls")
					}
					delete(fields, "risk_controls")
				}
			}
			// Public scalar fields retain their existing zero values when omitted.
			if fixture.Name == "missing" {
				var holder *Cardholder
				switch v := actual.(type) {
				case *Cardholder:
					holder = v
				case *ListCardholdersResponse:
					if len(v.Data) != 1 {
						t.Fatal(v)
					}
					holder = &v.Data[0]
				}
				if holder != nil && (holder.Gender != nil || holder.Nationality != nil) {
					t.Fatal("invented gender/nationality")
				}
			}
			if product, ok := actual.(*ListProductsResponse); ok && fixture.Name == "optional-absent" {
				if len(product.Data) != 1 || product.Data[0].ModeType != "" || product.Data[0].MaxCardQuota != 0 {
					t.Fatal(product)
				}
			}
			if status, ok := actual.(*CardStatusResponse); ok && fixture.Name == "reason-absent" && status.UpdateReason != nil {
				t.Fatal(status)
			}
			encoded, err := json.Marshal(actual)
			if err != nil {
				t.Fatal(err)
			}
			var got interface{}
			if err := json.Unmarshal(encoded, &got); err != nil {
				t.Fatal(err)
			}
			assertProvidedFields(t, fixture.Operation, got, want)
		})
	}
}
