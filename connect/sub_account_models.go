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
	JobTitleDirector                    = "DIRECTOR"
	JobTitleBeneficialOwner             = "BENEFICIAL_OWNER"
	JobTitleBeneficialOwnerAndDirector  = "BENEFICIAL_OWNER_AND_DIRECTOR"
	JobTitleAuthorisedPerson            = "AUTHORISED_PERSON"
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
	// EntityType must be "INDIVIDUAL" or "COMPANY"
	EntityType EntityType `json:"entity_type"`

	// Nickname is the account nickname (max 100 characters)
	Nickname string `json:"nickname"`

	// Inherit controls whether to inherit from master account. 1 = inherit, -1 = do not inherit.
	// Only applicable for COMPANY entity type.
	Inherit *int `json:"inherit,omitempty"`

	// IndividualInfo contains personal details. Required when EntityType is INDIVIDUAL.
	IndividualInfo *SubAccountIndividualInfo `json:"individual_info,omitempty"`

	// IdentityVerification contains identity document details. Required when EntityType is INDIVIDUAL.
	IdentityVerification *SubAccountIdentityVerification `json:"identity_verification,omitempty"`

	// ExpectedActivity contains expected business activity. Required when EntityType is INDIVIDUAL.
	ExpectedActivity *SubAccountExpectedActivity `json:"expected_activity,omitempty"`

	// ProofDocuments contains KYC proof documents. Required when EntityType is INDIVIDUAL.
	ProofDocuments *SubAccountProofDocuments `json:"proof_documents,omitempty"`

	// CompanyInfo contains company details. Required when EntityType is COMPANY (unless inherit=1).
	CompanyInfo *SubAccountCompanyInfo `json:"company_info,omitempty"`

	// CompanyAddress is the company registered address. Required when EntityType is COMPANY (unless inherit=1).
	CompanyAddress *SubAccountAddress `json:"company_address,omitempty"`

	// OwnershipDetails contains representatives and shareholder documents.
	// Required when EntityType is COMPANY (unless inherit=1).
	OwnershipDetails *SubAccountOwnershipDetails `json:"ownership_details,omitempty"`

	// BusinessDetails contains business operation details. Applicable for COMPANY entity type.
	BusinessDetails *SubAccountBusinessDetails `json:"business_details,omitempty"`

	// AdditionalDocuments contains extra required/optional documents. Applicable for COMPANY entity type.
	AdditionalDocuments *SubAccountAdditionalDocuments `json:"additional_documents,omitempty"`

	// TosAcceptance contains terms of service acceptance details. Required.
	TosAcceptance *SubAccountTosAcceptance `json:"tos_acceptance"`
}

// SubAccountIndividualInfo contains personal details for an individual sub-account.
type SubAccountIndividualInfo struct {
	// FirstNameEnglish is the first name in English
	FirstNameEnglish string `json:"first_name_english"`

	// LastNameEnglish is the last name in English
	LastNameEnglish string `json:"last_name_english"`

	// Nationality is the ISO 3166-1 alpha-2 country code of nationality
	Nationality string `json:"nationality"`

	// TaxNumber is the tax identification number
	TaxNumber string `json:"tax_number,omitempty"`

	// PhoneNumber is the phone number with country code
	PhoneNumber string `json:"phone_number"`

	// EmailAddress is the email address
	EmailAddress string `json:"email_address"`

	// DateOfBirth is the date of birth in YYYY-MM-DD format
	DateOfBirth string `json:"date_of_birth"`

	// CountryOrTerritory is the ISO 3166-1 alpha-2 country code of residence
	CountryOrTerritory string `json:"country_or_territory"`

	// StreetAddress is the street address
	StreetAddress string `json:"street_address"`

	// ApartmentSuiteOrFloor is the apartment, suite, or floor number
	ApartmentSuiteOrFloor string `json:"apartment_suite_or_floor,omitempty"`

	// City is the city name
	City string `json:"city"`

	// State is the state or province
	State string `json:"state,omitempty"`

	// PostalCode is the postal or ZIP code
	PostalCode string `json:"postal_code"`
}

// SubAccountIdentityVerification contains identity verification details.
type SubAccountIdentityVerification struct {
	// IdentificationType is the type of ID document: PASSPORT, DRIVERS_LICENSE, or NATIONAL_ID
	IdentificationType string `json:"identification_type"`

	// IdentificationValue is the ID document number
	IdentificationValue string `json:"identification_value"`

	// IdentityDocs is a list of identity document images (base64 encoded or file IDs)
	IdentityDocs []string `json:"identity_docs"`

	// FaceDocs is a list of face photo images (base64 encoded or file IDs).
	// Mandatory for individual accounts.
	FaceDocs []string `json:"face_docs"`
}

// SubAccountExpectedActivity contains expected business activity for the account.
type SubAccountExpectedActivity struct {
	// AccountPurpose is a list of intended account purposes.
	// Values: PURCHASE, BILL_PAYMENT, EDUCATIONAL_EXPENSES, PERSONAL_REMITTANCE,
	// CHARITABLE_DONATION, LOAN_REPAYMENT, INVESTMENT, OTHERS
	AccountPurpose []string `json:"account_purpose"`

	// OtherPurpose is required when AccountPurpose includes "OTHERS"
	OtherPurpose string `json:"other_purpose,omitempty"`

	// BankingCountries is a list of ISO 3166-1 alpha-2 country codes for banking operations
	BankingCountries []string `json:"banking_countries"`

	// BankingCurrencies is a list of ISO 4217 currency codes for banking operations
	BankingCurrencies []string `json:"banking_currencies"`

	// Internationally indicates if banking will be international. 0 = no, 1 = yes.
	Internationally int `json:"internationally"`

	// TurnoverMonthly is the estimated monthly turnover tier (TM001-TM005)
	TurnoverMonthly string `json:"turnover_monthly"`

	// TurnoverMonthlyCurrency is the ISO 4217 currency code for turnover
	TurnoverMonthlyCurrency string `json:"turnover_monthly_currency"`
}

// SubAccountProofDocuments contains KYC proof documents.
type SubAccountProofDocuments struct {
	// ProofOfAddress is a list of address proof documents (base64 or file IDs). Required.
	ProofOfAddress []string `json:"proof_of_address"`

	// SourceOfFunds is a list of source-of-funds documents (base64 or file IDs).
	// Required for Virtual Account applications.
	SourceOfFunds []string `json:"source_of_funds,omitempty"`

	// ProofOfPositionAndIncome is a list of position/income proof documents
	ProofOfPositionAndIncome []string `json:"proof_of_position_and_income,omitempty"`

	// OtherProof is a list of other supporting documents
	OtherProof []string `json:"other_proof,omitempty"`
}

// SubAccountCompanyInfo contains company details for a company sub-account.
type SubAccountCompanyInfo struct {
	// LegalBusinessName is the legal business name in local language
	LegalBusinessName string `json:"legal_business_name"`

	// LegalBusinessNameEnglish is the legal business name in English (ASCII, max 255)
	LegalBusinessNameEnglish string `json:"legal_business_name_english"`

	// CountryOfIncorporation is the ISO 3166-1 alpha-2 country code
	CountryOfIncorporation string `json:"country_of_incorporation"`

	// CompanyType is the type of company.
	// Values: SOLE_PROPRIETOR, LIMITED_COMPANY, PARTNERSHIP, LISTED, OTHERS
	CompanyType string `json:"company_type"`

	// PhoneNumber is the company phone number with country code
	PhoneNumber string `json:"phone_number"`

	// EmailAddress is the company email address
	EmailAddress string `json:"email_address"`

	// CompanyRegistrationNumber is the business registration number
	CompanyRegistrationNumber string `json:"company_registration_number"`

	// TaxType is the type of tax identifier: VAT, GST, or TAX
	TaxType string `json:"tax_type,omitempty"`

	// TaxNumber is the tax identification number
	TaxNumber string `json:"tax_number,omitempty"`

	// IncorporateDate is the date of incorporation in YYYY-MM-DD format
	IncorporateDate string `json:"incorparate_date"`

	// CertificationOfIncorporation is a list of incorporation certificate documents (base64 or file IDs)
	CertificationOfIncorporation []string `json:"certification_of_incorporation"`
}

// SubAccountAddress represents a physical address for sub-account requests.
type SubAccountAddress struct {
	// StreetAddress is the street address
	StreetAddress string `json:"street_address"`

	// ApartmentSuiteOrFloor is the apartment, suite, or floor number
	ApartmentSuiteOrFloor string `json:"apartment_suite_or_floor,omitempty"`

	// City is the city name
	City string `json:"city"`

	// State is the state or province
	State string `json:"state,omitempty"`

	// PostalCode is the postal or ZIP code
	PostalCode string `json:"postal_code"`
}

// SubAccountOwnershipDetails contains company ownership information.
type SubAccountOwnershipDetails struct {
	// Representatives is a list of company representatives/directors/UBOs
	Representatives []SubAccountRepresentative `json:"representatives"`

	// ShareholderDocs is a list of shareholder documents (base64 or file IDs)
	ShareholderDocs []string `json:"shareholder_docs"`
}

// SubAccountRepresentative represents a company representative/director/UBO.
type SubAccountRepresentative struct {
	// LegalFirstNameEnglish is the first name in English
	LegalFirstNameEnglish string `json:"legal_first_name_english"`

	// LegalLastNameEnglish is the last name in English
	LegalLastNameEnglish string `json:"legal_last_name_english"`

	// NameInOtherLanguage is the name in local/other language
	NameInOtherLanguage string `json:"name_in_other_language,omitempty"`

	// EmailAddress is the representative's email
	EmailAddress string `json:"email_address"`

	// IsApplicant indicates if this representative is the applicant. "0" or "1".
	// Only one representative can have IsApplicant = "1".
	IsApplicant string `json:"is_applicant"`

	// JobTitle is the representative's role.
	// Values: DIRECTOR, BENEFICIAL_OWNER, BENEFICIAL_OWNER_AND_DIRECTOR, AUTHORISED_PERSON
	JobTitle string `json:"job_title"`

	// OwnershipPercentage is the ownership percentage
	OwnershipPercentage string `json:"ownership_percentage,omitempty"`

	// Nationality is the ISO 3166-1 alpha-2 country code
	Nationality string `json:"nationality"`

	// TaxNumber is the tax identification number
	TaxNumber string `json:"tax_number,omitempty"`

	// PhoneNumber is the phone number with country code
	PhoneNumber string `json:"phone_number"`

	// DateOfBirth is the date of birth in YYYY-MM-DD format
	DateOfBirth string `json:"date_of_birth"`

	// CountryOrTerritory is the ISO 3166-1 alpha-2 country code of residence
	CountryOrTerritory string `json:"country_or_territory"`

	// StreetAddress is the street address
	StreetAddress string `json:"street_address"`

	// ApartmentSuiteOrFloor is the apartment, suite, or floor number
	ApartmentSuiteOrFloor string `json:"apartment_suite_or_floor,omitempty"`

	// City is the city name
	City string `json:"city"`

	// State is the state or province
	State string `json:"state,omitempty"`

	// PostalCode is the postal or ZIP code
	PostalCode string `json:"postal_code"`

	// IdentificationType is the type of ID document
	IdentificationType string `json:"identification_type"`

	// IdentificationValue is the ID document number
	IdentificationValue string `json:"identification_value"`

	// IdentityDocs is a list of identity document images (base64 or file IDs)
	IdentityDocs []string `json:"identity_docs"`

	// OtherDocuments is a list of additional documents
	OtherDocuments []SubAccountRepresentativeDocument `json:"other_documents,omitempty"`

	// FaceDocs is a list of face photo images (base64 or file IDs).
	// Required for DIRECTOR and BENEFICIAL_OWNER roles.
	FaceDocs []string `json:"face_docs,omitempty"`
}

// SubAccountRepresentativeDocument represents an additional document for a representative.
type SubAccountRepresentativeDocument struct {
	// Type is the document type identifier
	Type string `json:"type"`

	// DocStr is the document content (base64 encoded or file ID)
	DocStr string `json:"doc_str"`
}

// SubAccountBusinessDetails contains business operation details for a company.
type SubAccountBusinessDetails struct {
	// CountryOrTerritory is the ISO 3166-1 alpha-2 country code for business operations
	CountryOrTerritory string `json:"country_or_territory,omitempty"`

	// StreetAddress is the business operating street address
	StreetAddress string `json:"street_address,omitempty"`

	// City is the city name
	City string `json:"city,omitempty"`

	// State is the state or province
	State string `json:"state,omitempty"`

	// PostalCode is the postal or ZIP code
	PostalCode string `json:"postal_code,omitempty"`

	// Industry is the numeric industry classification code
	Industry string `json:"industry,omitempty"`

	// TurnoverMonthly is the estimated monthly turnover tier (TM001-TM005)
	TurnoverMonthly string `json:"turnover_monthly,omitempty"`

	// NumberOfEmployee is the employee count tier (BS001-BS005)
	NumberOfEmployee string `json:"number_of_employee,omitempty"`

	// WebsiteURL is the business website URL (max 100 characters)
	WebsiteURL string `json:"website_url,omitempty"`

	// CompanyDescription is a description of the company's business
	CompanyDescription string `json:"company_description,omitempty"`

	// AccountPurpose is a list of intended account purposes
	AccountPurpose []string `json:"account_purpose,omitempty"`

	// BankingCurrencies is a list of ISO 4217 currency codes for banking operations
	BankingCurrencies []string `json:"banking_currencies,omitempty"`

	// BankingCountries is a list of ISO 3166-1 alpha-2 country codes for banking operations
	BankingCountries []string `json:"banking_countries,omitempty"`

	// IssuingCountries is a list of ISO 3166-1 alpha-2 country codes for card issuing.
	// Required for ISSUING business.
	IssuingCountries []string `json:"issuing_countries,omitempty"`

	// IssuingMonthly is the estimated monthly issuing volume tier (TM001-TM005).
	// Required for ISSUING business.
	IssuingMonthly string `json:"issuing_monthly,omitempty"`
}

// SubAccountAdditionalDocuments contains extra documents for the sub-account application.
type SubAccountAdditionalDocuments struct {
	// RequiredDocs is a list of required additional documents
	RequiredDocs []SubAccountAdditionalDocument `json:"required_docs,omitempty"`

	// OptionDocs is a list of optional additional documents
	OptionDocs []SubAccountAdditionalDocument `json:"option_docs,omitempty"`
}

// SubAccountAdditionalDocument represents a single additional document.
type SubAccountAdditionalDocument struct {
	// ProfileKey is the document profile key identifier
	ProfileKey string `json:"profile_key"`

	// DocStr is the document content (base64 encoded or file ID)
	DocStr string `json:"doc_str"`
}

// SubAccountTosAcceptance contains terms of service acceptance details.
type SubAccountTosAcceptance struct {
	// IP is the IPv4 address of the user accepting the ToS
	IP string `json:"ip"`

	// Date is the acceptance date in ISO 8601 format
	Date string `json:"date"`

	// UserAgent is the browser user agent string
	UserAgent string `json:"user_agent,omitempty"`

	// TosAgreement controls auto-signing of the TPSP agreement. Set to 1 to auto-sign.
	TosAgreement int `json:"tos_agreement,omitempty"`
}

// ============================================================
// Response Models
// ============================================================

// CreateSubAccountResponse represents the response from POST /v1/accounts/create_accounts.
type CreateSubAccountResponse struct {
	// AccountID is the unique identifier for the created account (UUID)
	AccountID string `json:"account_id"`

	// ShortReferenceID is a short reference identifier for the account
	ShortReferenceID string `json:"short_reference_id"`

	// Status is the account status: ACTIVE, PROCESSING, INACTIVE, or CLOSED
	Status string `json:"status"`

	// VerificationStatus is the KYC verification status: APPROVED, PENDING, REJECT, EXPIRED, or RETURN
	VerificationStatus string `json:"verification_status"`
}
