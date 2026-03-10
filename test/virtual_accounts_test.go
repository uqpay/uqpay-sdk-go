package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/banking"
)

func TestVirtualAccountsCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	resp, err := client.Banking.VirtualAccounts.Create(ctx, &banking.CreateVirtualAccountRequest{
		Currency:      "USD",
		PaymentMethod: "LOCAL",
	})
	if err != nil {
		t.Fatalf("Failed to create virtual account: %v", err)
	}

	t.Logf("Create response: Message=%s", resp.Message)
	if resp.Message != "SUCCESS" {
		t.Errorf("Expected message SUCCESS, got %s", resp.Message)
	}
}

func TestVirtualAccountsList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	resp, err := client.Banking.VirtualAccounts.List(ctx, &banking.ListVirtualAccountsRequest{
		PageSize: 10, PageNumber: 1,
	})
	if err != nil {
		t.Fatalf("Failed to list virtual accounts: %v", err)
	}

	// Validate response pagination fields
	if resp.TotalPages < 1 {
		t.Errorf("Expected total_pages >= 1, got %d", resp.TotalPages)
	}
	if resp.TotalItems < 1 {
		t.Errorf("Expected total_items >= 1, got %d", resp.TotalItems)
	}
	if len(resp.Data) == 0 {
		t.Fatal("Expected at least one virtual account in response data")
	}

	// Validate VirtualAccount struct fields are deserialized correctly
	a := resp.Data[0]
	if a.AccountNumber == "" {
		t.Error("Expected account_number to be populated")
	}
	if a.Currency == "" {
		t.Error("Expected currency to be populated")
	}
	if a.Status == "" {
		t.Error("Expected status to be populated")
	}
	if a.Status != "ACTIVE" && a.Status != "INACTIVE" && a.Status != "CLOSED" {
		t.Errorf("Unexpected status value: %s (expected ACTIVE, INACTIVE, or CLOSED)", a.Status)
	}

	t.Logf("Found %d virtual accounts (total: %d, pages: %d)", len(resp.Data), resp.TotalItems, resp.TotalPages)
	for i, va := range resp.Data {
		if i >= 5 {
			break
		}
		t.Logf("  %d: Currency=%s AccountNumber=%s BankName=%s Status=%s", i+1, va.Currency, va.AccountNumber, va.BankName, va.Status)
		if va.Capability != nil {
			t.Logf("     Capability: PaymentMethod=%s", va.Capability.PaymentMethod)
		}
		if va.ClearingSystem != nil {
			t.Logf("     ClearingSystem: Type=%s Value=%s", va.ClearingSystem.Type, va.ClearingSystem.Value)
		}
	}
}

func TestVirtualAccountsListByCurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	resp, err := client.Banking.VirtualAccounts.List(ctx, &banking.ListVirtualAccountsRequest{
		PageSize: 10, PageNumber: 1, Currency: "USD",
	})
	if err != nil {
		t.Fatalf("Failed to list virtual accounts by currency: %v", err)
	}

	t.Logf("Found %d USD virtual accounts (total: %d)", len(resp.Data), resp.TotalItems)
	for i, a := range resp.Data {
		if i >= 3 {
			break
		}
		t.Logf("  %d: %s - %s, Status=%s", i+1, a.AccountNumber, a.BankName, a.Status)
	}
}
