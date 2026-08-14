package supporting

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/uqpay/uqpay-sdk-go/v2/common"
)

// FilesClient handles file operations
type FilesClient struct {
	client *common.APIClient
}

// UploadFileParams represents file upload parameters
type UploadFileParams struct {
	File     io.Reader
	FileName string
	Notes    string
}

// UploadFileResponse represents file upload response
type UploadFileResponse struct {
	CreateTime string `json:"create_time"`
	FileID     string `json:"file_id"`
	FileName   string `json:"file_name"`
	FileType   string `json:"file_type"`
	Size       int    `json:"size"`
	Notes      string `json:"notes"`
}

// DownloadLinksRequest represents download links request
type DownloadLinksRequest struct {
	FileIDs []string `json:"file_ids"` // required
}

// FileDownloadInfo represents file download information
type FileDownloadInfo struct {
	FileID   string `json:"file_id"`
	FileType string `json:"file_type"`
	FileName string `json:"file_name"`
	Size     int    `json:"size"`
	URL      string `json:"url"`
}

// DownloadLinksResponse represents download links response
type DownloadLinksResponse struct {
	Files       []FileDownloadInfo `json:"files"`
	AbsentFiles []string           `json:"absent_files"`
}

// Upload uploads a file to UQPAY
// POST /v1/files/upload
// Maximum file size: 20MB
// Supported types: jpeg, png, jpg, doc, docx, pdf
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

	var resp UploadFileResponse
	if err := c.client.PostMultipartWithOptions(
		ctx,
		"/v1/files/upload",
		&buf,
		writer.FormDataContentType(),
		&resp,
		firstRequestOptions(opts),
	); err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	return &resp, nil
}

// GetDownloadLinks retrieves download links for specified file IDs
// POST /v1/files/download_links
func (c *FilesClient) GetDownloadLinks(ctx context.Context, req *DownloadLinksRequest, opts ...*common.RequestOptions) (*DownloadLinksResponse, error) {
	var resp DownloadLinksResponse
	if err := c.client.PostWithOptions(ctx, "/v1/files/download_links", req, &resp, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to get download links: %w", err)
	}
	return &resp, nil
}
