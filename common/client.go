package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/uqpay/uqpay-sdk-go/configuration"
)

// TokenProvider provides auth tokens
type TokenProvider interface {
	GetToken() (string, error)
}

// APIClient handles HTTP requests to UQPAY API
type APIClient struct {
	Config        *configuration.Configuration
	TokenProvider TokenProvider
	HTTPClient    *http.Client
}

// NewAPIClient creates a new API client
func NewAPIClient(config *configuration.Configuration, tokenProvider TokenProvider) *APIClient {
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &APIClient{
		Config:        config,
		TokenProvider: tokenProvider,
		HTTPClient:    httpClient,
	}
}

// Do executes an HTTP request
func (c *APIClient) Do(ctx context.Context, method, path string, body, response interface{}) error {
	url := c.Config.Environment.BaseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Get auth token
	token, err := c.TokenProvider.GetToken()
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-auth-token", "Bearer "+token)
	req.Header.Set("x-idempotency-key", uuid.New().String())

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
		apiErr.StatusCode = resp.StatusCode
		return &apiErr
	}

	// Decode response
	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Get sends a GET request
func (c *APIClient) Get(ctx context.Context, path string, response interface{}) error {
	return c.Do(ctx, "GET", path, nil, response)
}

// Post sends a POST request
func (c *APIClient) Post(ctx context.Context, path string, body, response interface{}) error {
	return c.Do(ctx, "POST", path, body, response)
}

// Put sends a PUT request
func (c *APIClient) Put(ctx context.Context, path string, body, response interface{}) error {
	return c.Do(ctx, "PUT", path, body, response)
}

// Delete sends a DELETE request
func (c *APIClient) Delete(ctx context.Context, path string, response interface{}) error {
	return c.Do(ctx, "DELETE", path, nil, response)
}

// GetRaw sends a GET request and returns raw bytes (for file downloads)
func (c *APIClient) GetRaw(ctx context.Context, path string) ([]byte, error) {
	url := c.Config.Environment.BaseURL + path

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Get auth token
	token, err := c.TokenProvider.GetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/octet-stream")
	req.Header.Set("x-auth-token", "Bearer "+token)
	req.Header.Set("x-idempotency-key", uuid.New().String())

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
		apiErr.StatusCode = resp.StatusCode
		return nil, &apiErr
	}

	// Read response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}
