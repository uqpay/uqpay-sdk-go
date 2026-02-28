package connect

import (
	"context"
	"fmt"

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

// GetAdditionalDocumentsRequest represents a request for additional documents
type GetAdditionalDocumentsRequest struct {
	AccountID string `json:"account_id"`
}

// Document represents a required document
type Document struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Status      string `json:"status,omitempty"`
}

// GetAdditionalDocumentsResponse represents additional documents response
type GetAdditionalDocumentsResponse struct {
	AccountID string     `json:"account_id"`
	Documents []Document `json:"documents"`
}

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

// GetAdditionalDocuments retrieves additional required documents for an account
func (c *AccountsClient) GetAdditionalDocuments(ctx context.Context, accountID string) (*GetAdditionalDocumentsResponse, error) {
	var resp GetAdditionalDocumentsResponse
	path := fmt.Sprintf("/v1/accounts/get_additional?account_id=%s", accountID)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get additional documents: %w", err)
	}
	return &resp, nil
}

// Create creates a new account using the legacy API endpoint
func (c *AccountsClient) Create(ctx context.Context, req *CreateAccountRequest) (*Account, error) {
	// Validate discriminated union
	if req.EntityType == EntityTypeIndividual && req.Individual == nil {
		return nil, fmt.Errorf("individual details required for INDIVIDUAL entity type")
	}
	if req.EntityType == EntityTypeCompany && req.Company == nil {
		return nil, fmt.Errorf("company details required for COMPANY entity type")
	}

	var account Account
	if err := c.client.Post(ctx, "/v1/accounts", req, &account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}
	return &account, nil
}

// List lists accounts with optional filters
func (c *AccountsClient) List(ctx context.Context, req *ListAccountsRequest) (*ListAccountsResponse, error) {
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

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}
	return &resp, nil
}

// Update updates an existing account
func (c *AccountsClient) Update(ctx context.Context, accountID string, req *UpdateAccountRequest) (*Account, error) {
	var account Account
	path := fmt.Sprintf("/v1/accounts/%s", accountID)
	if err := c.client.Post(ctx, path, req, &account); err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}
	return &account, nil
}

// Get retrieves an account by ID
func (c *AccountsClient) Get(ctx context.Context, accountID string) (*Account, error) {
	var account Account
	path := fmt.Sprintf("/v1/accounts/%s", accountID)
	if err := c.client.Get(ctx, path, &account); err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return &account, nil
}
