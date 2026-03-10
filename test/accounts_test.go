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

		t.Logf("Found %d accounts (total: %d)", len(resp.Data), resp.TotalItems)

		if len(resp.Data) > 0 {
			account := resp.Data[0]
			t.Logf("First account: ID=%s, Type=%s, Status=%s, VerificationStatus=%s",
				account.AccountID, account.EntityType, account.Status, account.VerificationStatus)

			if account.ShortReferenceID != "" {
				t.Logf("  ShortReferenceID: %s", account.ShortReferenceID)
			}
			if account.AccountName != "" {
				t.Logf("  AccountName: %s", account.AccountName)
			}
			if account.Country != "" {
				t.Logf("  Country: %s", account.Country)
			}
			if len(account.BusinessCode) > 0 {
				t.Logf("  BusinessCode: %v", account.BusinessCode)
			}
			if account.ReviewReason != "" {
				t.Logf("  ReviewReason: %s", account.ReviewReason)
			}
			if account.ContactDetails != nil {
				t.Logf("  Contact: Email=%s, Phone=%s", account.ContactDetails.Email, account.ContactDetails.Phone)
			}

			// Log entity-specific details
			if account.EntityType == connect.EntityTypeCompany && account.BusinessDetails != nil {
				t.Logf("  Company: %s", account.BusinessDetails.LegalEntityNameEnglish)
			} else if account.EntityType == connect.EntityTypeIndividual && account.PersonDetails != nil {
				t.Logf("  Individual: %s %s",
					account.PersonDetails.FirstNameEnglish, account.PersonDetails.LastNameEnglish)
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

		// Test Get with default business code
		account, err := client.Connect.Accounts.Get(ctx, accountID, "")
		if err != nil {
			t.Logf("Get account returned: %v", err)
			return
		}

		t.Logf("Retrieved account: ID=%s, Type=%s, Status=%s", account.AccountID, account.EntityType, account.Status)

		// Verify the account ID matches
		if account.AccountID != accountID {
			t.Errorf("Expected account ID %s, got %s", accountID, account.AccountID)
		}

		if account.ShortReferenceID != "" {
			t.Logf("  ShortReferenceID: %s", account.ShortReferenceID)
		}
		if account.VerificationStatus != "" {
			t.Logf("  VerificationStatus: %s", account.VerificationStatus)
		}
		if account.ContactDetails != nil {
			t.Logf("  Contact: Email=%s, Phone=%s", account.ContactDetails.Email, account.ContactDetails.Phone)
		}
		if account.BusinessDetails != nil {
			t.Logf("  BusinessDetails: %s", account.BusinessDetails.LegalEntityNameEnglish)
		}
		if len(account.Representatives) > 0 {
			t.Logf("  Representatives: %d", len(account.Representatives))
		}
		if len(account.Documents) > 0 {
			t.Logf("  Documents: %d", len(account.Documents))
		}
		if account.TosAcceptance != nil {
			t.Logf("  TosAcceptance: IP=%s, Date=%s", account.TosAcceptance.IP, account.TosAcceptance.Date)
		}
	})

	t.Run("GetAdditionalDocuments", func(t *testing.T) {
		// Test GetAdditionalDocuments with country and business code
		docs, err := client.Connect.Accounts.GetAdditionalDocuments(ctx, "SG", "BANKING")
		if err != nil {
			t.Logf("GetAdditionalDocuments returned: %v", err)
			return
		}

		t.Logf("Retrieved %d additional document profiles", len(docs))

		for i, doc := range docs {
			required := "optional"
			if doc.ProfileOption == 1 {
				required = "required"
			}
			t.Logf("  Document %d: Key=%s, Name=%s, %s",
				i+1, doc.ProfileKey, doc.ProfileName, required)
		}
	})
}
