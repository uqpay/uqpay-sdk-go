package test

import (
	"context"
	"testing"

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
			AccountNumber:   "GC71950018692652646591",
			BankName:        "DBS Bank",
			SwiftCode:       "DBSSSGSG",
			BankCountryCode: "SG",
			BankAddress:     "22 Merina Boulevard, Singapore 028982",
			Currency:        "SGD",
		}

		resp, err := client.Payment.BankAccounts.Create(ctx, req)
		if err != nil {
			t.Fatalf("Create bank account failed: %v", err)
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
			t.Fatalf("Create bank account with ABA failed: %v", err)
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
		// First create a bank account to retrieve
		createReq := &payment.CreateBankAccountRequest{
			AccountNumber:   "GB91950018692652646591",
			BankName:        "DBS Bank",
			SwiftCode:       "DBSSSGSG",
			BankCountryCode: "SG",
			BankAddress:     "13 Marina Boulevard, Singapore 018182",
			Currency:        "SGD",
		}

		created, err := client.Payment.BankAccounts.Create(ctx, createReq)
		if err != nil {
			t.Fatalf("Create bank account failed: %v", err)
		}

		// Now retrieve it by ID
		resp, err := client.Payment.BankAccounts.Get(ctx, created.ID)
		if err != nil {
			t.Fatalf("Get bank account failed: %v", err)
		}

		// Assertions
		if resp.ID != created.ID {
			t.Errorf("ID mismatch: got %s, want %s", resp.ID, created.ID)
		}
		if resp.AccountNumber != createReq.AccountNumber {
			t.Errorf("AccountNumber mismatch: got %s, want %s", resp.AccountNumber, createReq.AccountNumber)
		}
		if resp.BankName != createReq.BankName {
			t.Errorf("BankName mismatch: got %s, want %s", resp.BankName, createReq.BankName)
		}
		if resp.SwiftCode != createReq.SwiftCode {
			t.Errorf("SwiftCode mismatch: got %s, want %s", resp.SwiftCode, createReq.SwiftCode)
		}
		if resp.BankCountryCode != createReq.BankCountryCode {
			t.Errorf("BankCountryCode mismatch: got %s, want %s", resp.BankCountryCode, createReq.BankCountryCode)
		}
		if resp.Currency != createReq.Currency {
			t.Errorf("Currency mismatch: got %s, want %s", resp.Currency, createReq.Currency)
		}
		if resp.AccountStatus == "" {
			t.Error("AccountStatus should not be empty")
		}

		t.Logf("Bank account retrieved successfully")
		t.Logf("   ID: %s", resp.ID)
		t.Logf("   Account Number: %s", resp.AccountNumber)
		t.Logf("   Account Name: %s", resp.AccountName)
		t.Logf("   Bank Name: %s", resp.BankName)
		t.Logf("   SWIFT Code: %s", resp.SwiftCode)
		t.Logf("   Country: %s", resp.BankCountryCode)
		t.Logf("   Currency: %s", resp.Currency)
		t.Logf("   Status: %s", resp.AccountStatus)
	})

	t.Run("Update", func(t *testing.T) {
		// First create a bank account to update
		createReq := &payment.CreateBankAccountRequest{
			AccountNumber:   "GB71350018692652646591",
			BankName:        "DBS Bank",
			SwiftCode:       "DBSSSGSG",
			BankCountryCode: "SG",
			BankAddress:     "15 Merina Boulevard, Singapore 017982",
			Currency:        "SGD",
		}

		created, err := client.Payment.BankAccounts.Create(ctx, createReq)
		if err != nil {
			t.Fatalf("Create bank account failed: %v", err)
		}

		// Now update the bank account
		updateReq := &payment.UpdateBankAccountRequest{
			AccountNumber:   "GB71950018692652646551",
			BankName:        "DBS Bank Updated",
			SwiftCode:       "DBSSSGSG",
			BankCountryCode: "SG",
			BankAddress:     "18 Marinea Boulevard, Singapore 016983",
		}

		resp, err := client.Payment.BankAccounts.Update(ctx, created.ID, updateReq)
		if err != nil {
			t.Fatalf("Update bank account failed: %v", err)
		}

		// Assertions
		if resp.ID != created.ID {
			t.Errorf("ID mismatch: got %s, want %s", resp.ID, created.ID)
		}
		if resp.AccountNumber != updateReq.AccountNumber {
			t.Errorf("AccountNumber mismatch: got %s, want %s", resp.AccountNumber, updateReq.AccountNumber)
		}
		if resp.BankName != updateReq.BankName {
			t.Errorf("BankName mismatch: got %s, want %s", resp.BankName, updateReq.BankName)
		}
		if resp.BankAddress != updateReq.BankAddress {
			t.Errorf("BankAddress mismatch: got %s, want %s", resp.BankAddress, updateReq.BankAddress)
		}
		if resp.AccountStatus == "" {
			t.Error("AccountStatus should not be empty")
		}

		t.Logf("Bank account updated successfully")
		t.Logf("   ID: %s", resp.ID)
		t.Logf("   Account Number: %s", resp.AccountNumber)
		t.Logf("   Account Name: %s", resp.AccountName)
		t.Logf("   Bank Name: %s", resp.BankName)
		t.Logf("   Bank Address: %s", resp.BankAddress)
		t.Logf("   SWIFT Code: %s", resp.SwiftCode)
		t.Logf("   Country: %s", resp.BankCountryCode)
		t.Logf("   Currency: %s", resp.Currency)
		t.Logf("   Status: %s", resp.AccountStatus)
	})

	t.Run("List", func(t *testing.T) {
		// First create a bank account to ensure there's at least one
		createReq := &payment.CreateBankAccountRequest{
			AccountNumber:   "GB71820018692652646598",
			BankName:        "DBS Bank",
			SwiftCode:       "DBSSSGSG",
			BankCountryCode: "SG",
			BankAddress:     "17 Marini Boulevard, Singapore 014982",
			Currency:        "SGD",
		}

		_, err := client.Payment.BankAccounts.Create(ctx, createReq)
		if err != nil {
			t.Fatalf("Create bank account failed: %v", err)
		}

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

		// Log first bank account details
		if len(resp.Data) > 0 {
			acc := resp.Data[0]
			t.Logf("   First bank account:")
			t.Logf("      ID: %s", acc.ID)
			t.Logf("      Account Number: %s", acc.AccountNumber)
			t.Logf("      Bank Name: %s", acc.BankName)
			t.Logf("      Currency: %s", acc.Currency)
			t.Logf("      Status: %s", acc.AccountStatus)
		}
	})
}
