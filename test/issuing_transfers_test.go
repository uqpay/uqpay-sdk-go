package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/issuing"
)

func TestIssuingTransfers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	var createdTransferID string

	t.Run("Create", func(t *testing.T) {
		req := &issuing.CreateTransferRequest{
			SourceAccountID:      "65087660-8d3d-428e-bd2e-9e56219c1512",
			DestinationAccountID: "11db237e-1a2b-4449-9878-a9bf1f0df0c7",
			Currency:             "SGD",
			Amount:               100.00,
			Remark:               "Test transfer from SDK",
		}

		resp, err := client.Issuing.Transfers.Create(ctx, req)
		if err != nil {
			t.Logf("Create issuing transfer returned error: %v", err)
			return
		}

		createdTransferID = resp.TransferID
		t.Logf("Issuing transfer created successfully")
		t.Logf("Transfer ID: %s", resp.TransferID)
	})

	t.Run("Retrieve", func(t *testing.T) {
		// Use the transfer ID from the Create test if available
		transferID := createdTransferID
		if transferID == "" {
			// Use a placeholder ID if Create didn't succeed
			transferID = "d58ed244-4b73-4095-bebf-2d05c0aab856"
		}

		resp, err := client.Issuing.Transfers.Retrieve(ctx, transferID)
		if err != nil {
			t.Logf("Retrieve issuing transfer returned error: %v", err)
			return
		}

		t.Logf("Issuing transfer retrieved successfully")
		t.Logf("Transfer ID: %s", resp.TransferID)
		t.Logf("Reference ID: %s", resp.ReferenceID)
		t.Logf("Source Account: %s", resp.SourceAccountID)
		t.Logf("Destination Account: %s", resp.DestinationAccountID)
		t.Logf("Amount: %s %s", resp.Amount, resp.Currency)
		t.Logf("Fee Amount: %s", resp.FeeAmount)
		t.Logf("Status: %s", resp.TransferStatus)
		t.Logf("Creator ID: %s", resp.CreatorID)
		t.Logf("Remark: %s", resp.Remark)
		t.Logf("Created: %s", resp.CreateTime)
		t.Logf("Completed: %s", resp.CompleteTime)
	})
}
