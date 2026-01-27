package test

import (
	"context"
	"testing"
)

func TestIssuingDownloadCenter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	t.Run("Download", func(t *testing.T) {
		// Note: Replace with a valid report ID from your test environment
		reportID := "test-report-id"

		resp, err := client.Issuing.DownloadCenter.Download(ctx, reportID)
		if err != nil {
			t.Logf("Download report returned error: %v", err)
			return
		}

		t.Logf("Report downloaded successfully")
		t.Logf("Data length: %d bytes", len(resp.Data))
	})

	t.Run("Download_EmptyID", func(t *testing.T) {
		_, err := client.Issuing.DownloadCenter.Download(ctx, "")
		if err == nil {
			t.Error("Expected error for empty report ID, got nil")
			return
		}
		t.Logf("Got expected error for empty report ID: %v", err)
	})
}
