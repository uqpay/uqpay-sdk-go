package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/uqpay/uqpay-sdk-go/banking"
	"github.com/uqpay/uqpay-sdk-go/common"
)

func TestPayouts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetBankingTestClient(t)
	ctx := context.Background()

	// Step 1: List beneficiaries to get an existing beneficiary_id
	beneficiaries, err := client.Banking.Beneficiaries.List(ctx, &banking.ListBeneficiariesRequest{
		PageSize: 1, PageNumber: 1,
	})
	if err != nil {
		t.Fatalf("Failed to list beneficiaries: %v", err)
	}
	if len(beneficiaries.Data) == 0 {
		t.Fatal("No beneficiaries in sandbox — cannot test Create Payout")
	}
	beneficiary := beneficiaries.Data[0]
	t.Logf("Using beneficiary: ID=%s, Type=%s, Method=%s",
		beneficiary.BeneficiaryID, beneficiary.EntityType, beneficiary.PaymentMethod)

	// Step 2: Create a payout using the existing beneficiary
	var createdPayoutID string
	t.Run("Create", func(t *testing.T) {
		payoutDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
		ref := fmt.Sprintf("SDK-TEST-%d", time.Now().UnixNano()%1000000)

		req := &banking.CreatePayoutRequest{
			Currency:        "USD",
			Amount:          "10.00",
			PurposeCode:     "GOODS_PURCHASED",
			PayoutReference: ref,
			FeePaidBy:       "OURS",
			PayoutDate:      payoutDate,
			BeneficiaryID:   beneficiary.BeneficiaryID,
		}

		resp, err := client.Banking.Payouts.Create(ctx, req, &common.RequestOptions{})
		if err != nil {
			t.Fatalf("Create payout failed: %v", err)
		}

		// Assert all CreatePayoutResponse fields deserialized
		if resp.PayoutID == "" {
			t.Error("CreatePayoutResponse.PayoutID is empty")
		}
		if resp.ShortReferenceID == "" {
			t.Error("CreatePayoutResponse.ShortReferenceID is empty")
		}
		if resp.PayoutStatus == "" {
			t.Error("CreatePayoutResponse.PayoutStatus is empty")
		}

		createdPayoutID = resp.PayoutID
		t.Logf("Created: ID=%s, Ref=%s, Status=%s",
			resp.PayoutID, resp.ShortReferenceID, resp.PayoutStatus)
	})

	// Step 3: List payouts — verify Payout struct fields
	t.Run("List", func(t *testing.T) {
		resp, err := client.Banking.Payouts.List(ctx, &banking.ListPayoutsRequest{
			PageSize: 10, PageNumber: 1,
		})
		if err != nil {
			t.Fatalf("List payouts failed: %v", err)
		}

		// Assert ListPayoutsResponse pagination fields
		if resp.TotalItems < 1 {
			t.Fatalf("Expected TotalItems >= 1, got %d", resp.TotalItems)
		}
		if resp.TotalPages < 1 {
			t.Fatalf("Expected TotalPages >= 1, got %d", resp.TotalPages)
		}
		if len(resp.Data) == 0 {
			t.Fatal("Expected at least 1 payout in Data")
		}

		// Assert all required Payout struct fields deserialized
		p := resp.Data[0]
		assertNonEmpty(t, "Payout.PayoutID", p.PayoutID)
		assertNonEmpty(t, "Payout.ShortReferenceID", p.ShortReferenceID)
		assertNonEmpty(t, "Payout.PayoutCurrency", p.PayoutCurrency)
		assertNonEmpty(t, "Payout.PayoutAmount", p.PayoutAmount)
		assertNonEmpty(t, "Payout.FeePaidBy", p.FeePaidBy)
		assertNonEmpty(t, "Payout.FeeCurrency", p.FeeCurrency)
		assertNonEmpty(t, "Payout.FeeAmount", p.FeeAmount)
		assertNonEmpty(t, "Payout.PayoutDate", p.PayoutDate)
		assertNonEmpty(t, "Payout.PayoutReference", p.PayoutReference)
		assertNonEmpty(t, "Payout.PurposeCode", p.PurposeCode)
		assertNonEmpty(t, "Payout.PayoutStatus", p.PayoutStatus)
		assertNonEmpty(t, "Payout.CreateTime", p.CreateTime)

		t.Logf("List OK: ID=%s, Amount=%s %s, Status=%s, Fee=%s %s",
			p.PayoutID, p.PayoutAmount, p.PayoutCurrency,
			p.PayoutStatus, p.FeeAmount, p.FeeCurrency)
		t.Logf("  PurposeCode=%s, Date=%s, Ref=%s, Method=%s",
			p.PurposeCode, p.PayoutDate, p.PayoutReference, p.PayoutMethod)
		if p.Conversion != nil {
			assertNonEmpty(t, "Payout.Conversion.CurrencyPair", p.Conversion.CurrencyPair)
			assertNonEmpty(t, "Payout.Conversion.ClientRate", p.Conversion.ClientRate)
			t.Logf("  Conversion: %s @ %s", p.Conversion.CurrencyPair, p.Conversion.ClientRate)
		}
	})

	// Step 4: List with all status filters — verify 200 OK for each
	t.Run("ListWithFilters", func(t *testing.T) {
		for _, status := range []string{"PENDING", "READY_TO_SEND", "COMPLETED", "FAILED", "REJECTED"} {
			resp, err := client.Banking.Payouts.List(ctx, &banking.ListPayoutsRequest{
				PageSize: 10, PageNumber: 1, PayoutStatus: status,
			})
			if err != nil {
				t.Errorf("List with status=%s failed: %v", status, err)
				continue
			}
			t.Logf("  %s: %d found (pages=%d)", status, resp.TotalItems, resp.TotalPages)
		}
	})

	// Step 5: Get payout by ID — verify PayoutDetailResponse fields
	t.Run("Get", func(t *testing.T) {
		if createdPayoutID == "" {
			t.Fatal("No payout ID from Create step")
		}

		resp, err := client.Banking.Payouts.Get(ctx, createdPayoutID)
		if err != nil {
			t.Fatalf("Get payout failed: %v", err)
		}

		// Assert embedded Payout fields
		if resp.PayoutID != createdPayoutID {
			t.Errorf("PayoutID mismatch: got %s, want %s", resp.PayoutID, createdPayoutID)
		}
		assertNonEmpty(t, "Get.PayoutID", resp.PayoutID)
		assertNonEmpty(t, "Get.ShortReferenceID", resp.ShortReferenceID)
		assertNonEmpty(t, "Get.PayoutCurrency", resp.PayoutCurrency)
		assertNonEmpty(t, "Get.PayoutAmount", resp.PayoutAmount)
		assertNonEmpty(t, "Get.FeePaidBy", resp.FeePaidBy)
		assertNonEmpty(t, "Get.FeeCurrency", resp.FeeCurrency)
		assertNonEmpty(t, "Get.FeeAmount", resp.FeeAmount)
		assertNonEmpty(t, "Get.PayoutDate", resp.PayoutDate)
		assertNonEmpty(t, "Get.PayoutReference", resp.PayoutReference)
		assertNonEmpty(t, "Get.PurposeCode", resp.PurposeCode)
		assertNonEmpty(t, "Get.PayoutStatus", resp.PayoutStatus)
		assertNonEmpty(t, "Get.CreateTime", resp.CreateTime)
		assertNonEmpty(t, "Get.PayoutMethod", resp.PayoutMethod)

		// Assert PayoutDetailResponse-specific fields
		assertNonEmpty(t, "Get.SourceCurrency", resp.SourceCurrency)
		assertNonEmpty(t, "Get.SourceAmount", resp.SourceAmount)
		assertNonEmpty(t, "Get.AmountPayerPays", resp.AmountPayerPays)
		assertNonEmpty(t, "Get.AmountBeneficiaryReceives", resp.AmountBeneficiaryReceives)

		// Assert nested Payer object
		if resp.Payer == nil {
			t.Error("Get.Payer is nil — expected payer details")
		} else {
			assertNonEmpty(t, "Get.Payer.EntityType", resp.Payer.EntityType)
			t.Logf("  Payer: ID=%s, Type=%s, Country=%s",
				resp.Payer.PayerID, resp.Payer.EntityType, resp.Payer.Country)
		}

		// Assert nested Beneficiary object
		if resp.Beneficiary == nil {
			t.Error("Get.Beneficiary is nil — expected beneficiary details")
		} else {
			assertNonEmpty(t, "Get.Beneficiary.BeneficiaryID", resp.Beneficiary.BeneficiaryID)
			assertNonEmpty(t, "Get.Beneficiary.EntityType", resp.Beneficiary.EntityType)
			t.Logf("  Beneficiary: ID=%s, Type=%s",
				resp.Beneficiary.BeneficiaryID, resp.Beneficiary.EntityType)
		}

		t.Logf("Get OK: ID=%s, Amount=%s %s, Status=%s, Method=%s",
			resp.PayoutID, resp.PayoutAmount, resp.PayoutCurrency,
			resp.PayoutStatus, resp.PayoutMethod)
		t.Logf("  Source=%s %s, PayerPays=%s, BeneficiaryReceives=%s",
			resp.SourceAmount, resp.SourceCurrency,
			resp.AmountPayerPays, resp.AmountBeneficiaryReceives)

		if resp.Conversion != nil {
			assertNonEmpty(t, "Get.Conversion.CurrencyPair", resp.Conversion.CurrencyPair)
			assertNonEmpty(t, "Get.Conversion.ClientRate", resp.Conversion.ClientRate)
			t.Logf("  Conversion: %s @ %s", resp.Conversion.CurrencyPair, resp.Conversion.ClientRate)
		}
	})
}

// assertNonEmpty fails the test if the value is empty
func assertNonEmpty(t *testing.T, field, value string) {
	t.Helper()
	if value == "" {
		t.Errorf("%s is empty — field not deserialized from API response", field)
	}
}
