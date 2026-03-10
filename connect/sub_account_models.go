package connect

// ============================================================
// Constants - Company Types
// ============================================================

const (
	CompanyTypeSoleProprietor = "SOLE_PROPRIETOR"
	CompanyTypeLimitedCompany = "LIMITED_COMPANY"
	CompanyTypePartnership    = "PARTNERSHIP"
	CompanyTypeListed         = "LISTED"
	CompanyTypeOthers         = "OTHERS"
)

// ============================================================
// Constants - Job Titles (Representative Roles)
// ============================================================

const (
	JobTitleDirector                   = "DIRECTOR"
	JobTitleBeneficialOwner            = "BENEFICIAL_OWNER"
	JobTitleBeneficialOwnerAndDirector = "BENEFICIAL_OWNER_AND_DIRECTOR"
	JobTitleAuthorisedPerson           = "AUTHORISED_PERSON"
)

// ============================================================
// Constants - Turnover Monthly Tiers
// ============================================================

const (
	TurnoverMonthlyTM001 = "TM001" // < $50,000
	TurnoverMonthlyTM002 = "TM002" // $50,000 - $100,000
	TurnoverMonthlyTM003 = "TM003" // $100,000 - $250,000
	TurnoverMonthlyTM004 = "TM004" // $250,000 - $500,000
	TurnoverMonthlyTM005 = "TM005" // > $500,000
)

// ============================================================
// Constants - Number of Employees
// ============================================================

const (
	NumberOfEmployeeBS001 = "BS001" // 0-1
	NumberOfEmployeeBS002 = "BS002" // 2-10
	NumberOfEmployeeBS003 = "BS003" // 11-50
	NumberOfEmployeeBS004 = "BS004" // 51-200
	NumberOfEmployeeBS005 = "BS005" // > 200
)

// ============================================================
// Constants - Account Purpose
// ============================================================

const (
	SubAccountPurposePurchase           = "PURCHASE"
	SubAccountPurposeBillPayment        = "BILL_PAYMENT"
	SubAccountPurposeEducationalExpense = "EDUCATIONAL_EXPENSES"
	SubAccountPurposePersonalRemittance = "PERSONAL_REMITTANCE"
	SubAccountPurposeCharitableDonation = "CHARITABLE_DONATION"
	SubAccountPurposeLoanRepayment      = "LOAN_REPAYMENT"
	SubAccountPurposeInvestment         = "INVESTMENT"
	SubAccountPurposeOthers             = "OTHERS"
)

// ============================================================
// Constants - Identification Types
// ============================================================

const (
	SubAccountIDTypePassport       = "PASSPORT"
	SubAccountIDTypeDriversLicense = "DRIVERS_LICENSE"
	SubAccountIDTypeNationalID     = "NATIONAL_ID"
)

// ============================================================
// Constants - Tax Types
// ============================================================

const (
	TaxTypeVAT = "VAT"
	TaxTypeGST = "GST"
	TaxTypeTAX = "TAX"
)

// ============================================================
// Constants - SubAccount Status (Response)
// ============================================================

const (
	SubAccountStatusActive     = "ACTIVE"
	SubAccountStatusProcessing = "PROCESSING"
	SubAccountStatusInactive   = "INACTIVE"
	SubAccountStatusClosed     = "CLOSED"
)

// ============================================================
// Constants - Verification Status (Response)
// ============================================================

const (
	SubAccountVerificationApproved = "APPROVED"
	SubAccountVerificationPending  = "PENDING"
	SubAccountVerificationReject   = "REJECT"
	SubAccountVerificationExpired  = "EXPIRED"
	SubAccountVerificationReturn   = "RETURN"
)

// ============================================================
// Request Models
// ============================================================

// CreateSubAccountRequest represents the request body for POST /v1/accounts/create_accounts.
// This is a discriminated union: populate IndividualInfo for INDIVIDUAL entity type,
// or CompanyInfo + OwnershipDetails + BusinessDetails for COMPANY entity type.
type CreateSubAccountRequest struct {
	BusinessType         string                          `json:"business_type"`                   // Required. Business line: BANKING, ACQUIRING, or ISSUING
	EntityType           EntityType                      `json:"entity_type"`                     // Required. INDIVIDUAL or COMPANY
	Nickname             string                          `json:"nickname"`                        // Required. Account display name, max 100 chars
	Inherit              *int                            `json:"inherit,omitempty"`               // Optional. COMPANY only: 1 = inherit from master, -1 = do not inherit
	IndividualInfo       *SubAccountIndividualInfo       `json:"individual_info,omitempty"`       // Required for INDIVIDUAL entity type
	IdentityVerification *SubAccountIdentityVerification `json:"identity_verification,omitempty"` // Required for INDIVIDUAL entity type
	ExpectedActivity     *SubAccountExpectedActivity     `json:"expected_activity,omitempty"`     // Required for INDIVIDUAL entity type
	ProofDocuments       *SubAccountProofDocuments       `json:"proof_documents,omitempty"`       // Required for INDIVIDUAL entity type
	CompanyInfo          *SubAccountCompanyInfo          `json:"company_info,omitempty"`          // Required for COMPANY (unless inherit=1)
	CompanyAddress       *SubAccountAddress              `json:"company_address,omitempty"`       // Required for COMPANY (unless inherit=1)
	OwnershipDetails     *SubAccountOwnershipDetails     `json:"ownership_details,omitempty"`     // Required for COMPANY (unless inherit=1)
	BusinessDetails      *SubAccountBusinessDetails      `json:"business_details,omitempty"`      // Required for COMPANY entity type
	AdditionalDocuments  *SubAccountAdditionalDocuments  `json:"additional_documents,omitempty"`  // Optional. Extra required/optional documents for COMPANY
	TosAcceptance        *SubAccountTosAcceptance        `json:"tos_acceptance"`                  // Required. Terms of service acceptance
}

// SubAccountIndividualInfo contains personal details for an individual sub-account.
type SubAccountIndividualInfo struct {
	FirstNameEnglish      string `json:"first_name_english"`                 // Required. Given name in English, max 100 chars
	LastNameEnglish       string `json:"last_name_english"`                  // Required. Family name in English, max 100 chars
	Nationality           string `json:"nationality"`                        // Required. ISO 3166-1 alpha-2 country code
	TaxNumber             string `json:"tax_number,omitempty"`               // Optional. Tax identification number, max 100 chars
	PhoneNumber           string `json:"phone_number"`                       // Required. With country code, max 25 chars, e.g. "+12345678"
	EmailAddress          string `json:"email_address"`                      // Required. Valid email, max 100 chars
	DateOfBirth           string `json:"date_of_birth"`                      // Required. Format: YYYY-MM-DD
	CountryOrTerritory    string `json:"country_or_territory"`               // Required. ISO 3166-1 alpha-2 residence country code
	StreetAddress         string `json:"street_address"`                     // Required. Street address, max 100 chars
	ApartmentSuiteOrFloor string `json:"apartment_suite_or_floor,omitempty"` // Optional. Unit/suite/floor, max 100 chars
	City                  string `json:"city"`                               // Required. City name, max 100 chars
	State                 string `json:"state,omitempty"`                    // Optional. State/province code
	PostalCode            string `json:"postal_code"`                        // Required. ZIP/postal code, max 100 chars
}

// SubAccountIdentityVerification contains identity verification details.
type SubAccountIdentityVerification struct {
	IdentificationType  string   `json:"identification_type"`  // Required. PASSPORT, DRIVERS_LICENSE, or NATIONAL_ID
	IdentificationValue string   `json:"identification_value"` // Required. ID document number, max 100 chars
	IdentityDocs        []string `json:"identity_docs"`        // Required. Identity document images (base64 or file IDs)
	FaceDocs            []string `json:"face_docs"`            // Required for individuals. Facial verification images (base64 or file IDs)
}

// SubAccountExpectedActivity contains expected business activity for the account.
type SubAccountExpectedActivity struct {
	AccountPurpose          []string `json:"account_purpose"`           // Required. PURCHASE, BILL_PAYMENT, EDUCATIONAL_EXPENSES, PERSONAL_REMITTANCE, CHARITABLE_DONATION, LOAN_REPAYMENT, INVESTMENT, OTHERS
	OtherPurpose            string   `json:"other_purpose,omitempty"`   // Required if account_purpose includes OTHERS
	BankingCountries        []string `json:"banking_countries"`         // Required. ISO 3166-1 alpha-2 country codes
	BankingCurrencies       []string `json:"banking_currencies"`        // Required. ISO 4217 currency codes
	Internationally         int      `json:"internationally"`           // Required. 0 = domestic only, 1 = international
	TurnoverMonthly         string   `json:"turnover_monthly"`          // Required. TM001-TM005 (monthly revenue bracket)
	TurnoverMonthlyCurrency string   `json:"turnover_monthly_currency"` // Required. ISO 4217 currency code for turnover estimate
}

// SubAccountProofDocuments contains KYC proof documents.
type SubAccountProofDocuments struct {
	ProofOfAddress           []string `json:"proof_of_address"`                       // Optional. Address verification docs (base64 or file IDs)
	SourceOfFunds            []string `json:"source_of_funds,omitempty"`              // Optional. Source-of-funds docs; required for Virtual Account applications
	ProofOfPositionAndIncome []string `json:"proof_of_position_and_income,omitempty"` // Optional. Employment/income proof docs (base64 or file IDs)
	OtherProof               []string `json:"other_proof,omitempty"`                  // Optional. Miscellaneous supporting docs (base64 or file IDs)
}

// SubAccountCompanyInfo contains company details for a company sub-account.
type SubAccountCompanyInfo struct {
	LegalBusinessName            string   `json:"legal_business_name"`            // Required. Business name in local language, max 100 chars
	LegalBusinessNameEnglish     string   `json:"legal_business_name_english"`    // Required. Business name in English (ASCII only), max 100 chars
	CountryOfIncorporation       string   `json:"country_of_incorporation"`       // Required. ISO 3166-1 alpha-2 country code
	CompanyType                  string   `json:"company_type"`                   // Required. SOLE_PROPRIETOR, LIMITED_COMPANY, PARTNERSHIP, LISTED, or OTHERS
	PhoneNumber                  string   `json:"phone_number"`                   // Required. With country code, max 25 chars, e.g. "+12345678"
	EmailAddress                 string   `json:"email_address"`                  // Required. Valid email, max 100 chars
	CompanyRegistrationNumber    string   `json:"company_registration_number"`    // Required. Official registration ID, max 100 chars
	TaxType                      string   `json:"tax_type,omitempty"`             // Optional. VAT, GST, or TAX
	TaxNumber                    string   `json:"tax_number,omitempty"`           // Optional. Tax ID number, max 100 chars
	IncorporateDate              string   `json:"incorparate_date"`               // Required. Format: YYYY-MM-DD
	CertificationOfIncorporation []string `json:"certification_of_incorporation"` // Required. Incorporation certificate docs (base64 or file IDs)
}

// SubAccountAddress represents a physical address for sub-account requests.
type SubAccountAddress struct {
	StreetAddress         string `json:"street_address"`                     // Required. Street address, max 100 chars
	ApartmentSuiteOrFloor string `json:"apartment_suite_or_floor,omitempty"` // Optional. Unit/suite/floor, max 100 chars
	City                  string `json:"city"`                               // Required. City name, max 100 chars
	State                 string `json:"state,omitempty"`                    // Required. State/province code
	PostalCode            string `json:"postal_code"`                        // Required. ZIP/postal code, max 100 chars
}

// SubAccountOwnershipDetails contains company ownership information.
type SubAccountOwnershipDetails struct {
	Representatives []SubAccountRepresentative `json:"representatives"`  // Required. Company directors, UBOs, and authorised persons
	ShareholderDocs []string                   `json:"shareholder_docs"` // Required. Shareholder documents (base64 or file IDs)
}

// SubAccountRepresentative represents a company representative/director/UBO.
type SubAccountRepresentative struct {
	LegalFirstNameEnglish string                             `json:"legal_first_name_english"`           // Required. Given name in English, max 100 chars
	LegalLastNameEnglish  string                             `json:"legal_last_name_english"`            // Required. Family name in English, max 100 chars
	NameInOtherLanguage   string                             `json:"name_in_other_language,omitempty"`   // Optional. Name in non-English language, max 100 chars
	EmailAddress          string                             `json:"email_address"`                      // Optional. Contact email, max 100 chars
	IsApplicant           string                             `json:"is_applicant"`                       // Required. "0" or "1"; only one representative may be "1"
	JobTitle              string                             `json:"job_title"`                          // Required. DIRECTOR, BENEFICIAL_OWNER, BENEFICIAL_OWNER_AND_DIRECTOR, or AUTHORISED_PERSON
	OwnershipPercentage   float64                            `json:"ownership_percentage,omitempty"`     // Optional. Ownership stake percentage, e.g. 15.5
	Nationality           string                             `json:"nationality"`                        // Required. ISO 3166-1 alpha-2 country code
	TaxNumber             string                             `json:"tax_number,omitempty"`               // Optional. Tax identification number, max 100 chars
	PhoneNumber           string                             `json:"phone_number"`                       // Optional. With country code, max 25 chars
	DateOfBirth           string                             `json:"date_of_birth"`                      // Required. Format: YYYY-MM-DD
	CountryOrTerritory    string                             `json:"country_or_territory"`               // Required. ISO 3166-1 alpha-2 residence country code
	StreetAddress         string                             `json:"street_address"`                     // Required. Residential street address, max 100 chars
	ApartmentSuiteOrFloor string                             `json:"apartment_suite_or_floor,omitempty"` // Optional. Unit/suite/floor, max 100 chars
	City                  string                             `json:"city"`                               // Required. Residential city, max 100 chars
	State                 string                             `json:"state,omitempty"`                    // Required. State/province code
	PostalCode            string                             `json:"postal_code"`                        // Required. ZIP/postal code, max 100 chars
	IdentificationType    string                             `json:"identification_type"`                // Required. PASSPORT, DRIVERS_LICENSE, or NATIONAL_ID
	IdentificationValue   string                             `json:"identification_value"`               // Required. ID document number, max 100 chars
	IdentityDocs          []string                           `json:"identity_docs"`                      // Required. Identity document images (base64 or file IDs)
	OtherDocuments        []SubAccountRepresentativeDocument `json:"other_documents,omitempty"`          // Optional. Additional supporting documents
	FaceDocs              []string                           `json:"face_docs,omitempty"`                // Optional. Face photos; required for at least one DIRECTOR/BENEFICIAL_OWNER
}

// SubAccountRepresentativeDocument represents an additional document for a representative.
type SubAccountRepresentativeDocument struct {
	Type   string `json:"type"`    // Required. PROOF_OF_ADDRESS or OTHERS
	DocStr string `json:"doc_str"` // Required. Base64-encoded document or file ID
}

// SubAccountBusinessDetails contains business operation details for a company.
type SubAccountBusinessDetails struct {
	CountryOrTerritory string   `json:"country_or_territory,omitempty"` // Required. ISO 3166-1 alpha-2 operating jurisdiction
	StreetAddress      string   `json:"street_address,omitempty"`       // Required. Business street address, max 100 chars
	City               string   `json:"city,omitempty"`                 // Required. Business city, max 100 chars
	State              string   `json:"state,omitempty"`                // Required. State/province code
	PostalCode         string   `json:"postal_code,omitempty"`          // Required. ZIP/postal code, max 100 chars
	Industry           string   `json:"industry,omitempty"`             // Required. Numeric industry classification code (MCC)
	TurnoverMonthly    string   `json:"turnover_monthly,omitempty"`     // Required. TM001-TM005 (monthly revenue bracket)
	NumberOfEmployee   string   `json:"number_of_employee,omitempty"`   // Required. BS001-BS005 (employee count bracket)
	WebsiteURL         string   `json:"website_url,omitempty"`          // Optional. Business website URL, max 100 chars
	CompanyDescription string   `json:"company_description,omitempty"`  // Optional. Business summary
	AccountPurpose     []string `json:"account_purpose,omitempty"`      // Optional. BUSINESS_PAYMENT, BILL_PAYMENT, CHARITABLE_DONATION, LOAN_REPAYMENT, INVESTMENT, COLLECTION_OF_BUSINESS, OTHERS
	BankingCurrencies  []string `json:"banking_currencies,omitempty"`   // Optional. ISO 4217 currency codes, e.g. "USD", "EUR"
	BankingCountries   []string `json:"banking_countries,omitempty"`    // Optional. ISO 3166-1 alpha-2 country codes
	IssuingCountries   []string `json:"issuing_countries,omitempty"`    // Optional. ISO 3166-1 alpha-2 codes; required for ISSUING business
	IssuingMonthly     string   `json:"issuing_monthly,omitempty"`      // Optional. TM001-TM005; required for ISSUING business
}

// SubAccountAdditionalDocuments contains extra documents for the sub-account application.
type SubAccountAdditionalDocuments struct {
	RequiredDocs []SubAccountAdditionalDocument `json:"required_docs,omitempty"` // Mandatory supporting documents (retrieve keys via Get Additional Documents)
	OptionDocs   []SubAccountAdditionalDocument `json:"option_docs,omitempty"`   // Optional supporting documents
}

// SubAccountAdditionalDocument represents a single additional document.
type SubAccountAdditionalDocument struct {
	ProfileKey string `json:"profile_key"` // Required. Document profile key, e.g. "business_articles_of_association"
	DocStr     string `json:"doc_str"`     // Required. Base64-encoded document or file ID
}

// SubAccountTosAcceptance contains terms of service acceptance details.
type SubAccountTosAcceptance struct {
	IP           string `json:"ip"`                      // Required. IPv4 address of the accepting user
	Date         string `json:"date"`                    // Required. Acceptance timestamp in ISO 8601 format
	UserAgent    string `json:"user_agent,omitempty"`    // Optional. Browser user agent string
	TosAgreement int    `json:"tos_agreement,omitempty"` // Optional. Set to 1 to auto-sign TPSP agreement (TPSP master accounts only)
}

// ============================================================
// Response Models
// ============================================================

// CreateSubAccountResponse represents the response from POST /v1/accounts/create_accounts.
type CreateSubAccountResponse struct {
	AccountID          string `json:"account_id"`          // Unique account identifier (UUID)
	ShortReferenceID   string `json:"short_reference_id"`  // Human-readable reference, e.g. "P220406-LLCVLRM"
	Status             string `json:"status"`              // ACTIVE, PROCESSING, INACTIVE, or CLOSED
	VerificationStatus string `json:"verification_status"` // APPROVED, PENDING, REJECT, EXPIRED, or RETURN
}
