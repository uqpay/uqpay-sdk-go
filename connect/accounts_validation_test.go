package connect

import (
	"context"
	"strings"
	"testing"
)

func completeCompanySubAccountRequest() *CreateSubAccountRequest {
	inherit := -1
	return &CreateSubAccountRequest{
		EntityType:     EntityTypeCompany,
		Inherit:        &inherit,
		CompanyInfo:    &SubAccountCompanyInfo{},
		CompanyAddress: &SubAccountAddress{},
		OwnershipDetails: &SubAccountOwnershipDetails{
			Representatives: []SubAccountRepresentative{{
				EmailAddress:        "representative@example.com",
				DateOfBirth:         "1985-03-20",
				OwnershipPercentage: "0",
			}},
		},
		BusinessDetails: &SubAccountBusinessDetails{
			AccountPurpose:        []SubAccountCompanyPurpose{CompanyPurposePaymentCollection},
			BankingCurrencies:     []string{"SGD"},
			BankingCountries:      []string{"SG"},
			ArticlesOfAssociation: []string{"file-id"},
		},
		TosAcceptance: &SubAccountTosAcceptance{},
	}
}

func TestValidateCreateSubAccountRequestCompanyContract(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*CreateSubAccountRequest)
		wantErr string
	}{
		{
			name: "missing representatives",
			mutate: func(req *CreateSubAccountRequest) {
				req.OwnershipDetails.Representatives = nil
			},
			wantErr: "ownership_details.representatives required",
		},
		{
			name: "missing representative email",
			mutate: func(req *CreateSubAccountRequest) {
				req.OwnershipDetails.Representatives[0].EmailAddress = ""
			},
			wantErr: "representatives[0].email_address is required",
		},
		{
			name: "missing representative date of birth",
			mutate: func(req *CreateSubAccountRequest) {
				req.OwnershipDetails.Representatives[0].DateOfBirth = ""
			},
			wantErr: "representatives[0].date_of_birth is required",
		},
		{
			name: "missing representative ownership percentage",
			mutate: func(req *CreateSubAccountRequest) {
				req.OwnershipDetails.Representatives[0].OwnershipPercentage = ""
			},
			wantErr: "representatives[0].ownership_percentage is required",
		},
		{
			name: "missing business details",
			mutate: func(req *CreateSubAccountRequest) {
				req.BusinessDetails = nil
			},
			wantErr: "business_details required",
		},
		{
			name: "missing account purpose",
			mutate: func(req *CreateSubAccountRequest) {
				req.BusinessDetails.AccountPurpose = nil
			},
			wantErr: "business_details.account_purpose is required",
		},
		{
			name: "missing banking currencies",
			mutate: func(req *CreateSubAccountRequest) {
				req.BusinessDetails.BankingCurrencies = nil
			},
			wantErr: "business_details.banking_currencies is required",
		},
		{
			name: "missing banking countries",
			mutate: func(req *CreateSubAccountRequest) {
				req.BusinessDetails.BankingCountries = nil
			},
			wantErr: "business_details.banking_countries is required",
		},
		{
			name: "missing articles of association",
			mutate: func(req *CreateSubAccountRequest) {
				req.BusinessDetails.ArticlesOfAssociation = nil
			},
			wantErr: "business_details.articles_of_association is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := completeCompanySubAccountRequest()
			tt.mutate(req)
			err := validateCreateSubAccountRequest(req)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("validateCreateSubAccountRequest() error = %v, want error containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCreateSubAccountRequestCompanyAcceptsCompleteContract(t *testing.T) {
	if err := validateCreateSubAccountRequest(completeCompanySubAccountRequest()); err != nil {
		t.Fatalf("validateCreateSubAccountRequest() returned an error: %v", err)
	}
}

func TestCreateSubAccountRejectsInvalidCompanyBeforeHTTP(t *testing.T) {
	req := completeCompanySubAccountRequest()
	req.OwnershipDetails.Representatives[0].EmailAddress = ""

	client := &AccountsClient{}
	_, err := client.CreateSubAccount(context.Background(), req)
	if err == nil || !strings.Contains(err.Error(), "representatives[0].email_address is required") {
		t.Fatalf("CreateSubAccount() error = %v, want representative email validation error", err)
	}
}

func TestValidateCreateSubAccountRequestCompanyInheritBypassesCompanyDetails(t *testing.T) {
	inherit := 1
	req := &CreateSubAccountRequest{
		EntityType:    EntityTypeCompany,
		Inherit:       &inherit,
		TosAcceptance: &SubAccountTosAcceptance{},
	}

	if err := validateCreateSubAccountRequest(req); err != nil {
		t.Fatalf("validateCreateSubAccountRequest() returned an error: %v", err)
	}
}
