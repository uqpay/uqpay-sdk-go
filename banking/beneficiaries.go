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
	StreetAddress string `json:"street_address"`         // required
	City          string `json:"city"`                   // required
	State         string `json:"state"`                  // required
	PostalCode    string `json:"postal_code"`            // required
	Country       string `json:"country"`                // required, ISO 3166-1 alpha-2
	Nationality   string `json:"nationality,omitempty"`  // optional, ISO 3166-1 alpha-2
}

// BankDetails represents beneficiary bank account details
type BankDetails struct {
	AccountNumber       string `json:"account_number"`                  // required
	AccountHolder       string `json:"account_holder"`                  // required
	AccountCurrencyCode string `json:"account_currency_code"`           // required, ISO 4217
	BankName            string `json:"bank_name"`                       // required
	BankAddress         string `json:"bank_address"`                    // required
	BankCountryCode     string `json:"bank_country_code"`               // required, ISO 3166-1 alpha-2
	SwiftCode           string `json:"swift_code"`                      // required
	ClearingSystem      string `json:"clearing_system"`                 // required, e.g. ACH, FEDWIRE, FASTER_PAYMENTS, SEPA
	RoutingCodeType1    string `json:"routing_code_type1,omitempty"`    // e.g. ach, sort_code, iban
	RoutingCodeValue1   string `json:"routing_code_value1,omitempty"`   // routing code value
	RoutingCodeType2    string `json:"routing_code_type2,omitempty"`    // optional second routing code
	RoutingCodeValue2   string `json:"routing_code_value2,omitempty"`   // optional second routing value
	IBAN                string `json:"iban,omitempty"`                  // conditional, required for European countries
}

// AdditionalInfo represents extra beneficiary information
type AdditionalInfo struct {
	OrganizationCode string `json:"organization_code,omitempty"` // Unified Social Credit Identifier for China
	ProxyID          string `json:"proxy_id,omitempty"`          // PayNow proxy identifier for SGD
	IDType           string `json:"id_type,omitempty"`           // PASSPORT, NATIONAL_ID, or DRIVERS_LICENSE
	IDNumber         string `json:"id_number,omitempty"`         // identification number for individuals
	TaxID            string `json:"tax_id,omitempty"`            // tax ID for company beneficiaries
	MSISDN           string `json:"msisdn,omitempty"`            // mobile number in +[code][number] format
}

// Beneficiary represents a beneficiary
type Beneficiary struct {
	BeneficiaryID  string          `json:"beneficiary_id"`
	EntityType     string          `json:"entity_type"`              // INDIVIDUAL or COMPANY
	FirstName      string          `json:"first_name,omitempty"`     // present if INDIVIDUAL
	LastName       string          `json:"last_name,omitempty"`      // present if INDIVIDUAL
	CompanyName    string          `json:"company_name,omitempty"`   // present if COMPANY
	IDNumber       string          `json:"id_number,omitempty"`      // present when account currency = COP
	Nickname       string          `json:"nickname,omitempty"`
	PaymentMethod  string          `json:"payment_method"`
	BankDetails    *BankDetails    `json:"bank_details"`
	Address        *Address        `json:"address"`
	AdditionalInfo *AdditionalInfo `json:"additional_info,omitempty"`
	Email          string          `json:"email,omitempty"`
	CreateTime     string          `json:"created_time"`
	UpdateTime     string          `json:"updated_time"`
	Status         string          `json:"status"` // active, inactive, deleted
}

// BeneficiaryCreationRequest represents a beneficiary creation request
type BeneficiaryCreationRequest struct {
	EntityType     string          `json:"entity_type"`                // required: INDIVIDUAL or COMPANY
	FirstName      string          `json:"first_name,omitempty"`      // required if INDIVIDUAL
	LastName       string          `json:"last_name,omitempty"`       // required if INDIVIDUAL
	CompanyName    string          `json:"company_name,omitempty"`    // required if COMPANY
	IDNumber       string          `json:"id_number,omitempty"`       // required when account currency = COP
	Nickname       string          `json:"nickname,omitempty"`        // optional, max 120 chars
	Currency       string          `json:"currency,omitempty"`        // optional
	Country        string          `json:"country,omitempty"`         // optional, ISO 3166-1 alpha-2
	PaymentMethod  string          `json:"payment_method"`            // required: LOCAL or SWIFT
	BankDetails    *BankDetails    `json:"bank_details"`              // required
	Address        *Address        `json:"address"`                   // required
	AdditionalInfo *AdditionalInfo `json:"additional_info,omitempty"` // optional
	Email          string          `json:"email,omitempty"`           // optional
}

// BeneficiaryCreationResponse represents a beneficiary creation response
type BeneficiaryCreationResponse struct {
	BeneficiaryID    string `json:"beneficiary_id"`
	ShortReferenceID string `json:"short_reference_id,omitempty"`
	Status           string `json:"status,omitempty"`
}

// ListBeneficiariesRequest represents a beneficiary list request
type ListBeneficiariesRequest struct {
	PageSize   int    `json:"page_size"`             // required, 10-100
	PageNumber int    `json:"page_number"`           // required, >=1
	Currency   string `json:"currency,omitempty"`    // optional
	Country    string `json:"country,omitempty"`     // optional, ISO 3166-1 alpha-2
	Status     string `json:"status,omitempty"`      // optional: active, inactive, deleted
	EntityType string `json:"entity_type,omitempty"` // optional: INDIVIDUAL, COMPANY
}

// ListBeneficiariesResponse represents a beneficiary list response
type ListBeneficiariesResponse struct {
	TotalPages int           `json:"total_pages"`
	TotalItems int           `json:"total_items"`
	Data       []Beneficiary `json:"data"`
}

// BeneficiaryCheckAdditionalInfo represents additional info for beneficiary check
type BeneficiaryCheckAdditionalInfo struct {
	ProxyID string `json:"proxy_id,omitempty"` // PayNow proxy identifier (SGD), supports UEN, phone number, or VPA
}

// BeneficiaryCheckRequest represents a beneficiary check request
type BeneficiaryCheckRequest struct {
	EntityType     string                           `json:"entity_type"`                // required: INDIVIDUAL or COMPANY
	PaymentMethod  string                           `json:"payment_method"`             // required: LOCAL or SWIFT
	AccountNumber  string                           `json:"account_number"`             // required
	Currency       string                           `json:"currency"`                   // required, ISO 4217
	FirstName      string                           `json:"first_name,omitempty"`       // required if INDIVIDUAL
	LastName       string                           `json:"last_name,omitempty"`        // required if INDIVIDUAL
	CompanyName    string                           `json:"company_name,omitempty"`     // required if COMPANY
	ClearingSystem string                           `json:"clearing_system,omitempty"`  // e.g. LOCAL, ACH, FEDWIRE
	IBAN           string                           `json:"iban,omitempty"`             // conditional, e.g. country code for validation
	AdditionalInfo *BeneficiaryCheckAdditionalInfo  `json:"additional_info,omitempty"`  // optional
}

// PaymentMethod represents an available payment method
type PaymentMethod struct {
	ClearingSystems string   `json:"clearing_systems"`
	Currency        string   `json:"currency"`
	Country         string   `json:"country"`
	PaymentMethod   string   `json:"payment_method"`
	ValidationField []string `json:"validation_field"`
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
	if req.Country != "" {
		path += fmt.Sprintf("&country=%s", req.Country)
	}
	if req.Status != "" {
		path += fmt.Sprintf("&status=%s", req.Status)
	}
	if req.EntityType != "" {
		path += fmt.Sprintf("&entity_type=%s", req.EntityType)
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
func (c *BeneficiariesClient) Update(ctx context.Context, beneficiaryID string, req *BeneficiaryCreationRequest) (*Beneficiary, error) {
	var resp Beneficiary
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
