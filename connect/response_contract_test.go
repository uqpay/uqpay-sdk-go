package connect

import (
	"encoding/json"
	"testing"
)

func TestAccountResponseContract(t *testing.T) {
	for _, raw := range []string{
		`{"entity_type":"COMPANY","business_details":{"legal_entity_name":"Example","registration_number":"reg-1","incorporation_date":"2020-01-01","merchant_category_code":"5734"}}`,
		`{"entity_type":"COMPANY","business_details":{"legal_entity_name_english":"Example","account_purpose":["CARD_ISSUING"],"product_description":null}}`,
		`{"entity_type":"INDIVIDUAL","person_details":{"first_name":"Test"},"residential_address":{"country":"SG","street_address":"Test Street"}}`,
	} {
		var account Account
		if err := json.Unmarshal([]byte(raw), &account); err != nil {
			t.Fatal(err)
		}
		out, _ := json.Marshal(account)
		var got, want map[string]interface{}
		_ = json.Unmarshal(out, &got)
		_ = json.Unmarshal([]byte(raw), &want)
		for k, v := range want {
			a, _ := json.Marshal(got[k])
			z, _ := json.Marshal(v)
			if string(a) != string(z) {
				t.Errorf("%s lost: %s want %s", k, a, z)
			}
		}
	}
}
