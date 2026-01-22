package test

import (
	"encoding/json"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/webhook"
)

// ============================================================================
// Test Fixtures - Exact JSON from UQPAY Documentation
// ============================================================================

// accountCreateWebhookJSON is the exact payload from docs.uqpay.com
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

// Minimal valid webhook for edge case testing
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

// Real webhook JSON from production environment with full fields
const realAccountCreateWebhookJSON = `{
	"version":"V1.6.0",
	"event_name":"ONBOARDING",
	"event_type":"onboarding.account.create",
	"event_id":"77edd072-b096-4fad-8df2-017d5bd9bb65",
	"source_id":"7cbf9ed8-cbbb-466a-8c5a-aa30cf6cf431",
	"data":{
		"account_id":"7cbf9ed8-cbbb-466a-8c5a-aa30cf6cf431",
		"account_name":"HelloTest1",
		"business_address":[{
			"city":"Clementi",
			"country":"AU",
			"line1":"IMM Outlet Mall, 2 Jurong East Street 21, Singapore 609601",
			"line2":"High Street, Houston, Texas, 19683",
			"postal_code":"749711",
			"state":"North Region"
		}],
		"business_details":{
			"account_purpose":["bill_payment"],
			"banking_countries":["AF"],
			"banking_currencies":["AED"],
			"business_address":[{
				"city":"Clementi",
				"country":"AU",
				"line1":"IMM Outlet Mall, 2 Jurong East Street 21, Singapore 609601",
				"line2":"High Street, Houston, Texas, 19683",
				"postal_code":"749711",
				"state":"North Region"
			}],
			"business_structure":"sole_proprietor",
			"country":"FR",
			"estimated_worker_count":"BS001",
			"identification_expiry_date":"--",
			"identification_issue_date":"",
			"identifier":{
				"number":"GB004760879",
				"type":"TAX"
			},
			"incorporation_date":"2001-01-01",
			"industry_code":"ICCV3_1106XX",
			"legal_entity_name":"自动化测试09:40:54",
			"legal_entity_name_english":"AutoTest Jhin",
			"mcc":"4411",
			"merchant_category_code":"4411",
			"monthly_estimated_revenue":{
				"amount":"TM001",
				"currency":"AED"
			},
			"product_description":"AutoTest business_scope",
			"registration_address":{
				"city":"Marine Parade+Modify",
				"country":"FR",
				"line1":"East Coast Park, 1220 East Coast Parkway, Singapore 468960+Modify",
				"line2":"Main Street, Dallas, Texas, 63028+Modify",
				"postal_code":"351008",
				"state":"North Region"
			},
			"registration_number":"FKIBGB2V",
			"website_url":"http://9bauuw0.com"
		},
		"contact_details":{
			"email":"gqvsbm42427584812@uqpay.mil",
			"phone":"+86 15002588508"
		},
		"country":"FR",
		"direct_id":"65087660-8d3d-428e-bd2e-9e56219c1512",
		"email":"gqvsbm42427584812@uqpay.mil",
		"entity_type":"COMPANY",
		"idv_status":"",
		"registration_address":{
			"city":"Marine Parade+Modify",
			"country":"FR",
			"line1":"East Coast Park, 1220 East Coast Parkway, Singapore 468960+Modify",
			"line2":"Main Street, Dallas, Texas, 63028+Modify",
			"postal_code":"351008",
			"state":"North Region"
		},
		"representatives":[{
			"area_code":"+86",
			"citizenship_status":0,
			"date_of_birth":"2002-01-01",
			"email":"kvqenzkyi41615768541@uqpay.biz",
			"first_name":"Jayce",
			"identification":{
				"documents":{
					"back":"account/CB2699475968/20251111/6f29b7627163403e9176d423a3e5b316.png",
					"front":"",
					"type":"1002"
				},
				"id_number":"223675196107079046",
				"identification_expiry_date":"",
				"identification_issue_date":"",
				"type":"NATIONAL_ID"
			},
			"idv_id":"",
			"idv_status":"verified",
			"is_applicant":true,
			"last_name":"Camille",
			"local_name":"Camille",
			"nationality":"GB",
			"phone_number":"17438084425",
			"representative_id":"2e2f97f1-de60-47fb-8d2d-4fd9017137f5",
			"residential_address":{
				"city":"Jurong",
				"country":"CH",
				"line1":"Singapore Art Museum, 61 Stamford Road, Singapore 178892",
				"line2":"Elm Street, Chicago, Ohio, 99683",
				"postal_code":"986674",
				"state":"North-East Region"
			},
			"roles":"UBO",
			"share_percentage":"66.71",
			"tax_number":"GB036734195"
		},{
			"area_code":"+86",
			"citizenship_status":0,
			"date_of_birth":"2002-01-01",
			"email":"lexwylxicj84231928811@uqpay.org",
			"first_name":"Graves",
			"identification":{
				"documents":{
					"back":"account/CB2699475968/20251111/ca08d02d65b54474acf1c61f648ce8ee.png",
					"front":"",
					"type":"1002"
				},
				"id_number":"530650197312200388",
				"identification_expiry_date":"",
				"identification_issue_date":"",
				"type":"NATIONAL_ID"
			},
			"idv_id":"",
			"idv_status":"verified",
			"is_applicant":false,
			"last_name":"Anivia",
			"local_name":"Anivia",
			"nationality":"CH",
			"phone_number":"13190897002",
			"representative_id":"ad39fe4e-1cd7-4e55-a2bb-412136d3908f",
			"residential_address":{
				"city":"Yishun",
				"country":"FR",
				"line1":"Jewel Changi Airport, 78 Airport Boulevard, Singapore 819666",
				"line2":"Broadway, San Diego, Michigan, 32225",
				"postal_code":"515460",
				"state":"West Region"
			},
			"roles":"UBO",
			"share_percentage":"66.71",
			"tax_number":"GB642122680"
		},{
			"area_code":"+86",
			"citizenship_status":0,
			"date_of_birth":"2002-01-01",
			"email":"kmnwmvc42085161233@uqpay.edu",
			"first_name":"Ekko",
			"identification":{
				"documents":{
					"back":"account/CB2699475968/20251111/788acf401e344402a9e554fd71d8174d.png",
					"front":"",
					"type":"1002"
				},
				"id_number":"129867199510210219",
				"identification_expiry_date":"",
				"identification_issue_date":"",
				"type":"NATIONAL_ID"
			},
			"idv_id":"",
			"idv_status":"verified",
			"is_applicant":false,
			"last_name":"Cassiopeia",
			"local_name":"Cassiopeia",
			"nationality":"JP",
			"phone_number":"15472875484",
			"representative_id":"f1bc84d6-611d-4d4f-9e5b-2b8870a3da5c",
			"residential_address":{
				"city":"Yishun",
				"country":"CN",
				"line1":"The Battle Box, 2 Cox Terrace, Singapore 179622",
				"line2":"Park Road, San Diego, Ohio, 08919",
				"postal_code":"717030",
				"state":"North Region"
			},
			"roles":"UBO",
			"share_percentage":"66.71",
			"tax_number":"GB981604705"
		}],
		"review_reason":"",
		"short_reference_id":"CB9618267136",
		"source":"web",
		"status":"ACTIVE",
		"verification_status":"APPROVED"
	}
}`

// ============================================================================
// Event Parsing Tests
// ============================================================================

func TestParseEvent(t *testing.T) {
	var event webhook.Event
	err := json.Unmarshal([]byte(accountCreateWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse event: %v", err)
	}

	// Verify envelope fields
	if event.Version != "V1.6.0" {
		t.Errorf("Version mismatch: got %s, want V1.6.0", event.Version)
	}
	if event.EventName != "ONBOARDING" {
		t.Errorf("EventName mismatch: got %s, want ONBOARDING", event.EventName)
	}
	if event.EventType != "onboarding.account.create" {
		t.Errorf("EventType mismatch: got %s, want onboarding.account.create", event.EventType)
	}
	if event.EventID != "8a78af1e-de83-43a5-b177-ecbc6a8a9fc6" {
		t.Errorf("EventID mismatch: got %s", event.EventID)
	}
	if event.SourceID != "f5bb6498-552e-40a5-b14b-616aa04ac1c1" {
		t.Errorf("SourceID mismatch: got %s", event.SourceID)
	}

	t.Logf("Event parsed successfully")
	t.Logf("   Version: %s", event.Version)
	t.Logf("   EventName: %s", event.EventName)
	t.Logf("   EventType: %s", event.EventType)
	t.Logf("   EventID: %s", event.EventID)
	t.Logf("   SourceID: %s", event.SourceID)
}

func TestEventHelperMethods(t *testing.T) {
	var event webhook.Event
	err := json.Unmarshal([]byte(accountCreateWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse event: %v", err)
	}

	// Test helper methods
	if !event.IsOnboardingEvent() {
		t.Error("IsOnboardingEvent should return true")
	}
	if !event.IsAccountCreateEvent() {
		t.Error("IsAccountCreateEvent should return true")
	}
	if event.IsAccountUpdateEvent() {
		t.Error("IsAccountUpdateEvent should return false for create event")
	}

	t.Logf("Helper methods verified")
	t.Logf("   IsOnboardingEvent: %v", event.IsOnboardingEvent())
	t.Logf("   IsAccountCreateEvent: %v", event.IsAccountCreateEvent())
	t.Logf("   IsAccountUpdateEvent: %v", event.IsAccountUpdateEvent())
}

// ============================================================================
// Account Data Parsing Tests
// ============================================================================

func TestParseAccountData_Create(t *testing.T) {
	var event webhook.Event
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

	t.Logf("Account data parsed successfully")
	t.Logf("   AccountID: %s", account.AccountID)
	t.Logf("   AccountName: %s", account.AccountName)
	t.Logf("   Email: %s", account.Email)
	t.Logf("   Country: %s", account.Country)
	t.Logf("   Status: %s", account.Status)
	t.Logf("   EntityType: %s", account.EntityType)
}

func TestParseAccountData_ContactDetails(t *testing.T) {
	var event webhook.Event
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

	t.Logf("Contact details verified")
	t.Logf("   Email: %s", account.ContactDetails.Email)
	t.Logf("   Phone: %s", account.ContactDetails.Phone)
}

func TestParseAccountData_BusinessDetails(t *testing.T) {
	var event webhook.Event
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

	t.Logf("Business details verified")
	t.Logf("   LegalEntityName: %s", bd.LegalEntityName)
	t.Logf("   LegalEntityNameEnglish: %s", bd.LegalEntityNameEnglish)
	t.Logf("   BusinessStructure: %s", bd.BusinessStructure)
	t.Logf("   MCC: %s", bd.MerchantCategoryCode)
	t.Logf("   AccountPurposes: %v", bd.AccountPurpose)
}

func TestParseAccountData_Representatives(t *testing.T) {
	var event webhook.Event
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

	t.Logf("Representative verified")
	t.Logf("   Name: %s %s", rep.FirstName, rep.LastName)
	t.Logf("   Role: %s", rep.Roles)
	t.Logf("   Nationality: %s", rep.Nationality)
	t.Logf("   SharePercentage: %s%%", rep.SharePercentage)
}

func TestParseAccountData_Addresses(t *testing.T) {
	var event webhook.Event
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

	t.Logf("Addresses verified")
	t.Logf("   Registration: %s, %s %s",
		account.RegistrationAddress.City,
		account.RegistrationAddress.State,
		account.RegistrationAddress.PostalCode)
	t.Logf("   Business: %s, %s",
		account.BusinessAddress[0].City,
		account.BusinessAddress[0].Country)
}

// ============================================================================
// Account Update Event Tests
// ============================================================================

func TestParseAccountData_Update(t *testing.T) {
	var event webhook.Event
	err := json.Unmarshal([]byte(accountUpdateWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse event: %v", err)
	}

	if event.EventType != webhook.EventTypeAccountUpdate {
		t.Errorf("EventType mismatch: got %s, want %s", event.EventType, webhook.EventTypeAccountUpdate)
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

	t.Logf("Account update event parsed")
	t.Logf("   Status: %s", account.Status)
	t.Logf("   Email: %s", account.Email)
}

// ============================================================================
// Edge Cases and Error Handling Tests
// ============================================================================

func TestParseAccountData_MinimalPayload(t *testing.T) {
	var event webhook.Event
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

	t.Logf("Minimal payload handled correctly")
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

	var event webhook.Event
	err := json.Unmarshal([]byte(wrongTypeJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse event: %v", err)
	}

	_, err = event.ParseAccountData()
	if err == nil {
		t.Error("ParseAccountData should fail for non-account event type")
	}

	t.Logf("Wrong event type correctly rejected: %v", err)
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

	var event webhook.Event
	err := json.Unmarshal([]byte(malformedJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse event envelope: %v", err)
	}

	_, err = event.ParseAccountData()
	if err == nil {
		t.Error("ParseAccountData should fail for malformed data")
	}

	t.Logf("Malformed data correctly rejected: %v", err)
}

func TestParseAccountData_InvalidJSON(t *testing.T) {
	invalidJSON := `{not valid json}`

	var event webhook.Event
	err := json.Unmarshal([]byte(invalidJSON), &event)
	if err == nil {
		t.Error("Should fail to parse invalid JSON")
	}

	t.Logf("Invalid JSON correctly rejected: %v", err)
}

// ============================================================================
// Constants Verification Tests
// ============================================================================

func TestEventConstants(t *testing.T) {
	// Verify event name constants match expected values
	if webhook.EventNameOnboarding != "ONBOARDING" {
		t.Errorf("EventNameOnboarding mismatch: got %s", webhook.EventNameOnboarding)
	}

	// Verify event type constants
	if webhook.EventTypeAccountCreate != "onboarding.account.create" {
		t.Errorf("EventTypeAccountCreate mismatch: got %s", webhook.EventTypeAccountCreate)
	}
	if webhook.EventTypeAccountUpdate != "onboarding.account.update" {
		t.Errorf("EventTypeAccountUpdate mismatch: got %s", webhook.EventTypeAccountUpdate)
	}

	t.Logf("Event constants verified")
	t.Logf("   EventNameOnboarding: %s", webhook.EventNameOnboarding)
	t.Logf("   EventTypeAccountCreate: %s", webhook.EventTypeAccountCreate)
	t.Logf("   EventTypeAccountUpdate: %s", webhook.EventTypeAccountUpdate)
}

func TestStatusConstants(t *testing.T) {
	// Account status
	if webhook.AccountStatusProcessing != "PROCESSING" {
		t.Errorf("AccountStatusProcessing mismatch")
	}
	if webhook.AccountStatusActive != "ACTIVE" {
		t.Errorf("AccountStatusActive mismatch")
	}

	// Verification status
	if webhook.VerificationStatusPending != "PENDING" {
		t.Errorf("VerificationStatusPending mismatch")
	}
	if webhook.VerificationStatusApproved != "APPROVED" {
		t.Errorf("VerificationStatusApproved mismatch")
	}

	// Entity type
	if webhook.EntityTypeCompany != "COMPANY" {
		t.Errorf("EntityTypeCompany mismatch")
	}
	if webhook.EntityTypeIndividual != "INDIVIDUAL" {
		t.Errorf("EntityTypeIndividual mismatch")
	}

	t.Logf("Status constants verified")
}

// ============================================================================
// Table-Driven Tests for Multiple Scenarios
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
			var event webhook.Event
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

			t.Logf("Verified: Status=%s, EntityType=%s", account.Status, account.EntityType)
		})
	}
}

// ============================================================================
// Real Production Webhook Tests
// ============================================================================

func TestParseRealAccountCreateWebhook(t *testing.T) {
	var event webhook.Event
	err := json.Unmarshal([]byte(realAccountCreateWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse real webhook event: %v", err)
	}

	// Verify envelope fields
	if event.Version != "V1.6.0" {
		t.Errorf("Version mismatch: got %s, want V1.6.0", event.Version)
	}
	if event.EventName != "ONBOARDING" {
		t.Errorf("EventName mismatch: got %s, want ONBOARDING", event.EventName)
	}
	if event.EventType != "onboarding.account.create" {
		t.Errorf("EventType mismatch: got %s", event.EventType)
	}
	if event.EventID != "77edd072-b096-4fad-8df2-017d5bd9bb65" {
		t.Errorf("EventID mismatch: got %s", event.EventID)
	}
	if event.SourceID != "7cbf9ed8-cbbb-466a-8c5a-aa30cf6cf431" {
		t.Errorf("SourceID mismatch: got %s", event.SourceID)
	}

	t.Logf("Real webhook event parsed successfully")
}

func TestParseRealAccountData(t *testing.T) {
	var event webhook.Event
	json.Unmarshal([]byte(realAccountCreateWebhookJSON), &event)

	account, err := event.ParseAccountData()
	if err != nil {
		t.Fatalf("Failed to parse real account data: %v", err)
	}

	// Verify core fields
	if account.AccountID != "7cbf9ed8-cbbb-466a-8c5a-aa30cf6cf431" {
		t.Errorf("AccountID mismatch: got %s", account.AccountID)
	}
	if account.AccountName != "HelloTest1" {
		t.Errorf("AccountName mismatch: got %s", account.AccountName)
	}
	if account.Email != "gqvsbm42427584812@uqpay.mil" {
		t.Errorf("Email mismatch: got %s", account.Email)
	}
	if account.Country != "FR" {
		t.Errorf("Country mismatch: got %s", account.Country)
	}
	if account.Status != "ACTIVE" {
		t.Errorf("Status mismatch: got %s", account.Status)
	}
	if account.EntityType != "COMPANY" {
		t.Errorf("EntityType mismatch: got %s", account.EntityType)
	}
	if account.VerificationStatus != "APPROVED" {
		t.Errorf("VerificationStatus mismatch: got %s", account.VerificationStatus)
	}
	if account.Source != "web" {
		t.Errorf("Source mismatch: got %s", account.Source)
	}
	if account.ShortReferenceID != "CB9618267136" {
		t.Errorf("ShortReferenceID mismatch: got %s", account.ShortReferenceID)
	}
	if account.DirectID != "65087660-8d3d-428e-bd2e-9e56219c1512" {
		t.Errorf("DirectID mismatch: got %s", account.DirectID)
	}

	t.Logf("Real account data parsed successfully")
	t.Logf("   AccountID: %s", account.AccountID)
	t.Logf("   AccountName: %s", account.AccountName)
	t.Logf("   Status: %s", account.Status)
	t.Logf("   Source: %s", account.Source)
}

func TestParseRealBusinessDetails(t *testing.T) {
	var event webhook.Event
	json.Unmarshal([]byte(realAccountCreateWebhookJSON), &event)

	account, err := event.ParseAccountData()
	if err != nil {
		t.Fatalf("Failed to parse account data: %v", err)
	}

	if account.BusinessDetails == nil {
		t.Fatal("BusinessDetails should not be nil")
	}

	bd := account.BusinessDetails

	// Verify new fields
	if bd.Country != "FR" {
		t.Errorf("BusinessDetails.Country mismatch: got %s", bd.Country)
	}
	if bd.MCC != "4411" {
		t.Errorf("BusinessDetails.MCC mismatch: got %s", bd.MCC)
	}
	if bd.MerchantCategoryCode != "4411" {
		t.Errorf("BusinessDetails.MerchantCategoryCode mismatch: got %s", bd.MerchantCategoryCode)
	}
	if bd.IndustryCode != "ICCV3_1106XX" {
		t.Errorf("BusinessDetails.IndustryCode mismatch: got %s", bd.IndustryCode)
	}
	if bd.BusinessStructure != "sole_proprietor" {
		t.Errorf("BusinessDetails.BusinessStructure mismatch: got %s", bd.BusinessStructure)
	}
	if bd.LegalEntityName != "自动化测试09:40:54" {
		t.Errorf("BusinessDetails.LegalEntityName mismatch: got %s", bd.LegalEntityName)
	}
	if bd.LegalEntityNameEnglish != "AutoTest Jhin" {
		t.Errorf("BusinessDetails.LegalEntityNameEnglish mismatch: got %s", bd.LegalEntityNameEnglish)
	}
	if bd.WebsiteURL != "http://9bauuw0.com" {
		t.Errorf("BusinessDetails.WebsiteURL mismatch: got %s", bd.WebsiteURL)
	}

	// Verify banking fields
	if len(bd.BankingCountries) != 1 || bd.BankingCountries[0] != "AF" {
		t.Errorf("BusinessDetails.BankingCountries mismatch: got %v", bd.BankingCountries)
	}
	if len(bd.BankingCurrencies) != 1 || bd.BankingCurrencies[0] != "AED" {
		t.Errorf("BusinessDetails.BankingCurrencies mismatch: got %v", bd.BankingCurrencies)
	}

	// Verify account purpose
	if len(bd.AccountPurpose) != 1 || bd.AccountPurpose[0] != "bill_payment" {
		t.Errorf("BusinessDetails.AccountPurpose mismatch: got %v", bd.AccountPurpose)
	}

	// Verify registration address within business details
	if bd.RegistrationAddress == nil {
		t.Error("BusinessDetails.RegistrationAddress should not be nil")
	} else if bd.RegistrationAddress.City != "Marine Parade+Modify" {
		t.Errorf("BusinessDetails.RegistrationAddress.City mismatch: got %s", bd.RegistrationAddress.City)
	}

	// Verify business address within business details
	if len(bd.BusinessAddress) != 1 {
		t.Errorf("BusinessDetails.BusinessAddress length mismatch: got %d", len(bd.BusinessAddress))
	} else if bd.BusinessAddress[0].City != "Clementi" {
		t.Errorf("BusinessDetails.BusinessAddress[0].City mismatch: got %s", bd.BusinessAddress[0].City)
	}

	t.Logf("Real business details parsed successfully")
	t.Logf("   Country: %s", bd.Country)
	t.Logf("   MCC: %s", bd.MCC)
	t.Logf("   IndustryCode: %s", bd.IndustryCode)
	t.Logf("   BankingCountries: %v", bd.BankingCountries)
	t.Logf("   BankingCurrencies: %v", bd.BankingCurrencies)
}

func TestParseRealRepresentatives(t *testing.T) {
	var event webhook.Event
	json.Unmarshal([]byte(realAccountCreateWebhookJSON), &event)

	account, err := event.ParseAccountData()
	if err != nil {
		t.Fatalf("Failed to parse account data: %v", err)
	}

	if len(account.Representatives) != 3 {
		t.Fatalf("Representatives length mismatch: got %d, want 3", len(account.Representatives))
	}

	// Test first representative (the applicant)
	rep := account.Representatives[0]

	if rep.RepresentativeID != "2e2f97f1-de60-47fb-8d2d-4fd9017137f5" {
		t.Errorf("RepresentativeID mismatch: got %s", rep.RepresentativeID)
	}
	if rep.FirstName != "Jayce" {
		t.Errorf("FirstName mismatch: got %s", rep.FirstName)
	}
	if rep.LastName != "Camille" {
		t.Errorf("LastName mismatch: got %s", rep.LastName)
	}
	if rep.LocalName != "Camille" {
		t.Errorf("LocalName mismatch: got %s", rep.LocalName)
	}
	if rep.Roles != "UBO" {
		t.Errorf("Roles mismatch: got %s", rep.Roles)
	}
	if rep.AreaCode != "+86" {
		t.Errorf("AreaCode mismatch: got %s", rep.AreaCode)
	}
	if rep.PhoneNumber != "17438084425" {
		t.Errorf("PhoneNumber mismatch: got %s", rep.PhoneNumber)
	}
	if rep.Email != "kvqenzkyi41615768541@uqpay.biz" {
		t.Errorf("Email mismatch: got %s", rep.Email)
	}
	if rep.TaxNumber != "GB036734195" {
		t.Errorf("TaxNumber mismatch: got %s", rep.TaxNumber)
	}
	if rep.SharePercentage != "66.71" {
		t.Errorf("SharePercentage mismatch: got %s", rep.SharePercentage)
	}
	if !rep.IsApplicant {
		t.Error("IsApplicant should be true for first representative")
	}
	if rep.IDVStatus != "verified" {
		t.Errorf("IDVStatus mismatch: got %s", rep.IDVStatus)
	}
	if rep.CitizenshipStatus != 0 {
		t.Errorf("CitizenshipStatus mismatch: got %d", rep.CitizenshipStatus)
	}

	// Verify identification with documents
	if rep.Identification == nil {
		t.Fatal("Identification should not be nil")
	}
	if rep.Identification.Type != "NATIONAL_ID" {
		t.Errorf("Identification.Type mismatch: got %s", rep.Identification.Type)
	}
	if rep.Identification.IDNumber != "223675196107079046" {
		t.Errorf("Identification.IDNumber mismatch: got %s", rep.Identification.IDNumber)
	}
	if rep.Identification.Documents == nil {
		t.Fatal("Identification.Documents should not be nil")
	}
	if rep.Identification.Documents.Back != "account/CB2699475968/20251111/6f29b7627163403e9176d423a3e5b316.png" {
		t.Errorf("Identification.Documents.Back mismatch: got %s", rep.Identification.Documents.Back)
	}
	if rep.Identification.Documents.Type != "1002" {
		t.Errorf("Identification.Documents.Type mismatch: got %s", rep.Identification.Documents.Type)
	}

	// Verify second and third representatives are not applicants
	if account.Representatives[1].IsApplicant {
		t.Error("Second representative should not be applicant")
	}
	if account.Representatives[2].IsApplicant {
		t.Error("Third representative should not be applicant")
	}

	t.Logf("Real representatives parsed successfully")
	t.Logf("   Representative 1: %s %s (Applicant: %v)", rep.FirstName, rep.LastName, rep.IsApplicant)
	t.Logf("   Representative 2: %s %s", account.Representatives[1].FirstName, account.Representatives[1].LastName)
	t.Logf("   Representative 3: %s %s", account.Representatives[2].FirstName, account.Representatives[2].LastName)
}

// ============================================================================
// Payment Intent Webhook Tests
// ============================================================================

// Real payment intent webhook JSON from production
const paymentIntentCreatedWebhookJSON = `{
	"version":"V1.6.0",
	"event_name":"ACQUIRING",
	"event_type":"acquiring.payment_intent.created",
	"event_id":"5008b4da-a5e0-4a07-86f1-42cf6931afed",
	"source_id":"PI2013833849980588032",
	"data":{
		"amount":"101",
		"cancel_time":null,
		"cancellation_reason":"",
		"complete_time":null,
		"create_time":"2026-01-21T12:39:39.716846826+08:00",
		"currency":"USD",
		"description":"Test payment intent",
		"intent_status":"REQUIRES_PAYMENT_METHOD",
		"merchant_order_id":"test-order-002",
		"metadata":{"test":"true"},
		"payment_intent_id":"PI2013833849980588032",
		"payment_method":null
	}
}`

func TestParsePaymentIntentCreatedWebhook(t *testing.T) {
	var event webhook.Event
	err := json.Unmarshal([]byte(paymentIntentCreatedWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse payment intent webhook event: %v", err)
	}

	// Verify envelope fields
	if event.Version != "V1.6.0" {
		t.Errorf("Version mismatch: got %s, want V1.6.0", event.Version)
	}
	if event.EventName != "ACQUIRING" {
		t.Errorf("EventName mismatch: got %s, want ACQUIRING", event.EventName)
	}
	if event.EventType != "acquiring.payment_intent.created" {
		t.Errorf("EventType mismatch: got %s", event.EventType)
	}
	if event.EventID != "5008b4da-a5e0-4a07-86f1-42cf6931afed" {
		t.Errorf("EventID mismatch: got %s", event.EventID)
	}
	if event.SourceID != "PI2013833849980588032" {
		t.Errorf("SourceID mismatch: got %s", event.SourceID)
	}

	// Verify helper methods
	if !event.IsAcquiringEvent() {
		t.Error("IsAcquiringEvent should return true")
	}
	if !event.IsPaymentIntentEvent() {
		t.Error("IsPaymentIntentEvent should return true")
	}
	if event.IsOnboardingEvent() {
		t.Error("IsOnboardingEvent should return false")
	}

	t.Logf("Payment intent webhook event parsed successfully")
	t.Logf("   EventName: %s", event.EventName)
	t.Logf("   EventType: %s", event.EventType)
	t.Logf("   SourceID: %s", event.SourceID)
}

func TestParsePaymentIntentData(t *testing.T) {
	var event webhook.Event
	json.Unmarshal([]byte(paymentIntentCreatedWebhookJSON), &event)

	paymentIntent, err := event.ParsePaymentIntentData()
	if err != nil {
		t.Fatalf("Failed to parse payment intent data: %v", err)
	}

	// Verify core fields
	if paymentIntent.PaymentIntentID != "PI2013833849980588032" {
		t.Errorf("PaymentIntentID mismatch: got %s", paymentIntent.PaymentIntentID)
	}
	if paymentIntent.Amount != "101" {
		t.Errorf("Amount mismatch: got %s", paymentIntent.Amount)
	}
	if paymentIntent.Currency != "USD" {
		t.Errorf("Currency mismatch: got %s", paymentIntent.Currency)
	}
	if paymentIntent.Description != "Test payment intent" {
		t.Errorf("Description mismatch: got %s", paymentIntent.Description)
	}
	if paymentIntent.IntentStatus != "REQUIRES_PAYMENT_METHOD" {
		t.Errorf("IntentStatus mismatch: got %s", paymentIntent.IntentStatus)
	}
	if paymentIntent.MerchantOrderID != "test-order-002" {
		t.Errorf("MerchantOrderID mismatch: got %s", paymentIntent.MerchantOrderID)
	}
	if paymentIntent.CreateTime != "2026-01-21T12:39:39.716846826+08:00" {
		t.Errorf("CreateTime mismatch: got %s", paymentIntent.CreateTime)
	}

	// Verify metadata
	if paymentIntent.Metadata == nil {
		t.Fatal("Metadata should not be nil")
	}
	if paymentIntent.Metadata["test"] != "true" {
		t.Errorf("Metadata[test] mismatch: got %s", paymentIntent.Metadata["test"])
	}

	// Verify nullable fields
	if paymentIntent.PaymentMethod != nil {
		t.Error("PaymentMethod should be nil")
	}
	if paymentIntent.CancelTime != nil {
		t.Error("CancelTime should be nil")
	}
	if paymentIntent.CompleteTime != nil {
		t.Error("CompleteTime should be nil")
	}
	if paymentIntent.CancellationReason != "" {
		t.Errorf("CancellationReason should be empty, got %s", paymentIntent.CancellationReason)
	}

	t.Logf("Payment intent data parsed successfully")
	t.Logf("   PaymentIntentID: %s", paymentIntent.PaymentIntentID)
	t.Logf("   Amount: %s %s", paymentIntent.Amount, paymentIntent.Currency)
	t.Logf("   Status: %s", paymentIntent.IntentStatus)
	t.Logf("   MerchantOrderID: %s", paymentIntent.MerchantOrderID)
	t.Logf("   Metadata: %v", paymentIntent.Metadata)
}

func TestPaymentIntentEventConstants(t *testing.T) {
	// Verify event name constant
	if webhook.EventNameAcquiring != "ACQUIRING" {
		t.Errorf("EventNameAcquiring mismatch: got %s", webhook.EventNameAcquiring)
	}

	// Verify event type constants
	if webhook.EventTypePaymentIntentCreated != "acquiring.payment_intent.created" {
		t.Errorf("EventTypePaymentIntentCreated mismatch: got %s", webhook.EventTypePaymentIntentCreated)
	}
	if webhook.EventTypePaymentIntentSucceeded != "acquiring.payment_intent.succeeded" {
		t.Errorf("EventTypePaymentIntentSucceeded mismatch: got %s", webhook.EventTypePaymentIntentSucceeded)
	}
	if webhook.EventTypePaymentIntentFailed != "acquiring.payment_intent.failed" {
		t.Errorf("EventTypePaymentIntentFailed mismatch: got %s", webhook.EventTypePaymentIntentFailed)
	}
	if webhook.EventTypePaymentIntentCanceled != "acquiring.payment_intent.canceled" {
		t.Errorf("EventTypePaymentIntentCanceled mismatch: got %s", webhook.EventTypePaymentIntentCanceled)
	}

	// Verify intent status constants
	if webhook.IntentStatusRequiresPaymentMethod != "REQUIRES_PAYMENT_METHOD" {
		t.Errorf("IntentStatusRequiresPaymentMethod mismatch")
	}
	if webhook.IntentStatusSucceeded != "SUCCEEDED" {
		t.Errorf("IntentStatusSucceeded mismatch")
	}
	if webhook.IntentStatusCanceled != "CANCELED" {
		t.Errorf("IntentStatusCanceled mismatch")
	}
	if webhook.IntentStatusFailed != "FAILED" {
		t.Errorf("IntentStatusFailed mismatch")
	}

	t.Logf("Payment intent constants verified")
}

func TestParsePaymentIntentData_WrongEventType(t *testing.T) {
	// Try to parse account event as payment intent
	var event webhook.Event
	json.Unmarshal([]byte(accountCreateWebhookJSON), &event)

	_, err := event.ParsePaymentIntentData()
	if err == nil {
		t.Error("ParsePaymentIntentData should fail for account event type")
	}

	t.Logf("Wrong event type correctly rejected: %v", err)
}

// Payment intent failed webhook JSON
const paymentIntentFailedWebhookJSON = `{
	"version":"V1.6.0",
	"event_name":"ACQUIRING",
	"event_type":"acquiring.payment_intent.failed",
	"event_id":"1e77dbee-4b0b-4b5b-9599-9bbbcf148fd4",
	"source_id":"PI2013833849980588032",
	"data":{
		"amount":"101",
		"cancel_time":null,
		"cancellation_reason":"",
		"complete_time":"2026-01-21T13:09:47.727+08:00",
		"create_time":"2026-01-21T12:39:39.717+08:00",
		"currency":"USD",
		"description":"Test payment intent",
		"intent_status":"FAILED",
		"merchant_order_id":"test-order-002",
		"metadata":{"test":"true"},
		"payment_intent_id":"PI2013833849980588032",
		"payment_method":null
	}
}`

// Payment intent succeeded webhook JSON
const paymentIntentSucceededWebhookJSON = `{
	"version":"V1.6.0",
	"event_name":"ACQUIRING",
	"event_type":"acquiring.payment_intent.succeeded",
	"event_id":"2a88cf2e-5c1d-4c6c-a6aa-1cccdf259e5e",
	"source_id":"PI2013833849980588032",
	"data":{
		"amount":"101",
		"cancel_time":null,
		"cancellation_reason":"",
		"complete_time":"2026-01-21T13:15:22.456+08:00",
		"create_time":"2026-01-21T12:39:39.717+08:00",
		"currency":"USD",
		"description":"Test payment intent",
		"intent_status":"SUCCEEDED",
		"merchant_order_id":"test-order-002",
		"metadata":{"test":"true"},
		"payment_intent_id":"PI2013833849980588032",
		"payment_method":{
			"type":"card",
			"card":{
				"brand":"visa",
				"last4":"4242",
				"exp_month":12,
				"exp_year":2027,
				"funding":"credit",
				"country":"US"
			}
		}
	}
}`

// Payment intent canceled webhook JSON
const paymentIntentCanceledWebhookJSON = `{
	"version":"V1.6.0",
	"event_name":"ACQUIRING",
	"event_type":"acquiring.payment_intent.canceled",
	"event_id":"3b99df3f-6d2e-5d7d-b7bb-2dddef36af6f",
	"source_id":"PI2013833849980588032",
	"data":{
		"amount":"101",
		"cancel_time":"2026-01-21T13:20:15.123+08:00",
		"cancellation_reason":"requested_by_customer",
		"complete_time":null,
		"create_time":"2026-01-21T12:39:39.717+08:00",
		"currency":"USD",
		"description":"Test payment intent",
		"intent_status":"CANCELED",
		"merchant_order_id":"test-order-002",
		"metadata":{"test":"true"},
		"payment_intent_id":"PI2013833849980588032",
		"payment_method":null
	}
}`

func TestParsePaymentIntentFailedWebhook(t *testing.T) {
	var event webhook.Event
	err := json.Unmarshal([]byte(paymentIntentFailedWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse payment intent failed webhook: %v", err)
	}

	// Verify event type
	if event.EventType != webhook.EventTypePaymentIntentFailed {
		t.Errorf("EventType mismatch: got %s, want %s", event.EventType, webhook.EventTypePaymentIntentFailed)
	}

	// Should still be recognized as payment intent event
	if !event.IsPaymentIntentEvent() {
		t.Error("IsPaymentIntentEvent should return true for failed event")
	}

	// Parse data using the same struct
	paymentIntent, err := event.ParsePaymentIntentData()
	if err != nil {
		t.Fatalf("Failed to parse payment intent data: %v", err)
	}

	// Verify status is FAILED
	if paymentIntent.IntentStatus != webhook.IntentStatusFailed {
		t.Errorf("IntentStatus mismatch: got %s, want %s", paymentIntent.IntentStatus, webhook.IntentStatusFailed)
	}

	// Verify complete_time is now set (unlike created event)
	if paymentIntent.CompleteTime == nil {
		t.Error("CompleteTime should not be nil for failed event")
	} else if *paymentIntent.CompleteTime != "2026-01-21T13:09:47.727+08:00" {
		t.Errorf("CompleteTime mismatch: got %s", *paymentIntent.CompleteTime)
	}

	t.Logf("Payment intent failed webhook parsed successfully")
	t.Logf("   EventType: %s", event.EventType)
	t.Logf("   IntentStatus: %s", paymentIntent.IntentStatus)
	t.Logf("   CompleteTime: %s", *paymentIntent.CompleteTime)
}

func TestParsePaymentIntentSucceededWebhook(t *testing.T) {
	var event webhook.Event
	err := json.Unmarshal([]byte(paymentIntentSucceededWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse payment intent succeeded webhook: %v", err)
	}

	// Verify event type
	if event.EventType != webhook.EventTypePaymentIntentSucceeded {
		t.Errorf("EventType mismatch: got %s, want %s", event.EventType, webhook.EventTypePaymentIntentSucceeded)
	}

	if !event.IsPaymentIntentEvent() {
		t.Error("IsPaymentIntentEvent should return true for succeeded event")
	}

	paymentIntent, err := event.ParsePaymentIntentData()
	if err != nil {
		t.Fatalf("Failed to parse payment intent data: %v", err)
	}

	// Verify status is SUCCEEDED
	if paymentIntent.IntentStatus != webhook.IntentStatusSucceeded {
		t.Errorf("IntentStatus mismatch: got %s, want %s", paymentIntent.IntentStatus, webhook.IntentStatusSucceeded)
	}

	// Verify complete_time is set
	if paymentIntent.CompleteTime == nil {
		t.Error("CompleteTime should not be nil for succeeded event")
	}

	// Verify payment method is populated
	if paymentIntent.PaymentMethod == nil {
		t.Fatal("PaymentMethod should not be nil for succeeded event")
	}
	if paymentIntent.PaymentMethod.Type != "card" {
		t.Errorf("PaymentMethod.Type mismatch: got %s", paymentIntent.PaymentMethod.Type)
	}
	if paymentIntent.PaymentMethod.Card == nil {
		t.Fatal("PaymentMethod.Card should not be nil")
	}
	if paymentIntent.PaymentMethod.Card.Brand != "visa" {
		t.Errorf("Card.Brand mismatch: got %s", paymentIntent.PaymentMethod.Card.Brand)
	}
	if paymentIntent.PaymentMethod.Card.Last4 != "4242" {
		t.Errorf("Card.Last4 mismatch: got %s", paymentIntent.PaymentMethod.Card.Last4)
	}
	if paymentIntent.PaymentMethod.Card.ExpMonth != 12 {
		t.Errorf("Card.ExpMonth mismatch: got %d", paymentIntent.PaymentMethod.Card.ExpMonth)
	}
	if paymentIntent.PaymentMethod.Card.ExpYear != 2027 {
		t.Errorf("Card.ExpYear mismatch: got %d", paymentIntent.PaymentMethod.Card.ExpYear)
	}
	if paymentIntent.PaymentMethod.Card.Funding != "credit" {
		t.Errorf("Card.Funding mismatch: got %s", paymentIntent.PaymentMethod.Card.Funding)
	}
	if paymentIntent.PaymentMethod.Card.Country != "US" {
		t.Errorf("Card.Country mismatch: got %s", paymentIntent.PaymentMethod.Card.Country)
	}

	// Verify cancel fields are not set
	if paymentIntent.CancelTime != nil {
		t.Error("CancelTime should be nil for succeeded event")
	}
	if paymentIntent.CancellationReason != "" {
		t.Errorf("CancellationReason should be empty, got %s", paymentIntent.CancellationReason)
	}

	t.Logf("Payment intent succeeded webhook parsed successfully")
	t.Logf("   EventType: %s", event.EventType)
	t.Logf("   IntentStatus: %s", paymentIntent.IntentStatus)
	t.Logf("   CompleteTime: %s", *paymentIntent.CompleteTime)
	t.Logf("   PaymentMethod: %s (****%s)", paymentIntent.PaymentMethod.Card.Brand, paymentIntent.PaymentMethod.Card.Last4)
}

func TestParsePaymentIntentCanceledWebhook(t *testing.T) {
	var event webhook.Event
	err := json.Unmarshal([]byte(paymentIntentCanceledWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse payment intent canceled webhook: %v", err)
	}

	// Verify event type
	if event.EventType != webhook.EventTypePaymentIntentCanceled {
		t.Errorf("EventType mismatch: got %s, want %s", event.EventType, webhook.EventTypePaymentIntentCanceled)
	}

	if !event.IsPaymentIntentEvent() {
		t.Error("IsPaymentIntentEvent should return true for canceled event")
	}

	paymentIntent, err := event.ParsePaymentIntentData()
	if err != nil {
		t.Fatalf("Failed to parse payment intent data: %v", err)
	}

	// Verify status is CANCELED
	if paymentIntent.IntentStatus != webhook.IntentStatusCanceled {
		t.Errorf("IntentStatus mismatch: got %s, want %s", paymentIntent.IntentStatus, webhook.IntentStatusCanceled)
	}

	// Verify cancel_time is set
	if paymentIntent.CancelTime == nil {
		t.Error("CancelTime should not be nil for canceled event")
	} else if *paymentIntent.CancelTime != "2026-01-21T13:20:15.123+08:00" {
		t.Errorf("CancelTime mismatch: got %s", *paymentIntent.CancelTime)
	}

	// Verify cancellation_reason is set
	if paymentIntent.CancellationReason != "requested_by_customer" {
		t.Errorf("CancellationReason mismatch: got %s", paymentIntent.CancellationReason)
	}

	// Verify complete_time is NOT set (canceled before completion)
	if paymentIntent.CompleteTime != nil {
		t.Error("CompleteTime should be nil for canceled event")
	}

	t.Logf("Payment intent canceled webhook parsed successfully")
	t.Logf("   EventType: %s", event.EventType)
	t.Logf("   IntentStatus: %s", paymentIntent.IntentStatus)
	t.Logf("   CancelTime: %s", *paymentIntent.CancelTime)
	t.Logf("   CancellationReason: %s", paymentIntent.CancellationReason)
}

// Test handling all payment intent event types with same struct
func TestPaymentIntentAllEventTypes(t *testing.T) {
	testCases := []struct {
		name             string
		json             string
		expectedType     string
		expectedStatus   string
		hasCompleteTime  bool
		hasCancelTime    bool
		hasPaymentMethod bool
		hasCancelReason  bool
	}{
		{
			name:             "Created",
			json:             paymentIntentCreatedWebhookJSON,
			expectedType:     webhook.EventTypePaymentIntentCreated,
			expectedStatus:   webhook.IntentStatusRequiresPaymentMethod,
			hasCompleteTime:  false,
			hasCancelTime:    false,
			hasPaymentMethod: false,
			hasCancelReason:  false,
		},
		{
			name:             "Succeeded",
			json:             paymentIntentSucceededWebhookJSON,
			expectedType:     webhook.EventTypePaymentIntentSucceeded,
			expectedStatus:   webhook.IntentStatusSucceeded,
			hasCompleteTime:  true,
			hasCancelTime:    false,
			hasPaymentMethod: true,
			hasCancelReason:  false,
		},
		{
			name:             "Failed",
			json:             paymentIntentFailedWebhookJSON,
			expectedType:     webhook.EventTypePaymentIntentFailed,
			expectedStatus:   webhook.IntentStatusFailed,
			hasCompleteTime:  true,
			hasCancelTime:    false,
			hasPaymentMethod: false,
			hasCancelReason:  false,
		},
		{
			name:             "Canceled",
			json:             paymentIntentCanceledWebhookJSON,
			expectedType:     webhook.EventTypePaymentIntentCanceled,
			expectedStatus:   webhook.IntentStatusCanceled,
			hasCompleteTime:  false,
			hasCancelTime:    true,
			hasPaymentMethod: false,
			hasCancelReason:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var event webhook.Event
			err := json.Unmarshal([]byte(tc.json), &event)
			if err != nil {
				t.Fatalf("Failed to parse event: %v", err)
			}

			if event.EventType != tc.expectedType {
				t.Errorf("EventType mismatch: got %s, want %s", event.EventType, tc.expectedType)
			}

			if !event.IsPaymentIntentEvent() {
				t.Error("IsPaymentIntentEvent should return true")
			}

			if !event.IsAcquiringEvent() {
				t.Error("IsAcquiringEvent should return true")
			}

			paymentIntent, err := event.ParsePaymentIntentData()
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			if paymentIntent.IntentStatus != tc.expectedStatus {
				t.Errorf("IntentStatus mismatch: got %s, want %s", paymentIntent.IntentStatus, tc.expectedStatus)
			}

			// Verify optional fields based on event type
			if tc.hasCompleteTime && paymentIntent.CompleteTime == nil {
				t.Error("CompleteTime should not be nil")
			}
			if !tc.hasCompleteTime && paymentIntent.CompleteTime != nil {
				t.Error("CompleteTime should be nil")
			}

			if tc.hasCancelTime && paymentIntent.CancelTime == nil {
				t.Error("CancelTime should not be nil")
			}
			if !tc.hasCancelTime && paymentIntent.CancelTime != nil {
				t.Error("CancelTime should be nil")
			}

			if tc.hasPaymentMethod && paymentIntent.PaymentMethod == nil {
				t.Error("PaymentMethod should not be nil")
			}
			if !tc.hasPaymentMethod && paymentIntent.PaymentMethod != nil {
				t.Error("PaymentMethod should be nil")
			}

			if tc.hasCancelReason && paymentIntent.CancellationReason == "" {
				t.Error("CancellationReason should not be empty")
			}
			if !tc.hasCancelReason && paymentIntent.CancellationReason != "" {
				t.Error("CancellationReason should be empty")
			}

			t.Logf("Verified: EventType=%s, IntentStatus=%s", event.EventType, paymentIntent.IntentStatus)
		})
	}
}

// ============================================================================
// Payment Attempt Webhook Tests
// ============================================================================

// Payment attempt created webhook JSON
const paymentAttemptCreatedWebhookJSON = `{
	"version":"V1.6.0",
	"event_name":"ACQUIRING",
	"event_type":"acquiring.payment_attempt.created",
	"event_id":"3b1e591c-c9cd-4e26-be19-d83d7ec907e1",
	"source_id":"PA2013848038472159232",
	"data":{
		"amount":"0.01",
		"attempt_status":"INITIATED",
		"cancel_time":null,
		"cancellation_reason":"",
		"captured_amount":"0.01",
		"complete_time":null,
		"create_time":"2026-01-21T13:36:02.516375218+08:00",
		"currency":"USD",
		"failure_code":"",
		"merchant_order_id":"test-AlipayCN_QRCode-001",
		"payment_attempt_id":"PA2013848038472159232",
		"payment_intent_id":"PI2013848035972354048",
		"payment_method":{
			"alipaycn":{
				"flow":"qrcode",
				"os_type":""
			},
			"type":"alipaycn"
		},
		"refunded_amount":"0"
	}
}`

// Payment attempt capture requested webhook JSON
const paymentAttemptCaptureRequestedWebhookJSON = `{
	"version":"V1.6.0",
	"event_name":"ACQUIRING",
	"event_type":"acquiring.payment_attempt.capture_requested",
	"event_id":"afd89fc8-752b-4905-a4f1-c37bb718cfcd",
	"source_id":"PA2013848038472159232",
	"data":{
		"amount":"0.01",
		"attempt_status":"CAPTURE_REQUESTED",
		"cancel_time":null,
		"cancellation_reason":"",
		"captured_amount":"0.01",
		"complete_time":null,
		"create_time":"2026-01-21T13:36:02.516+08:00",
		"currency":"USD",
		"failure_code":"",
		"merchant_order_id":"test-AlipayCN_QRCode-001",
		"payment_attempt_id":"PA2013848038472159232",
		"payment_intent_id":"PI2013848035972354048",
		"payment_method":{
			"alipaycn":{
				"flow":"qrcode",
				"os_type":""
			},
			"type":"alipaycn"
		},
		"refunded_amount":"0"
	}
}`

// Payment attempt succeeded webhook JSON
const paymentAttemptSucceededWebhookJSON = `{
	"version":"V1.6.0",
	"event_name":"ACQUIRING",
	"event_type":"acquiring.payment_attempt.succeeded",
	"event_id":"4c2f602d-d0de-5f37-cf2a-e94e8fd018f2",
	"source_id":"PA2013848038472159232",
	"data":{
		"amount":"0.01",
		"attempt_status":"SUCCEEDED",
		"cancel_time":null,
		"cancellation_reason":"",
		"captured_amount":"0.01",
		"complete_time":"2026-01-21T13:40:15.789+08:00",
		"create_time":"2026-01-21T13:36:02.516+08:00",
		"currency":"USD",
		"failure_code":"",
		"merchant_order_id":"test-AlipayCN_QRCode-001",
		"payment_attempt_id":"PA2013848038472159232",
		"payment_intent_id":"PI2013848035972354048",
		"payment_method":{
			"alipaycn":{
				"flow":"qrcode",
				"os_type":""
			},
			"type":"alipaycn"
		},
		"refunded_amount":"0"
	}
}`

// Payment attempt failed webhook JSON
const paymentAttemptFailedWebhookJSON = `{
	"version":"V1.6.0",
	"event_name":"ACQUIRING",
	"event_type":"acquiring.payment_attempt.failed",
	"event_id":"5d3g713e-e1ef-6g48-dg3b-f05f9ge129g3",
	"source_id":"PA2013848038472159232",
	"data":{
		"amount":"0.01",
		"attempt_status":"FAILED",
		"cancel_time":null,
		"cancellation_reason":"",
		"captured_amount":"0",
		"complete_time":"2026-01-21T13:42:30.123+08:00",
		"create_time":"2026-01-21T13:36:02.516+08:00",
		"currency":"USD",
		"failure_code":"payment_declined",
		"merchant_order_id":"test-AlipayCN_QRCode-001",
		"payment_attempt_id":"PA2013848038472159232",
		"payment_intent_id":"PI2013848035972354048",
		"payment_method":{
			"alipaycn":{
				"flow":"qrcode",
				"os_type":""
			},
			"type":"alipaycn"
		},
		"refunded_amount":"0"
	}
}`

// Payment attempt canceled webhook JSON
const paymentAttemptCanceledWebhookJSON = `{
	"version":"V1.6.0",
	"event_name":"ACQUIRING",
	"event_type":"acquiring.payment_attempt.canceled",
	"event_id":"6e4h824f-f2fg-7h59-eh4c-g16g0hf230h4",
	"source_id":"PA2013848038472159232",
	"data":{
		"amount":"0.01",
		"attempt_status":"CANCELED",
		"cancel_time":"2026-01-21T13:45:00.456+08:00",
		"cancellation_reason":"requested_by_customer",
		"captured_amount":"0",
		"complete_time":null,
		"create_time":"2026-01-21T13:36:02.516+08:00",
		"currency":"USD",
		"failure_code":"",
		"merchant_order_id":"test-AlipayCN_QRCode-001",
		"payment_attempt_id":"PA2013848038472159232",
		"payment_intent_id":"PI2013848035972354048",
		"payment_method":{
			"alipaycn":{
				"flow":"qrcode",
				"os_type":""
			},
			"type":"alipaycn"
		},
		"refunded_amount":"0"
	}
}`

func TestParsePaymentAttemptCreatedWebhook(t *testing.T) {
	var event webhook.Event
	err := json.Unmarshal([]byte(paymentAttemptCreatedWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse payment attempt created webhook: %v", err)
	}

	// Verify event type
	if event.EventType != webhook.EventTypePaymentAttemptCreated {
		t.Errorf("EventType mismatch: got %s, want %s", event.EventType, webhook.EventTypePaymentAttemptCreated)
	}

	if !event.IsPaymentAttemptEvent() {
		t.Error("IsPaymentAttemptEvent should return true")
	}

	if !event.IsAcquiringEvent() {
		t.Error("IsAcquiringEvent should return true")
	}

	// Should NOT be a payment intent event
	if event.IsPaymentIntentEvent() {
		t.Error("IsPaymentIntentEvent should return false for attempt event")
	}

	attempt, err := event.ParsePaymentAttemptData()
	if err != nil {
		t.Fatalf("Failed to parse payment attempt data: %v", err)
	}

	// Verify core fields
	if attempt.PaymentAttemptID != "PA2013848038472159232" {
		t.Errorf("PaymentAttemptID mismatch: got %s", attempt.PaymentAttemptID)
	}
	if attempt.PaymentIntentID != "PI2013848035972354048" {
		t.Errorf("PaymentIntentID mismatch: got %s", attempt.PaymentIntentID)
	}
	if attempt.Amount != "0.01" {
		t.Errorf("Amount mismatch: got %s", attempt.Amount)
	}
	if attempt.Currency != "USD" {
		t.Errorf("Currency mismatch: got %s", attempt.Currency)
	}
	if attempt.AttemptStatus != webhook.AttemptStatusInitiated {
		t.Errorf("AttemptStatus mismatch: got %s", attempt.AttemptStatus)
	}
	if attempt.MerchantOrderID != "test-AlipayCN_QRCode-001" {
		t.Errorf("MerchantOrderID mismatch: got %s", attempt.MerchantOrderID)
	}
	if attempt.CapturedAmount != "0.01" {
		t.Errorf("CapturedAmount mismatch: got %s", attempt.CapturedAmount)
	}
	if attempt.RefundedAmount != "0" {
		t.Errorf("RefundedAmount mismatch: got %s", attempt.RefundedAmount)
	}

	// Verify payment method
	if attempt.PaymentMethod == nil {
		t.Fatal("PaymentMethod should not be nil")
	}
	if attempt.PaymentMethod.Type != "alipaycn" {
		t.Errorf("PaymentMethod.Type mismatch: got %s", attempt.PaymentMethod.Type)
	}
	if attempt.PaymentMethod.AlipayCN == nil {
		t.Fatal("PaymentMethod.AlipayCN should not be nil")
	}
	if attempt.PaymentMethod.AlipayCN.Flow != "qrcode" {
		t.Errorf("AlipayCN.Flow mismatch: got %s", attempt.PaymentMethod.AlipayCN.Flow)
	}

	t.Logf("Payment attempt created webhook parsed successfully")
	t.Logf("   PaymentAttemptID: %s", attempt.PaymentAttemptID)
	t.Logf("   PaymentIntentID: %s", attempt.PaymentIntentID)
	t.Logf("   AttemptStatus: %s", attempt.AttemptStatus)
	t.Logf("   PaymentMethod: %s (%s)", attempt.PaymentMethod.Type, attempt.PaymentMethod.AlipayCN.Flow)
}

func TestParsePaymentAttemptCaptureRequestedWebhook(t *testing.T) {
	var event webhook.Event
	err := json.Unmarshal([]byte(paymentAttemptCaptureRequestedWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse payment attempt capture requested webhook: %v", err)
	}

	if event.EventType != webhook.EventTypePaymentAttemptCaptureRequested {
		t.Errorf("EventType mismatch: got %s", event.EventType)
	}

	if !event.IsPaymentAttemptEvent() {
		t.Error("IsPaymentAttemptEvent should return true")
	}

	attempt, err := event.ParsePaymentAttemptData()
	if err != nil {
		t.Fatalf("Failed to parse payment attempt data: %v", err)
	}

	if attempt.AttemptStatus != webhook.AttemptStatusCaptureRequested {
		t.Errorf("AttemptStatus mismatch: got %s", attempt.AttemptStatus)
	}

	t.Logf("Payment attempt capture requested webhook parsed successfully")
	t.Logf("   AttemptStatus: %s", attempt.AttemptStatus)
}

func TestParsePaymentAttemptSucceededWebhook(t *testing.T) {
	var event webhook.Event
	err := json.Unmarshal([]byte(paymentAttemptSucceededWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse payment attempt succeeded webhook: %v", err)
	}

	if event.EventType != webhook.EventTypePaymentAttemptSucceeded {
		t.Errorf("EventType mismatch: got %s", event.EventType)
	}

	attempt, err := event.ParsePaymentAttemptData()
	if err != nil {
		t.Fatalf("Failed to parse payment attempt data: %v", err)
	}

	if attempt.AttemptStatus != webhook.AttemptStatusSucceeded {
		t.Errorf("AttemptStatus mismatch: got %s", attempt.AttemptStatus)
	}

	// Succeeded should have complete_time
	if attempt.CompleteTime == nil {
		t.Error("CompleteTime should not be nil for succeeded attempt")
	}

	t.Logf("Payment attempt succeeded webhook parsed successfully")
	t.Logf("   AttemptStatus: %s", attempt.AttemptStatus)
	t.Logf("   CompleteTime: %s", *attempt.CompleteTime)
}

func TestParsePaymentAttemptFailedWebhook(t *testing.T) {
	var event webhook.Event
	err := json.Unmarshal([]byte(paymentAttemptFailedWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse payment attempt failed webhook: %v", err)
	}

	if event.EventType != webhook.EventTypePaymentAttemptFailed {
		t.Errorf("EventType mismatch: got %s", event.EventType)
	}

	attempt, err := event.ParsePaymentAttemptData()
	if err != nil {
		t.Fatalf("Failed to parse payment attempt data: %v", err)
	}

	if attempt.AttemptStatus != webhook.AttemptStatusFailed {
		t.Errorf("AttemptStatus mismatch: got %s", attempt.AttemptStatus)
	}

	// Failed should have failure_code and complete_time
	if attempt.FailureCode != "payment_declined" {
		t.Errorf("FailureCode mismatch: got %s", attempt.FailureCode)
	}
	if attempt.CompleteTime == nil {
		t.Error("CompleteTime should not be nil for failed attempt")
	}

	t.Logf("Payment attempt failed webhook parsed successfully")
	t.Logf("   AttemptStatus: %s", attempt.AttemptStatus)
	t.Logf("   FailureCode: %s", attempt.FailureCode)
}

func TestParsePaymentAttemptCanceledWebhook(t *testing.T) {
	var event webhook.Event
	err := json.Unmarshal([]byte(paymentAttemptCanceledWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse payment attempt canceled webhook: %v", err)
	}

	if event.EventType != webhook.EventTypePaymentAttemptCanceled {
		t.Errorf("EventType mismatch: got %s", event.EventType)
	}

	attempt, err := event.ParsePaymentAttemptData()
	if err != nil {
		t.Fatalf("Failed to parse payment attempt data: %v", err)
	}

	if attempt.AttemptStatus != webhook.AttemptStatusCanceled {
		t.Errorf("AttemptStatus mismatch: got %s", attempt.AttemptStatus)
	}

	// Canceled should have cancel_time and cancellation_reason
	if attempt.CancelTime == nil {
		t.Error("CancelTime should not be nil for canceled attempt")
	}
	if attempt.CancellationReason != "requested_by_customer" {
		t.Errorf("CancellationReason mismatch: got %s", attempt.CancellationReason)
	}

	t.Logf("Payment attempt canceled webhook parsed successfully")
	t.Logf("   AttemptStatus: %s", attempt.AttemptStatus)
	t.Logf("   CancelTime: %s", *attempt.CancelTime)
	t.Logf("   CancellationReason: %s", attempt.CancellationReason)
}

// Test handling all payment attempt event types with same struct
func TestPaymentAttemptAllEventTypes(t *testing.T) {
	testCases := []struct {
		name            string
		json            string
		expectedType    string
		expectedStatus  string
		hasCompleteTime bool
		hasCancelTime   bool
		hasFailureCode  bool
		hasCancelReason bool
	}{
		{
			name:            "Created",
			json:            paymentAttemptCreatedWebhookJSON,
			expectedType:    webhook.EventTypePaymentAttemptCreated,
			expectedStatus:  webhook.AttemptStatusInitiated,
			hasCompleteTime: false,
			hasCancelTime:   false,
			hasFailureCode:  false,
			hasCancelReason: false,
		},
		{
			name:            "CaptureRequested",
			json:            paymentAttemptCaptureRequestedWebhookJSON,
			expectedType:    webhook.EventTypePaymentAttemptCaptureRequested,
			expectedStatus:  webhook.AttemptStatusCaptureRequested,
			hasCompleteTime: false,
			hasCancelTime:   false,
			hasFailureCode:  false,
			hasCancelReason: false,
		},
		{
			name:            "Succeeded",
			json:            paymentAttemptSucceededWebhookJSON,
			expectedType:    webhook.EventTypePaymentAttemptSucceeded,
			expectedStatus:  webhook.AttemptStatusSucceeded,
			hasCompleteTime: true,
			hasCancelTime:   false,
			hasFailureCode:  false,
			hasCancelReason: false,
		},
		{
			name:            "Failed",
			json:            paymentAttemptFailedWebhookJSON,
			expectedType:    webhook.EventTypePaymentAttemptFailed,
			expectedStatus:  webhook.AttemptStatusFailed,
			hasCompleteTime: true,
			hasCancelTime:   false,
			hasFailureCode:  true,
			hasCancelReason: false,
		},
		{
			name:            "Canceled",
			json:            paymentAttemptCanceledWebhookJSON,
			expectedType:    webhook.EventTypePaymentAttemptCanceled,
			expectedStatus:  webhook.AttemptStatusCanceled,
			hasCompleteTime: false,
			hasCancelTime:   true,
			hasFailureCode:  false,
			hasCancelReason: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var event webhook.Event
			err := json.Unmarshal([]byte(tc.json), &event)
			if err != nil {
				t.Fatalf("Failed to parse event: %v", err)
			}

			if event.EventType != tc.expectedType {
				t.Errorf("EventType mismatch: got %s, want %s", event.EventType, tc.expectedType)
			}

			if !event.IsPaymentAttemptEvent() {
				t.Error("IsPaymentAttemptEvent should return true")
			}

			if !event.IsAcquiringEvent() {
				t.Error("IsAcquiringEvent should return true")
			}

			attempt, err := event.ParsePaymentAttemptData()
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			if attempt.AttemptStatus != tc.expectedStatus {
				t.Errorf("AttemptStatus mismatch: got %s, want %s", attempt.AttemptStatus, tc.expectedStatus)
			}

			// Verify optional fields based on event type
			if tc.hasCompleteTime && attempt.CompleteTime == nil {
				t.Error("CompleteTime should not be nil")
			}
			if !tc.hasCompleteTime && attempt.CompleteTime != nil {
				t.Error("CompleteTime should be nil")
			}

			if tc.hasCancelTime && attempt.CancelTime == nil {
				t.Error("CancelTime should not be nil")
			}
			if !tc.hasCancelTime && attempt.CancelTime != nil {
				t.Error("CancelTime should be nil")
			}

			if tc.hasFailureCode && attempt.FailureCode == "" {
				t.Error("FailureCode should not be empty")
			}
			if !tc.hasFailureCode && attempt.FailureCode != "" {
				t.Error("FailureCode should be empty")
			}

			if tc.hasCancelReason && attempt.CancellationReason == "" {
				t.Error("CancellationReason should not be empty")
			}
			if !tc.hasCancelReason && attempt.CancellationReason != "" {
				t.Error("CancellationReason should be empty")
			}

			t.Logf("Verified: EventType=%s, AttemptStatus=%s", event.EventType, attempt.AttemptStatus)
		})
	}
}

func TestPaymentAttemptEventConstants(t *testing.T) {
	// Verify event type constants
	if webhook.EventTypePaymentAttemptCreated != "acquiring.payment_attempt.created" {
		t.Errorf("EventTypePaymentAttemptCreated mismatch")
	}
	if webhook.EventTypePaymentAttemptCaptureRequested != "acquiring.payment_attempt.capture_requested" {
		t.Errorf("EventTypePaymentAttemptCaptureRequested mismatch")
	}
	if webhook.EventTypePaymentAttemptSucceeded != "acquiring.payment_attempt.succeeded" {
		t.Errorf("EventTypePaymentAttemptSucceeded mismatch")
	}
	if webhook.EventTypePaymentAttemptFailed != "acquiring.payment_attempt.failed" {
		t.Errorf("EventTypePaymentAttemptFailed mismatch")
	}
	if webhook.EventTypePaymentAttemptCanceled != "acquiring.payment_attempt.canceled" {
		t.Errorf("EventTypePaymentAttemptCanceled mismatch")
	}

	// Verify attempt status constants
	if webhook.AttemptStatusInitiated != "INITIATED" {
		t.Errorf("AttemptStatusInitiated mismatch")
	}
	if webhook.AttemptStatusCaptureRequested != "CAPTURE_REQUESTED" {
		t.Errorf("AttemptStatusCaptureRequested mismatch")
	}
	if webhook.AttemptStatusSucceeded != "SUCCEEDED" {
		t.Errorf("AttemptStatusSucceeded mismatch")
	}
	if webhook.AttemptStatusFailed != "FAILED" {
		t.Errorf("AttemptStatusFailed mismatch")
	}
	if webhook.AttemptStatusCanceled != "CANCELED" {
		t.Errorf("AttemptStatusCanceled mismatch")
	}

	t.Logf("Payment attempt constants verified")
}

func TestParsePaymentAttemptData_WrongEventType(t *testing.T) {
	// Try to parse payment intent event as payment attempt
	var event webhook.Event
	json.Unmarshal([]byte(paymentIntentCreatedWebhookJSON), &event)

	_, err := event.ParsePaymentAttemptData()
	if err == nil {
		t.Error("ParsePaymentAttemptData should fail for payment intent event type")
	}

	t.Logf("Wrong event type correctly rejected: %v", err)
}
