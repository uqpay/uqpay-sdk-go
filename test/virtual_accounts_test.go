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

	t.Logf("Found %d virtual accounts (total: %d, pages: %d)", len(resp.Data), resp.TotalItems, resp.TotalPages)
	for i, a := range resp.Data {
		if i >= 5 {
			break
		}
		t.Logf("  %d: %s - %s (%s), Status=%s", i+1, a.Currency, a.AccountNumber, a.BankName, a.Status)
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
