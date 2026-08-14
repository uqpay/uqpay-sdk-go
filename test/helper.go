package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/uqpay/uqpay-sdk-go/v2"
	"github.com/uqpay/uqpay-sdk-go/v2/configuration"
)

func init() {
	// Try to load .env file from project root
	// Look for .env in current directory and parent directories
	dirs := []string{
		".env",
		"../.env",
		filepath.Join("..", ".env"),
	}
	for _, path := range dirs {
		if err := godotenv.Load(path); err == nil {
			break
		}
	}
}

// GetTestClient creates a test client with credentials from environment variables
func GetTestClient(t *testing.T) *uqpay.Client {
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
