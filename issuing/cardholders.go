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

// DeliveryAddress represents a cardholder's delivery address
type DeliveryAddress struct {
	City       string `json:"city"`            // Required. City/district/suburb/town/village, max 128 chars
	Country    string `json:"country"`         // Required. ISO 3166-1 alpha-2 country code
	Line1      string `json:"line1"`           // Required. Street, PO Box, or company name, max 255 chars
	State      string `json:"state,omitempty"` // Optional. State/county/province/region, max 128 chars
	Line2      string `json:"line2,omitempty"` // Optional. Apartment/suite/unit/building, max 255 chars
	PostalCode string `json:"postal_code"`     // Required. ZIP or postal code, max 16 chars
}

// CreateCardholderRequest represents a cardholder creation request
type CreateCardholderRequest struct {
	Email           string           `json:"email"`                      // Required. Cardholder's email address
	FirstName       string           `json:"first_name"`                 // Required. Alphabetic and spaces only, 1-40 chars
	LastName        string           `json:"last_name"`                  // Required. Alphabetic and spaces only, 1-40 chars
	CountryCode     string           `json:"country_code"`               // Required. ISO 3166-1 alpha-2 code (e.g. "SG")
	PhoneNumber     string           `json:"phone_number"`               // Required. Cardholder's phone number
	DateOfBirth     string           `json:"date_of_birth,omitempty"`    // Optional. Format: yyyy-mm-dd
	DeliveryAddress *DeliveryAddress `json:"delivery_address,omitempty"` // Optional. Physical mailing address
	DocumentType    string           `json:"document_type,omitempty"`    // Optional. Allowed: pdf, png, jpg, jpeg
	Document        string           `json:"document,omitempty"`         // Optional. Base64-encoded ID document, max 2MB
}

// CreateCardholderResponse represents a cardholder creation/update response
type CreateCardholderResponse struct {
	CardholderID     string `json:"cardholder_id"`     // UUID. Unique cardholder identifier
	CardholderStatus string `json:"cardholder_status"` // FAILED, PENDING, SUCCESS, or INCOMPLETE
}

// UpdateCardholderRequest represents a cardholder update request (all fields optional)
type UpdateCardholderRequest struct {
	CountryCode     string           `json:"country_code,omitempty"`     // Optional. ISO 3166-1 alpha-2 code
	Email           string           `json:"email,omitempty"`            // Optional. Cardholder's email address
	PhoneNumber     string           `json:"phone_number,omitempty"`     // Optional. Cardholder's phone number
	DeliveryAddress *DeliveryAddress `json:"delivery_address,omitempty"` // Optional. Physical mailing address
	DocumentType    string           `json:"document_type,omitempty"`    // Optional. Allowed: pdf, png, jpg, jpeg
	Document        string           `json:"document,omitempty"`         // Optional. Base64-encoded ID document, max 2MB
	DateOfBirth     string           `json:"date_of_birth,omitempty"`    // Optional. Format: yyyy-mm-dd
}

// Cardholder represents a cardholder in list/retrieve responses
type Cardholder struct {
	CardholderID     string           `json:"cardholder_id"`              // UUID. Unique cardholder identifier
	Email            string           `json:"email"`                      // Cardholder's email address
	NumberOfCards    int              `json:"number_of_cards"`            // Total cards associated with this cardholder
	FirstName        string           `json:"first_name"`                 // Alphabetic and spaces only, 1-40 chars
	LastName         string           `json:"last_name"`                  // Alphabetic and spaces only, 1-40 chars
	CreateTime       string           `json:"create_time"`                // Format: YYYY-MM-DD HH:MM:SS
	CardholderStatus string           `json:"cardholder_status"`          // FAILED, PENDING, SUCCESS, or INCOMPLETE
	DateOfBirth      string           `json:"date_of_birth,omitempty"`    // Format: yyyy-mm-dd
	CountryCode      string           `json:"country_code"`               // ISO 3166-1 alpha-2 code
	PhoneNumber      string           `json:"phone_number"`               // Cardholder's phone number
	DeliveryAddress  *DeliveryAddress `json:"delivery_address,omitempty"` // Physical mailing address
	ReviewStatus     string           `json:"review_status,omitempty"`    // Reserved for future use
}

// ListCardholdersRequest represents a cardholder list request
type ListCardholdersRequest struct {
	PageSize         int    `json:"page_size"`                   // Required. Items per page, min 10, max 100
	PageNumber       int    `json:"page_number"`                 // Required. Page to retrieve, min 1
	CardholderStatus string `json:"cardholder_status,omitempty"` // Optional. Filter: PENDING, SUCCESS, INCOMPLETE, or FAILED
}

// ListCardholdersResponse represents a cardholder list response
type ListCardholdersResponse struct {
	TotalPages int          `json:"total_pages"` // Total available pages
	TotalItems int          `json:"total_items"` // Total count of cardholders
	Data       []Cardholder `json:"data"`        // List of cardholder objects
}

// Create creates a new cardholder
func (c *CardholdersClient) Create(ctx context.Context, req *CreateCardholderRequest) (*CreateCardholderResponse, error) {
	var resp CreateCardholderResponse
	if err := c.client.Post(ctx, "/v1/issuing/cardholders", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create cardholder: %w", err)
	}
	return &resp, nil
}

// Update updates an existing cardholder
func (c *CardholdersClient) Update(ctx context.Context, cardholderID string, req *UpdateCardholderRequest) (*CreateCardholderResponse, error) {
	var resp CreateCardholderResponse
	path := fmt.Sprintf("/v1/issuing/cardholders/%s", cardholderID)
	if err := c.client.Post(ctx, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to update cardholder: %w", err)
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

// List lists cardholders
func (c *CardholdersClient) List(ctx context.Context, req *ListCardholdersRequest) (*ListCardholdersResponse, error) {
	var resp ListCardholdersResponse
	path := fmt.Sprintf("/v1/issuing/cardholders?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)
	if req.CardholderStatus != "" {
		path += fmt.Sprintf("&cardholder_status=%s", req.CardholderStatus)
	}
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list cardholders: %w", err)
	}
	return &resp, nil
}
