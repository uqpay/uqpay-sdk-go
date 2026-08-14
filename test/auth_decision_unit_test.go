package test

import (
	"testing"

	uqpay "github.com/uqpay/uqpay-sdk-go/v2"
	"github.com/uqpay/uqpay-sdk-go/v2/configuration"
)

func TestAuthDecisionIsAvailableOnIssuingClient(t *testing.T) {
	client, err := uqpay.NewClient("test-client", "test-key", configuration.Sandbox())
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if client.Issuing.AuthDecision == nil {
		t.Fatal("Issuing.AuthDecision is nil")
	}
	if client.Issuing.MerchantBrands == nil {
		t.Fatal("Issuing.MerchantBrands was removed while adding AuthDecision")
	}
}
