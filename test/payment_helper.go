package test

import (
	"os"
	"testing"

	"github.com/uqpay/uqpay-sdk-go"
	"github.com/uqpay/uqpay-sdk-go/configuration"
)

// GetPaymentTestClient creates a test client for Payment API tests
func GetPaymentTestClient(t *testing.T) *uqpay.Client {
	t.Helper()

	// Skip integration tests in CI environment
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test in CI environment")
	}

	clientID := os.Getenv("UQPAY_CLIENT_ID")
	apiKey := os.Getenv("UQPAY_API_KEY")

	if clientID == "" || apiKey == "" {
		t.Skip("Skipping test: UQPAY_CLIENT_ID and UQPAY_API_KEY environment variables not set")
	}

	client, err := uqpay.NewClient(clientID, apiKey, configuration.Sandbox())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return client
}
