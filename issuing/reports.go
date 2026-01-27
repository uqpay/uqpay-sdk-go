package issuing

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// ReportsClient handles issuing report operations
type ReportsClient struct {
	client *common.APIClient
}

// ============================================================================
// Request Structures
// ============================================================================

// CreateReportRequest represents a request to create an issuing report
type CreateReportRequest struct {
	ReportType string `json:"report_type"` // required - The type of report file, only SETTLEMENT and LEDGER can be accepted
	StartTime  string `json:"start_time"`  // required - The earliest time for transaction
	EndTime    string `json:"end_time"`    // required - The latest timestamp for transaction
}

// ============================================================================
// Response Structures
// ============================================================================

// CreateReportResponse represents a response from creating an issuing report
type CreateReportResponse struct {
	ReportID string `json:"report_id"` // Unique identifier for the report
}

// ============================================================================
// API Methods
// ============================================================================

// Create creates a new issuing report for account transactions, card transactions or transaction settlements
func (c *ReportsClient) Create(ctx context.Context, req *CreateReportRequest) (*CreateReportResponse, error) {
	var resp CreateReportResponse
	if err := c.client.Post(ctx, "/v1/issuing/reports", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create issuing report: %w", err)
	}
	return &resp, nil
}
