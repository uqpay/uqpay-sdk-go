package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// TokenResponse represents the auth token response
type TokenResponse struct {
	AuthToken string `json:"auth_token"`
	ExpiredAt int64  `json:"expired_at"`
}

// TokenProvider automatically manages and refreshes auth tokens
type TokenProvider struct {
	mu            sync.RWMutex
	baseURL       string
	clientID      string
	apiKey        string
	httpClient    *http.Client
	currentToken  string
	expiresAt     time.Time
	refreshBuffer time.Duration
}

// NewTokenProvider creates a new token provider
func NewTokenProvider(baseURL, clientID, apiKey string, httpClient *http.Client) *TokenProvider {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	return &TokenProvider{
		baseURL:       baseURL,
		clientID:      clientID,
		apiKey:        apiKey,
		httpClient:    httpClient,
		refreshBuffer: 5 * time.Minute,
	}
}

// GetToken returns a valid token, refreshing if necessary
func (p *TokenProvider) GetToken() (string, error) {
	p.mu.RLock()
	if time.Now().Add(p.refreshBuffer).Before(p.expiresAt) {
		token := p.currentToken
		p.mu.RUnlock()
		return token, nil
	}
	p.mu.RUnlock()

	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check after acquiring write lock
	if time.Now().Add(p.refreshBuffer).Before(p.expiresAt) {
		return p.currentToken, nil
	}

	// Refresh token
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := p.baseURL + "/v1/connect/token"
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-client-id", p.clientID)
	req.Header.Set("x-api-key", p.apiKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get token: status %d, body: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	p.currentToken = tokenResp.AuthToken
	p.expiresAt = time.Unix(tokenResp.ExpiredAt, 0)

	return p.currentToken, nil
}
