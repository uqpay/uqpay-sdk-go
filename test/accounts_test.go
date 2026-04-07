package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/connect"
)

func TestAccounts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		req := &connect.ListAccountsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Connect.Accounts.List(ctx, req)
		if err != nil {
			t.Logf("List accounts returned: %v", err)
			return
		}

		t.Logf("✅ Found %d accounts (total: %d)", len(resp.Data), resp.TotalItems)

		if len(resp.Data) > 0 {
			account := resp.Data[0]
			t.Logf("First account: ID=%s, Type=%s, Status=%s, PayoutsEnabled=%t, ChargesEnabled=%t",
				account.AccountID, account.EntityType, account.Status, account.PayoutsEnabled, account.ChargesEnabled)

			// If there are requirements, log them
			if account.Requirements != nil {
				if len(account.Requirements.CurrentlyDue) > 0 {
					t.Logf("  Currently due requirements: %v", account.Requirements.CurrentlyDue)
				}
				if len(account.Requirements.EventuallyDue) > 0 {
					t.Logf("  Eventually due requirements: %v", account.Requirements.EventuallyDue)
				}
				if account.Requirements.Disabled {
					t.Logf("  Account disabled: %s", account.Requirements.DisabledReason)
				}
			}

			// Log entity-specific details
			if account.EntityType == connect.EntityTypeIndividual && account.Individual != nil {
				t.Logf("  Individual: %s %s, Email=%s",
					account.Individual.FirstName, account.Individual.LastName, account.Individual.ContactInfo.Email)
			} else if account.EntityType == connect.EntityTypeCompany && account.Company != nil {
				t.Logf("  Company: %s, Type=%s, Email=%s",
					account.Company.LegalName, account.Company.BusinessType, account.Company.ContactInfo.Email)
			}
		}
	})

	t.Run("Get", func(t *testing.T) {
		// First, list accounts to get a valid account ID
		listReq := &connect.ListAccountsRequest{
			PageSize:   1,
			PageNumber: 1,
		}

		listResp, err := client.Connect.Accounts.List(ctx, listReq)
		if err != nil {
			t.Logf("List accounts failed, skipping Get test: %v", err)
			return
		}

		if len(listResp.Data) == 0 {
			t.Log("No accounts available, skipping Get test")
			return
		}

		accountID := listResp.Data[0].AccountID

		// Test Get
		account, err := client.Connect.Accounts.Get(ctx, accountID)
		if err != nil {
			t.Logf("Get account returned: %v", err)
			return
		}

		t.Logf("✅ Retrieved account: ID=%s, Type=%s, Status=%s", account.AccountID, account.EntityType, account.Status)

		// Verify the account ID matches
		if account.AccountID != accountID {
			t.Errorf("Expected account ID %s, got %s", accountID, account.AccountID)
		}

		// Log detailed information
		t.Logf("  Payouts Enabled: %t", account.PayoutsEnabled)
		t.Logf("  Charges Enabled: %t", account.ChargesEnabled)
		t.Logf("  Created: %s", account.CreateTime)
		if account.UpdateTime != "" {
			t.Logf("  Updated: %s", account.UpdateTime)
		}

		// Log metadata if present
		if len(account.Metadata) > 0 {
			t.Logf("  Metadata: %v", account.Metadata)
		}
	})

	t.Run("GetAdditionalDocuments", func(t *testing.T) {
		// Retrieves required/optional document types for company sub-accounts by country
		docs, err := client.Connect.Accounts.GetAdditionalDocuments(ctx, "GB", "BANKING")
		if err != nil {
			t.Logf("GetAdditionalDocuments returned: %v", err)
			return
		}

		t.Logf("✅ Retrieved additional documents for GB/BANKING: %d items", len(docs))
		for i, doc := range docs {
			required := "optional"
			if doc.ProfileOption == 1 {
				required = "required"
			}
			t.Logf("  [%d] %s (%s) - %s", i+1, doc.ProfileKey, doc.ProfileName, required)
		}
	})
}
