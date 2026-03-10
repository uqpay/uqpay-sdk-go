package supporting

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// FilesClient handles file operations
type FilesClient struct {
	client *common.APIClient
}

// UploadFileParams represents file upload parameters
type UploadFileParams struct {
	File     io.Reader // Required. Binary file content. Max 20MB. Supported: jpeg, png, jpg, doc, docx, pdf
	FileName string    // Required. Original file name
	Notes    string    // Optional. File annotations, max 50 characters
}

// UploadFileResponse represents file upload response
type UploadFileResponse struct {
	CreateTime string `json:"create_time"` // Timestamp when the file was created, e.g. "2024-08-22T17:12:58+08:00"
	FileID     string `json:"file_id"`     // Unique file identifier, used in attachment references
	FileName   string `json:"file_name"`   // Original file name
	FileType   string `json:"file_type"`   // File extension, e.g. "png", "pdf"
	Size       int    `json:"size"`        // File size in bytes
	Notes      string `json:"notes"`       // File annotations
}

// DownloadLinksRequest represents download links request
type DownloadLinksRequest struct {
	FileIDs []string `json:"file_ids"` // Required. List of file IDs (UUID format) to generate download links for
}

// FileDownloadInfo represents file download information
type FileDownloadInfo struct {
	FileID   string `json:"file_id"`   // Unique file identifier (UUID)
	FileType string `json:"file_type"` // File extension, e.g. "png"
	FileName string `json:"file_name"` // Complete file name with extension
	Size     int    `json:"size"`      // File size in bytes, max 20MB
	URL      string `json:"url"`       // Direct download link for the file
}

// DownloadLinksResponse represents download links response
type DownloadLinksResponse struct {
	Files       []FileDownloadInfo `json:"files"`        // List of file download details
	AbsentFiles []string           `json:"absent_files"` // File IDs that were not found in the system
}

// Upload uploads a file to UQPAY
// POST /v1/files/upload
// Maximum file size: 20MB
// Supported types: jpeg, png, jpg, doc, docx, pdf
// Optional RequestOptions can be provided to set custom headers like x-on-behalf-of or x-idempotency-key
func (c *FilesClient) Upload(ctx context.Context, params *UploadFileParams, opts ...*common.RequestOptions) (*UploadFileResponse, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Create form file
	part, err := writer.CreateFormFile("file", params.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file content
	if _, err := io.Copy(part, params.File); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Add notes if provided
	if params.Notes != "" {
		if err := writer.WriteField("notes", params.Notes); err != nil {
			return nil, fmt.Errorf("failed to write notes field: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	var resp UploadFileResponse
	if err := c.client.PostMultipart(ctx, "/v1/files/upload", &buf, writer.FormDataContentType(), &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	return &resp, nil
}

// GetDownloadLinks retrieves download links for specified file IDs
// POST /v1/files/download_links
// Optional RequestOptions can be provided to set custom headers like x-on-behalf-of or x-idempotency-key
func (c *FilesClient) GetDownloadLinks(ctx context.Context, req *DownloadLinksRequest, opts ...*common.RequestOptions) (*DownloadLinksResponse, error) {
	var resp DownloadLinksResponse
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.PostWithOptions(ctx, "/v1/files/download_links", req, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to get download links: %w", err)
	}
	return &resp, nil
}
