package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/connect"
)

func TestCreateSubAccount(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	t.Run("Individual", func(t *testing.T) {
		req := &connect.CreateSubAccountRequest{
			EntityType: connect.EntityTypeIndividual,
			Nickname:   "SDK Test Individual",
			IndividualInfo: &connect.SubAccountIndividualInfo{
				FirstNameEnglish:   "John",
				LastNameEnglish:    "Doe",
				Nationality:        "GB",
				PhoneNumber:        "+447911123456",
				EmailAddress:       "john.doe.sdktest@example.com",
				DateOfBirth:        "1990-01-15",
				CountryOrTerritory: "GB",
				StreetAddress:      "123 Baker Street",
				City:               "London",
				State:              "England",
				PostalCode:         "W1U 6RS",
			},
			IdentityVerification: &connect.SubAccountIdentityVerification{
				IdentificationType:  connect.SubAccountIDTypePassport,
				IdentificationValue: "P12345678",
				IdentityDocs:        []string{"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=="},
				FaceDocs:            []string{"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=="},
			},
			ExpectedActivity: &connect.SubAccountExpectedActivity{
				AccountPurpose:      []string{connect.SubAccountPurposePurchase, connect.SubAccountPurposeBillPayment},
				BankingCountries:    []string{"GB", "US"},
				BankingCurrencies:   []string{"GBP", "USD"},
				Internationally:     1,
				TurnoverMonthly:     connect.TurnoverMonthlyTM002,
				TurnoverMonthlyCurrency: "GBP",
			},
			ProofDocuments: &connect.SubAccountProofDocuments{
				ProofOfAddress: []string{"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=="},
			},
			TosAcceptance: &connect.SubAccountTosAcceptance{
				IP:           "192.168.1.1",
				Date:         "2026-02-28T00:00:00Z",
				UserAgent:    "uqpay-sdk-go/test",
				TosAgreement: 1,
			},
		}

		resp, err := client.Connect.Accounts.CreateSubAccount(ctx, req)
		if err != nil {
			t.Fatalf("CreateSubAccount Individual failed: %v", err)
		}

		t.Logf("✅ Created Individual sub-account:")
		t.Logf("  AccountID: %s", resp.AccountID)
		t.Logf("  ShortReferenceID: %s", resp.ShortReferenceID)
		t.Logf("  Status: %s", resp.Status)
		t.Logf("  VerificationStatus: %s", resp.VerificationStatus)

		if resp.AccountID == "" {
			t.Error("AccountID should not be empty")
		}
		if resp.Status == "" {
			t.Error("Status should not be empty")
		}
	})

	t.Run("Company", func(t *testing.T) {
		inheritNo := -1
		req := &connect.CreateSubAccountRequest{
			EntityType: connect.EntityTypeCompany,
			Nickname:   "SDK Test Company",
			Inherit:    &inheritNo,
			CompanyInfo: &connect.SubAccountCompanyInfo{
				LegalBusinessName:        "SDK Test Ltd",
				LegalBusinessNameEnglish: "SDK Test Ltd",
				CountryOfIncorporation:   "GB",
				CompanyType:              connect.CompanyTypeLimitedCompany,
				PhoneNumber:              "+442071234567",
				EmailAddress:             "company.sdktest@example.com",
				CompanyRegistrationNumber: "12345678",
				TaxType:                  connect.TaxTypeVAT,
				TaxNumber:               "GB123456789",
				IncorporateDate:          "2020-06-15",
				CertificationOfIncorporation: []string{"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=="},
			},
			CompanyAddress: &connect.SubAccountAddress{
				StreetAddress: "456 Oxford Street",
				City:          "London",
				State:         "England",
				PostalCode:    "W1C 1AP",
			},
			OwnershipDetails: &connect.SubAccountOwnershipDetails{
				Representatives: []connect.SubAccountRepresentative{
					{
						LegalFirstNameEnglish: "Jane",
						LegalLastNameEnglish:  "Smith",
						EmailAddress:          "jane.smith@example.com",
						IsApplicant:           "1",
						JobTitle:              connect.JobTitleBeneficialOwnerAndDirector,
						OwnershipPercentage:   100,
						Nationality:           "GB",
						PhoneNumber:           "+447911654321",
						DateOfBirth:           "1985-03-20",
						CountryOrTerritory:    "GB",
						StreetAddress:         "789 Park Lane",
						City:                  "London",
						State:                 "England",
						PostalCode:            "W1K 1PN",
						IdentificationType:    connect.SubAccountIDTypePassport,
						IdentificationValue:   "P87654321",
						IdentityDocs:          []string{"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=="},
						FaceDocs:              []string{"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=="},
					},
				},
				ShareholderDocs: []string{"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=="},
			},
			BusinessDetails: &connect.SubAccountBusinessDetails{
				CountryOrTerritory: "GB",
				StreetAddress:      "456 Oxford Street",
				City:               "London",
				State:              "England",
				PostalCode:         "W1C 1AP",
				Industry:           "7372",
				TurnoverMonthly:    connect.TurnoverMonthlyTM003,
				NumberOfEmployee:   connect.NumberOfEmployeeBS002,
				WebsiteURL:         "https://sdktest.example.com",
				CompanyDescription: "Software development and testing",
				AccountPurpose:     []string{connect.SubAccountPurposeBillPayment},
				BankingCurrencies:  []string{"GBP", "USD", "EUR"},
				BankingCountries:   []string{"GB", "US", "DE"},
			},
			TosAcceptance: &connect.SubAccountTosAcceptance{
				IP:           "192.168.1.1",
				Date:         "2026-02-28T00:00:00Z",
				UserAgent:    "uqpay-sdk-go/test",
				TosAgreement: 1,
			},
		}

		resp, err := client.Connect.Accounts.CreateSubAccount(ctx, req)
		if err != nil {
			t.Fatalf("CreateSubAccount Company failed: %v", err)
		}

		t.Logf("✅ Created Company sub-account:")
		t.Logf("  AccountID: %s", resp.AccountID)
		t.Logf("  ShortReferenceID: %s", resp.ShortReferenceID)
		t.Logf("  Status: %s", resp.Status)
		t.Logf("  VerificationStatus: %s", resp.VerificationStatus)

		if resp.AccountID == "" {
			t.Error("AccountID should not be empty")
		}
		if resp.Status == "" {
			t.Error("Status should not be empty")
		}
	})

	t.Run("Validation_MissingTos", func(t *testing.T) {
		req := &connect.CreateSubAccountRequest{
			EntityType: connect.EntityTypeIndividual,
			Nickname:   "Missing TOS Test",
			IndividualInfo: &connect.SubAccountIndividualInfo{
				FirstNameEnglish:   "Test",
				LastNameEnglish:    "User",
				Nationality:        "GB",
				PhoneNumber:        "+447911000000",
				EmailAddress:       "test@example.com",
				DateOfBirth:        "1990-01-01",
				CountryOrTerritory: "GB",
				StreetAddress:      "1 Test Street",
				City:               "London",
				PostalCode:         "W1A 1AA",
			},
			IdentityVerification: &connect.SubAccountIdentityVerification{
				IdentificationType:  connect.SubAccountIDTypePassport,
				IdentificationValue: "P00000000",
				IdentityDocs:        []string{"doc1"},
				FaceDocs:            []string{"face1"},
			},
			ExpectedActivity: &connect.SubAccountExpectedActivity{
				AccountPurpose:          []string{connect.SubAccountPurposePurchase},
				BankingCountries:        []string{"GB"},
				BankingCurrencies:       []string{"GBP"},
				Internationally:         0,
				TurnoverMonthly:         connect.TurnoverMonthlyTM001,
				TurnoverMonthlyCurrency: "GBP",
			},
			ProofDocuments: &connect.SubAccountProofDocuments{
				ProofOfAddress: []string{"doc1"},
			},
			// TosAcceptance intentionally nil
		}

		_, err := client.Connect.Accounts.CreateSubAccount(ctx, req)
		if err == nil {
			t.Error("Expected error for missing TosAcceptance, got nil")
		} else {
			t.Logf("✅ Correctly rejected missing TosAcceptance: %v", err)
		}
	})

	t.Run("Validation_IndividualMissingFields", func(t *testing.T) {
		req := &connect.CreateSubAccountRequest{
			EntityType: connect.EntityTypeIndividual,
			Nickname:   "Missing Fields Test",
			// Missing required fields for Individual
			TosAcceptance: &connect.SubAccountTosAcceptance{
				IP:   "192.168.1.1",
				Date: "2026-02-28T00:00:00Z",
			},
		}

		_, err := client.Connect.Accounts.CreateSubAccount(ctx, req)
		if err == nil {
			t.Error("Expected error for missing Individual fields, got nil")
		} else {
			t.Logf("✅ Correctly rejected missing Individual fields: %v", err)
		}
	})
}
