package test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/uqpay/uqpay-sdk-go/common"
	"github.com/uqpay/uqpay-sdk-go/payment"
)

// ============================================================================
// Bank Account Tests
// ============================================================================

func TestBankAccounts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		req := &payment.CreateBankAccountRequest{
			AccountNumber:   "1234567890",
			BankName:        "DBS Bank",
			SwiftCode:       "DBSSSGSG",
			BankCountryCode: "SG",
			BankAddress:     "12 Marina Boulevard, Singapore 018982",
			Currency:        "USD",
		}

		resp, err := client.Payment.BankAccounts.Create(ctx, req, &common.RequestOptions{
			IdempotencyKey: uuid.New().String(),
		})
		if err != nil {
			t.Logf("Create bank account returned: %v", err)
			return
		}

		// Assertions
		if resp.ID == "" {
			t.Error("ID should not be empty")
		}
		if resp.AccountNumber != req.AccountNumber {
			t.Errorf("AccountNumber mismatch: got %s, want %s", resp.AccountNumber, req.AccountNumber)
		}
		if resp.BankName != req.BankName {
			t.Errorf("BankName mismatch: got %s, want %s", resp.BankName, req.BankName)
		}
		if resp.Currency != req.Currency {
			t.Errorf("Currency mismatch: got %s, want %s", resp.Currency, req.Currency)
		}
		if resp.AccountStatus == "" {
			t.Error("AccountStatus should not be empty")
		}

		t.Logf("Bank account created successfully")
		t.Logf("   ID: %s", resp.ID)
		t.Logf("   Account Number: %s", resp.AccountNumber)
		t.Logf("   Account Name: %s", resp.AccountName)
		t.Logf("   Bank Name: %s", resp.BankName)
		t.Logf("   SWIFT Code: %s", resp.SwiftCode)
		t.Logf("   Country: %s", resp.BankCountryCode)
		t.Logf("   Currency: %s", resp.Currency)
		t.Logf("   Status: %s", resp.AccountStatus)
	})

	t.Run("CreateWithBankCode_USD_ABA", func(t *testing.T) {
		req := &payment.CreateBankAccountRequest{
			AccountNumber:   "121456789011",
			BankName:        "Bank of America",
			SwiftCode:       "BOFAUS3N",
			BankCountryCode: "US",
			BankAddress:     "110 N Tryon St, Charlotte, NC 18255",
			Currency:        "USD",
			BankCodeType:    "aba",
			BankCodeValue:   "021000021",
		}

		resp, err := client.Payment.BankAccounts.Create(ctx, req)
		if err != nil {
			t.Logf("Create bank account with ABA returned: %v", err)
			return
		}

		// Assertions
		if resp.ID == "" {
			t.Error("ID should not be empty")
		}
		if resp.BankCodeType != "aba" {
			t.Errorf("BankCodeType mismatch: got %s, want aba", resp.BankCodeType)
		}
		if resp.BankCodeValue != req.BankCodeValue {
			t.Errorf("BankCodeValue mismatch: got %s, want %s", resp.BankCodeValue, req.BankCodeValue)
		}

		t.Logf("Bank account with ABA created successfully")
		t.Logf("   ID: %s", resp.ID)
		t.Logf("   Bank Code Type: %s", resp.BankCodeType)
		t.Logf("   Bank Code Value: %s", resp.BankCodeValue)
		t.Logf("   Status: %s", resp.AccountStatus)
	})

	t.Run("Get", func(t *testing.T) {
		// List existing bank accounts and use the first one
		listResp, err := client.Payment.BankAccounts.List(ctx, &payment.ListBankAccountsRequest{
			PageNumber: 1,
			PageSize:   1,
		})
		if err != nil {
			t.Logf("List bank accounts failed, skipping Get test: %v", err)
			return
		}
		if len(listResp.Data) == 0 {
			t.Log("No bank accounts available, skipping Get test")
			return
		}

		id := listResp.Data[0].ID
		resp, err := client.Payment.BankAccounts.Get(ctx, id)
		if err != nil {
			t.Fatalf("Get bank account failed: %v", err)
		}

		if resp.ID != id {
			t.Errorf("ID mismatch: got %s, want %s", resp.ID, id)
		}
		if resp.AccountStatus == "" {
			t.Error("AccountStatus should not be empty")
		}

		t.Logf("✅ Bank account retrieved: ID=%s, Currency=%s, Status=%s", resp.ID, resp.Currency, resp.AccountStatus)
	})

	t.Run("Update", func(t *testing.T) {
		// List existing bank accounts and use the first one
		listResp, err := client.Payment.BankAccounts.List(ctx, &payment.ListBankAccountsRequest{
			PageNumber: 1,
			PageSize:   1,
		})
		if err != nil {
			t.Logf("List bank accounts failed, skipping Update test: %v", err)
			return
		}
		if len(listResp.Data) == 0 {
			t.Log("No bank accounts available, skipping Update test")
			return
		}

		id := listResp.Data[0].ID
		updateReq := &payment.UpdateBankAccountRequest{
			AccountNumber:   "1234567890",
			BankName:        "DBS Bank Ltd",
			SwiftCode:       "DBSSSGSG",
			BankCountryCode: "SG",
			BankAddress:     "12 Marina Boulevard, Singapore 018982",
		}

		resp, err := client.Payment.BankAccounts.Update(ctx, id, updateReq)
		if err != nil {
			t.Fatalf("Update bank account failed: %v", err)
		}

		if resp.ID == "" {
			t.Error("ID should not be empty")
		}
		if resp.AccountStatus == "" {
			t.Error("AccountStatus should not be empty")
		}

		t.Logf("✅ Bank account updated: ID=%s, BankName=%s, Status=%s", resp.ID, resp.BankName, resp.AccountStatus)
	})

	t.Run("List", func(t *testing.T) {
		// Now list bank accounts
		listReq := &payment.ListBankAccountsRequest{
			PageNumber: 1,
			PageSize:   10,
		}

		resp, err := client.Payment.BankAccounts.List(ctx, listReq)
		if err != nil {
			t.Fatalf("List bank accounts failed: %v", err)
		}

		// Assertions
		if resp.TotalItems < 1 {
			t.Error("TotalItems should be at least 1")
		}
		if resp.TotalPages < 1 {
			t.Error("TotalPages should be at least 1")
		}
		if len(resp.Data) < 1 {
			t.Error("Data should contain at least 1 bank account")
		}

		t.Logf("Bank accounts listed successfully")
		t.Logf("   Total Pages: %d", resp.TotalPages)
		t.Logf("   Total Items: %d", resp.TotalItems)
		t.Logf("   Items in page: %d", len(resp.Data))

		// Log all bank account details
		for i, acc := range resp.Data {
			t.Logf("   Bank Account #%d:", i+1)
			t.Logf("      ID: %s", acc.ID)
			t.Logf("      Account Number: %s", acc.AccountNumber)
			t.Logf("      Bank Name: %s", acc.BankName)
			t.Logf("      Currency: %s", acc.Currency)
			t.Logf("      Status: %s", acc.AccountStatus)
		}
	})
}
