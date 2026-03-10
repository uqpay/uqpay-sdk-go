package connect

import (
	"context"
	"fmt"
	"net/url"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// AccountsClient handles account operations
type AccountsClient struct {
	client *common.APIClient
}

// EntityType represents the type of account entity
type EntityType string

const (
	EntityTypeIndividual EntityType = "INDIVIDUAL"
	EntityTypeCompany    EntityType = "COMPANY"
)

// ============================================================
// Shared Request Types
// ============================================================

// Address represents a physical address used in requests
type Address struct {
	Line1      string `json:"line1"`           // Street address, max 100 chars
	Line2      string `json:"line2,omitempty"` // Suite/apt/building, max 100 chars
	City       string `json:"city"`            // City name, max 100 chars
	State      string `json:"state,omitempty"` // State/province, max 2 chars
	PostalCode string `json:"postal_code"`     // ZIP/postal code, max 16 chars
	Country    string `json:"country"`         // ISO 3166-1 alpha-2 country code
}

// TosAcceptance represents terms of service acceptance
type TosAcceptance struct {
	Date      string `json:"date"`                 // Acceptance timestamp in ISO 8601 format
	IP        string `json:"ip"`                   // IPv4 address of the accepting user
	UserAgent string `json:"user_agent,omitempty"` // Browser user agent string
}

// ============================================================
// Create Account (Legacy) - POST /v1/accounts
// Deprecated: Use CreateSubAccount instead.
// ============================================================

// CreateAccountContactDetails represents contact details for the Create Account request.
type CreateAccountContactDetails struct {
	Email string `json:"email"` // Required. Email address
	Phone string `json:"phone"` // Required. Phone with country code
}

// MonthlyEstimatedRevenue represents a monthly revenue estimate.
type MonthlyEstimatedRevenue struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency,omitempty"`
}

// PersonIdentification represents identification in a Create Account request.
type PersonIdentification struct {
	Type     string `json:"type"`                // PASSPORT, DRIVERS_LICENSE, or NATIONAL_ID
	IDNumber string `json:"id_number,omitempty"` // ID document number
}

// CreateAccountPersonDetails represents person details for the Create Account request.
type CreateAccountPersonDetails struct {
	FirstNameEnglish        string                   `json:"first_name_english,omitempty"`
	LastNameEnglish         string                   `json:"last_name_english,omitempty"`
	FirstName               string                   `json:"first_name,omitempty"`
	LastName                string                   `json:"last_name,omitempty"`
	Nationality             string                   `json:"nationality,omitempty"`
	DateOfBirth             string                   `json:"date_of_birth,omitempty"`
	TaxNumber               string                   `json:"tax_number,omitempty"`
	Internationally         *int                     `json:"internationally,omitempty"`
	MonthlyEstimatedRevenue *MonthlyEstimatedRevenue `json:"monthly_estimated_revenue,omitempty"`
	AccountPurpose          []string                 `json:"account_purpose,omitempty"`
	OtherPurpose            string                   `json:"other_purpose,omitempty"`
	Identification          *PersonIdentification    `json:"identification,omitempty"`
}

// CreateAccountDocument represents a document in the Create Account request.
type CreateAccountDocument struct {
	Type  string `json:"type"`            // Document type
	Front string `json:"front,omitempty"` // Front image (base64 or URL)
	Back  string `json:"back,omitempty"`  // Back image (base64 or URL)
}

// CreateAccountRequest represents an account creation request (legacy endpoint).
// Deprecated: Use CreateSubAccountRequest instead.
type CreateAccountRequest struct {
	EntityType         EntityType                   `json:"entity_type"`                   // Required. INDIVIDUAL or COMPANY
	Name               string                       `json:"name"`                          // Required. Nickname displayed in Dashboard
	ContactDetails     *CreateAccountContactDetails `json:"contact_details"`               // Required. Contact information
	PersonDetails      *CreateAccountPersonDetails  `json:"person_details,omitempty"`      // Required for INDIVIDUAL
	ResidentialAddress *Address                     `json:"residential_address,omitempty"` // Required for INDIVIDUAL
	Documents          []CreateAccountDocument      `json:"documents,omitempty"`           // Required. Identity/KYC documents
	TosAcceptance      *TosAcceptance               `json:"tos_acceptance"`                // Required
}

// CreateAccountResponse represents the response from POST /v1/accounts.
type CreateAccountResponse struct {
	AccountID          string `json:"account_id"`          // Unique account identifier (UUID)
	ShortReferenceID   string `json:"short_reference_id"`  // Short reference ID
	Status             string `json:"status"`              // ACTIVE, PROCESSING, INACTIVE, or CLOSED
	VerificationStatus string `json:"verification_status"` // APPROVED, PENDING, REJECT, EXPIRED, or RETURN
}

// ============================================================
// Get Additional Documents - GET /v1/accounts/get_additional
// ============================================================

// AdditionalDocument represents a document profile returned by Get Additional Documents.
type AdditionalDocument struct {
	ProfileKey    string `json:"profile_key"`    // Unique key for the document type
	ProfileName   string `json:"profile_name"`   // Description of the file
	ProfileOption int    `json:"profile_option"` // 1 = required, 0 = optional
}

// ============================================================
// List Connected Accounts - GET /v1/accounts
// ============================================================

// ListAccountsRequest represents an accounts list request.
type ListAccountsRequest struct {
	PageSize   int `json:"page_size"`   // Required. Items per page, 10-100 (default 10)
	PageNumber int `json:"page_number"` // Required. Page number, >= 1 (default 1)
}

// ListAccountsResponse represents an accounts list response.
type ListAccountsResponse struct {
	TotalPages int                `json:"total_pages"` // Total number of pages
	TotalItems int                `json:"total_items"` // Total number of accounts
	Data       []ConnectedAccount `json:"data"`        // Array of account objects
}

// ============================================================
// Update Account - POST /v1/accounts/{id}
// Note: This endpoint has not been audited against API docs.
// ============================================================

// ContactDetails represents contact information for update requests.
type ContactDetails struct {
	Email       string `json:"email"`        // Email address
	PhoneNumber string `json:"phone_number"` // Phone number with country code
}

// IndividualDetails represents individual account details for update requests.
type IndividualDetails struct {
	FirstName     string         `json:"first_name"`
	LastName      string         `json:"last_name"`
	DateOfBirth   string         `json:"date_of_birth"`
	SSNLast4      string         `json:"ssn_last4,omitempty"`
	Address       Address        `json:"address"`
	ContactInfo   ContactDetails `json:"contact_info"`
	TosAcceptance TosAcceptance  `json:"tos_acceptance"`
}

// CompanyDetails represents company account details for update requests.
type CompanyDetails struct {
	LegalName       string           `json:"legal_name"`
	TaxID           string           `json:"tax_id,omitempty"`
	BusinessType    string           `json:"business_type"`
	Address         Address          `json:"address"`
	ContactInfo     ContactDetails   `json:"contact_info"`
	TosAcceptance   TosAcceptance    `json:"tos_acceptance"`
	Representatives []Representative `json:"representatives,omitempty"`
}

// Representative represents a company representative for update requests.
type Representative struct {
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	DateOfBirth string  `json:"date_of_birth"`
	Email       string  `json:"email"`
	Address     Address `json:"address"`
	SSNLast4    string  `json:"ssn_last4,omitempty"`
}

// UpdateAccountRequest represents an account update request.
// Note: This endpoint has not been audited against API docs.
type UpdateAccountRequest struct {
	Individual *IndividualDetails `json:"individual,omitempty"`
	Company    *CompanyDetails    `json:"company,omitempty"`
	Metadata   map[string]string  `json:"metadata,omitempty"`
}

// ============================================================
// Response Models (shared by List and Retrieve endpoints)
// ============================================================

// AccountContactDetails represents contact details in API responses.
type AccountContactDetails struct {
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"` // Phone with country code, e.g. "+6588880000"
}

// AccountIdentifier represents a tax/business identifier.
type AccountIdentifier struct {
	Type   string `json:"type,omitempty"`
	Number string `json:"number,omitempty"`
}

// AccountBusinessDetails represents business details in API responses.
type AccountBusinessDetails struct {
	LegalEntityName         string                   `json:"legal_entity_name,omitempty"`
	LegalEntityNameEnglish  string                   `json:"legal_entity_name_english,omitempty"`
	IncorporationDate       string                   `json:"incorporation_date,omitempty"`
	RegistrationNumber      string                   `json:"registration_number,omitempty"`
	BusinessStructure       string                   `json:"business_structure,omitempty"`
	ProductDescription      string                   `json:"product_description,omitempty"`
	MerchantCategoryCode    string                   `json:"merchant_category_code,omitempty"`
	EstimatedWorkerCount    string                   `json:"estimated_worker_count,omitempty"`
	MonthlyEstimatedRevenue *MonthlyEstimatedRevenue `json:"monthly_estimated_revenue,omitempty"`
	AccountPurpose          []string                 `json:"account_purpose,omitempty"`
	Identifier              *AccountIdentifier       `json:"identifier,omitempty"`
	WebsiteURL              string                   `json:"website_url,omitempty"`
}

// IdentificationDocuments represents document URLs/IDs for identification.
type IdentificationDocuments struct {
	Front       string `json:"front,omitempty"`
	FrontFileID string `json:"front_file_id,omitempty"`
	Back        string `json:"back,omitempty"`
	BackFileID  string `json:"back_file_id,omitempty"`
}

// AccountIdentification represents identification details in API responses.
type AccountIdentification struct {
	Type      string                   `json:"type,omitempty"`
	IDNumber  string                   `json:"id_number,omitempty"`
	Remark    string                   `json:"remark,omitempty"`
	Documents *IdentificationDocuments `json:"documents,omitempty"`
}

// AccountPersonDetails represents person details in API responses.
type AccountPersonDetails struct {
	FirstNameEnglish        string                   `json:"first_name_english,omitempty"`
	LastNameEnglish         string                   `json:"last_name_english,omitempty"`
	FirstName               string                   `json:"first_name,omitempty"`
	LastName                string                   `json:"last_name,omitempty"`
	Nationality             string                   `json:"nationality,omitempty"`
	DateOfBirth             string                   `json:"date_of_birth,omitempty"`
	TaxNumber               string                   `json:"tax_number,omitempty"`
	MonthlyEstimatedRevenue *MonthlyEstimatedRevenue `json:"monthly_estimated_revenue,omitempty"`
	OtherPurpose            string                   `json:"other_purpose,omitempty"`
	AccountPurpose          []string                 `json:"account_purpose,omitempty"`
	Internationally         *int                     `json:"internationally,omitempty"`
	Identification          *AccountIdentification   `json:"identification,omitempty"`
}

// AccountAddress represents an address in API responses.
type AccountAddress struct {
	City       string `json:"city,omitempty"`
	Country    string `json:"country,omitempty"`
	Line1      string `json:"line1,omitempty"`
	State      string `json:"state,omitempty"`
	Line2      string `json:"line2,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}

// RepresentativeOtherDocument represents a document attached to a representative.
type RepresentativeOtherDocument struct {
	Type  string `json:"type,omitempty"`
	Front string `json:"front,omitempty"`
}

// AccountRepresentative represents a representative in API responses.
type AccountRepresentative struct {
	RepresentativeID   string                        `json:"representative_id,omitempty"`
	Roles              string                        `json:"roles,omitempty"`
	FirstName          string                        `json:"first_name,omitempty"`
	LastName           string                        `json:"last_name,omitempty"`
	Nationality        string                        `json:"nationality,omitempty"`
	DateOfBirth        string                        `json:"date_of_birth,omitempty"`
	SharePercentage    string                        `json:"share_percentage,omitempty"`
	Identification     *AccountIdentification        `json:"identification,omitempty"`
	ResidentialAddress *AccountAddress               `json:"residential_address,omitempty"`
	AsApplicant        bool                          `json:"as_applicant,omitempty"`
	OtherDocuments     []RepresentativeOtherDocument `json:"other_doucments,omitempty"` // Note: API has typo "doucments"
}

// AccountDocument represents a document in API responses.
type AccountDocument struct {
	Type  string `json:"type,omitempty"`
	Front string `json:"front,omitempty"`
	Back  string `json:"back,omitempty"`
}

// ConnectedAccount represents an account as returned by the List and Retrieve API endpoints.
type ConnectedAccount struct {
	AccountID          string                  `json:"account_id"`
	ShortReferenceID   string                  `json:"short_reference_id,omitempty"`
	BusinessCode       []string                `json:"business_code,omitempty"`
	Email              string                  `json:"email,omitempty"`
	AccountName        string                  `json:"account_name,omitempty"`
	Country            string                  `json:"country,omitempty"`
	Status             string                  `json:"status"`
	VerificationStatus string                  `json:"verification_status,omitempty"`
	EntityType         EntityType              `json:"entity_type,omitempty"`
	ReviewReason       string                  `json:"review_reason,omitempty"`
	ContactDetails     *AccountContactDetails  `json:"contact_details,omitempty"`
	BusinessDetails    *AccountBusinessDetails `json:"business_details,omitempty"`
	PersonDetails      *AccountPersonDetails   `json:"person_details,omitempty"`
	// Fields present only in Retrieve Account response
	RegistrationAddress *AccountAddress         `json:"registration_address,omitempty"`
	BusinessAddress     []AccountAddress        `json:"business_address,omitempty"`
	Representatives     []AccountRepresentative `json:"representatives,omitempty"`
	Documents           []AccountDocument       `json:"documents,omitempty"`
	TosAcceptance       *TosAcceptance          `json:"tos_acceptance,omitempty"`
}

// ============================================================
// Methods
// ============================================================

// CreateSubAccount creates a new sub-account using the new API endpoint.
// For INDIVIDUAL accounts, populate IndividualInfo, IdentityVerification, ExpectedActivity, and ProofDocuments.
// For COMPANY accounts, populate CompanyInfo, CompanyAddress, OwnershipDetails, and BusinessDetails.
func (c *AccountsClient) CreateSubAccount(ctx context.Context, req *CreateSubAccountRequest) (*CreateSubAccountResponse, error) {
	if req.EntityType == EntityTypeIndividual {
		if req.IndividualInfo == nil {
			return nil, fmt.Errorf("individual_info required for INDIVIDUAL entity type")
		}
		if req.IdentityVerification == nil {
			return nil, fmt.Errorf("identity_verification required for INDIVIDUAL entity type")
		}
		if req.ExpectedActivity == nil {
			return nil, fmt.Errorf("expected_activity required for INDIVIDUAL entity type")
		}
		if req.ProofDocuments == nil {
			return nil, fmt.Errorf("proof_documents required for INDIVIDUAL entity type")
		}
	}
	if req.EntityType == EntityTypeCompany {
		if req.Inherit == nil || *req.Inherit != 1 {
			if req.CompanyInfo == nil {
				return nil, fmt.Errorf("company_info required for COMPANY entity type when inherit != 1")
			}
			if req.CompanyAddress == nil {
				return nil, fmt.Errorf("company_address required for COMPANY entity type when inherit != 1")
			}
			if req.OwnershipDetails == nil {
				return nil, fmt.Errorf("ownership_details required for COMPANY entity type when inherit != 1")
			}
		}
	}
	if req.TosAcceptance == nil {
		return nil, fmt.Errorf("tos_acceptance is required")
	}

	var resp CreateSubAccountResponse
	if err := c.client.Post(ctx, "/v1/accounts/create_accounts", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create sub-account: %w", err)
	}
	return &resp, nil
}

// GetAdditionalDocuments retrieves required/optional document types for a company sub-account
// based on the specified country and business code.
func (c *AccountsClient) GetAdditionalDocuments(ctx context.Context, country, businessCode string) ([]AdditionalDocument, error) {
	var resp []AdditionalDocument
	path := fmt.Sprintf("/v1/accounts/get_additional?country=%s&business_code=%s",
		url.QueryEscape(country), url.QueryEscape(businessCode))
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get additional documents: %w", err)
	}
	return resp, nil
}

// Create creates a new account using the legacy API endpoint.
// Deprecated: Use CreateSubAccount instead.
func (c *AccountsClient) Create(ctx context.Context, req *CreateAccountRequest) (*CreateAccountResponse, error) {
	var resp CreateAccountResponse
	if err := c.client.Post(ctx, "/v1/accounts", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}
	return &resp, nil
}

// List lists connected accounts with pagination.
func (c *AccountsClient) List(ctx context.Context, req *ListAccountsRequest) (*ListAccountsResponse, error) {
	var resp ListAccountsResponse
	path := fmt.Sprintf("/v1/accounts?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}
	return &resp, nil
}

// Update updates an existing account.
// Note: This endpoint has not been audited against API docs.
func (c *AccountsClient) Update(ctx context.Context, accountID string, req *UpdateAccountRequest) (*ConnectedAccount, error) {
	var account ConnectedAccount
	path := fmt.Sprintf("/v1/accounts/%s", accountID)
	if err := c.client.Post(ctx, path, req, &account); err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}
	return &account, nil
}

// Get retrieves an account by ID with an optional business code filter.
// Pass an empty string for businessCode to use the API default (BANKING).
func (c *AccountsClient) Get(ctx context.Context, accountID string, businessCode string) (*ConnectedAccount, error) {
	var account ConnectedAccount
	path := fmt.Sprintf("/v1/accounts/%s", accountID)
	if businessCode != "" {
		path += fmt.Sprintf("?business_code=%s", url.QueryEscape(businessCode))
	}
	if err := c.client.Get(ctx, path, &account); err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return &account, nil
}
