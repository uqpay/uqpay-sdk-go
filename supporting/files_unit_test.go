package supporting

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/v2/common"
	"github.com/uqpay/uqpay-sdk-go/v2/configuration"
)

type staticTokenProvider struct {
	token string
}

func (p *staticTokenProvider) GetToken() (string, error) {
	return p.token, nil
}

type uploadRequest struct {
	header   http.Header
	fileName string
	fileData string
	notes    string
}

func TestFilesUploadForwardsRequestOptions(t *testing.T) {
	requests := make(chan uploadRequest, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := uploadRequest{header: r.Header.Clone()}
		if err := r.ParseMultipartForm(1 << 20); err == nil {
			got.notes = r.FormValue("notes")
			file, header, fileErr := r.FormFile("file")
			if fileErr == nil {
				defer file.Close()
				data, _ := io.ReadAll(file)
				got.fileName = header.Filename
				got.fileData = string(data)
			}
		}
		requests <- got
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"file_id":"file_123","file_name":"evidence.txt"}`))
	}))
	defer server.Close()

	config := &configuration.Configuration{
		Environment: &configuration.Environment{BaseURL: server.URL},
		HTTPClient:  server.Client(),
	}
	apiClient := common.NewAPIClient(config, &staticTokenProvider{token: "default-token"})
	client := NewClient(apiClient)

	resp, err := client.Files.Upload(context.Background(), &UploadFileParams{
		File:     bytes.NewBufferString("supporting evidence"),
		FileName: "evidence.txt",
		Notes:    "sub-account evidence",
	}, &common.RequestOptions{
		OnBehalfOf:     "account_sub_789",
		ClientID:       "client_789",
		IdempotencyKey: "idempotency_789",
		AuthToken:      "request-token",
	})
	if err != nil {
		t.Fatalf("Files.Upload returned an error: %v", err)
	}
	if resp.FileID != "file_123" {
		t.Errorf("response FileID = %q, want %q", resp.FileID, "file_123")
	}

	got := <-requests
	if got.fileName != "evidence.txt" {
		t.Errorf("uploaded filename = %q, want %q", got.fileName, "evidence.txt")
	}
	if got.fileData != "supporting evidence" {
		t.Errorf("uploaded content = %q, want %q", got.fileData, "supporting evidence")
	}
	if got.notes != "sub-account evidence" {
		t.Errorf("uploaded notes = %q, want %q", got.notes, "sub-account evidence")
	}
	wantHeaders := map[string]string{
		"x-on-behalf-of":    "account_sub_789",
		"x-client-id":       "client_789",
		"x-idempotency-key": "idempotency_789",
		"x-auth-token":      "Bearer request-token",
	}
	for name, want := range wantHeaders {
		if gotValue := got.header.Get(name); gotValue != want {
			t.Errorf("%s = %q, want %q", name, gotValue, want)
		}
	}
	if gotValue := got.header.Get("Content-Type"); gotValue == "" {
		t.Error("Content-Type is empty, want multipart/form-data")
	}
}
