package supporting

import "github.com/uqpay/uqpay-sdk-go/common"

// Client represents the Supporting Services API client
type Client struct {
	Files *FilesClient
}

// NewClient creates a new Supporting Services API client
func NewClient(apiClient *common.APIClient) *Client {
	return &Client{
		Files: &FilesClient{client: apiClient},
	}
}

func firstRequestOptions(opts []*common.RequestOptions) *common.RequestOptions {
	if len(opts) == 0 {
		return nil
	}
	return opts[0]
}
