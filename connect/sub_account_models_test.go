package connect

import (
	"encoding/json"
	"testing"
)

// TestSubAccountIndividualInfo_RequiredFieldsSerialize verifies that the
// individual_info object serializes the breaking-change required fields
// (effective 2026-03-19 and 2026-07-02) into the JSON request body with the
// correct keys. Before the gender/annual_income fields were added these keys
// were absent from the marshalled body (and the struct literal would not even
// compile), so this test pins the contract.
func TestSubAccountIndividualInfo_RequiredFieldsSerialize(t *testing.T) {
	info := SubAccountIndividualInfo{
		FirstNameEnglish:      "John",
		LastNameEnglish:       "Doe",
		Nationality:           "SG",
		PhoneNumber:           "+6591234567",
		EmailAddress:          "john.doe@example.com",
		DateOfBirth:           "1990-01-15",
		CountryOrTerritory:    "SG",
		StreetAddress:         "1 Raffles Place",
		ApartmentSuiteOrFloor: "#12-01",
		City:                  "Singapore",
		State:                 "Singapore",
		PostalCode:            "048616",
		EmploymentStatus:      EmploymentStatusEmployed,
		Industry:              "Information Technology/IT",
		JobTitle:              "Business and administration professionals",
		CompanyName:           "Acme Corp.",
		Gender:                GenderMale,
		AnnualIncome:          "85000",
	}

	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("failed to marshal SubAccountIndividualInfo: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to unmarshal serialized body: %v", err)
	}

	// All breaking-change required keys must be present in the request body.
	wantKeys := map[string]string{
		"employment_status": "Employed",
		"industry":          "Information Technology/IT",
		"job_title":         "Business and administration professionals",
		"company_name":      "Acme Corp.",
		"gender":            "MALE",
		"annual_income":     "85000",
		// state is in the spec's required list; it must always be emitted.
		"state": "Singapore",
		// optional but supported field
		"apartment_suite_or_floor": "#12-01",
	}

	for key, want := range wantKeys {
		v, ok := got[key]
		if !ok {
			t.Errorf("expected JSON key %q to be present in request body, but it was missing", key)
			continue
		}
		if s, _ := v.(string); s != want {
			t.Errorf("JSON key %q = %q, want %q", key, s, want)
		}
	}
}

// TestSubAccountIndividualInfo_StateAlwaysEmitted ensures state is treated as a
// required field (no omitempty): an empty value must still serialize the key so
// the API surfaces a clear "state required" error rather than the SDK silently
// dropping it.
func TestSubAccountIndividualInfo_StateAlwaysEmitted(t *testing.T) {
	info := SubAccountIndividualInfo{}

	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	for _, key := range []string{"state", "gender", "annual_income"} {
		if _, ok := got[key]; !ok {
			t.Errorf("required key %q must be emitted even when empty (no omitempty)", key)
		}
	}
}
