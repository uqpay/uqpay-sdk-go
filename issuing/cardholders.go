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
	City       string `json:"city"`
	Country    string `json:"country"` // ISO 3166-1 alpha-2
	Line1      string `json:"line1"`
	State      string `json:"state,omitempty"`
	Line2      string `json:"line2,omitempty"`
	PostalCode string `json:"postal_code"`
}

// CreateCardholderRequest represents a cardholder creation request
type CreateCardholderRequest struct {
	Email           string           `json:"email"`
	FirstName       string           `json:"first_name"`
	LastName        string           `json:"last_name"`
	CountryCode     string           `json:"country_code"`
	PhoneNumber     string           `json:"phone_number"`
	DateOfBirth     string           `json:"date_of_birth,omitempty"`
	DeliveryAddress *DeliveryAddress `json:"delivery_address,omitempty"`
	DocumentType    string           `json:"document_type,omitempty"` // file format: pdf, png, jpg, jpeg
	Document        string           `json:"document,omitempty"`      // base64 encoded document content
}

// CreateCardholderResponse represents a cardholder creation/update response
type CreateCardholderResponse struct {
	CardholderID     string `json:"cardholder_id"`
	CardholderStatus string `json:"cardholder_status"`
}

// UpdateCardholderRequest represents a cardholder update request
type UpdateCardholderRequest struct {
	CountryCode     string           `json:"country_code,omitempty"`
	Email           string           `json:"email,omitempty"`
	PhoneNumber     string           `json:"phone_number,omitempty"`
	DeliveryAddress *DeliveryAddress `json:"delivery_address,omitempty"`
	DocumentType    string           `json:"document_type,omitempty"` // file format: pdf, png, jpg, jpeg
	Document        string           `json:"document,omitempty"`      // base64 encoded document content
	DateOfBirth     string           `json:"date_of_birth,omitempty"`
}

// Cardholder represents a cardholder in list/retrieve responses
type Cardholder struct {
	CardholderID     string           `json:"cardholder_id"`
	Email            string           `json:"email"`
	NumberOfCards    int              `json:"number_of_cards"`
	FirstName        string           `json:"first_name"`
	LastName         string           `json:"last_name"`
	CreateTime       string           `json:"create_time"`
	CardholderStatus string           `json:"cardholder_status"`
	DateOfBirth      string           `json:"date_of_birth,omitempty"`
	CountryCode      string           `json:"country_code"`
	PhoneNumber      string           `json:"phone_number"`
	DeliveryAddress  *DeliveryAddress `json:"delivery_address,omitempty"`
	ReviewStatus     string           `json:"review_status,omitempty"`
}

// ListCardholdersRequest represents a cardholder list request
type ListCardholdersRequest struct {
	PageSize         int    `json:"page_size"`
	PageNumber       int    `json:"page_number"`
	CardholderStatus string `json:"cardholder_status,omitempty"` // optional filter: PENDING, SUCCESS, INCOMPLETE, FAILED
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
