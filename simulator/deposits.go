package simulator

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

type DepositsClient struct{ client *common.APIClient }

type CreateDepositRequest struct {
	Amount                float64 `json:"amount"`
	Currency              string  `json:"currency"`
	ReceiverAccountNumber string  `json:"receiver_account_number,omitempty"`
	SenderSwiftCode       string  `json:"sender_swift_code"`
	SenderAccountNumber   string  `json:"sender_account_number,omitempty"`
	SenderCountry         string  `json:"sender_country,omitempty"`
	SenderName            string  `json:"sender_name,omitempty"`
}

type DepositSender struct {
	SenderName          string `json:"sender_name"`
	SenderCountry       string `json:"sender_country"`
	SenderAccountNumber string `json:"sender_account_number"`
	SenderSwiftCode     string `json:"sender_swift_code"`
}

type CreateDepositResponse struct {
	DepositID             string        `json:"deposit_id"`
	ShortReferenceID      string        `json:"short_reference_id"`
	Amount                string        `json:"amount"`
	Currency              string        `json:"currency"`
	DepositStatus         string        `json:"deposit_status"`
	CreateTime            string        `json:"create_time"`
	CompleteTime          string        `json:"complete_time"`
	ReceiverAccountNumber string        `json:"receiver_account_number"`
	Sender                DepositSender `json:"sender"`
}

func (c *DepositsClient) Create(ctx context.Context, req *CreateDepositRequest, opts ...*common.RequestOptions) (*CreateDepositResponse, error) {
	var resp CreateDepositResponse
	if err := c.client.PostWithOptions(ctx, "/v1/simulation/deposit", req, &resp, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to simulate deposit: %w", err)
	}
	return &resp, nil
}
