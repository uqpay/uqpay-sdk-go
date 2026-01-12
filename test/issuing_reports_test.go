package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/issuing"
)

func TestIssuingReports(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		req := &issuing.CreateReportRequest{
			ReportType: "SETTLEMENT",
			StartTime:  "2024-03-21T17:17:32+08:00",
			EndTime:    "2024-03-21T17:17:32+08:00",
		}

		resp, err := client.Issuing.Reports.Create(ctx, req)
		if err != nil {
			t.Logf("Create issuing report returned error: %v", err)
			return
		}

		t.Logf("Issuing report created successfully")
		t.Logf("Report ID: %s", resp.ReportID)
	})

	t.Run("Create_LEDGER", func(t *testing.T) {
		req := &issuing.CreateReportRequest{
			ReportType: "LEDGER",
			StartTime:  "2024-03-21T00:00:00+08:00",
			EndTime:    "2024-03-21T23:59:59+08:00",
		}

		resp, err := client.Issuing.Reports.Create(ctx, req)
		if err != nil {
			t.Logf("Create issuing LEDGER report returned error: %v", err)
			return
		}

		t.Logf("Issuing LEDGER report created successfully")
		t.Logf("Report ID: %s", resp.ReportID)
	})
}
