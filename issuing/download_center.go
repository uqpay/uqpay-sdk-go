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
	Data []byte // The raw file data
}

// ============================================================================
// API Methods
// ============================================================================

// Download downloads a report file by its ID
func (c *DownloadCenterClient) Download(ctx context.Context, reportID string) (*DownloadReportResponse, error) {
	if reportID == "" {
		return nil, fmt.Errorf("report ID is required")
	}

	path := fmt.Sprintf("/v1/issuing/reports/%s", reportID)
	data, err := c.client.GetRaw(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to download report: %w", err)
	}

	return &DownloadReportResponse{
		Data: data,
	}, nil
}
