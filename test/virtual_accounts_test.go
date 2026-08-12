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
		Country:       "BH",
		Currency:      "USD",
		PaymentMethod: "SWIFT",
	})
	if err != nil {
		t.Fatalf("Failed to create virtual account: %v", err)
	}

	t.Logf("Created virtual account:")
	t.Logf("  ApplicationID=%s, Version=%d, Currency=%s, Country=%s, Status=%s", account.Data.ApplicationID, account.Data.PublicVersion, account.Data.Currency, account.Data.Country, account.Data.Status)
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
		Country:       "BH",
		Currency:      "USD",
		PaymentMethod: "SWIFT",
	})
	if err != nil {
		t.Fatalf("Failed to create virtual account: %v", err)
	}
	t.Logf("Created application: ID=%s, Currency=%s", account.Data.ApplicationID, account.Data.Currency)

	listResp, err := client.Banking.VirtualAccounts.List(ctx, &banking.ListVirtualAccountsRequest{
		PageSize: 10, PageNumber: 1,
	})
	if err != nil {
		t.Fatalf("Failed to list virtual accounts: %v", err)
	}

	found := false
	for _, a := range listResp.Data {
		for _, result := range account.Data.Results {
			for _, bank := range result.VirtualAccounts {
				if a.AccountBankID == bank.AccountBankID {
					found = true
					t.Logf("Found in list: BankID=%s, Status=%s", a.AccountBankID, a.Status)
				}
			}
		}
	}
	if !found {
		t.Log("Note: Created account not found in first page")
	}
}

func TestVirtualAccountsOmittedPaymentMethod(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	account, err := client.Banking.VirtualAccounts.Create(ctx, &banking.CreateVirtualAccountRequest{
		Country:  "SG",
		Currency: "USD",
	})
	if err != nil {
		t.Fatalf("Failed to create multi-currency virtual account: %v", err)
	}

	t.Logf("Created application evaluating methods: ID=%s, Status=%s", account.Data.ApplicationID, account.Data.Status)
}
