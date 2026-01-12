package test

import (
	"os"
	"testing"
	"time"

	"github.com/uqpay/uqpay-sdk-go/auth"
	"github.com/uqpay/uqpay-sdk-go/configuration"
)

// minInt returns the smaller of two integers (Go 1.19 compatible)
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ============================================================================
// Auth Token Provider Tests
// ============================================================================

func TestTokenProvider(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip integration tests in CI environment
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test in CI environment")
	}

	clientID := os.Getenv("UQPAY_CLIENT_ID")
	apiKey := os.Getenv("UQPAY_API_KEY")

	if clientID == "" || apiKey == "" {
		t.Skip("Skipping test: UQPAY_CLIENT_ID and UQPAY_API_KEY environment variables not set")
	}

	env := configuration.Sandbox()

	t.Run("NewTokenProvider", func(t *testing.T) {
		provider := auth.NewTokenProvider(env.BaseURL, clientID, apiKey, nil)

		if provider == nil {
			t.Fatal("NewTokenProvider returned nil")
		}

		t.Log("TokenProvider created successfully")
	})

	t.Run("GetToken_ValidCredentials", func(t *testing.T) {
		provider := auth.NewTokenProvider(env.BaseURL, clientID, apiKey, nil)

		token, err := provider.GetToken()
		if err != nil {
			t.Fatalf("GetToken failed: %v", err)
		}

		// Assertions
		if token == "" {
			t.Error("Token should not be empty")
		}

		// Token should be a reasonable length (typical JWT/token length)
		if len(token) < 20 {
			t.Errorf("Token length seems too short: %d characters", len(token))
		}

		t.Logf("Token retrieved successfully")
		t.Logf("   Token length: %d characters", len(token))
		t.Logf("   Token preview: %s...", token[:minInt(20, len(token))])
	})

	t.Run("GetToken_Caching", func(t *testing.T) {
		provider := auth.NewTokenProvider(env.BaseURL, clientID, apiKey, nil)

		// Get token first time
		token1, err := provider.GetToken()
		if err != nil {
			t.Fatalf("First GetToken failed: %v", err)
		}

		// Get token second time (should be cached)
		token2, err := provider.GetToken()
		if err != nil {
			t.Fatalf("Second GetToken failed: %v", err)
		}

		// Both tokens should be the same (cached)
		if token1 != token2 {
			t.Error("Cached token should be returned on second call")
		}

		t.Log("Token caching works correctly")
		t.Logf("   First token:  %s...", token1[:minInt(20, len(token1))])
		t.Logf("   Second token: %s...", token2[:minInt(20, len(token2))])
	})

	t.Run("GetToken_MultipleCallsPerformance", func(t *testing.T) {
		provider := auth.NewTokenProvider(env.BaseURL, clientID, apiKey, nil)

		// First call (will make HTTP request)
		start := time.Now()
		_, err := provider.GetToken()
		if err != nil {
			t.Fatalf("First GetToken failed: %v", err)
		}
		firstCallDuration := time.Since(start)

		// Subsequent calls (should be cached and fast)
		start = time.Now()
		for i := 0; i < 100; i++ {
			_, err := provider.GetToken()
			if err != nil {
				t.Fatalf("GetToken call %d failed: %v", i, err)
			}
		}
		cachedCallsDuration := time.Since(start)

		t.Logf("Performance test completed")
		t.Logf("   First call (HTTP): %v", firstCallDuration)
		t.Logf("   100 cached calls:  %v", cachedCallsDuration)
		t.Logf("   Avg cached call:   %v", cachedCallsDuration/100)

		// Cached calls should be significantly faster than first call
		if cachedCallsDuration > firstCallDuration {
			t.Log("   Warning: Cached calls slower than expected, but test passed")
		}
	})
}

func TestTokenProvider_InvalidCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip integration tests in CI environment
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test in CI environment")
	}

	env := configuration.Sandbox()

	t.Run("GetToken_InvalidClientID", func(t *testing.T) {
		provider := auth.NewTokenProvider(env.BaseURL, "invalid-client-id", "invalid-api-key", nil)

		_, err := provider.GetToken()
		if err == nil {
			t.Error("Expected error for invalid credentials, got nil")
		}

		t.Logf("Invalid credentials correctly rejected: %v", err)
	})

	t.Run("GetToken_EmptyCredentials", func(t *testing.T) {
		provider := auth.NewTokenProvider(env.BaseURL, "", "", nil)

		_, err := provider.GetToken()
		if err == nil {
			t.Error("Expected error for empty credentials, got nil")
		}

		t.Logf("Empty credentials correctly rejected: %v", err)
	})
}

func TestTokenProvider_InvalidBaseURL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip integration tests in CI environment
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test in CI environment")
	}

	t.Run("GetToken_InvalidURL", func(t *testing.T) {
		provider := auth.NewTokenProvider("https://invalid-url.example.com/api", "client", "key", nil)

		_, err := provider.GetToken()
		if err == nil {
			t.Error("Expected error for invalid URL, got nil")
		}

		t.Logf("Invalid URL correctly rejected: %v", err)
	})
}
