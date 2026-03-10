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
	ReportType string `json:"report_type"` // Required. SETTLEMENT or LEDGER
	StartTime  string `json:"start_time"`  // Required. Earliest transaction time, ISO 8601 format
	EndTime    string `json:"end_time"`    // Required. Latest transaction time, ISO 8601 format
}

// ============================================================================
// Response Structures
// ============================================================================

// CreateReportResponse represents a response from creating an issuing report
type CreateReportResponse struct {
	ReportID string `json:"report_id"` // UUID. Unique identifier for the created report
}

// ============================================================================
// API Methods
// ============================================================================

// Create creates a new issuing report for account transactions, card transactions or transaction settlements
// Optional RequestOptions can be provided to set custom headers like x-on-behalf-of or x-idempotency-key
func (c *ReportsClient) Create(ctx context.Context, req *CreateReportRequest, opts ...*common.RequestOptions) (*CreateReportResponse, error) {
	var resp CreateReportResponse
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.PostWithOptions(ctx, "/v1/issuing/reports", req, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to create issuing report: %w", err)
	}
	return &resp, nil
}
