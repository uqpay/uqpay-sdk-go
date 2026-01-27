package webhook

// AccountData represents the account information in onboarding webhook events.
// This is returned in the data field for onboarding.account.create and
// onboarding.account.update events.
type AccountData struct {
	// AccountID is the unique identifier for the account
	AccountID string `json:"account_id"`

	// DirectID is the direct account identifier
	DirectID string `json:"direct_id,omitempty"`

	// ShortReferenceID is a short reference identifier for the account
	ShortReferenceID string `json:"short_reference_id,omitempty"`

	// Email is the primary email address for the account
	Email string `json:"email"`

	// AccountName is the display name of the account
	AccountName string `json:"account_name"`

	// Country is the ISO 3166-1 alpha-2 country code
	Country string `json:"country"`

	// Status is the current account status (e.g., "PROCESSING", "ACTIVE")
	Status string `json:"status"`

	// IDVStatus is the identity verification status (e.g., "PENDING", "APPROVED")
	IDVStatus string `json:"idv_status,omitempty"`

	// VerificationStatus is the overall verification status
	VerificationStatus string `json:"verification_status,omitempty"`

	// ReviewReason provides the reason for the current review status
	ReviewReason string `json:"review_reason,omitempty"`

	// EntityType is the type of entity (e.g., "COMPANY", "INDIVIDUAL")
	EntityType string `json:"entity_type"`

	// ContactDetails contains contact information for the account
	ContactDetails *ContactDetails `json:"contact_details,omitempty"`

	// BusinessDetails contains business-specific information
	BusinessDetails *BusinessDetails `json:"business_details,omitempty"`

	// RegistrationAddress is the official registration address
	RegistrationAddress *Address `json:"registration_address,omitempty"`

	// BusinessAddress is a list of business operating addresses
	BusinessAddress []Address `json:"business_address,omitempty"`

	// Representatives is a list of account representatives/directors
	Representatives []Representative `json:"representatives,omitempty"`

	// Source indicates where the account was created from (e.g., "web", "api")
	Source string `json:"source,omitempty"`
}

// ContactDetails represents contact information for an account
type ContactDetails struct {
	// Email is the contact email address
	Email string `json:"email,omitempty"`

	// Phone is the contact phone number with country code
	Phone string `json:"phone,omitempty"`
}

// BusinessDetails contains detailed business information
type BusinessDetails struct {
	// LegalEntityName is the legal name in the local language
	LegalEntityName string `json:"legal_entity_name,omitempty"`

	// LegalEntityNameEnglish is the legal name in English
	LegalEntityNameEnglish string `json:"legal_entity_name_english,omitempty"`

	// IncorporationDate is the date of incorporation (YYYY-MM-DD)
	IncorporationDate string `json:"incorporation_date,omitempty"`

	// RegistrationNumber is the business registration number
	RegistrationNumber string `json:"registration_number,omitempty"`

	// BusinessStructure is the type of business structure (e.g., "sole_proprietor")
	BusinessStructure string `json:"business_structure,omitempty"`

	// ProductDescription describes the products or services offered
	ProductDescription string `json:"product_description,omitempty"`

	// MerchantCategoryCode is the MCC code for the business
	MerchantCategoryCode string `json:"merchant_category_code,omitempty"`

	// MCC is an alternative field for merchant category code
	MCC string `json:"mcc,omitempty"`

	// EstimatedWorkerCount is the estimated number of workers
	EstimatedWorkerCount string `json:"estimated_worker_count,omitempty"`

	// MonthlyEstimatedRevenue contains revenue estimates
	MonthlyEstimatedRevenue *MonthlyEstimatedRevenue `json:"monthly_estimated_revenue,omitempty"`

	// AccountPurpose lists the intended purposes for the account
	AccountPurpose []string `json:"account_purpose,omitempty"`

	// Identifier contains business identification details
	Identifier *Identifier `json:"identifier,omitempty"`

	// WebsiteURL is the business website URL
	WebsiteURL string `json:"website_url,omitempty"`

	// Country is the ISO 3166-1 alpha-2 country code for the business
	Country string `json:"country,omitempty"`

	// IndustryCode is the industry classification code
	IndustryCode string `json:"industry_code,omitempty"`

	// BankingCountries is a list of countries for banking operations
	BankingCountries []string `json:"banking_countries,omitempty"`

	// BankingCurrencies is a list of currencies for banking operations
	BankingCurrencies []string `json:"banking_currencies,omitempty"`

	// BusinessAddress is a list of business operating addresses within business details
	BusinessAddress []Address `json:"business_address,omitempty"`

	// RegistrationAddress is the registration address within business details
	RegistrationAddress *Address `json:"registration_address,omitempty"`

	// IdentificationExpiryDate is the expiry date of the business identification
	IdentificationExpiryDate string `json:"identification_expiry_date,omitempty"`

	// IdentificationIssueDate is the issue date of the business identification
	IdentificationIssueDate string `json:"identification_issue_date,omitempty"`
}

// MonthlyEstimatedRevenue represents monthly revenue estimates
type MonthlyEstimatedRevenue struct {
	// Amount is the estimated monthly revenue amount or tier code
	Amount string `json:"amount"`

	// Currency is the ISO 4217 currency code
	Currency string `json:"currency"`
}

// Identifier represents a business identification number
type Identifier struct {
	// Type is the type of identifier (e.g., "VAT", "EIN", "SSN")
	Type string `json:"type"`

	// Number is the identifier number
	Number string `json:"number"`
}

// Address represents a physical address
type Address struct {
	// Line1 is the primary address line
	Line1 string `json:"line1"`

	// Line2 is the secondary address line (optional)
	Line2 string `json:"line2,omitempty"`

	// City is the city name
	City string `json:"city"`

	// State is the state or province
	State string `json:"state,omitempty"`

	// PostalCode is the postal or ZIP code
	PostalCode string `json:"postal_code,omitempty"`

	// Country is the ISO 3166-1 alpha-2 country code
	Country string `json:"country,omitempty"`
}

// Representative represents an account representative or director
type Representative struct {
	// RepresentativeID is the unique identifier for the representative
	RepresentativeID string `json:"representative_id,omitempty"`

	// Roles is the role of the representative (e.g., "DIRECTOR", "OWNER", "UBO")
	Roles string `json:"roles,omitempty"`

	// FirstName is the representative's first name
	FirstName string `json:"first_name"`

	// LastName is the representative's last name
	LastName string `json:"last_name"`

	// LocalName is the representative's name in local language
	LocalName string `json:"local_name,omitempty"`

	// Nationality is the ISO 3166-1 alpha-2 country code of nationality
	Nationality string `json:"nationality,omitempty"`

	// DateOfBirth is the date of birth (YYYY-MM-DD)
	DateOfBirth string `json:"date_of_birth,omitempty"`

	// SharePercentage is the ownership percentage
	SharePercentage string `json:"share_percentage,omitempty"`

	// AreaCode is the phone area/country code (e.g., "+86")
	AreaCode string `json:"area_code,omitempty"`

	// PhoneNumber is the phone number without area code
	PhoneNumber string `json:"phone_number,omitempty"`

	// Email is the representative's email address
	Email string `json:"email,omitempty"`

	// TaxNumber is the tax identification number
	TaxNumber string `json:"tax_number,omitempty"`

	// Identification contains ID document information
	Identification *Identification `json:"identification,omitempty"`

	// ResidentialAddress is the representative's residential address
	ResidentialAddress *Address `json:"residential_address,omitempty"`

	// IsApplicant indicates if this representative is the applicant
	IsApplicant bool `json:"is_applicant,omitempty"`

	// AsApplicant is an alternative field for applicant status (deprecated, use IsApplicant)
	AsApplicant bool `json:"as_applicant,omitempty"`

	// IDVStatus is the identity verification status for this representative
	IDVStatus string `json:"idv_status,omitempty"`

	// IDVID is the identity verification ID
	IDVID string `json:"idv_id,omitempty"`

	// CitizenshipStatus is the citizenship status code
	CitizenshipStatus int `json:"citizenship_status,omitempty"`
}

// Identification represents identity document information
type Identification struct {
	// Type is the document type (e.g., "PASSPORT", "NATIONAL_ID")
	Type string `json:"type"`

	// IDNumber is the document number
	IDNumber string `json:"id_number"`

	// Documents contains the document images/files
	Documents *IdentificationDocuments `json:"documents,omitempty"`

	// IdentificationExpiryDate is the expiry date of the identification document
	IdentificationExpiryDate string `json:"identification_expiry_date,omitempty"`

	// IdentificationIssueDate is the issue date of the identification document
	IdentificationIssueDate string `json:"identification_issue_date,omitempty"`
}

// IdentificationDocuments contains document image references
type IdentificationDocuments struct {
	// Front is the path/URL to the front of the document
	Front string `json:"front,omitempty"`

	// Back is the path/URL to the back of the document
	Back string `json:"back,omitempty"`

	// Type is the document type code
	Type string `json:"type,omitempty"`
}

// Account status constants
const (
	AccountStatusProcessing = "PROCESSING"
	AccountStatusActive     = "ACTIVE"
	AccountStatusSuspended  = "SUSPENDED"
	AccountStatusClosed     = "CLOSED"
)

// Verification status constants
const (
	VerificationStatusPending  = "PENDING"
	VerificationStatusApproved = "APPROVED"
	VerificationStatusRejected = "REJECTED"
)

// Entity type constants
const (
	EntityTypeCompany    = "COMPANY"
	EntityTypeIndividual = "INDIVIDUAL"
)

// Business structure constants
const (
	BusinessStructureSoleProprietor   = "SOLE_PROPRIETOR"
	BusinessStructurePartnership      = "PARTNERSHIP"
	BusinessStructureCorporation      = "CORPORATION"
	BusinessStructureLLC              = "LLC"
	BusinessStructureNonProfit        = "NON_PROFIT"
	BusinessStructureGovernmentEntity = "GOVERNMENT_ENTITY"
	BusinessStructurePubliclyTraded   = "PUBLICLY_TRADED"
	BusinessStructurePrivatelyHeld    = "PRIVATELY_HELD"
)

// Representative role constants
const (
	RoleDirector       = "DIRECTOR"
	RoleOwner          = "OWNER"
	RoleShareholder    = "SHAREHOLDER"
	RoleAuthorizedUser = "AUTHORIZED_USER"
	RoleUBO            = "UBO" // Ultimate Beneficial Owner
)

// Identification type constants
const (
	IdentificationTypePassport      = "PASSPORT"
	IdentificationTypeNationalID    = "NATIONAL_ID"
	IdentificationTypeDriverLicense = "DRIVER_LICENSE"
)

// Account purpose constants
const (
	AccountPurposeCollection  = "COLLECTION"
	AccountPurposePayout      = "PAYOUT"
	AccountPurposeBillPayment = "bill_payment"
)
