package uqpay

import (
	"net/http"

	"github.com/uqpay/uqpay-sdk-go/auth"
	"github.com/uqpay/uqpay-sdk-go/banking"
	"github.com/uqpay/uqpay-sdk-go/common"
	"github.com/uqpay/uqpay-sdk-go/configuration"
	"github.com/uqpay/uqpay-sdk-go/connect"
	"github.com/uqpay/uqpay-sdk-go/issuing"
	"github.com/uqpay/uqpay-sdk-go/payment"
	"github.com/uqpay/uqpay-sdk-go/supporting"
)

// Client is the main UQPAY SDK client
type Client struct {
	Issuing    *issuing.Client
	Banking    *banking.Client
	Connect    *connect.Client
	Supporting *supporting.Client
	Payment    *payment.Client
}

// NewClient creates a new UQPAY client
func NewClient(clientID, apiKey string, env *configuration.Environment) (*Client, error) {
	config := &configuration.Configuration{
		ClientID:    clientID,
		APIKey:      apiKey,
		Environment: env,
		HTTPClient:  &http.Client{},
	}

	// Create token provider
	tokenProvider := auth.NewTokenProvider(
		env.BaseURL,
		clientID,
		apiKey,
		config.HTTPClient,
	)

	// Create API client for main APIs
	apiClient := common.NewAPIClient(config, tokenProvider)

	// Create separate configuration for Files API (different base URL)
	filesConfig := &configuration.Configuration{
		ClientID:    clientID,
		APIKey:      apiKey,
		Environment: &configuration.Environment{BaseURL: env.FilesBaseURL},
		HTTPClient:  &http.Client{},
	}
	filesTokenProvider := auth.NewTokenProvider(
		env.FilesBaseURL,
		clientID,
		apiKey,
		filesConfig.HTTPClient,
	)
	filesAPIClient := common.NewAPIClient(filesConfig, filesTokenProvider)

	// Initialize service clients
	return &Client{
		Issuing:    issuing.NewClient(apiClient),
		Banking:    banking.NewClient(apiClient),
		Connect:    connect.NewClient(apiClient),
		Supporting: supporting.NewClient(filesAPIClient), // Use separate client for Files API
		Payment:    payment.NewClient(apiClient),
	}, nil
}
