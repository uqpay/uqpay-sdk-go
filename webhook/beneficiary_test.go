package webhook

import (
	"encoding/json"
	"testing"
)

// ============================================================================
// Test Fixtures - Beneficiary Webhooks
// ============================================================================

const beneficiarySuccessfulWebhookJSON = `{
	"version": "V1.6.0",
	"event_name": "BENEFICIARY",
	"event_type": "beneficiary.successful",
	"event_id": "0d314cfb-188f-4f5e-a213-7805032c7ec5",
	"source_id": "cd69af8d-4bff-4e12-99fb-f27fa14255a0",
	"data": {
		"account_currency_code": "MYR",
		"account_id": "65087660-8d3d-428e-bd2e-9e56219c1512",
		"account_number": "3413423141212",
		"bank_country_code": "MY",
		"beneficiary_address": "{\"nationality\":\"MY\",\"country_code\":\"MY\",\"city\":\"KL\",\"street_address\":\"34234rwdasd , adsadas\",\"postal_code\":\"234312\",\"state\":\"KL\"}",
		"beneficiary_bank_details": "{\"bank_name\":\"MALAYAN BANKING BERHAD (MAYBANK)\",\"bank_address\":\"MENARAMAYBANKFLOOR8100 JALAN TUN PERAK KUALA LUMPUR KUALALUMPUR50050POB12010 MALAYSIA\",\"bank_country_code\":\"MY\",\"account_holder\":\"Alpha ac\",\"account_currency_code\":\"MYR\",\"account_number\":\"3413423141212\",\"iban\":\"\",\"swift_code\":\"MBBEMYKLXXX\",\"clearing_system\":\"LOCAL\",\"routing_code_type1\":\"\",\"routing_code_value1\":\"\",\"routing_code_type2\":\"\",\"routing_code_value2\":\"\"}",
		"beneficiary_company_name": "Alpha ac",
		"beneficiary_email": "",
		"beneficiary_entity_type": "INDIVIDUAL",
		"beneficiary_first_name": "Zu",
		"beneficiary_id": "cd69af8d-4bff-4e12-99fb-f27fa14255a0",
		"beneficiary_last_name": "Wa",
		"beneficiary_nickname": "Zu Wa",
		"beneficiary_status": "ACTIVE",
		"direct_id": "0",
		"payment_type": "LOCAL",
		"short_reference_id": "BF260126-4ZIYF4QC"
	}
}`

const beneficiaryCompanyWebhookJSON = `{
	"version": "V1.6.0",
	"event_name": "BENEFICIARY",
	"event_type": "beneficiary.successful",
	"event_id": "1e425dfc-299g-5g6f-b324-8916143d8e6a",
	"source_id": "de7abf9e-5cgg-5f23-aagb-g38gb25366b1",
	"data": {
		"account_currency_code": "USD",
		"account_id": "76198771-9e4e-539f-ce3f-af67320d2623",
		"account_number": "9876543210",
		"bank_country_code": "US",
		"beneficiary_address": "{\"nationality\":\"US\",\"country_code\":\"US\",\"city\":\"New York\",\"street_address\":\"123 Wall Street\",\"postal_code\":\"10005\",\"state\":\"NY\"}",
		"beneficiary_bank_details": "{\"bank_name\":\"BANK OF AMERICA\",\"bank_address\":\"100 N Tryon St, Charlotte, NC\",\"bank_country_code\":\"US\",\"account_holder\":\"ACME Corp\",\"account_currency_code\":\"USD\",\"account_number\":\"9876543210\",\"iban\":\"\",\"swift_code\":\"BOFAUS3N\",\"clearing_system\":\"LOCAL\",\"routing_code_type1\":\"ABA\",\"routing_code_value1\":\"026009593\",\"routing_code_type2\":\"\",\"routing_code_value2\":\"\"}",
		"beneficiary_company_name": "ACME Corp",
		"beneficiary_email": "payments@acme.com",
		"beneficiary_entity_type": "COMPANY",
		"beneficiary_first_name": "",
		"beneficiary_id": "de7abf9e-5cgg-5f23-aagb-g38gb25366b1",
		"beneficiary_last_name": "",
		"beneficiary_nickname": "ACME Corp",
		"beneficiary_status": "ACTIVE",
		"direct_id": "0",
		"payment_type": "INTERNATIONAL",
		"short_reference_id": "BF260126-5AJZG5RD"
	}
}`

// ============================================================================
// Beneficiary Successful Event Tests
// ============================================================================

func TestParseBeneficiarySuccessfulWebhook(t *testing.T) {
	var event Event
	err := json.Unmarshal([]byte(beneficiarySuccessfulWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse beneficiary successful webhook: %v", err)
	}

	// Verify envelope fields
	if event.Version != "V1.6.0" {
		t.Errorf("Version mismatch: got %s, want V1.6.0", event.Version)
	}
	if event.EventName != EventNameBeneficiary {
		t.Errorf("EventName mismatch: got %s, want %s", event.EventName, EventNameBeneficiary)
	}
	if event.EventType != EventTypeBeneficiarySuccessful {
		t.Errorf("EventType mismatch: got %s, want %s", event.EventType, EventTypeBeneficiarySuccessful)
	}
	if event.EventID != "0d314cfb-188f-4f5e-a213-7805032c7ec5" {
		t.Errorf("EventID mismatch: got %s", event.EventID)
	}
	if event.SourceID != "cd69af8d-4bff-4e12-99fb-f27fa14255a0" {
		t.Errorf("SourceID mismatch: got %s", event.SourceID)
	}

	// Verify helper methods
	if !event.IsBeneficiaryEvent() {
		t.Error("IsBeneficiaryEvent should return true")
	}
	if !event.IsBeneficiarySuccessfulEvent() {
		t.Error("IsBeneficiarySuccessfulEvent should return true")
	}
	if event.IsAcquiringEvent() {
		t.Error("IsAcquiringEvent should return false")
	}
}

func TestParseBeneficiaryData(t *testing.T) {
	var event Event
	json.Unmarshal([]byte(beneficiarySuccessfulWebhookJSON), &event)

	beneficiary, err := event.ParseBeneficiaryData()
	if err != nil {
		t.Fatalf("Failed to parse beneficiary data: %v", err)
	}

	// Verify account information
	if beneficiary.AccountID != "65087660-8d3d-428e-bd2e-9e56219c1512" {
		t.Errorf("AccountID mismatch: got %s", beneficiary.AccountID)
	}
	if beneficiary.AccountNumber != "3413423141212" {
		t.Errorf("AccountNumber mismatch: got %s", beneficiary.AccountNumber)
	}
	if beneficiary.AccountCurrencyCode != "MYR" {
		t.Errorf("AccountCurrencyCode mismatch: got %s", beneficiary.AccountCurrencyCode)
	}

	// Verify beneficiary identification
	if beneficiary.BeneficiaryID != "cd69af8d-4bff-4e12-99fb-f27fa14255a0" {
		t.Errorf("BeneficiaryID mismatch: got %s", beneficiary.BeneficiaryID)
	}
	if beneficiary.BeneficiaryFirstName != "Zu" {
		t.Errorf("BeneficiaryFirstName mismatch: got %s", beneficiary.BeneficiaryFirstName)
	}
	if beneficiary.BeneficiaryLastName != "Wa" {
		t.Errorf("BeneficiaryLastName mismatch: got %s", beneficiary.BeneficiaryLastName)
	}
	if beneficiary.BeneficiaryNickname != "Zu Wa" {
		t.Errorf("BeneficiaryNickname mismatch: got %s", beneficiary.BeneficiaryNickname)
	}

	// Verify beneficiary details
	if beneficiary.BeneficiaryCompanyName != "Alpha ac" {
		t.Errorf("BeneficiaryCompanyName mismatch: got %s", beneficiary.BeneficiaryCompanyName)
	}
	if beneficiary.BeneficiaryEntityType != BeneficiaryEntityTypeIndividual {
		t.Errorf("BeneficiaryEntityType mismatch: got %s", beneficiary.BeneficiaryEntityType)
	}
	if beneficiary.BeneficiaryStatus != BeneficiaryStatusActive {
		t.Errorf("BeneficiaryStatus mismatch: got %s", beneficiary.BeneficiaryStatus)
	}

	// Verify payment information
	if beneficiary.PaymentType != PaymentTypeLocal {
		t.Errorf("PaymentType mismatch: got %s", beneficiary.PaymentType)
	}
	if beneficiary.ShortReferenceID != "BF260126-4ZIYF4QC" {
		t.Errorf("ShortReferenceID mismatch: got %s", beneficiary.ShortReferenceID)
	}
	if beneficiary.BankCountryCode != "MY" {
		t.Errorf("BankCountryCode mismatch: got %s", beneficiary.BankCountryCode)
	}
}

// ============================================================================
// Nested JSON Parsing Tests
// ============================================================================

func TestParseBeneficiaryAddress(t *testing.T) {
	var event Event
	json.Unmarshal([]byte(beneficiarySuccessfulWebhookJSON), &event)

	beneficiary, _ := event.ParseBeneficiaryData()

	address, err := beneficiary.GetBeneficiaryAddress()
	if err != nil {
		t.Fatalf("Failed to parse beneficiary address: %v", err)
	}

	if address.Nationality != "MY" {
		t.Errorf("Nationality mismatch: got %s", address.Nationality)
	}
	if address.CountryCode != "MY" {
		t.Errorf("CountryCode mismatch: got %s", address.CountryCode)
	}
	if address.City != "KL" {
		t.Errorf("City mismatch: got %s", address.City)
	}
	if address.State != "KL" {
		t.Errorf("State mismatch: got %s", address.State)
	}
	if address.StreetAddress != "34234rwdasd , adsadas" {
		t.Errorf("StreetAddress mismatch: got %s", address.StreetAddress)
	}
	if address.PostalCode != "234312" {
		t.Errorf("PostalCode mismatch: got %s", address.PostalCode)
	}
}

func TestParseBeneficiaryBankDetails(t *testing.T) {
	var event Event
	json.Unmarshal([]byte(beneficiarySuccessfulWebhookJSON), &event)

	beneficiary, _ := event.ParseBeneficiaryData()

	bankDetails, err := beneficiary.GetBeneficiaryBankDetails()
	if err != nil {
		t.Fatalf("Failed to parse beneficiary bank details: %v", err)
	}

	if bankDetails.BankName != "MALAYAN BANKING BERHAD (MAYBANK)" {
		t.Errorf("BankName mismatch: got %s", bankDetails.BankName)
	}
	if bankDetails.BankCountryCode != "MY" {
		t.Errorf("BankCountryCode mismatch: got %s", bankDetails.BankCountryCode)
	}
	if bankDetails.AccountHolder != "Alpha ac" {
		t.Errorf("AccountHolder mismatch: got %s", bankDetails.AccountHolder)
	}
	if bankDetails.AccountCurrencyCode != "MYR" {
		t.Errorf("AccountCurrencyCode mismatch: got %s", bankDetails.AccountCurrencyCode)
	}
	if bankDetails.AccountNumber != "3413423141212" {
		t.Errorf("AccountNumber mismatch: got %s", bankDetails.AccountNumber)
	}
	if bankDetails.SwiftCode != "MBBEMYKLXXX" {
		t.Errorf("SwiftCode mismatch: got %s", bankDetails.SwiftCode)
	}
	if bankDetails.ClearingSystem != "LOCAL" {
		t.Errorf("ClearingSystem mismatch: got %s", bankDetails.ClearingSystem)
	}
}

// ============================================================================
// Helper Method Tests
// ============================================================================

func TestBeneficiaryHelperMethods(t *testing.T) {
	var event Event
	json.Unmarshal([]byte(beneficiarySuccessfulWebhookJSON), &event)

	beneficiary, _ := event.ParseBeneficiaryData()

	// Test GetFullName
	fullName := beneficiary.GetFullName()
	if fullName != "Zu Wa" {
		t.Errorf("GetFullName mismatch: got %s, want Zu Wa", fullName)
	}

	// Test IsIndividual
	if !beneficiary.IsIndividual() {
		t.Error("IsIndividual should return true")
	}
	if beneficiary.IsCompany() {
		t.Error("IsCompany should return false")
	}

	// Test IsActive
	if !beneficiary.IsActive() {
		t.Error("IsActive should return true")
	}

	// Test IsLocalPayment
	if !beneficiary.IsLocalPayment() {
		t.Error("IsLocalPayment should return true")
	}
	if beneficiary.IsInternationalPayment() {
		t.Error("IsInternationalPayment should return false")
	}
}

// ============================================================================
// Company Beneficiary Tests
// ============================================================================

func TestParseBeneficiaryCompany(t *testing.T) {
	var event Event
	err := json.Unmarshal([]byte(beneficiaryCompanyWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse company beneficiary webhook: %v", err)
	}

	beneficiary, err := event.ParseBeneficiaryData()
	if err != nil {
		t.Fatalf("Failed to parse beneficiary data: %v", err)
	}

	// Verify company-specific fields
	if beneficiary.BeneficiaryEntityType != BeneficiaryEntityTypeCompany {
		t.Errorf("BeneficiaryEntityType mismatch: got %s", beneficiary.BeneficiaryEntityType)
	}
	if beneficiary.BeneficiaryCompanyName != "ACME Corp" {
		t.Errorf("BeneficiaryCompanyName mismatch: got %s", beneficiary.BeneficiaryCompanyName)
	}
	if beneficiary.BeneficiaryEmail != "payments@acme.com" {
		t.Errorf("BeneficiaryEmail mismatch: got %s", beneficiary.BeneficiaryEmail)
	}

	// Verify helper methods for company
	if beneficiary.IsIndividual() {
		t.Error("IsIndividual should return false for company")
	}
	if !beneficiary.IsCompany() {
		t.Error("IsCompany should return true")
	}

	// Verify international payment type
	if beneficiary.IsLocalPayment() {
		t.Error("IsLocalPayment should return false")
	}
	if !beneficiary.IsInternationalPayment() {
		t.Error("IsInternationalPayment should return true")
	}

	// GetFullName should return empty for company with no first/last name
	if beneficiary.GetFullName() != "" {
		t.Errorf("GetFullName should return empty for company without names, got %s", beneficiary.GetFullName())
	}

	// Verify bank details with routing codes
	bankDetails, err := beneficiary.GetBeneficiaryBankDetails()
	if err != nil {
		t.Fatalf("Failed to parse bank details: %v", err)
	}
	if bankDetails.RoutingCodeType1 != "ABA" {
		t.Errorf("RoutingCodeType1 mismatch: got %s", bankDetails.RoutingCodeType1)
	}
	if bankDetails.RoutingCodeValue1 != "026009593" {
		t.Errorf("RoutingCodeValue1 mismatch: got %s", bankDetails.RoutingCodeValue1)
	}
}

// ============================================================================
// Constants Tests
// ============================================================================

func TestBeneficiaryStatusConstants(t *testing.T) {
	if BeneficiaryStatusActive != "ACTIVE" {
		t.Errorf("BeneficiaryStatusActive mismatch")
	}
	if BeneficiaryStatusInactive != "INACTIVE" {
		t.Errorf("BeneficiaryStatusInactive mismatch")
	}
	if BeneficiaryStatusPending != "PENDING" {
		t.Errorf("BeneficiaryStatusPending mismatch")
	}
	if BeneficiaryStatusRejected != "REJECTED" {
		t.Errorf("BeneficiaryStatusRejected mismatch")
	}
}

func TestBeneficiaryEntityTypeConstants(t *testing.T) {
	if BeneficiaryEntityTypeIndividual != "INDIVIDUAL" {
		t.Errorf("BeneficiaryEntityTypeIndividual mismatch")
	}
	if BeneficiaryEntityTypeCompany != "COMPANY" {
		t.Errorf("BeneficiaryEntityTypeCompany mismatch")
	}
}

func TestPaymentTypeConstants(t *testing.T) {
	if PaymentTypeLocal != "LOCAL" {
		t.Errorf("PaymentTypeLocal mismatch")
	}
	if PaymentTypeInternational != "INTERNATIONAL" {
		t.Errorf("PaymentTypeInternational mismatch")
	}
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestParseBeneficiaryData_WrongEventType(t *testing.T) {
	wrongTypeJSON := `{
		"version": "V1.6.0",
		"event_name": "ONBOARDING",
		"event_type": "onboarding.account.create",
		"event_id": "test-id",
		"source_id": "test-source",
		"data": {}
	}`

	var event Event
	json.Unmarshal([]byte(wrongTypeJSON), &event)

	_, err := event.ParseBeneficiaryData()
	if err == nil {
		t.Error("ParseBeneficiaryData should fail for non-beneficiary event type")
	}
}

func TestGetBeneficiaryAddress_EmptyAddress(t *testing.T) {
	beneficiary := &BeneficiaryData{
		BeneficiaryAddressRaw: "",
	}

	address, err := beneficiary.GetBeneficiaryAddress()
	if err != nil {
		t.Errorf("GetBeneficiaryAddress should not error for empty address: %v", err)
	}
	if address != nil {
		t.Error("GetBeneficiaryAddress should return nil for empty address")
	}
}

func TestGetBeneficiaryBankDetails_EmptyDetails(t *testing.T) {
	beneficiary := &BeneficiaryData{
		BeneficiaryBankDetailsRaw: "",
	}

	details, err := beneficiary.GetBeneficiaryBankDetails()
	if err != nil {
		t.Errorf("GetBeneficiaryBankDetails should not error for empty details: %v", err)
	}
	if details != nil {
		t.Error("GetBeneficiaryBankDetails should return nil for empty details")
	}
}
