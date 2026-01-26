package webhook

import (
	"encoding/json"
	"testing"
)

// ============================================================================
// Test Fixtures - Account Webhooks
// ============================================================================

const accountCreateWebhookJSON = `{
    "version": "V1.6.0",
    "event_name": "ONBOARDING",
    "event_type": "onboarding.account.create",
    "event_id": "8a78af1e-de83-43a5-b177-ecbc6a8a9fc6",
    "source_id": "f5bb6498-552e-40a5-b14b-616aa04ac1c1",
    "data": {
        "account_id": "f5bb6498-552e-40a5-b14b-616aa04ac1c1",
        "direct_id": "0",
        "short_reference_id": "B1020563456",
        "email": "example@company.com",
        "account_name": "UQPAY PTE LTD.",
        "country": "US",
        "status": "PROCESSING",
        "idv_status": "PENDING",
        "verification_status": "APPROVED",
        "review_reason": "No risk.",
        "entity_type": "COMPANY",
        "contact_details": {
            "email": "example@company.com",
            "phone": "+6588880000"
        },
        "business_details": {
            "legal_entity_name": "ソフトバンクインターナショナル株式会社",
            "legal_entity_name_english": "UQPAY PTE LTD.",
            "incorporation_date": "2013-04-15",
            "registration_number": "201901754K",
            "business_structure": "SOLE_PROPRIETOR",
            "product_description": "Design and produce high-tech wireless headphones.",
            "merchant_category_code": "5733",
            "estimated_worker_count": "BS001",
            "monthly_estimated_revenue": {
                "amount": "TM001",
                "currency": "SGD"
            },
            "account_purpose": [
                "COLLECTION",
                "PAYOUT"
            ],
            "identifier": {
                "type": "VAT",
                "number": "XXX-XX-XXXX"
            },
            "website_url": "https://yourcompany.com"
        },
        "registration_address": {
            "city": "San Francisco",
            "line1": "9 N Buona Vista Dr",
            "state": "CA",
            "line2": "THE METROPOLIS",
            "postal_code": "94103"
        },
        "business_address": [
            {
                "city": "Singapore",
                "country": "SG",
                "line1": "9 N Buona Vista Dr",
                "state": "SG",
                "line2": "THE METROPOLIS",
                "postal_code": "138666"
            }
        ],
        "representatives": [
            {
                "roles": "DIRECTOR",
                "first_name": "Mock",
                "last_name": "Toy",
                "nationality": "SG",
                "date_of_birth": "2024-01-28",
                "share_percentage": "20.01",
                "identification": {
                    "type": "PASSPORT",
                    "id_number": "27738277K"
                },
                "residential_address": {
                    "city": "Singapore",
                    "country": "SG",
                    "line1": "9 N Buona Vista Dr",
                    "state": "SG",
                    "line2": "THE METROPOLIS",
                    "postal_code": "138666"
                },
                "as_applicant": false
            }
        ]
    }
}`

const accountUpdateWebhookJSON = `{
    "version": "V1.6.0",
    "event_name": "ONBOARDING",
    "event_type": "onboarding.account.update",
    "event_id": "9b89bf2f-ef94-54b6-c288-fdcd7b9b0fd7",
    "source_id": "f5bb6498-552e-40a5-b14b-616aa04ac1c1",
    "data": {
        "account_id": "f5bb6498-552e-40a5-b14b-616aa04ac1c1",
        "email": "updated@company.com",
        "account_name": "UQPAY PTE LTD.",
        "country": "US",
        "status": "ACTIVE",
        "verification_status": "APPROVED",
        "entity_type": "COMPANY"
    }
}`

const minimalAccountWebhookJSON = `{
    "version": "V1.6.0",
    "event_name": "ONBOARDING",
    "event_type": "onboarding.account.create",
    "event_id": "minimal-event-id",
    "source_id": "minimal-source-id",
    "data": {
        "account_id": "minimal-account-id",
        "email": "test@test.com",
        "account_name": "Test",
        "country": "US",
        "status": "PROCESSING",
        "entity_type": "INDIVIDUAL"
    }
}`

// ============================================================================
// Account Create Tests
// ============================================================================

func TestParseAccountData_Create(t *testing.T) {
	var event Event
	err := json.Unmarshal([]byte(accountCreateWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse event: %v", err)
	}

	account, err := event.ParseAccountData()
	if err != nil {
		t.Fatalf("Failed to parse account data: %v", err)
	}

	// Verify core fields
	if account.AccountID != "f5bb6498-552e-40a5-b14b-616aa04ac1c1" {
		t.Errorf("AccountID mismatch: got %s", account.AccountID)
	}
	if account.Email != "example@company.com" {
		t.Errorf("Email mismatch: got %s", account.Email)
	}
	if account.AccountName != "UQPAY PTE LTD." {
		t.Errorf("AccountName mismatch: got %s", account.AccountName)
	}
	if account.Country != "US" {
		t.Errorf("Country mismatch: got %s", account.Country)
	}
	if account.Status != "PROCESSING" {
		t.Errorf("Status mismatch: got %s", account.Status)
	}
	if account.EntityType != "COMPANY" {
		t.Errorf("EntityType mismatch: got %s", account.EntityType)
	}
	if account.IDVStatus != "PENDING" {
		t.Errorf("IDVStatus mismatch: got %s", account.IDVStatus)
	}
	if account.VerificationStatus != "APPROVED" {
		t.Errorf("VerificationStatus mismatch: got %s", account.VerificationStatus)
	}
}

func TestParseAccountData_ContactDetails(t *testing.T) {
	var event Event
	json.Unmarshal([]byte(accountCreateWebhookJSON), &event)

	account, err := event.ParseAccountData()
	if err != nil {
		t.Fatalf("Failed to parse account data: %v", err)
	}

	if account.ContactDetails == nil {
		t.Fatal("ContactDetails should not be nil")
	}

	if account.ContactDetails.Email != "example@company.com" {
		t.Errorf("ContactDetails.Email mismatch: got %s", account.ContactDetails.Email)
	}
	if account.ContactDetails.Phone != "+6588880000" {
		t.Errorf("ContactDetails.Phone mismatch: got %s", account.ContactDetails.Phone)
	}
}

func TestParseAccountData_BusinessDetails(t *testing.T) {
	var event Event
	json.Unmarshal([]byte(accountCreateWebhookJSON), &event)

	account, err := event.ParseAccountData()
	if err != nil {
		t.Fatalf("Failed to parse account data: %v", err)
	}

	if account.BusinessDetails == nil {
		t.Fatal("BusinessDetails should not be nil")
	}

	bd := account.BusinessDetails

	if bd.LegalEntityName != "ソフトバンクインターナショナル株式会社" {
		t.Errorf("LegalEntityName mismatch: got %s", bd.LegalEntityName)
	}
	if bd.LegalEntityNameEnglish != "UQPAY PTE LTD." {
		t.Errorf("LegalEntityNameEnglish mismatch: got %s", bd.LegalEntityNameEnglish)
	}
	if bd.IncorporationDate != "2013-04-15" {
		t.Errorf("IncorporationDate mismatch: got %s", bd.IncorporationDate)
	}
	if bd.RegistrationNumber != "201901754K" {
		t.Errorf("RegistrationNumber mismatch: got %s", bd.RegistrationNumber)
	}
	if bd.BusinessStructure != "SOLE_PROPRIETOR" {
		t.Errorf("BusinessStructure mismatch: got %s", bd.BusinessStructure)
	}
	if bd.MerchantCategoryCode != "5733" {
		t.Errorf("MerchantCategoryCode mismatch: got %s", bd.MerchantCategoryCode)
	}
	if bd.WebsiteURL != "https://yourcompany.com" {
		t.Errorf("WebsiteURL mismatch: got %s", bd.WebsiteURL)
	}

	// Verify monthly estimated revenue
	if bd.MonthlyEstimatedRevenue == nil {
		t.Fatal("MonthlyEstimatedRevenue should not be nil")
	}
	if bd.MonthlyEstimatedRevenue.Currency != "SGD" {
		t.Errorf("MonthlyEstimatedRevenue.Currency mismatch: got %s", bd.MonthlyEstimatedRevenue.Currency)
	}

	// Verify account purposes
	if len(bd.AccountPurpose) != 2 {
		t.Errorf("AccountPurpose length mismatch: got %d, want 2", len(bd.AccountPurpose))
	}

	// Verify identifier
	if bd.Identifier == nil {
		t.Fatal("Identifier should not be nil")
	}
	if bd.Identifier.Type != "VAT" {
		t.Errorf("Identifier.Type mismatch: got %s", bd.Identifier.Type)
	}
}

func TestParseAccountData_Representatives(t *testing.T) {
	var event Event
	json.Unmarshal([]byte(accountCreateWebhookJSON), &event)

	account, err := event.ParseAccountData()
	if err != nil {
		t.Fatalf("Failed to parse account data: %v", err)
	}

	if len(account.Representatives) != 1 {
		t.Fatalf("Representatives length mismatch: got %d, want 1", len(account.Representatives))
	}

	rep := account.Representatives[0]

	if rep.Roles != "DIRECTOR" {
		t.Errorf("Roles mismatch: got %s", rep.Roles)
	}
	if rep.FirstName != "Mock" {
		t.Errorf("FirstName mismatch: got %s", rep.FirstName)
	}
	if rep.LastName != "Toy" {
		t.Errorf("LastName mismatch: got %s", rep.LastName)
	}
	if rep.Nationality != "SG" {
		t.Errorf("Nationality mismatch: got %s", rep.Nationality)
	}
	if rep.DateOfBirth != "2024-01-28" {
		t.Errorf("DateOfBirth mismatch: got %s", rep.DateOfBirth)
	}
	if rep.SharePercentage != "20.01" {
		t.Errorf("SharePercentage mismatch: got %s", rep.SharePercentage)
	}
	if rep.AsApplicant != false {
		t.Errorf("AsApplicant mismatch: got %v", rep.AsApplicant)
	}

	// Verify identification
	if rep.Identification == nil {
		t.Fatal("Identification should not be nil")
	}
	if rep.Identification.Type != "PASSPORT" {
		t.Errorf("Identification.Type mismatch: got %s", rep.Identification.Type)
	}
	if rep.Identification.IDNumber != "27738277K" {
		t.Errorf("Identification.IDNumber mismatch: got %s", rep.Identification.IDNumber)
	}

	// Verify residential address
	if rep.ResidentialAddress == nil {
		t.Fatal("ResidentialAddress should not be nil")
	}
	if rep.ResidentialAddress.City != "Singapore" {
		t.Errorf("ResidentialAddress.City mismatch: got %s", rep.ResidentialAddress.City)
	}
}

func TestParseAccountData_Addresses(t *testing.T) {
	var event Event
	json.Unmarshal([]byte(accountCreateWebhookJSON), &event)

	account, err := event.ParseAccountData()
	if err != nil {
		t.Fatalf("Failed to parse account data: %v", err)
	}

	// Registration address
	if account.RegistrationAddress == nil {
		t.Fatal("RegistrationAddress should not be nil")
	}
	if account.RegistrationAddress.City != "San Francisco" {
		t.Errorf("RegistrationAddress.City mismatch: got %s", account.RegistrationAddress.City)
	}
	if account.RegistrationAddress.State != "CA" {
		t.Errorf("RegistrationAddress.State mismatch: got %s", account.RegistrationAddress.State)
	}
	if account.RegistrationAddress.PostalCode != "94103" {
		t.Errorf("RegistrationAddress.PostalCode mismatch: got %s", account.RegistrationAddress.PostalCode)
	}

	// Business addresses
	if len(account.BusinessAddress) != 1 {
		t.Fatalf("BusinessAddress length mismatch: got %d, want 1", len(account.BusinessAddress))
	}
	if account.BusinessAddress[0].Country != "SG" {
		t.Errorf("BusinessAddress.Country mismatch: got %s", account.BusinessAddress[0].Country)
	}
}

// ============================================================================
// Account Update Tests
// ============================================================================

func TestParseAccountData_Update(t *testing.T) {
	var event Event
	err := json.Unmarshal([]byte(accountUpdateWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse event: %v", err)
	}

	if event.EventType != EventTypeAccountUpdate {
		t.Errorf("EventType mismatch: got %s, want %s", event.EventType, EventTypeAccountUpdate)
	}
	if !event.IsAccountUpdateEvent() {
		t.Error("IsAccountUpdateEvent should return true")
	}

	account, err := event.ParseAccountData()
	if err != nil {
		t.Fatalf("Failed to parse account data: %v", err)
	}

	if account.Status != "ACTIVE" {
		t.Errorf("Status mismatch: got %s, want ACTIVE", account.Status)
	}
	if account.Email != "updated@company.com" {
		t.Errorf("Email mismatch: got %s", account.Email)
	}
}

// ============================================================================
// Edge Cases and Error Handling Tests
// ============================================================================

func TestParseAccountData_MinimalPayload(t *testing.T) {
	var event Event
	err := json.Unmarshal([]byte(minimalAccountWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse minimal event: %v", err)
	}

	account, err := event.ParseAccountData()
	if err != nil {
		t.Fatalf("Failed to parse minimal account data: %v", err)
	}

	if account.AccountID != "minimal-account-id" {
		t.Errorf("AccountID mismatch: got %s", account.AccountID)
	}
	if account.EntityType != "INDIVIDUAL" {
		t.Errorf("EntityType mismatch: got %s", account.EntityType)
	}

	// Optional fields should be nil/empty
	if account.BusinessDetails != nil {
		t.Error("BusinessDetails should be nil for minimal payload")
	}
	if account.ContactDetails != nil {
		t.Error("ContactDetails should be nil for minimal payload")
	}
}

func TestParseAccountData_WrongEventType(t *testing.T) {
	wrongTypeJSON := `{
		"version": "V1.6.0",
		"event_name": "PAYMENT",
		"event_type": "payment.completed",
		"event_id": "test-id",
		"source_id": "test-source",
		"data": {}
	}`

	var event Event
	err := json.Unmarshal([]byte(wrongTypeJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse event: %v", err)
	}

	_, err = event.ParseAccountData()
	if err == nil {
		t.Error("ParseAccountData should fail for non-account event type")
	}
}

func TestParseAccountData_MalformedData(t *testing.T) {
	malformedJSON := `{
		"version": "V1.6.0",
		"event_name": "ONBOARDING",
		"event_type": "onboarding.account.create",
		"event_id": "test-id",
		"source_id": "test-source",
		"data": "not an object"
	}`

	var event Event
	err := json.Unmarshal([]byte(malformedJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse event envelope: %v", err)
	}

	_, err = event.ParseAccountData()
	if err == nil {
		t.Error("ParseAccountData should fail for malformed data")
	}
}

// ============================================================================
// Onboarding Status Constants Tests
// ============================================================================

func TestOnboardingStatusConstants(t *testing.T) {
	// Account status
	if AccountStatusProcessing != "PROCESSING" {
		t.Errorf("AccountStatusProcessing mismatch")
	}
	if AccountStatusActive != "ACTIVE" {
		t.Errorf("AccountStatusActive mismatch")
	}

	// Verification status
	if VerificationStatusPending != "PENDING" {
		t.Errorf("VerificationStatusPending mismatch")
	}
	if VerificationStatusApproved != "APPROVED" {
		t.Errorf("VerificationStatusApproved mismatch")
	}

	// Entity type
	if EntityTypeCompany != "COMPANY" {
		t.Errorf("EntityTypeCompany mismatch")
	}
	if EntityTypeIndividual != "INDIVIDUAL" {
		t.Errorf("EntityTypeIndividual mismatch")
	}
}

// ============================================================================
// Table-Driven Tests
// ============================================================================

func TestParseAccountData_TableDriven(t *testing.T) {
	testCases := []struct {
		name           string
		json           string
		expectError    bool
		expectedStatus string
		expectedType   string
	}{
		{
			name:           "AccountCreate_Full",
			json:           accountCreateWebhookJSON,
			expectError:    false,
			expectedStatus: "PROCESSING",
			expectedType:   "COMPANY",
		},
		{
			name:           "AccountUpdate",
			json:           accountUpdateWebhookJSON,
			expectError:    false,
			expectedStatus: "ACTIVE",
			expectedType:   "COMPANY",
		},
		{
			name:           "Minimal",
			json:           minimalAccountWebhookJSON,
			expectError:    false,
			expectedStatus: "PROCESSING",
			expectedType:   "INDIVIDUAL",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var event Event
			err := json.Unmarshal([]byte(tc.json), &event)
			if err != nil {
				t.Fatalf("Failed to parse event: %v", err)
			}

			account, err := event.ParseAccountData()
			if tc.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if account.Status != tc.expectedStatus {
				t.Errorf("Status mismatch: got %s, want %s", account.Status, tc.expectedStatus)
			}
			if account.EntityType != tc.expectedType {
				t.Errorf("EntityType mismatch: got %s, want %s", account.EntityType, tc.expectedType)
			}
		})
	}
}
