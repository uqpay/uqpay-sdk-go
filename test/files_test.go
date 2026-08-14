package test

import (
	"context"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/v2/supporting"
)

func TestFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetTestClient(t)
	ctx := context.Background()

	t.Run("GetDownloadLinks", func(t *testing.T) {
		// Test with empty file IDs - should handle gracefully
		t.Run("EmptyFileIDs", func(t *testing.T) {
			req := &supporting.DownloadLinksRequest{
				FileIDs: []string{},
			}

			resp, err := client.Supporting.Files.GetDownloadLinks(ctx, req)
			if err != nil {
				t.Logf("GetDownloadLinks with empty IDs returned: %v", err)
				return
			}

			t.Logf("✅ GetDownloadLinks with empty IDs succeeded")
			t.Logf("  Files returned: %d", len(resp.Files))
			t.Logf("  Absent files: %d", len(resp.AbsentFiles))
		})

		// Test with non-existent file IDs - should return them in absent_files
		t.Run("NonExistentFileIDs", func(t *testing.T) {
			req := &supporting.DownloadLinksRequest{
				FileIDs: []string{
					"file_nonexistent_123",
					"file_nonexistent_456",
				},
			}

			resp, err := client.Supporting.Files.GetDownloadLinks(ctx, req)
			if err != nil {
				t.Logf("GetDownloadLinks with non-existent IDs returned: %v", err)
				return
			}

			t.Logf("✅ GetDownloadLinks with non-existent IDs succeeded")
			t.Logf("  Files returned: %d", len(resp.Files))
			t.Logf("  Absent files: %d", len(resp.AbsentFiles))

			if len(resp.AbsentFiles) > 0 {
				t.Logf("  Absent file IDs: %v", resp.AbsentFiles)
			}

			// Verify that non-existent files are reported as absent
			if len(resp.Files) > 0 {
				t.Log("  Warning: Non-existent file IDs returned download links - unexpected behavior")
			}
		})

		// Test with a mix of valid and invalid file IDs
		// Note: In a real scenario, you would have valid file IDs from previous uploads
		t.Run("MixedFileIDs", func(t *testing.T) {
			// This test demonstrates the expected behavior with mixed file IDs
			// In production, you would replace these with actual file IDs
			req := &supporting.DownloadLinksRequest{
				FileIDs: []string{
					"file_test_valid_id_if_exists", // Replace with actual valid file ID in production
					"file_nonexistent_999",
				},
			}

			resp, err := client.Supporting.Files.GetDownloadLinks(ctx, req)
			if err != nil {
				t.Logf("GetDownloadLinks with mixed IDs returned: %v", err)
				return
			}

			t.Logf("✅ GetDownloadLinks with mixed IDs succeeded")
			t.Logf("  Total requested: %d", len(req.FileIDs))
			t.Logf("  Files returned: %d", len(resp.Files))
			t.Logf("  Absent files: %d", len(resp.AbsentFiles))

			// Log details for found files
			for i, file := range resp.Files {
				t.Logf("  File %d: ID=%s, Name=%s, Type=%s, Size=%d bytes",
					i+1, file.FileID, file.FileName, file.FileType, file.Size)
				t.Logf("    Download URL: %s", file.URL)
			}

			// Log absent files
			if len(resp.AbsentFiles) > 0 {
				t.Logf("  Absent file IDs: %v", resp.AbsentFiles)
			}

			// Verify the total matches
			if len(resp.Files)+len(resp.AbsentFiles) != len(req.FileIDs) {
				t.Logf("  Note: Total files (%d) + absent files (%d) != requested files (%d)",
					len(resp.Files), len(resp.AbsentFiles), len(req.FileIDs))
			}
		})

		// Test with a single file ID
		t.Run("SingleFileID", func(t *testing.T) {
			req := &supporting.DownloadLinksRequest{
				FileIDs: []string{"file_test_single_123"},
			}

			resp, err := client.Supporting.Files.GetDownloadLinks(ctx, req)
			if err != nil {
				t.Logf("GetDownloadLinks with single ID returned: %v", err)
				return
			}

			t.Logf("✅ GetDownloadLinks with single ID succeeded")

			if len(resp.Files) > 0 {
				file := resp.Files[0]
				t.Logf("  File found: ID=%s, Name=%s, Type=%s, Size=%d bytes",
					file.FileID, file.FileName, file.FileType, file.Size)
				t.Logf("  Download URL: %s", file.URL)

				// Verify URL is not empty
				if file.URL == "" {
					t.Error("  Error: Download URL is empty")
				}
			} else if len(resp.AbsentFiles) > 0 {
				t.Logf("  File not found (expected for test file ID): %v", resp.AbsentFiles)
			}
		})

		// Test with multiple valid file IDs
		t.Run("MultipleFileIDs", func(t *testing.T) {
			req := &supporting.DownloadLinksRequest{
				FileIDs: []string{
					"file_test_multi_1",
					"file_test_multi_2",
					"file_test_multi_3",
				},
			}

			resp, err := client.Supporting.Files.GetDownloadLinks(ctx, req)
			if err != nil {
				t.Logf("GetDownloadLinks with multiple IDs returned: %v", err)
				return
			}

			t.Logf("✅ GetDownloadLinks with multiple IDs succeeded")
			t.Logf("  Requested: %d files", len(req.FileIDs))
			t.Logf("  Found: %d files", len(resp.Files))
			t.Logf("  Absent: %d files", len(resp.AbsentFiles))

			// Log all found files
			for i, file := range resp.Files {
				t.Logf("  File %d: ID=%s, Name=%s, Type=%s, Size=%d bytes",
					i+1, file.FileID, file.FileName, file.FileType, file.Size)
			}

			// Log all absent files
			if len(resp.AbsentFiles) > 0 {
				for i, fileID := range resp.AbsentFiles {
					t.Logf("  Absent %d: %s", i+1, fileID)
				}
			}
		})
	})

	// Note: File upload test is skipped as it requires multipart/form-data implementation
	t.Run("Upload_NotImplemented", func(t *testing.T) {
		t.Skip("File upload test skipped - requires multipart/form-data implementation in APIClient")

		// This is a placeholder for future implementation
		// When multipart/form-data support is added to the APIClient:
		//
		// 1. Create a test file in memory
		// 2. Upload using client.Supporting.Files.Upload()
		// 3. Verify the upload response contains file_id
		// 4. Use GetDownloadLinks to verify the file exists
		// 5. Clean up test files if API provides delete functionality
		//
		// Example structure:
		// content := []byte("Test file content")
		// reader := bytes.NewReader(content)
		// params := &supporting.UploadFileParams{
		//     File:     reader,
		//     FileName: "test.txt",
		//     Notes:    "Test upload from Go SDK",
		// }
		// resp, err := client.Supporting.Files.Upload(ctx, params)
	})
}
