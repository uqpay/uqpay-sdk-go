package connect

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/v4/common"
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

// Address represents a physical address
type Address struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// ContactDetails represents contact information
type ContactDetails struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

// TosAcceptance represents terms of service acceptance
type TosAcceptance struct {
	Date      string `json:"date"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent,omitempty"`
}

// IndividualDetails represents individual account details
type IndividualDetails struct {
	FirstName     string         `json:"first_name"`
	LastName      string         `json:"last_name"`
	DateOfBirth   string         `json:"date_of_birth"`
	SSNLast4      string         `json:"ssn_last4,omitempty"`
	Address       Address        `json:"address"`
	ContactInfo   ContactDetails `json:"contact_info"`
	TosAcceptance TosAcceptance  `json:"tos_acceptance"`
}

// CompanyDetails represents company account details
type CompanyDetails struct {
	LegalName       string           `json:"legal_name"`
	TaxID           string           `json:"tax_id,omitempty"`
	BusinessType    string           `json:"business_type"`
	Address         Address          `json:"address"`
	ContactInfo     ContactDetails   `json:"contact_info"`
	TosAcceptance   TosAcceptance    `json:"tos_acceptance"`
	Representatives []Representative `json:"representatives,omitempty"`
}

// Representative represents a company representative
type Representative struct {
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	DateOfBirth string  `json:"date_of_birth"`
	Email       string  `json:"email"`
	Address     Address `json:"address"`
	SSNLast4    string  `json:"ssn_last4,omitempty"`
}

// CreateAccountRequest represents an account creation request
// This struct handles the discriminated union for INDIVIDUAL vs COMPANY entity types
type CreateAccountRequest struct {
	EntityType EntityType         `json:"entity_type"`
	Individual *IndividualDetails `json:"individual,omitempty"`
	Company    *CompanyDetails    `json:"company,omitempty"`
	Metadata   map[string]string  `json:"metadata,omitempty"`
}

// Account represents a Connect account
type Account struct {
	// Optional response fields are shared by list summaries and account detail.
	ShortReferenceID    string                   `json:"short_reference_id,omitempty"`
	BusinessCode        []string                 `json:"business_code,omitempty"`
	Email               string                   `json:"email,omitempty"`
	AccountName         string                   `json:"account_name,omitempty"`
	Country             string                   `json:"country,omitempty"`
	VerificationStatus  string                   `json:"verification_status,omitempty"`
	ReviewReason        string                   `json:"review_reason,omitempty"`
	ContactDetails      map[string]interface{}   `json:"contact_details,omitempty"`
	BusinessDetails     map[string]interface{}   `json:"business_details,omitempty"`
	PersonDetails       map[string]interface{}   `json:"person_details,omitempty"`
	RegistrationAddress map[string]interface{}   `json:"registration_address,omitempty"`
	ResidentialAddress  map[string]interface{}   `json:"residential_address,omitempty"`
	BusinessAddress     []map[string]interface{} `json:"business_address,omitempty"`
	Representatives     []map[string]interface{} `json:"representatives,omitempty"`
	Documents           []map[string]interface{} `json:"documents,omitempty"`
	TosAcceptance       *TosAcceptance           `json:"tos_acceptance,omitempty"`

	AccountID      string               `json:"account_id"`
	EntityType     EntityType           `json:"entity_type"`
	Individual     *IndividualDetails   `json:"individual,omitempty"`
	Company        *CompanyDetails      `json:"company,omitempty"`
	Status         string               `json:"status"`
	PayoutsEnabled bool                 `json:"payouts_enabled"`
	ChargesEnabled bool                 `json:"charges_enabled"`
	Requirements   *AccountRequirements `json:"requirements,omitempty"`
	Metadata       map[string]string    `json:"metadata,omitempty"`
	CreateTime     string               `json:"create_time"`
	UpdateTime     string               `json:"update_time,omitempty"`
}

// AccountRequirements represents account verification requirements
type AccountRequirements struct {
	CurrentlyDue   []string `json:"currently_due,omitempty"`
	EventuallyDue  []string `json:"eventually_due,omitempty"`
	PastDue        []string `json:"past_due,omitempty"`
	Disabled       bool     `json:"disabled"`
	DisabledReason string   `json:"disabled_reason,omitempty"`
}

// ListAccountsRequest represents an accounts list request
type ListAccountsRequest struct {
	PageSize   int    `json:"page_size,omitempty"`
	PageNumber int    `json:"page_number,omitempty"`
	Status     string `json:"status,omitempty"`
}

// ListAccountsResponse represents an accounts list response
type ListAccountsResponse struct {
	TotalPages int       `json:"total_pages"`
	TotalItems int       `json:"total_items"`
	Data       []Account `json:"data"`
}

// UpdateAccountRequest represents an account update request
type UpdateAccountRequest struct {
	Individual *IndividualDetails `json:"individual,omitempty"`
	Company    *CompanyDetails    `json:"company,omitempty"`
	Metadata   map[string]string  `json:"metadata,omitempty"`
}

// RetrieveAccountResponse represents an account retrieval response
type RetrieveAccountResponse struct {
	Account
}

// AdditionalDocument represents a document type required or optional for company sub-account creation
type AdditionalDocument struct {
	ProfileKey    string `json:"profile_key"`    // Unique key representing the document type, e.g. "ARTICLES_OF_ASSOCIATION"
	ProfileName   string `json:"profile_name"`   // Human-readable description of the document
	ProfileOption int    `json:"profile_option"` // 1 = required, 0 = optional
}

// CreateSubAccount creates a new sub-account using the new API endpoint.
// For INDIVIDUAL accounts, populate IndividualInfo, IdentityVerification, ExpectedActivity, and ProofDocuments.
// For COMPANY accounts, populate CompanyInfo, CompanyAddress, OwnershipDetails, and BusinessDetails.
func (c *AccountsClient) CreateSubAccount(ctx context.Context, req *CreateSubAccountRequest, opts ...*common.RequestOptions) (*CreateSubAccountResponse, error) {
	if err := validateCreateSubAccountRequest(req); err != nil {
		return nil, err
	}

	var resp CreateSubAccountResponse
	if err := c.client.PostWithOptions(ctx, "/v1/accounts/create_accounts", req, &resp, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to create sub-account: %w", err)
	}
	return &resp, nil
}

func validateCreateSubAccountRequest(req *CreateSubAccountRequest) error {
	if req == nil {
		return fmt.Errorf("request is required")
	}
	if req.EntityType == EntityTypeIndividual {
		if req.IndividualInfo == nil {
			return fmt.Errorf("individual_info required for INDIVIDUAL entity type")
		}
		if req.IdentityVerification == nil {
			return fmt.Errorf("identity_verification required for INDIVIDUAL entity type")
		}
		if req.ExpectedActivity == nil {
			return fmt.Errorf("expected_activity required for INDIVIDUAL entity type")
		}
		if req.ProofDocuments == nil {
			return fmt.Errorf("proof_documents required for INDIVIDUAL entity type")
		}
	}
	if req.EntityType == EntityTypeCompany {
		if req.Inherit == nil || *req.Inherit != 1 {
			if req.CompanyInfo == nil {
				return fmt.Errorf("company_info required for COMPANY entity type when inherit != 1")
			}
			if req.CompanyAddress == nil {
				return fmt.Errorf("company_address required for COMPANY entity type when inherit != 1")
			}
			if req.OwnershipDetails == nil {
				return fmt.Errorf("ownership_details required for COMPANY entity type when inherit != 1")
			}
			if req.OwnershipDetails.Representatives == nil {
				return fmt.Errorf("ownership_details.representatives required for COMPANY entity type when inherit != 1")
			}
			for i, representative := range req.OwnershipDetails.Representatives {
				if representative.EmailAddress == "" {
					return fmt.Errorf("ownership_details.representatives[%d].email_address is required", i)
				}
				if representative.DateOfBirth == "" {
					return fmt.Errorf("ownership_details.representatives[%d].date_of_birth is required", i)
				}
				if representative.OwnershipPercentage == "" {
					return fmt.Errorf("ownership_details.representatives[%d].ownership_percentage is required", i)
				}
			}
			if req.BusinessDetails == nil {
				return fmt.Errorf("business_details required for COMPANY entity type when inherit != 1")
			}
			if req.BusinessDetails.AccountPurpose == nil {
				return fmt.Errorf("business_details.account_purpose is required")
			}
			if req.BusinessDetails.BankingCurrencies == nil {
				return fmt.Errorf("business_details.banking_currencies is required")
			}
			if req.BusinessDetails.BankingCountries == nil {
				return fmt.Errorf("business_details.banking_countries is required")
			}
			if req.BusinessDetails.ArticlesOfAssociation == nil {
				return fmt.Errorf("business_details.articles_of_association is required")
			}
		}
	}
	if req.TosAcceptance == nil {
		return fmt.Errorf("tos_acceptance is required")
	}
	return nil
}

// GetAdditionalDocuments retrieves the required and optional document types for creating
// a company-type sub-account based on the specified country and business code (e.g. "BANKING").
func (c *AccountsClient) GetAdditionalDocuments(ctx context.Context, country, businessCode string, opts ...*common.RequestOptions) ([]AdditionalDocument, error) {
	var resp []AdditionalDocument
	path := fmt.Sprintf("/v1/accounts/get_additional?country=%s&business_code=%s", country, businessCode)
	if err := c.client.GetWithOptions(ctx, path, &resp, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to get additional documents: %w", err)
	}
	return resp, nil
}

// Create creates a new account using the legacy API endpoint
func (c *AccountsClient) Create(ctx context.Context, req *CreateAccountRequest, opts ...*common.RequestOptions) (*Account, error) {
	// Validate discriminated union
	if req.EntityType == EntityTypeIndividual && req.Individual == nil {
		return nil, fmt.Errorf("individual details required for INDIVIDUAL entity type")
	}
	if req.EntityType == EntityTypeCompany && req.Company == nil {
		return nil, fmt.Errorf("company details required for COMPANY entity type")
	}

	var account Account
	if err := c.client.PostWithOptions(ctx, "/v1/accounts", req, &account, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}
	return &account, nil
}

// List lists accounts with optional filters
func (c *AccountsClient) List(ctx context.Context, req *ListAccountsRequest, opts ...*common.RequestOptions) (*ListAccountsResponse, error) {
	var resp ListAccountsResponse
	path := "/v1/accounts?"

	if req.PageSize > 0 {
		path += fmt.Sprintf("page_size=%d&", req.PageSize)
	}
	if req.PageNumber > 0 {
		path += fmt.Sprintf("page_number=%d&", req.PageNumber)
	}
	if req.Status != "" {
		path += fmt.Sprintf("status=%s&", req.Status)
	}

	// Remove trailing '?' or '&'
	if path[len(path)-1] == '?' || path[len(path)-1] == '&' {
		path = path[:len(path)-1]
	}

	if err := c.client.GetWithOptions(ctx, path, &resp, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}
	return &resp, nil
}

// Update updates an existing account
func (c *AccountsClient) Update(ctx context.Context, accountID string, req *UpdateAccountRequest, opts ...*common.RequestOptions) (*Account, error) {
	var account Account
	path := fmt.Sprintf("/v1/accounts/%s", accountID)
	if err := c.client.PostWithOptions(ctx, path, req, &account, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}
	return &account, nil
}

// Get retrieves an account by ID. An optional businessCode query parameter can be provided
// to filter by business type (e.g. "BANKING"). Omit or pass empty string to use the API default.
func (c *AccountsClient) Get(ctx context.Context, accountID string, businessCode ...string) (*Account, error) {
	return c.get(ctx, accountID, nil, businessCode...)
}

// GetWithOptions retrieves an account by ID with optional request headers.
// An optional businessCode query parameter can be provided after opts.
func (c *AccountsClient) GetWithOptions(ctx context.Context, accountID string, opts *common.RequestOptions, businessCode ...string) (*Account, error) {
	return c.get(ctx, accountID, opts, businessCode...)
}

func (c *AccountsClient) get(ctx context.Context, accountID string, opts *common.RequestOptions, businessCode ...string) (*Account, error) {
	var account Account
	path := fmt.Sprintf("/v1/accounts/%s", accountID)
	if len(businessCode) > 0 && businessCode[0] != "" {
		path += fmt.Sprintf("?business_code=%s", businessCode[0])
	}
	if err := c.client.GetWithOptions(ctx, path, &account, opts); err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return &account, nil
}
