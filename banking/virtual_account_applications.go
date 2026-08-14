package banking

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/uqpay/uqpay-sdk-go/v2/common"
)

const (
	VirtualAccountApplicationSubmitted          = "SUBMITTED"
	VirtualAccountApplicationPartiallyCompleted = "PARTIALLY_COMPLETED"
	VirtualAccountApplicationCompleted          = "COMPLETED"
	VirtualAccountApplicationFailed             = "FAILED"
	VirtualAccountApplicationClosed             = "CLOSED"
)

type VirtualAccountApplicationsClient struct{ client *common.APIClient }

type VirtualAccountApplicationClearingSystem struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type VirtualAccountApplicationBankDetail struct {
	AccountBankID  string                                  `json:"account_bank_id"`
	AccountHolder  string                                  `json:"account_holder"`
	AccountNumber  string                                  `json:"account_number"`
	CountryCode    string                                  `json:"country_code"`
	Currency       string                                  `json:"currency"`
	BankName       string                                  `json:"bank_name"`
	BankAddress    string                                  `json:"bank_address"`
	ClearingSystem VirtualAccountApplicationClearingSystem `json:"clearing_system"`
	Status         string                                  `json:"status"`
	CloseReason    string                                  `json:"close_reason"`
}

type VirtualAccountApplicationResultError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type VirtualAccountApplicationResult struct {
	PaymentMethod   string                                `json:"payment_method"`
	Status          string                                `json:"status"`
	VirtualAccounts []VirtualAccountApplicationBankDetail `json:"virtual_accounts"`
	Error           *VirtualAccountApplicationResultError `json:"error"`
}

type VirtualAccountApplication struct {
	ApplicationID string                            `json:"application_id"`
	PublicVersion int64                             `json:"public_version"`
	Country       string                            `json:"country"`
	Currency      string                            `json:"currency"`
	Status        string                            `json:"status"`
	Results       []VirtualAccountApplicationResult `json:"results"`
}

type VirtualAccountApplicationResponse struct {
	Data VirtualAccountApplication `json:"data"`
}

type VirtualAccountApplicationSummary struct {
	ApplicationID string `json:"application_id"`
	PublicVersion int64  `json:"public_version"`
	Country       string `json:"country"`
	Currency      string `json:"currency"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

type ListVirtualAccountApplicationsRequest struct {
	PageNumber int
	PageSize   int
	Status     string
	Country    string
	Currency   string
}

type ListVirtualAccountApplicationsResponse struct {
	TotalPages int64                              `json:"total_pages"`
	TotalItems int64                              `json:"total_items"`
	Data       []VirtualAccountApplicationSummary `json:"data"`
}

func (c *VirtualAccountApplicationsClient) List(ctx context.Context, req *ListVirtualAccountApplicationsRequest, opts ...*common.RequestOptions) (*ListVirtualAccountApplicationsResponse, error) {
	values := url.Values{}
	values.Set("page_number", strconv.Itoa(req.PageNumber))
	values.Set("page_size", strconv.Itoa(req.PageSize))
	if req.Status != "" {
		values.Set("status", req.Status)
	}
	if req.Country != "" {
		values.Set("country", req.Country)
	}
	if req.Currency != "" {
		values.Set("currency", req.Currency)
	}
	var resp ListVirtualAccountApplicationsResponse
	if err := c.client.GetWithOptions(ctx, "/v1/virtual/applications?"+values.Encode(), &resp, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to list virtual account applications: %w", err)
	}
	return &resp, nil
}

func (c *VirtualAccountApplicationsClient) Retrieve(ctx context.Context, applicationID string, opts ...*common.RequestOptions) (*VirtualAccountApplicationResponse, error) {
	var resp VirtualAccountApplicationResponse
	path := "/v1/virtual/applications/" + url.PathEscape(applicationID)
	if err := c.client.GetWithOptions(ctx, path, &resp, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to retrieve virtual account application: %w", err)
	}
	return &resp, nil
}
