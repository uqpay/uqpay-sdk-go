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

	account, err := client.Banking.VirtualAccounts.Create(ctx, &banking.CreateVirtualAccountRequest{
		Currency:      "USD",
		PaymentMethod: "LOCAL",
	})
	if err != nil {
		t.Fatalf("Failed to create virtual account: %v", err)
	}

	t.Logf("Created virtual account:")
	t.Logf("  BankID=%s, Holder=%s, Currency=%s", account.AccountBankID, account.AccountHolder, account.Currency)
	t.Logf("  AccountNumber=%s, Country=%s, Status=%s", account.AccountNumber, account.CountryCode, account.Status)
	t.Logf("  Bank=%s", account.BankName)
	if account.Capability != nil {
		t.Logf("  PaymentMethod=%s", account.Capability.PaymentMethod)
	}
	if account.ClearingSystem != nil {
		t.Logf("  Clearing=%s: %s", account.ClearingSystem.Type, account.ClearingSystem.Value)
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

	t.Logf("Found %d virtual accounts (total: %d, pages: %d)", len(resp.Data), resp.TotalItems, resp.TotalPages)
	for i, a := range resp.Data {
		if i >= 5 {
			break
		}
		t.Logf("  %d: %s - %s (%s), Status=%s", i+1, a.Currency, a.AccountNumber, a.BankName, a.Status)
	}
}

func TestVirtualAccountsCreateAndList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	account, err := client.Banking.VirtualAccounts.Create(ctx, &banking.CreateVirtualAccountRequest{
		Currency:      "USD",
		PaymentMethod: "LOCAL",
	})
	if err != nil {
		t.Fatalf("Failed to create virtual account: %v", err)
	}
	t.Logf("Created: BankID=%s, Currency=%s", account.AccountBankID, account.Currency)

	listResp, err := client.Banking.VirtualAccounts.List(ctx, &banking.ListVirtualAccountsRequest{
		PageSize: 10, PageNumber: 1,
	})
	if err != nil {
		t.Fatalf("Failed to list virtual accounts: %v", err)
	}

	found := false
	for _, a := range listResp.Data {
		if a.AccountBankID == account.AccountBankID {
			found = true
			t.Logf("Found in list: BankID=%s, Status=%s", a.AccountBankID, a.Status)
			break
		}
	}
	if !found {
		t.Log("Note: Created account not found in first page")
	}
}

func TestVirtualAccountsMultipleCurrencies(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	account, err := client.Banking.VirtualAccounts.Create(ctx, &banking.CreateVirtualAccountRequest{
		Currency:      "USD,EUR,GBP",
		PaymentMethod: "LOCAL",
	})
	if err != nil {
		t.Fatalf("Failed to create multi-currency virtual account: %v", err)
	}

	t.Logf("Created multi-currency account: BankID=%s, Status=%s", account.AccountBankID, account.Status)
}
