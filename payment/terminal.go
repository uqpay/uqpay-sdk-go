package payment

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/v2/common"
)

type TerminalsClient struct{ client *common.APIClient }

type RegisterTerminalRequest struct {
	FirmCode      string `json:"firm_code"`
	FirmSN        string `json:"firm_sn"`
	TerminalModel string `json:"terminal_model"`
}

type RegisterTerminalResponse struct {
	CreateTime string `json:"create_time"`
	FirmSN     string `json:"firm_sn"`
	TerminalID string `json:"terminal_id"`
}

type GetPINKeyRequest struct {
	TerminalID string `json:"terminal_id"`
	PrivateKey string `json:"prv_key"`
}

type GetPINKeyResponse struct {
	EncryptedPINKey string `json:"encrypt_pin_key"`
	PINKeyExpiresAt string `json:"pin_key_expire"`
	TerminalID      string `json:"terminal_id"`
}

func (c *TerminalsClient) Register(ctx context.Context, req *RegisterTerminalRequest, opts ...*common.RequestOptions) (*RegisterTerminalResponse, error) {
	var resp RegisterTerminalResponse
	opt := requestOptionsWithClientID(c.client.Config.ClientID, opts...)
	if err := c.client.PostWithOptions(ctx, "/v2/terminal/register", req, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to register terminal: %w", err)
	}
	return &resp, nil
}

func (c *TerminalsClient) GetPINKey(ctx context.Context, req *GetPINKeyRequest, opts ...*common.RequestOptions) (*GetPINKeyResponse, error) {
	var resp GetPINKeyResponse
	opt := requestOptionsWithClientID(c.client.Config.ClientID, opts...)
	if err := c.client.PostWithOptions(ctx, "/v2/terminal/getPinKey", req, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to get terminal PIN key: %w", err)
	}
	return &resp, nil
}
