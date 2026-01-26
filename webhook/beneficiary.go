package webhook

import "encoding/json"

// Beneficiary status constants
const (
	BeneficiaryStatusActive   = "ACTIVE"
	BeneficiaryStatusInactive = "INACTIVE"
	BeneficiaryStatusPending  = "PENDING"
	BeneficiaryStatusRejected = "REJECTED"
)

// Beneficiary entity type constants
const (
	BeneficiaryEntityTypeIndividual = "INDIVIDUAL"
	BeneficiaryEntityTypeCompany    = "COMPANY"
)

// Payment type constants
const (
	PaymentTypeLocal         = "LOCAL"
	PaymentTypeInternational = "INTERNATIONAL"
)

// BeneficiaryData represents the data payload for beneficiary webhook events.
type BeneficiaryData struct {
	// Account information
	AccountID           string `json:"account_id"`
	AccountNumber       string `json:"account_number"`
	AccountCurrencyCode string `json:"account_currency_code"`

	// Beneficiary identification
	BeneficiaryID        string `json:"beneficiary_id"`
	BeneficiaryFirstName string `json:"beneficiary_first_name"`
	BeneficiaryLastName  string `json:"beneficiary_last_name"`
	BeneficiaryNickname  string `json:"beneficiary_nickname"`

	// Beneficiary details
	BeneficiaryCompanyName string `json:"beneficiary_company_name"`
	BeneficiaryEmail       string `json:"beneficiary_email"`
	BeneficiaryEntityType  string `json:"beneficiary_entity_type"`
	BeneficiaryStatus      string `json:"beneficiary_status"`

	// Address and bank details (stored as JSON strings)
	BeneficiaryAddressRaw     string `json:"beneficiary_address"`
	BeneficiaryBankDetailsRaw string `json:"beneficiary_bank_details"`

	// Bank information
	BankCountryCode string `json:"bank_country_code"`

	// Payment information
	PaymentType      string `json:"payment_type"`
	ShortReferenceID string `json:"short_reference_id"`
	DirectID         string `json:"direct_id"`
}

// BeneficiaryAddress represents the parsed beneficiary address.
type BeneficiaryAddress struct {
	Nationality   string `json:"nationality"`
	CountryCode   string `json:"country_code"`
	City          string `json:"city"`
	State         string `json:"state"`
	StreetAddress string `json:"street_address"`
	PostalCode    string `json:"postal_code"`
}

// BeneficiaryBankDetails represents the parsed beneficiary bank details.
type BeneficiaryBankDetails struct {
	BankName            string `json:"bank_name"`
	BankAddress         string `json:"bank_address"`
	BankCountryCode     string `json:"bank_country_code"`
	AccountHolder       string `json:"account_holder"`
	AccountCurrencyCode string `json:"account_currency_code"`
	AccountNumber       string `json:"account_number"`
	IBAN                string `json:"iban"`
	SwiftCode           string `json:"swift_code"`
	ClearingSystem      string `json:"clearing_system"`
	RoutingCodeType1    string `json:"routing_code_type1"`
	RoutingCodeValue1   string `json:"routing_code_value1"`
	RoutingCodeType2    string `json:"routing_code_type2"`
	RoutingCodeValue2   string `json:"routing_code_value2"`
}

// GetBeneficiaryAddress parses and returns the beneficiary address.
// Returns nil if the address is empty or cannot be parsed.
func (b *BeneficiaryData) GetBeneficiaryAddress() (*BeneficiaryAddress, error) {
	if b.BeneficiaryAddressRaw == "" {
		return nil, nil
	}

	var address BeneficiaryAddress
	if err := json.Unmarshal([]byte(b.BeneficiaryAddressRaw), &address); err != nil {
		return nil, err
	}
	return &address, nil
}

// GetBeneficiaryBankDetails parses and returns the beneficiary bank details.
// Returns nil if the bank details are empty or cannot be parsed.
func (b *BeneficiaryData) GetBeneficiaryBankDetails() (*BeneficiaryBankDetails, error) {
	if b.BeneficiaryBankDetailsRaw == "" {
		return nil, nil
	}

	var details BeneficiaryBankDetails
	if err := json.Unmarshal([]byte(b.BeneficiaryBankDetailsRaw), &details); err != nil {
		return nil, err
	}
	return &details, nil
}

// GetFullName returns the beneficiary's full name (first + last).
func (b *BeneficiaryData) GetFullName() string {
	if b.BeneficiaryFirstName == "" && b.BeneficiaryLastName == "" {
		return ""
	}
	if b.BeneficiaryFirstName == "" {
		return b.BeneficiaryLastName
	}
	if b.BeneficiaryLastName == "" {
		return b.BeneficiaryFirstName
	}
	return b.BeneficiaryFirstName + " " + b.BeneficiaryLastName
}

// IsIndividual returns true if the beneficiary is an individual.
func (b *BeneficiaryData) IsIndividual() bool {
	return b.BeneficiaryEntityType == BeneficiaryEntityTypeIndividual
}

// IsCompany returns true if the beneficiary is a company.
func (b *BeneficiaryData) IsCompany() bool {
	return b.BeneficiaryEntityType == BeneficiaryEntityTypeCompany
}

// IsActive returns true if the beneficiary status is active.
func (b *BeneficiaryData) IsActive() bool {
	return b.BeneficiaryStatus == BeneficiaryStatusActive
}

// IsLocalPayment returns true if this is a local payment type.
func (b *BeneficiaryData) IsLocalPayment() bool {
	return b.PaymentType == PaymentTypeLocal
}

// IsInternationalPayment returns true if this is an international payment type.
func (b *BeneficiaryData) IsInternationalPayment() bool {
	return b.PaymentType == PaymentTypeInternational
}
