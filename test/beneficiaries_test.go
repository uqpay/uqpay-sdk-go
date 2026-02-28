package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/uqpay/uqpay-sdk-go/banking"
)

// newUSBankDetails creates a standard US ACH beneficiary bank details for testing
func newUSBankDetails(accountHolder string) *banking.BankDetails {
	// Use timestamp-based account number to avoid "beneficiary repeat addition" errors
	acctNum := fmt.Sprintf("99%d", time.Now().UnixNano()%100000000)
	return &banking.BankDetails{
		AccountNumber:       acctNum,
		AccountHolder:       accountHolder,
		AccountCurrencyCode: "USD",
		BankName:            "JPMorgan Chase",
		BankAddress:         "383 Madison Avenue, New York, NY 10179",
		BankCountryCode:     "US",
		SwiftCode:           "CHASUS33",
		ClearingSystem:      "ACH",
		RoutingCodeType1:    "ach",
		RoutingCodeValue1:   "021000021",
	}
}

// newUSAddress creates a standard US address for testing
func newUSAddress() *banking.Address {
	return &banking.Address{
		StreetAddress: "100 Test Street",
		City:          "New York",
		State:         "NY",
		PostalCode:    "10001",
		Country:       "US",
		CountryCode:   "US",
	}
}

func TestBeneficiaries(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		req := &banking.ListBeneficiariesRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Banking.Beneficiaries.List(ctx, req)
		if err != nil {
			t.Logf("List beneficiaries returned error: %v", err)
			return
		}

		t.Logf("Found %d beneficiaries (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			b := resp.Data[0]
			t.Logf("First beneficiary: ID=%s, Type=%s, PaymentMethod=%s",
				b.BeneficiaryID, b.EntityType, b.PaymentMethod)
			if b.BankDetails != nil {
				t.Logf("  Bank: %s (%s)", b.BankDetails.BankName, b.BankDetails.AccountNumber)
			}
		}
	})

	t.Run("ListWithFilters", func(t *testing.T) {
		req := &banking.ListBeneficiariesRequest{
			PageSize:   10,
			PageNumber: 1,
			EntityType: "INDIVIDUAL",
		}

		resp, err := client.Banking.Beneficiaries.List(ctx, req)
		if err != nil {
			t.Logf("List beneficiaries with filters returned error: %v", err)
			return
		}

		t.Logf("Found %d individual beneficiaries (total: %d)", len(resp.Data), resp.TotalItems)
	})

	t.Run("ListPaymentMethods", func(t *testing.T) {
		methods, err := client.Banking.Beneficiaries.ListPaymentMethods(ctx, "USD", "US")
		if err != nil {
			t.Logf("List payment methods returned error: %v", err)
			return
		}

		t.Logf("Found %d payment methods for USD/US", len(methods))
		for i, m := range methods {
			t.Logf("  Method %d: clearing=%s, payment_method=%s, currency=%s, country=%s",
				i+1, m.ClearingSystems, m.PaymentMethod, m.Currency, m.Country)
		}
	})

	t.Run("ListPaymentMethodsMultipleCurrencies", func(t *testing.T) {
		testCases := []struct {
			currency string
			country  string
		}{
			{"USD", "US"},
			{"GBP", "GB"},
			{"EUR", "DE"},
			{"SGD", "SG"},
		}

		for _, tc := range testCases {
			methods, err := client.Banking.Beneficiaries.ListPaymentMethods(ctx, tc.currency, tc.country)
			if err != nil {
				t.Logf("%s/%s: %v", tc.currency, tc.country, err)
				continue
			}
			t.Logf("%s/%s: %d payment methods available", tc.currency, tc.country, len(methods))
		}
	})

	t.Run("Create", func(t *testing.T) {
		req := &banking.BeneficiaryCreationRequest{
			EntityType:    "INDIVIDUAL",
			FirstName:     "SDK",
			LastName:      "Test",
			Currency:      "USD",
			Country:       "US",
			PaymentMethod: "LOCAL",
			BankDetails:   newUSBankDetails("SDK Test"),
			Address:       newUSAddress(),
			Email:         "sdk-test@example.com",
		}

		resp, err := client.Banking.Beneficiaries.Create(ctx, req)
		if err != nil {
			t.Fatalf("Create beneficiary failed: %v", err)
		}

		if resp.BeneficiaryID == "" {
			t.Error("BeneficiaryID should not be empty")
		}

		t.Logf("Beneficiary created: ID=%s, Ref=%s", resp.BeneficiaryID, resp.ShortReferenceID)
	})

	t.Run("Get", func(t *testing.T) {
		listResp, err := client.Banking.Beneficiaries.List(ctx, &banking.ListBeneficiariesRequest{
			PageSize: 10, PageNumber: 1,
		})
		if err != nil {
			t.Fatalf("List beneficiaries failed: %v", err)
		}
		if len(listResp.Data) == 0 {
			t.Skip("No beneficiaries available to test Get")
		}

		id := listResp.Data[0].BeneficiaryID
		resp, err := client.Banking.Beneficiaries.Get(ctx, id)
		if err != nil {
			t.Fatalf("Get beneficiary failed: %v", err)
		}

		if resp.BeneficiaryID != id {
			t.Errorf("ID mismatch: got %s, want %s", resp.BeneficiaryID, id)
		}

		t.Logf("Get OK: ID=%s, Type=%s, PaymentMethod=%s", resp.BeneficiaryID, resp.EntityType, resp.PaymentMethod)
	})

	t.Run("Delete", func(t *testing.T) {
		// Create one to delete
		created, err := client.Banking.Beneficiaries.Create(ctx, &banking.BeneficiaryCreationRequest{
			EntityType:    "INDIVIDUAL",
			FirstName:     "Delete",
			LastName:      "Test",
			Currency:      "USD",
			Country:       "US",
			PaymentMethod: "LOCAL",
			BankDetails:   newUSBankDetails("Delete Test"),
			Address:       newUSAddress(),
		})
		if err != nil {
			t.Fatalf("Create for delete test failed: %v", err)
		}

		t.Logf("Created beneficiary for delete: %s", created.BeneficiaryID)

		err = client.Banking.Beneficiaries.Delete(ctx, created.BeneficiaryID)
		if err != nil {
			t.Fatalf("Delete beneficiary failed: %v", err)
		}

		t.Logf("Deleted: %s", created.BeneficiaryID)
	})

	t.Run("FullLifecycle", func(t *testing.T) {
		// Step 1: Create
		created, err := client.Banking.Beneficiaries.Create(ctx, &banking.BeneficiaryCreationRequest{
			EntityType:    "INDIVIDUAL",
			FirstName:     "Lifecycle",
			LastName:      "Test",
			Currency:      "USD",
			Country:       "US",
			PaymentMethod: "LOCAL",
			BankDetails:   newUSBankDetails("Lifecycle Test"),
			Address:       newUSAddress(),
			Email:         "lifecycle@example.com",
		})
		if err != nil {
			t.Fatalf("Step 1 - Create failed: %v", err)
		}
		t.Logf("Step 1 - Created: %s", created.BeneficiaryID)

		// Step 2: Get
		fetched, err := client.Banking.Beneficiaries.Get(ctx, created.BeneficiaryID)
		if err != nil {
			t.Fatalf("Step 2 - Get failed: %v", err)
		}
		if fetched.BeneficiaryID != created.BeneficiaryID {
			t.Errorf("ID mismatch: got %s, want %s", fetched.BeneficiaryID, created.BeneficiaryID)
		}
		t.Logf("Step 2 - Get verified: %s %s", fetched.FirstName, fetched.LastName)

		// Step 3: Delete
		err = client.Banking.Beneficiaries.Delete(ctx, created.BeneficiaryID)
		if err != nil {
			t.Fatalf("Step 3 - Delete failed: %v", err)
		}
		t.Logf("Step 3 - Deleted: %s", created.BeneficiaryID)
	})
}
