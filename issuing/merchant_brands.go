package issuing

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/uqpay/uqpay-sdk-go/common"
)

type MerchantBrandsClient struct{ client *common.APIClient }

type MerchantBrand struct {
	MerchantCode string `json:"merchant_code"`
	DisplayName  string `json:"display_name"`
}

type ListMerchantBrandsRequest struct {
	DisplayName  string
	MerchantCode string
	PageNumber   int
	PageSize     int
}

type ListMerchantBrandsResponse struct {
	TotalItems int             `json:"total_items"`
	TotalPages int             `json:"total_pages"`
	Data       []MerchantBrand `json:"data"`
}

func (c *MerchantBrandsClient) List(ctx context.Context, req *ListMerchantBrandsRequest, opts ...*common.RequestOptions) (*ListMerchantBrandsResponse, error) {
	params := url.Values{}
	params.Set("page_number", strconv.Itoa(req.PageNumber))
	params.Set("page_size", strconv.Itoa(req.PageSize))
	if req.DisplayName != "" {
		params.Set("display_name", req.DisplayName)
	}
	if req.MerchantCode != "" {
		params.Set("merchant_code", req.MerchantCode)
	}
	var resp ListMerchantBrandsResponse
	if err := c.client.GetWithOptions(ctx, "/v1/issuing/merchant_brands?"+params.Encode(), &resp, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to list merchant brands: %w", err)
	}
	return &resp, nil
}
