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
	StreetAddress string `json:"street_address"`          // required
	City          string `json:"city"`                    // required
	State         string `json:"state,omitempty"`         // optional, ISO 3166-2
	PostalCode    string `json:"postal_code,omitempty"`   // optional
	Country       string `json:"country"`                 // required, ISO 3166-1 alpha-2
	CountryCode   string `json:"country_code,omitempty"`  // optional, ISO 3166-1 alpha-2
	Nationality   string `json:"nationality,omitempty"`   // optional
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
	IBAN                string `json:"iban,omitempty"`                  // optional, SEPA specific
	SortCode            string `json:"sort_code,omitempty"`             // optional, UK specific
	BIC                 string `json:"bic,omitempty"`                   // optional, SEPA specific
}

// Beneficiary represents a beneficiary
type Beneficiary struct {
	BeneficiaryID string       `json:"beneficiary_id"`
	EntityType    string       `json:"entity_type"`    // INDIVIDUAL or COMPANY
	FirstName     string       `json:"first_name"`     // required if INDIVIDUAL
	LastName      string       `json:"last_name"`      // required if INDIVIDUAL
	CompanyName   string       `json:"company_name"`   // required if COMPANY
	Currency      string       `json:"currency"`       // required
	Country       string       `json:"country"`        // required, ISO 3166-1 alpha-2
	PaymentMethod string       `json:"payment_method"` // required
	BankDetails   *BankDetails `json:"bank_details"`   // required
	Address       *Address     `json:"address"`        // required
	Email         string       `json:"email,omitempty"`
	PhoneNumber   string       `json:"phone_number,omitempty"`
	Reference     string       `json:"reference,omitempty"`
	CreateTime    string       `json:"create_time"`
	UpdateTime    string       `json:"update_time"`
	Status        string       `json:"status"` // active, inactive, deleted
}

// BeneficiaryCreationRequest represents a beneficiary creation request
type BeneficiaryCreationRequest struct {
	EntityType    string       `json:"entity_type"`    // required: INDIVIDUAL or COMPANY
	FirstName     string       `json:"first_name"`     // required if INDIVIDUAL
	LastName      string       `json:"last_name"`      // required if INDIVIDUAL
	CompanyName   string       `json:"company_name"`   // required if COMPANY
	Currency      string       `json:"currency"`       // required
	Country       string       `json:"country"`        // required, ISO 3166-1 alpha-2
	PaymentMethod string       `json:"payment_method"` // required
	BankDetails   *BankDetails `json:"bank_details"`   // required
	Address       *Address     `json:"address"`        // required
	Email         string       `json:"email,omitempty"`
	PhoneNumber   string       `json:"phone_number,omitempty"`
	Reference     string       `json:"reference,omitempty"`
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

// BeneficiaryCheckRequest represents a beneficiary check request
type BeneficiaryCheckRequest struct {
	Currency      string       `json:"currency"`       // required
	Country       string       `json:"country"`        // required, ISO 3166-1 alpha-2
	PaymentMethod string       `json:"payment_method"` // required
	BankDetails   *BankDetails `json:"bank_details"`   // required
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
