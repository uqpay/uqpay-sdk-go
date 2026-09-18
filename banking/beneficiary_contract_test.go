package banking

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

func TestFrozenBeneficiaryContract(t *testing.T) {
	raw, err := os.ReadFile("testdata/beneficiary-contract.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Name, Operation, Path string
		Request, Body         json.RawMessage
	}
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatal(err)
	}
	for _, fixture := range cases {
		t.Run(fixture.Name+"/"+fixture.Operation, func(t *testing.T) {
			type capturedRequest struct {
				Method, Path string
				Body         map[string]interface{}
			}
			requests := make(chan capturedRequest, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var body map[string]interface{}
				if r.Method == "POST" {
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Error(err)
					}
				}
				requests <- capturedRequest{r.Method, r.URL.Path, body}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(fixture.Body)
			}))
			defer server.Close()
			api := common.NewAPIClient(&configuration.Configuration{Environment: &configuration.Environment{BaseURL: server.URL}, HTTPClient: server.Client()}, &vaStaticTokenProvider{})
			client := NewClient(api).Beneficiaries
			ctx := context.Background()
			var actual interface{}
			switch fixture.Operation {
			case "check":
				var request BeneficiaryCheckRequest
				if err := json.Unmarshal(fixture.Request, &request); err != nil {
					t.Fatal(err)
				}
				actual, err = client.Check(ctx, &request)
			case "list":
				actual, err = client.List(ctx, &ListBeneficiariesRequest{PageSize: 10, PageNumber: 1})
			case "get":
				actual, err = client.Get(ctx, "beneficiary-1")
			default:
				t.Fatal(fixture.Operation)
			}
			if err != nil {
				t.Fatal(err)
			}
			req := <-requests
			method := "GET"
			if fixture.Operation == "check" {
				method = "POST"
				var want map[string]interface{}
				if err := json.Unmarshal(fixture.Request, &want); err != nil {
					t.Fatal(err)
				}
				// Public scalar omitempty merges explicit empty with absent; both select IBAN.
				if want["account_number"] == "" {
					delete(want, "account_number")
				}
				if !reflect.DeepEqual(req.Body, want) {
					t.Fatalf("request: %#v want %#v", req.Body, want)
				}
			}
			if req.Method != method || req.Path != fixture.Path {
				t.Fatal(req)
			}
			var expected interface{}
			if err := json.Unmarshal(fixture.Body, &expected); err != nil {
				t.Fatal(err)
			}
			if fixture.Operation != "check" {
				var beneficiary *Beneficiary
				fields := expected.(map[string]interface{})
				if page, ok := actual.(*ListBeneficiariesResponse); ok {
					if len(page.Data) != 1 {
						t.Fatal(page)
					}
					beneficiary = &page.Data[0]
					fields = fields["data"].([]interface{})[0].(map[string]interface{})
				} else {
					beneficiary = actual.(*Beneficiary)
				}
				address := fields["address"].(map[string]interface{})
				if beneficiary.Address == nil {
					t.Fatal("missing address")
				}
				wantNationality, _ := address["nationality"].(string)
				if beneficiary.Address.Nationality != wantNationality {
					t.Fatal("nationality lost")
				}
				delete(address, "nationality") // Explicit empty and absent share the existing scalar zero value.
			}
			encoded, err := json.Marshal(actual)
			if err != nil {
				t.Fatal(err)
			}
			var got interface{}
			if err := json.Unmarshal(encoded, &got); err != nil {
				t.Fatal(err)
			}
			assertBeneficiaryFields(t, fixture.Name, got, expected)
		})
	}
}
func assertBeneficiaryFields(t *testing.T, path string, got, want interface{}) {
	t.Helper()
	switch w := want.(type) {
	case map[string]interface{}:
		g, ok := got.(map[string]interface{})
		if !ok {
			t.Fatal(path, got)
		}
		for k, v := range w {
			actual, ok := g[k]
			if !ok {
				t.Fatal(path, k, "missing")
			}
			assertBeneficiaryFields(t, path+"."+k, actual, v)
		}
	case []interface{}:
		g, ok := got.([]interface{})
		if !ok || len(g) != len(w) {
			t.Fatal(path, "array mismatch")
		}
		for i, v := range w {
			assertBeneficiaryFields(t, path, g[i], v)
		}
	default:
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("%s got %#v want %#v", path, got, want)
		}
	}
}
