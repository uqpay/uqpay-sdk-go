package banking

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// BeneficiariesClient handles beneficiary operations
type BeneficiariesClient struct {
	client *common.APIClient
}

// Address represents a beneficiary address
type Address struct {
	StreetAddress string `json:"street_address"`        // required, max 255 chars
	City          string `json:"city"`                  // required, max 36 chars
	State         string `json:"state"`                 // required, max 96 chars
	PostalCode    string `json:"postal_code"`           // required, max 12 chars
	Country       string `json:"country"`               // required, ISO 3166-1 alpha-2 two-letter code
	Nationality   string `json:"nationality,omitempty"` // optional, ISO 3166-1 alpha-2 two-letter code
}

// BankDetails represents beneficiary bank account details
type BankDetails struct {
	AccountNumber       string `json:"account_number"`                // conditional, max 60 chars, alphanumeric only; required if IBAN absent
	AccountHolder       string `json:"account_holder"`                // required, max 240 chars
	AccountCurrencyCode string `json:"account_currency_code"`         // required, ISO 4217 three-letter currency code
	BankName            string `json:"bank_name"`                     // required, max 240 chars
	BankAddress         string `json:"bank_address"`                  // required, max 240 chars
	BankCountryCode     string `json:"bank_country_code"`             // required, ISO 3166-1 alpha-2 two-letter code
	SwiftCode           string `json:"swift_code"`                    // required, max 30 chars
	ClearingSystem      string `json:"clearing_system"`               // required, e.g. ACH, Fedwire, SWIFT, FAST, GIRO, RTGS, PayNow, FPS, EFT
	RoutingCodeType1    string `json:"routing_code_type1,omitempty"`  // conditional, e.g. ach, aba, bank_code, sort_code, bsb_code, ifsc, cnaps_number
	RoutingCodeValue1   string `json:"routing_code_value1,omitempty"` // conditional, max 48 chars; required if routing_code_type1 is set
	RoutingCodeType2    string `json:"routing_code_type2,omitempty"`  // conditional, e.g. branch_code
	RoutingCodeValue2   string `json:"routing_code_value2,omitempty"` // conditional, max 48 chars; required if routing_code_type2 is set
	IBAN                string `json:"iban,omitempty"`                // conditional, max 36 chars; mandatory for European countries
}

// AdditionalInfo represents extra beneficiary information
type AdditionalInfo struct {
	OrganizationCode string `json:"organization_code,omitempty"` // optional, Unified Social Credit Identifier for Mainland China entities
	ProxyID          string `json:"proxy_id,omitempty"`          // optional, PayNow proxy identifier (UEN, phone, or VPA) for SGD
	IDType           string `json:"id_type,omitempty"`           // conditional, PASSPORT | NATIONAL_ID | DRIVERS_LICENSE; required for COP/INDIVIDUAL
	IDNumber         string `json:"id_number,omitempty"`         // conditional, ID document number; required for COP/INDIVIDUAL
	TaxID            string `json:"tax_id,omitempty"`            // conditional, tax registration number; required for COP/COMPANY
	MSISDN           string `json:"msisdn,omitempty"`            // conditional, mobile number in +[country_code][number] format; required for COP, HKD/LOCAL
}

// Beneficiary represents a beneficiary
type Beneficiary struct {
	BeneficiaryID    string          `json:"beneficiary_id"`            // UUID v4 unique identifier
	EntityType       string          `json:"entity_type"`               // INDIVIDUAL | COMPANY
	FirstName        string          `json:"first_name,omitempty"`      // present if INDIVIDUAL, max 45 chars
	LastName         string          `json:"last_name,omitempty"`       // present if INDIVIDUAL, max 45 chars
	CompanyName      string          `json:"company_name,omitempty"`    // present if COMPANY, max 120 chars
	IDNumber         string          `json:"id_number,omitempty"`       // Mainland China resident ID; mandatory for CNH & LOCAL
	Nickname         string          `json:"nickname,omitempty"`        // optional, max 120 chars
	PaymentMethod    string          `json:"payment_method"`            // LOCAL | SWIFT
	BankDetails      *BankDetails    `json:"bank_details"`              // bank account details
	Address          *Address        `json:"address"`                   // beneficiary address
	AdditionalInfo   *AdditionalInfo `json:"additional_info,omitempty"` // supplementary data
	Email            string          `json:"email,omitempty"`           // contact email address
	ShortReferenceID string          `json:"short_reference_id"`        // system-generated reference identifier
	Summary          string          `json:"summary,omitempty"`         // beneficiary summary
	CreateTime       string          `json:"create_time"`               // ISO 8601 timestamp
	UpdateTime       string          `json:"update_time"`               // ISO 8601 timestamp
	Status           string          `json:"beneficiary_status"`        // ACTIVE | PENDING
}

// BeneficiaryCreationRequest represents a beneficiary creation request
type BeneficiaryCreationRequest struct {
	EntityType     string          `json:"entity_type"`               // required, INDIVIDUAL | COMPANY
	FirstName      string          `json:"first_name,omitempty"`      // required if INDIVIDUAL, max 45 chars
	LastName       string          `json:"last_name,omitempty"`       // required if INDIVIDUAL, max 45 chars
	CompanyName    string          `json:"company_name,omitempty"`    // required if COMPANY, max 120 chars
	IDNumber       string          `json:"id_number,omitempty"`       // conditional, mandatory for Mainland China residents with CNH & LOCAL
	Nickname       string          `json:"nickname,omitempty"`        // optional, max 120 chars
	PaymentMethod  string          `json:"payment_method"`            // required, LOCAL | SWIFT
	BankDetails    *BankDetails    `json:"bank_details"`              // required
	Address        *Address        `json:"address"`                   // required
	AdditionalInfo *AdditionalInfo `json:"additional_info,omitempty"` // optional
	Email          string          `json:"email,omitempty"`           // optional
}

// BeneficiaryCreationResponse represents a beneficiary creation response
type BeneficiaryCreationResponse struct {
	BeneficiaryID    string `json:"beneficiary_id"`               // UUID v4 unique identifier
	ShortReferenceID string `json:"short_reference_id,omitempty"` // system-generated reference identifier
}

// ListBeneficiariesRequest represents a beneficiary list request
type ListBeneficiariesRequest struct {
	PageSize    int    `json:"page_size"`              // required, min 1, max 100
	PageNumber  int    `json:"page_number"`            // required, min 1
	Currency    string `json:"currency,omitempty"`     // optional, ISO 4217 three-letter currency code
	EntityType  string `json:"entity_type,omitempty"`  // optional, INDIVIDUAL | COMPANY
	Nickname    string `json:"nickname,omitempty"`     // optional, max 120 chars
	CompanyName string `json:"company_name,omitempty"` // optional, max 120 chars
}

// ListBeneficiariesResponse represents a beneficiary list response
type ListBeneficiariesResponse struct {
	TotalPages int           `json:"total_pages"` // total number of pages available
	TotalItems int           `json:"total_items"` // total count of available items
	Data       []Beneficiary `json:"data"`        // array of beneficiary objects
}

// BeneficiaryCheckAdditionalInfo represents additional info for beneficiary check
type BeneficiaryCheckAdditionalInfo struct {
	ProxyID string `json:"proxy_id,omitempty"` // PayNow proxy identifier (SGD), supports UEN, phone number, or VPA
}

// BeneficiaryCheckRequest represents a beneficiary check request
type BeneficiaryCheckRequest struct {
	EntityType      string                          `json:"entity_type"`                 // required, INDIVIDUAL | COMPANY
	PaymentMethod   string                          `json:"payment_method"`              // required, LOCAL | SWIFT
	AccountNumber   string                          `json:"account_number"`              // required, max 60 chars, alphanumeric only
	Currency        string                          `json:"currency"`                    // required, ISO 4217 three-letter currency code
	BankCountryCode string                          `json:"bank_country_code,omitempty"` // optional, ISO 3166-1 alpha-2 two-letter code
	FirstName       string                          `json:"first_name,omitempty"`        // optional, for INDIVIDUAL beneficiaries, max 45 chars
	LastName        string                          `json:"last_name,omitempty"`         // optional, for INDIVIDUAL beneficiaries, max 45 chars
	CompanyName     string                          `json:"company_name,omitempty"`      // optional, for COMPANY beneficiaries, max 120 chars
	ClearingSystem  string                          `json:"clearing_system,omitempty"`   // optional, e.g. LOCAL, SWIFT, ACH, FAST, GIRO, Fedwire
	IBAN            string                          `json:"iban,omitempty"`              // optional, max 36 chars; mandatory for specified European countries
	AdditionalInfo  *BeneficiaryCheckAdditionalInfo `json:"additional_info,omitempty"`   // optional
}

// PaymentMethod represents an available payment method
type PaymentMethod struct {
	ClearingSystems string `json:"clearing_systems"` // e.g. ACH, FAST, GIRO, Fedwire, SWIFT, FPS, EFT, CHAPS, NPP, BPAY
	Currency        string `json:"currency"`         // ISO 4217 three-letter currency code
	Country         string `json:"country"`          // ISO 3166-1 alpha-2 two-letter country code
	PaymentMethod   string `json:"payment_method"`   // LOCAL | SWIFT
}

// ListPaymentMethodsResponse wraps the payment methods list response
type ListPaymentMethodsResponse struct {
	Data []PaymentMethod `json:"data"`
}

// Create creates a new beneficiary
func (c *BeneficiariesClient) Create(ctx context.Context, req *BeneficiaryCreationRequest) (*BeneficiaryCreationResponse, error) {
	var resp BeneficiaryCreationResponse
	if err := c.client.Post(ctx, "/v1/beneficiaries", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create beneficiary: %w", err)
	}
	return &resp, nil
}

// List lists beneficiaries with optional filters
func (c *BeneficiariesClient) List(ctx context.Context, req *ListBeneficiariesRequest) (*ListBeneficiariesResponse, error) {
	var resp ListBeneficiariesResponse
	path := fmt.Sprintf("/v1/beneficiaries?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)

	if req.Currency != "" {
		path += fmt.Sprintf("&currency=%s", req.Currency)
	}
	if req.EntityType != "" {
		path += fmt.Sprintf("&entity_type=%s", req.EntityType)
	}
	if req.Nickname != "" {
		path += fmt.Sprintf("&nickname=%s", req.Nickname)
	}
	if req.CompanyName != "" {
		path += fmt.Sprintf("&company_name=%s", req.CompanyName)
	}

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list beneficiaries: %w", err)
	}
	return &resp, nil
}

// Get retrieves a specific beneficiary by ID
func (c *BeneficiariesClient) Get(ctx context.Context, beneficiaryID string) (*Beneficiary, error) {
	var resp Beneficiary
	path := fmt.Sprintf("/v1/beneficiaries/%s", beneficiaryID)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get beneficiary: %w", err)
	}
	return &resp, nil
}

// Update updates an existing beneficiary
func (c *BeneficiariesClient) Update(ctx context.Context, beneficiaryID string, req *BeneficiaryCreationRequest) (*BeneficiaryCreationResponse, error) {
	var resp BeneficiaryCreationResponse
	path := fmt.Sprintf("/v1/beneficiaries/%s", beneficiaryID)
	if err := c.client.Post(ctx, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to update beneficiary: %w", err)
	}
	return &resp, nil
}

// Delete deletes a beneficiary
func (c *BeneficiariesClient) Delete(ctx context.Context, beneficiaryID string) error {
	path := fmt.Sprintf("/v1/beneficiaries/%s/delete", beneficiaryID)
	if err := c.client.Post(ctx, path, nil, nil); err != nil {
		return fmt.Errorf("failed to delete beneficiary: %w", err)
	}
	return nil
}

// ListPaymentMethods retrieves available payment methods for a currency and country
func (c *BeneficiariesClient) ListPaymentMethods(ctx context.Context, currency, country string) ([]PaymentMethod, error) {
	var resp ListPaymentMethodsResponse
	path := fmt.Sprintf("/v1/beneficiaries/paymentmethods?currency=%s&country=%s", currency, country)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list payment methods: %w", err)
	}
	return resp.Data, nil
}

// Check validates beneficiary details before creation
func (c *BeneficiariesClient) Check(ctx context.Context, req *BeneficiaryCheckRequest) (*Beneficiary, error) {
	var resp Beneficiary
	if err := c.client.Post(ctx, "/v1/beneficiaries/check", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to check beneficiary: %w", err)
	}
	return &resp, nil
}
