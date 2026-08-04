package connect

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/uqpay/uqpay-sdk-go/common"
)

type RFIsClient struct{ client *common.APIClient }

type RFIStatus string

const (
	RFIStatusSubmittedPending RFIStatus = "SUBMITTED_PENDING"
	RFIStatusRejected         RFIStatus = "REJECTED"
	RFIStatusApproved         RFIStatus = "APPROVED"
	RFIStatusActionRequired   RFIStatus = "ACTION_REQUIRED"
)

type RFIAnswerItem struct {
	Key         string   `json:"key"`
	Type        string   `json:"type"`
	Attachments []string `json:"attachments"`
}

type RFIQuestion struct {
	Key     string `json:"key"`
	Comment string `json:"comment,omitempty"`
	Type    string `json:"type"`
}

type RFIRequestItem struct {
	Question RFIQuestion    `json:"question"`
	Answer   *RFIAnswerItem `json:"answer,omitempty"`
}

type RFI struct {
	AccountID  string           `json:"account_id"`
	RFIID      string           `json:"rfi_id"`
	Status     RFIStatus        `json:"status"`
	CreateTime string           `json:"create_time"`
	UpdateTime string           `json:"update_time"`
	Request    []RFIRequestItem `json:"request"`
}

type ListRFIsRequest struct {
	PageSize   int
	PageNumber int
	Status     RFIStatus
}

type ListRFIsResponse struct {
	TotalPages int   `json:"total_pages"`
	TotalItems int   `json:"total_items"`
	Data       []RFI `json:"data"`
}

type AnswerRFIRequest struct {
	RFIID  string          `json:"rfi_id"`
	Answer []RFIAnswerItem `json:"answer"`
}

func (c *RFIsClient) List(ctx context.Context, req *ListRFIsRequest, opts ...*common.RequestOptions) (*ListRFIsResponse, error) {
	params := url.Values{}
	params.Set("page_size", strconv.Itoa(req.PageSize))
	params.Set("page_number", strconv.Itoa(req.PageNumber))
	if req.Status != "" {
		params.Set("status", string(req.Status))
	}
	var resp ListRFIsResponse
	if err := c.client.GetWithOptions(ctx, "/v1/rfis?"+params.Encode(), &resp, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to list RFIs: %w", err)
	}
	return &resp, nil
}

func (c *RFIsClient) Get(ctx context.Context, rfiID string, opts ...*common.RequestOptions) (*RFI, error) {
	var resp RFI
	if err := c.client.GetWithOptions(ctx, "/v1/rfis/"+url.PathEscape(rfiID), &resp, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to get RFI: %w", err)
	}
	return &resp, nil
}

func (c *RFIsClient) Answer(ctx context.Context, req *AnswerRFIRequest, opts ...*common.RequestOptions) (*RFI, error) {
	var resp RFI
	if err := c.client.PostWithOptions(ctx, "/v1/rfis/answer", req, &resp, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to answer RFI: %w", err)
	}
	return &resp, nil
}
