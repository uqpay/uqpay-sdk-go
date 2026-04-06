package issuing

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// CardholdersClient handles cardholder operations
type CardholdersClient struct {
	client *common.APIClient
}

// ResidentialAddress represents a cardholder's residential address
type ResidentialAddress struct {
	Country    string  `json:"country"`
	State      *string `json:"state,omitempty"`
	City       string  `json:"city"`
	District   *string `json:"district,omitempty"`
	Line1      string  `json:"line1"`
	Line2      *string `json:"line2,omitempty"`
	LineEn     *string `json:"line_en,omitempty"`
	PostalCode *string `json:"postal_code,omitempty"`
}

// Identity represents a cardholder's identity document
type Identity struct {
	Type      string  `json:"type"` // ID_CARD or PASSPORT
	Number    string  `json:"number"`
	FrontFile string  `json:"front_file"`
	BackFile  *string `json:"back_file,omitempty"`
	HandFile  *string `json:"hand_file,omitempty"`
}

// KycProof represents third-party KYC verification proof
type KycProof struct {
	Provider    string `json:"provider"`
	ReferenceID string `json:"reference_id"`
}

// KycVerification represents KYC verification information
type KycVerification struct {
	Method   string    `json:"method"` // THIRD_PARTY or SUMSUB_REDIRECT
	KycProof *KycProof `json:"kyc_proof,omitempty"`
}

// CreateCardholderRequest represents a cardholder creation request
type CreateCardholderRequest struct {
	Email               string              `json:"email"`
	PhoneNumber         string              `json:"phone_number"`
	FirstName           string              `json:"first_name"`
	LastName            string              `json:"last_name"`
	CountryCode         string              `json:"country_code"`
	DateOfBirth         *string             `json:"date_of_birth,omitempty"`
	Gender              *string             `json:"gender,omitempty"` // MALE or FEMALE
	Nationality         *string             `json:"nationality,omitempty"`
	ResidentialAddress  *ResidentialAddress `json:"residential_address,omitempty"`
	Identity            *Identity           `json:"identity,omitempty"`
	KycVerification     *KycVerification    `json:"kyc_verification,omitempty"`
	DocumentType        *string             `json:"document_type,omitempty"` // pdf, png, jpg, jpeg
	Document            *string             `json:"document,omitempty"`
}

// CreateCardholderResponse represents the response after creating a cardholder
type CreateCardholderResponse struct {
	CardholderID       string  `json:"cardholder_id"`
	CardholderStatus   string  `json:"cardholder_status"`
	IdvVerificationURL *string `json:"idv_verification_url,omitempty"`
	IdvURLExpiresAt    *string `json:"idv_url_expires_at,omitempty"`
}

// Cardholder represents a full cardholder object (returned by Get and List)
type Cardholder struct {
	CardholderID        string              `json:"cardholder_id"`
	Email               string              `json:"email"`
	FirstName           string              `json:"first_name"`
	LastName            string              `json:"last_name"`
	CountryCode         string              `json:"country_code"`
	CardholderStatus    string              `json:"cardholder_status"`
	CreateTime          string              `json:"create_time"`
	NumberOfCards       int                 `json:"number_of_cards"`
	DateOfBirth         *string             `json:"date_of_birth,omitempty"`
	PhoneNumber         *string             `json:"phone_number,omitempty"`
	Gender              *string             `json:"gender,omitempty"`
	Nationality         *string             `json:"nationality,omitempty"`
	ResidentialAddress  *ResidentialAddress `json:"residential_address,omitempty"`
	ReviewStatus        *string             `json:"review_status,omitempty"`
	IdvStatus           *string             `json:"idv_status,omitempty"`
	IdvVerificationURL  *string             `json:"idv_verification_url,omitempty"`
	IdvURLExpiresAt     *string             `json:"idv_url_expires_at,omitempty"`
}

// UpdateCardholderRequest represents a cardholder update request
// Note: first_name and last_name cannot be updated
type UpdateCardholderRequest struct {
	CountryCode        *string             `json:"country_code,omitempty"`
	Email              *string             `json:"email,omitempty"`
	PhoneNumber        *string             `json:"phone_number,omitempty"`
	DateOfBirth        *string             `json:"date_of_birth,omitempty"`
	Gender             *string             `json:"gender,omitempty"` // MALE or FEMALE
	Nationality        *string             `json:"nationality,omitempty"`
	ResidentialAddress *ResidentialAddress `json:"residential_address,omitempty"`
	Identity           *Identity           `json:"identity,omitempty"`
	KycVerification    *KycVerification    `json:"kyc_verification,omitempty"`
	DocumentType       *string             `json:"document_type,omitempty"` // pdf, png, jpg, jpeg
	Document           *string             `json:"document,omitempty"`
}

// UpdateCardholderResponse represents the response after updating a cardholder
type UpdateCardholderResponse struct {
	CardholderID       string  `json:"cardholder_id"`
	CardholderStatus   string  `json:"cardholder_status"`
	IdvVerificationURL *string `json:"idv_verification_url,omitempty"`
	IdvURLExpiresAt    *string `json:"idv_url_expires_at,omitempty"`
}

// ListCardholdersRequest represents a cardholder list request
type ListCardholdersRequest struct {
	PageSize   int `json:"page_size"`
	PageNumber int `json:"page_number"`
}

// ListCardholdersResponse represents a cardholder list response
type ListCardholdersResponse struct {
	TotalPages int          `json:"total_pages"`
	TotalItems int          `json:"total_items"`
	Data       []Cardholder `json:"data"`
}

// Create creates a new cardholder
func (c *CardholdersClient) Create(ctx context.Context, req *CreateCardholderRequest) (*CreateCardholderResponse, error) {
	var resp CreateCardholderResponse
	if err := c.client.Post(ctx, "/v1/issuing/cardholders", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create cardholder: %w", err)
	}
	return &resp, nil
}

// Get retrieves a cardholder by ID
func (c *CardholdersClient) Get(ctx context.Context, cardholderID string) (*Cardholder, error) {
	var cardholder Cardholder
	path := fmt.Sprintf("/v1/issuing/cardholders/%s", cardholderID)
	if err := c.client.Get(ctx, path, &cardholder); err != nil {
		return nil, fmt.Errorf("failed to get cardholder: %w", err)
	}
	return &cardholder, nil
}

// Update updates the specified cardholder
// Note: first_name and last_name cannot be updated
func (c *CardholdersClient) Update(ctx context.Context, cardholderID string, req *UpdateCardholderRequest) (*UpdateCardholderResponse, error) {
	var resp UpdateCardholderResponse
	path := fmt.Sprintf("/v1/issuing/cardholders/%s", cardholderID)
	if err := c.client.Post(ctx, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to update cardholder: %w", err)
	}
	return &resp, nil
}

// List lists cardholders
func (c *CardholdersClient) List(ctx context.Context, req *ListCardholdersRequest) (*ListCardholdersResponse, error) {
	var resp ListCardholdersResponse
	path := fmt.Sprintf("/v1/issuing/cardholders?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list cardholders: %w", err)
	}
	return &resp, nil
}
