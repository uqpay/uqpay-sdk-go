package issuing

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// DownloadCenterClient handles issuing download center operations
type DownloadCenterClient struct {
	client *common.APIClient
}

// ============================================================================
// Response Structures
// ============================================================================

// DownloadReportResponse represents the response from downloading a report file
type DownloadReportResponse struct {
	Data []byte // Raw binary file data (application/octet-stream)
}

// ============================================================================
// API Methods
// ============================================================================

// Download downloads a report file by its ID
// Optional RequestOptions can be provided to set custom headers like x-on-behalf-of or x-idempotency-key
func (c *DownloadCenterClient) Download(ctx context.Context, reportID string, opts ...*common.RequestOptions) (*DownloadReportResponse, error) {
	if reportID == "" {
		return nil, fmt.Errorf("report ID is required")
	}

	path := fmt.Sprintf("/v1/issuing/reports/%s", reportID)
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	data, err := c.client.GetRawWithOptions(ctx, path, "application/octet-stream", opt)
	if err != nil {
		return nil, fmt.Errorf("failed to download report: %w", err)
	}

	return &DownloadReportResponse{
		Data: data,
	}, nil
}

// DownloadAsJSON retrieves report metadata/status as JSON by its ID
// Optional RequestOptions can be provided to set custom headers like x-on-behalf-of or x-idempotency-key
func (c *DownloadCenterClient) DownloadAsJSON(ctx context.Context, reportID string, response interface{}, opts ...*common.RequestOptions) error {
	if reportID == "" {
		return fmt.Errorf("report ID is required")
	}

	path := fmt.Sprintf("/v1/issuing/reports/%s", reportID)
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	return c.client.GetWithOptions(ctx, path, response, opt)
}
