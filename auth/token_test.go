package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func tokenServer(t *testing.T, expiry time.Time) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	calls := new(atomic.Int32)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		call := calls.Add(1)
		if r.Method != http.MethodPost || r.URL.Path != "/v1/connect/token" {
			t.Errorf("unexpected token request: %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("x-client-id") != "test-client" || r.Header.Get("x-api-key") != "test-key" {
			t.Error("token request omitted credentials")
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Error("token request omitted JSON accept header")
		}
		token, expiresAt := "first-token", expiry
		if call > 1 {
			token, expiresAt = "refreshed-token", time.Now().Add(time.Hour)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(TokenResponse{AuthToken: token, ExpiredAt: expiresAt.Unix()}); err != nil {
			t.Errorf("write token response: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	return server, calls
}

func TestGetTokenCachesValidToken(t *testing.T) {
	server, calls := tokenServer(t, time.Now().Add(time.Hour))
	provider := NewTokenProvider(server.URL, "test-client", "test-key", server.Client())
	for i := 0; i < 3; i++ {
		token, err := provider.GetToken()
		if err != nil || token != "first-token" {
			t.Fatalf("GetToken() = %q, %v; want first-token, nil", token, err)
		}
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("token requests = %d; want 1", got)
	}
}

func TestGetTokenRefreshesExpiredOrExpiringToken(t *testing.T) {
	for _, tc := range []struct {
		name      string
		remaining time.Duration
	}{
		{"expired", -time.Hour},
		{"inside_refresh_buffer", time.Minute},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server, calls := tokenServer(t, time.Now().Add(tc.remaining))
			provider := NewTokenProvider(server.URL, "test-client", "test-key", server.Client())
			for _, want := range []string{"first-token", "refreshed-token", "refreshed-token"} {
				token, err := provider.GetToken()
				if err != nil || token != want {
					t.Fatalf("GetToken() = %q, %v; want %q, nil", token, err, want)
				}
			}
			if got := calls.Load(); got != 2 {
				t.Fatalf("token requests = %d; want 2", got)
			}
		})
	}
}

func TestGetTokenConcurrentCallsFetchOnce(t *testing.T) {
	server, calls := tokenServer(t, time.Now().Add(time.Hour))
	provider := NewTokenProvider(server.URL, "test-client", "test-key", server.Client())
	const workers = 20
	var ready, done sync.WaitGroup
	ready.Add(workers)
	done.Add(workers)
	start := make(chan struct{})
	for i := 0; i < workers; i++ {
		go func() {
			defer done.Done()
			ready.Done()
			<-start
			token, err := provider.GetToken()
			if err != nil || token != "first-token" {
				t.Errorf("GetToken() = %q, %v; want first-token, nil", token, err)
			}
		}()
	}
	ready.Wait()
	close(start)
	done.Wait()
	if got := calls.Load(); got != 1 {
		t.Fatalf("concurrent token requests = %d; want 1", got)
	}
}
