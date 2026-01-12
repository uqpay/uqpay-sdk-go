package issuing

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// TransfersClient handles issuing transfer operations
type TransfersClient struct {
	client *common.APIClient
}

// ============================================================================
// Request Structures
// ============================================================================

// CreateTransferRequest represents a request to create an issuing transfer
type CreateTransferRequest struct {
	SourceAccountID      string  `json:"source_account_id"`      // required - The account id that initiated the transfer
	DestinationAccountID string  `json:"destination_account_id"` // required - The account id that received the transfer
	Currency             string  `json:"currency"`               // required - Transfer currency
	Amount               float64 `json:"amount"`                 // required - Transfer amount (precision limited to two decimal places)
	Remark               string  `json:"remark,omitempty"`       // optional - The remark of the transfer
}

// ============================================================================
// Response Structures
// ============================================================================

// CreateTransferResponse represents a response from creating an issuing transfer
type CreateTransferResponse struct {
	TransferID string `json:"transfer_id"` // Unique identifier for transfer
}

// Transfer represents the full details of an issuing transfer
type Transfer struct {
	TransferID           string `json:"transfer_id"`            // Unique identifier for transfer
	ReferenceID          string `json:"reference_id"`           // Short reference id for the transfer
	SourceAccountID      string `json:"source_account_id"`      // The account id that initiated the transfer
	DestinationAccountID string `json:"destination_account_id"` // The account id that received the transfer
	Amount               string `json:"amount"`                 // Transfer amount
	FeeAmount            string `json:"fee_amount"`             // Transaction fee amount
	Currency             string `json:"currency"`               // Transfer currency
	TransferStatus       string `json:"transfer_status"`        // Transfer status: pending, failed, completed
	CreateTime           string `json:"create_time"`            // Transfer create time
	CompleteTime         string `json:"complete_time"`          // Transfer complete time
	CreatorID            string `json:"creator_id"`             // The account id that create the transfer
	Remark               string `json:"remark"`                 // The remark of the transfer
}

// ============================================================================
// API Methods
// ============================================================================

// Create creates a new issuing transfer
// This API is specifically designed for fund transfers between the master account
// and its sub-accounts, and does not apply to cross-business-line or external transfers.
func (c *TransfersClient) Create(ctx context.Context, req *CreateTransferRequest) (*CreateTransferResponse, error) {
	var resp CreateTransferResponse
	if err := c.client.Post(ctx, "/v1/issuing/transfers", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create issuing transfer: %w", err)
	}
	return &resp, nil
}

// Retrieve retrieves an issuing transfer with the provided transfer id
func (c *TransfersClient) Retrieve(ctx context.Context, transferID string) (*Transfer, error) {
	var resp Transfer
	path := fmt.Sprintf("/v1/issuing/transfers/%s", transferID)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to retrieve issuing transfer: %w", err)
	}
	return &resp, nil
}
