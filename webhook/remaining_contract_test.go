package webhook

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestRemainingSignedFields(t *testing.T) {
	raw, err := os.ReadFile("testdata/remaining-webhooks.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Kind string                 `json:"kind"`
		Data map[string]interface{} `json:"data"`
	}
	if err = json.Unmarshal(raw, &cases); err != nil {
		t.Fatal(err)
	}
	for _, c := range cases {
		t.Run(c.Kind, func(t *testing.T) {
			body, _ := json.Marshal(map[string]interface{}{"event_type": c.Kind, "data": c.Data})
			ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
			mac := hmac.New(sha512.New, []byte("offline-secret"))
			mac.Write(append(append([]byte{}, body...), []byte(ts)...))
			sig := hex.EncodeToString(mac.Sum(nil))
			verifier := NewVerifier("offline-secret")
			event, e := verifier.ConstructEvent(body, sig, ts)
			if e != nil {
				t.Fatal(e)
			}
			var got map[string]interface{}
			if e = json.Unmarshal(event.Data, &got); e != nil {
				t.Fatal(e)
			}
			if !reflect.DeepEqual(got, c.Data) {
				t.Fatal(got)
			}
			if _, e = verifier.ConstructEvent(append(body, ' '), sig, ts); e == nil {
				t.Fatal("changed bytes accepted")
			}
			var typed interface{}
			switch c.Kind {
			case "representative":
				typed = &Representative{}
			case "cardholder":
				typed = &Cardholder{}
			case "fee":
				typed = &NetworkProtectionFeeData{}
			case "transfer":
				typed = &IssuingTransferStatusChangedData{}
			case "transaction":
				typed = &CardTransactionData{}
			case "deposit":
				typed = &DepositData{}
			case "intent":
				typed = &PaymentIntentData{}
			case "rfi":
				return
			}
			if e = json.Unmarshal(event.Data, typed); e != nil {
				t.Fatal(e)
			}
			// Inspect public fields before omitempty serialization, including zero values.
			val := reflect.ValueOf(typed).Elem()
			typ := val.Type()
			for key, want := range c.Data {
				found := false
				for i := 0; i < typ.NumField(); i++ {
					tag := typ.Field(i).Tag.Get("json")
					if tag != key && tag != key+",omitempty" {
						continue
					}
					found = true
					b, _ := json.Marshal(val.Field(i).Interface())
					var actual interface{}
					json.Unmarshal(b, &actual)
					if key == "complete_time" && want == nil {
						if actual != "" {
							t.Fatal(actual)
						}
						continue
					}
					if !reflect.DeepEqual(actual, want) {
						t.Fatalf("%s: got %s want %#v", key, b, want)
					}
				}
				if !found {
					t.Fatalf("missing typed field %s", key)
				}
			}
		})
	}
}
