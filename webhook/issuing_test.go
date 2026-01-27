package webhook

import (
	"encoding/json"
	"testing"
)

func TestParseCardCreateSucceededEvent(t *testing.T) {
	jsonPayload := `{
		"version": "V1.6.0",
		"event_name": "ISSUING",
		"event_type": "card.create.succeeded",
		"event_id": "13ced737-9f47-481f-a054-d789078327d3",
		"source_id": "16f6c372-e69e-432c-9af3-14fa84800224",
		"data": {
			"card_available_balance": "10000",
			"card_bin": "49372418",
			"card_id": "a738d29b-3dd7-4fe4-9119-3a3024100f30",
			"card_number": "49372418****4306",
			"card_product_id": "5a9239b7-1618-41a5-989b-4a41fc2d4856",
			"card_scheme": "VISA",
			"card_status": "ACTIVE",
			"cardholder": {
				"cardholder_id": "a88465b4-f9f2-45f6-bc28-ecadaad1062f",
				"cardholder_status": "SUCCESS",
				"create_time": "2025-11-11T09:45:41+08:00",
				"email": "utddluuo62196793096@uqpay.gov",
				"first_name": "Corki",
				"last_name": "Jarvan IV"
			},
			"form_factor": "VIRTUAL",
			"metadata": {},
			"mode_type": "SINGLE",
			"risk_control": {},
			"spending_limits": [
				{
					"amount": "2500",
					"interval": "PER_TRANSACTION"
				}
			]
		}
	}`

	var event Event
	err := json.Unmarshal([]byte(jsonPayload), &event)
	if err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Verify event envelope
	if event.Version != "V1.6.0" {
		t.Errorf("Expected version V1.6.0, got %s", event.Version)
	}
	if event.EventName != EventNameIssuing {
		t.Errorf("Expected event name %s, got %s", EventNameIssuing, event.EventName)
	}
	if event.EventType != EventTypeCardCreateSucceeded {
		t.Errorf("Expected event type %s, got %s", EventTypeCardCreateSucceeded, event.EventType)
	}
	if event.EventID != "13ced737-9f47-481f-a054-d789078327d3" {
		t.Errorf("Expected event ID 13ced737-9f47-481f-a054-d789078327d3, got %s", event.EventID)
	}

	// Verify helper methods
	if !event.IsIssuingEvent() {
		t.Error("Expected IsIssuingEvent to return true")
	}
	if !event.IsCardEvent() {
		t.Error("Expected IsCardEvent to return true")
	}

	// Parse card data
	cardData, err := event.ParseCardData()
	if err != nil {
		t.Fatalf("Failed to parse card data: %v", err)
	}

	// Verify card data fields
	if cardData.CardID != "a738d29b-3dd7-4fe4-9119-3a3024100f30" {
		t.Errorf("Expected card ID a738d29b-3dd7-4fe4-9119-3a3024100f30, got %s", cardData.CardID)
	}
	if cardData.CardNumber != "49372418****4306" {
		t.Errorf("Expected card number 49372418****4306, got %s", cardData.CardNumber)
	}
	if cardData.CardBin != "49372418" {
		t.Errorf("Expected card bin 49372418, got %s", cardData.CardBin)
	}
	if cardData.CardScheme != CardSchemeVisa {
		t.Errorf("Expected card scheme %s, got %s", CardSchemeVisa, cardData.CardScheme)
	}
	if cardData.CardStatus != CardStatusActive {
		t.Errorf("Expected card status %s, got %s", CardStatusActive, cardData.CardStatus)
	}
	if cardData.CardAvailableBalance != "10000" {
		t.Errorf("Expected card available balance 10000, got %s", cardData.CardAvailableBalance)
	}
	if cardData.FormFactor != FormFactorVirtual {
		t.Errorf("Expected form factor %s, got %s", FormFactorVirtual, cardData.FormFactor)
	}
	if cardData.ModeType != ModeTypeSingle {
		t.Errorf("Expected mode type %s, got %s", ModeTypeSingle, cardData.ModeType)
	}

	// Verify cardholder
	if cardData.Cardholder == nil {
		t.Fatal("Expected cardholder to be present")
	}
	if cardData.Cardholder.CardholderID != "a88465b4-f9f2-45f6-bc28-ecadaad1062f" {
		t.Errorf("Expected cardholder ID a88465b4-f9f2-45f6-bc28-ecadaad1062f, got %s", cardData.Cardholder.CardholderID)
	}
	if cardData.Cardholder.CardholderStatus != CardholderStatusSuccess {
		t.Errorf("Expected cardholder status %s, got %s", CardholderStatusSuccess, cardData.Cardholder.CardholderStatus)
	}
	if cardData.Cardholder.FirstName != "Corki" {
		t.Errorf("Expected first name Corki, got %s", cardData.Cardholder.FirstName)
	}
	if cardData.Cardholder.LastName != "Jarvan IV" {
		t.Errorf("Expected last name Jarvan IV, got %s", cardData.Cardholder.LastName)
	}
	if cardData.Cardholder.Email != "utddluuo62196793096@uqpay.gov" {
		t.Errorf("Expected email utddluuo62196793096@uqpay.gov, got %s", cardData.Cardholder.Email)
	}

	// Verify spending limits
	if len(cardData.SpendingLimits) != 1 {
		t.Fatalf("Expected 1 spending limit, got %d", len(cardData.SpendingLimits))
	}
	if cardData.SpendingLimits[0].Amount != "2500" {
		t.Errorf("Expected spending limit amount 2500, got %s", cardData.SpendingLimits[0].Amount)
	}
	if cardData.SpendingLimits[0].Interval != SpendingIntervalPerTransaction {
		t.Errorf("Expected spending limit interval %s, got %s", SpendingIntervalPerTransaction, cardData.SpendingLimits[0].Interval)
	}
}

func TestParseCardStatusUpdateSucceededEvent(t *testing.T) {
	jsonPayload := `{
		"version": "V1.6.0",
		"event_name": "ISSUING",
		"event_type": "card.status.update.succeeded",
		"event_id": "b833b7bd-ea47-46f0-a9aa-cd611ecd1cee",
		"source_id": "c0aba11b-4b4f-4d3d-a7e1-9feb0211497c",
		"data": {
			"card_id": "a738d29b-3dd7-4fe4-9119-3a3024100f30",
			"card_number": "49372418****4306",
			"card_status": "BLOCKED",
			"update_reason": "test",
			"update_time": "2026-01-26T15:29:59+08:00"
		}
	}`

	var event Event
	err := json.Unmarshal([]byte(jsonPayload), &event)
	if err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Verify event envelope
	if event.Version != "V1.6.0" {
		t.Errorf("Expected version V1.6.0, got %s", event.Version)
	}
	if event.EventName != EventNameIssuing {
		t.Errorf("Expected event name %s, got %s", EventNameIssuing, event.EventName)
	}
	if event.EventType != EventTypeCardStatusUpdateSucceeded {
		t.Errorf("Expected event type %s, got %s", EventTypeCardStatusUpdateSucceeded, event.EventType)
	}

	// Verify helper methods
	if !event.IsIssuingEvent() {
		t.Error("Expected IsIssuingEvent to return true")
	}
	if !event.IsCardEvent() {
		t.Error("Expected IsCardEvent to return true")
	}
	if !event.IsCardStatusUpdateEvent() {
		t.Error("Expected IsCardStatusUpdateEvent to return true")
	}

	// Parse card status update data
	statusData, err := event.ParseCardStatusUpdateData()
	if err != nil {
		t.Fatalf("Failed to parse card status update data: %v", err)
	}

	// Verify card status update data fields
	if statusData.CardID != "a738d29b-3dd7-4fe4-9119-3a3024100f30" {
		t.Errorf("Expected card ID a738d29b-3dd7-4fe4-9119-3a3024100f30, got %s", statusData.CardID)
	}
	if statusData.CardNumber != "49372418****4306" {
		t.Errorf("Expected card number 49372418****4306, got %s", statusData.CardNumber)
	}
	if statusData.CardStatus != CardStatusBlocked {
		t.Errorf("Expected card status %s, got %s", CardStatusBlocked, statusData.CardStatus)
	}
	if statusData.UpdateReason != "test" {
		t.Errorf("Expected update reason 'test', got %s", statusData.UpdateReason)
	}
	if statusData.UpdateTime != "2026-01-26T15:29:59+08:00" {
		t.Errorf("Expected update time 2026-01-26T15:29:59+08:00, got %s", statusData.UpdateTime)
	}
}

func TestParseIssuingFeeCardEvent(t *testing.T) {
	jsonPayload := `{
		"version": "V1.6.0",
		"event_name": "ISSUING",
		"event_type": "issuing.fee.card",
		"event_id": "a3b5e621-2a3b-4f33-a066-9e1f98867927",
		"source_id": "f2115314-78b8-414a-9b46-bf6f3c39fd38",
		"data": {
			"billing_amount": "1",
			"billing_currency": "USD",
			"card_available_balance": "9999",
			"card_id": "a738d29b-3dd7-4fe4-9119-3a3024100f30",
			"card_number": "49372418****4306",
			"cardholder_id": "a88465b4-f9f2-45f6-bc28-ecadaad1062f",
			"posted_time": "2026-01-26T15:24:06+08:00",
			"reference_id": "f2115314-78b8-414a-9b46-bf6f3c39fd38",
			"remark": "deduct virtual card [4306] maintenance fee 1 USD",
			"short_reference_id": "CL260126-22YO6L67NMDC",
			"transaction_amount": "1",
			"transaction_currency": "USD",
			"transaction_status": "APPROVED",
			"transaction_time": "2026-01-26T15:24:06+08:00",
			"transaction_type": "FEE"
		}
	}`

	var event Event
	err := json.Unmarshal([]byte(jsonPayload), &event)
	if err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Verify event envelope
	if event.Version != "V1.6.0" {
		t.Errorf("Expected version V1.6.0, got %s", event.Version)
	}
	if event.EventName != EventNameIssuing {
		t.Errorf("Expected event name %s, got %s", EventNameIssuing, event.EventName)
	}
	if event.EventType != EventTypeIssuingFeeCard {
		t.Errorf("Expected event type %s, got %s", EventTypeIssuingFeeCard, event.EventType)
	}

	// Verify helper methods
	if !event.IsIssuingEvent() {
		t.Error("Expected IsIssuingEvent to return true")
	}
	if !event.IsCardTransactionEvent() {
		t.Error("Expected IsCardTransactionEvent to return true")
	}

	// Parse card transaction data
	txnData, err := event.ParseCardTransactionData()
	if err != nil {
		t.Fatalf("Failed to parse card transaction data: %v", err)
	}

	// Verify card transaction data fields
	if txnData.CardID != "a738d29b-3dd7-4fe4-9119-3a3024100f30" {
		t.Errorf("Expected card ID a738d29b-3dd7-4fe4-9119-3a3024100f30, got %s", txnData.CardID)
	}
	if txnData.CardNumber != "49372418****4306" {
		t.Errorf("Expected card number 49372418****4306, got %s", txnData.CardNumber)
	}
	if txnData.CardholderID != "a88465b4-f9f2-45f6-bc28-ecadaad1062f" {
		t.Errorf("Expected cardholder ID a88465b4-f9f2-45f6-bc28-ecadaad1062f, got %s", txnData.CardholderID)
	}
	if txnData.CardAvailableBalance != "9999" {
		t.Errorf("Expected card available balance 9999, got %s", txnData.CardAvailableBalance)
	}
	if txnData.TransactionAmount != "1" {
		t.Errorf("Expected transaction amount 1, got %s", txnData.TransactionAmount)
	}
	if txnData.TransactionCurrency != "USD" {
		t.Errorf("Expected transaction currency USD, got %s", txnData.TransactionCurrency)
	}
	if txnData.BillingAmount != "1" {
		t.Errorf("Expected billing amount 1, got %s", txnData.BillingAmount)
	}
	if txnData.BillingCurrency != "USD" {
		t.Errorf("Expected billing currency USD, got %s", txnData.BillingCurrency)
	}
	if txnData.TransactionStatus != TransactionStatusApproved {
		t.Errorf("Expected transaction status %s, got %s", TransactionStatusApproved, txnData.TransactionStatus)
	}
	if txnData.TransactionType != TransactionTypeFee {
		t.Errorf("Expected transaction type %s, got %s", TransactionTypeFee, txnData.TransactionType)
	}
	if txnData.ReferenceID != "f2115314-78b8-414a-9b46-bf6f3c39fd38" {
		t.Errorf("Expected reference ID f2115314-78b8-414a-9b46-bf6f3c39fd38, got %s", txnData.ReferenceID)
	}
	if txnData.ShortReferenceID != "CL260126-22YO6L67NMDC" {
		t.Errorf("Expected short reference ID CL260126-22YO6L67NMDC, got %s", txnData.ShortReferenceID)
	}
	if txnData.Remark != "deduct virtual card [4306] maintenance fee 1 USD" {
		t.Errorf("Expected remark 'deduct virtual card [4306] maintenance fee 1 USD', got %s", txnData.Remark)
	}
}

func TestParseCardUpdateSucceededEvent(t *testing.T) {
	jsonPayload := `{
		"version": "V1.6.0",
		"event_name": "ISSUING",
		"event_type": "card.update.succeeded",
		"event_id": "ec14696a-9867-474f-98e1-096010828bf8",
		"source_id": "e3c76528-8c3a-42c8-b1be-40bb16c6d9c0",
		"data": {
			"available_balance": "9999",
			"card_bin": "49372418",
			"card_currency": "USD",
			"card_id": "a738d29b-3dd7-4fe4-9119-3a3024100f30",
			"card_limit": "0",
			"card_number": "49372418****4306",
			"card_order_id": "e3c76528-8c3a-42c8-b1be-40bb16c6d9c0",
			"card_product_id": "5a9239b7-1618-41a5-989b-4a41fc2d4856",
			"card_scheme": "VISA",
			"card_status": "ACTIVE",
			"cardholder": {
				"cardholder_id": "a88465b4-f9f2-45f6-bc28-ecadaad1062f",
				"cardholder_status": "SUCCESS",
				"country_code": "SO",
				"create_time": "2025-11-11T09:45:41+08:00",
				"date_of_birth": "1998-10-11",
				"email": "utddluuo62196793096@uqpay.gov",
				"first_name": "Corki",
				"last_name": "Jarvan IV",
				"number_of_cards": 3,
				"phone_number": "612607870"
			},
			"form_factor": "VIRTUAL",
			"metadata": {
				"": ""
			},
			"mode_type": "SINGLE",
			"no_pin_payment_amount": "USD",
			"order_status": "success",
			"risk_control": {
				"allow_3ds_transactions": "N"
			},
			"spending_controls": [
				{
					"amount": "3500",
					"interval": "PER_TRANSACTION"
				}
			]
		}
	}`

	var event Event
	err := json.Unmarshal([]byte(jsonPayload), &event)
	if err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Verify event envelope
	if event.Version != "V1.6.0" {
		t.Errorf("Expected version V1.6.0, got %s", event.Version)
	}
	if event.EventName != EventNameIssuing {
		t.Errorf("Expected event name %s, got %s", EventNameIssuing, event.EventName)
	}
	if event.EventType != EventTypeCardUpdateSucceeded {
		t.Errorf("Expected event type %s, got %s", EventTypeCardUpdateSucceeded, event.EventType)
	}

	// Verify helper methods
	if !event.IsIssuingEvent() {
		t.Error("Expected IsIssuingEvent to return true")
	}
	if !event.IsCardEvent() {
		t.Error("Expected IsCardEvent to return true")
	}
	if !event.IsCardCreateOrUpdateEvent() {
		t.Error("Expected IsCardCreateOrUpdateEvent to return true")
	}

	// Parse card data
	cardData, err := event.ParseCardData()
	if err != nil {
		t.Fatalf("Failed to parse card data: %v", err)
	}

	// Verify card data fields
	if cardData.CardID != "a738d29b-3dd7-4fe4-9119-3a3024100f30" {
		t.Errorf("Expected card ID a738d29b-3dd7-4fe4-9119-3a3024100f30, got %s", cardData.CardID)
	}
	if cardData.CardNumber != "49372418****4306" {
		t.Errorf("Expected card number 49372418****4306, got %s", cardData.CardNumber)
	}
	if cardData.CardOrderID != "e3c76528-8c3a-42c8-b1be-40bb16c6d9c0" {
		t.Errorf("Expected card order ID e3c76528-8c3a-42c8-b1be-40bb16c6d9c0, got %s", cardData.CardOrderID)
	}
	if cardData.CardCurrency != "USD" {
		t.Errorf("Expected card currency USD, got %s", cardData.CardCurrency)
	}
	if cardData.CardLimit != "0" {
		t.Errorf("Expected card limit 0, got %s", cardData.CardLimit)
	}
	if cardData.AvailableBalance != "9999" {
		t.Errorf("Expected available balance 9999, got %s", cardData.AvailableBalance)
	}
	if cardData.OrderStatus != "success" {
		t.Errorf("Expected order status success, got %s", cardData.OrderStatus)
	}
	if cardData.NoPinPaymentAmount != "USD" {
		t.Errorf("Expected no pin payment amount USD, got %s", cardData.NoPinPaymentAmount)
	}

	// Verify GetAvailableBalance helper
	if cardData.GetAvailableBalance() != "9999" {
		t.Errorf("Expected GetAvailableBalance to return 9999, got %s", cardData.GetAvailableBalance())
	}

	// Verify cardholder with extended fields
	if cardData.Cardholder == nil {
		t.Fatal("Expected cardholder to be present")
	}
	if cardData.Cardholder.CountryCode != "SO" {
		t.Errorf("Expected country code SO, got %s", cardData.Cardholder.CountryCode)
	}
	if cardData.Cardholder.DateOfBirth != "1998-10-11" {
		t.Errorf("Expected date of birth 1998-10-11, got %s", cardData.Cardholder.DateOfBirth)
	}
	if cardData.Cardholder.NumberOfCards != 3 {
		t.Errorf("Expected number of cards 3, got %d", cardData.Cardholder.NumberOfCards)
	}
	if cardData.Cardholder.PhoneNumber != "612607870" {
		t.Errorf("Expected phone number 612607870, got %s", cardData.Cardholder.PhoneNumber)
	}

	// Verify spending controls
	if len(cardData.SpendingControls) != 1 {
		t.Fatalf("Expected 1 spending control, got %d", len(cardData.SpendingControls))
	}
	if cardData.SpendingControls[0].Amount != "3500" {
		t.Errorf("Expected spending control amount 3500, got %s", cardData.SpendingControls[0].Amount)
	}

	// Verify GetSpendingLimits helper returns spending_controls when available
	limits := cardData.GetSpendingLimits()
	if len(limits) != 1 || limits[0].Amount != "3500" {
		t.Errorf("Expected GetSpendingLimits to return spending_controls")
	}

	// Verify risk control
	if cardData.RiskControl == nil {
		t.Fatal("Expected risk control to be present")
	}
	if cardData.RiskControl.Allow3DSTransactions != "N" {
		t.Errorf("Expected allow_3ds_transactions N, got %s", cardData.RiskControl.Allow3DSTransactions)
	}
}

func TestParseCardRechargeSucceededEvent(t *testing.T) {
	jsonPayload := `{
		"version": "V1.6.0",
		"event_name": "ISSUING",
		"event_type": "card.recharge.succeeded",
		"event_id": "f7b130cd-efa9-4f2a-9178-8c98dfa59726",
		"source_id": "a8c73130-8bda-45f7-bc88-f5edc3fbaf74",
		"data": {
			"amount": "200",
			"card_available_balance": "10199",
			"card_currency": "USD",
			"card_id": "a738d29b-3dd7-4fe4-9119-3a3024100f30",
			"card_status": "ACTIVE",
			"complete_time": "2026-01-26T15:39:02+08:00",
			"order_status": "SUCCESS",
			"update_time": "2026-01-26T15:39:02+08:00"
		}
	}`

	var event Event
	err := json.Unmarshal([]byte(jsonPayload), &event)
	if err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Verify event envelope
	if event.Version != "V1.6.0" {
		t.Errorf("Expected version V1.6.0, got %s", event.Version)
	}
	if event.EventName != EventNameIssuing {
		t.Errorf("Expected event name %s, got %s", EventNameIssuing, event.EventName)
	}
	if event.EventType != EventTypeCardRechargeSucceeded {
		t.Errorf("Expected event type %s, got %s", EventTypeCardRechargeSucceeded, event.EventType)
	}

	// Verify helper methods
	if !event.IsIssuingEvent() {
		t.Error("Expected IsIssuingEvent to return true")
	}
	if !event.IsCardEvent() {
		t.Error("Expected IsCardEvent to return true")
	}
	if !event.IsCardRechargeEvent() {
		t.Error("Expected IsCardRechargeEvent to return true")
	}

	// Parse card recharge data
	rechargeData, err := event.ParseCardRechargeData()
	if err != nil {
		t.Fatalf("Failed to parse card recharge data: %v", err)
	}

	// Verify card recharge data fields
	if rechargeData.CardID != "a738d29b-3dd7-4fe4-9119-3a3024100f30" {
		t.Errorf("Expected card ID a738d29b-3dd7-4fe4-9119-3a3024100f30, got %s", rechargeData.CardID)
	}
	if rechargeData.Amount != "200" {
		t.Errorf("Expected amount 200, got %s", rechargeData.Amount)
	}
	if rechargeData.CardCurrency != "USD" {
		t.Errorf("Expected card currency USD, got %s", rechargeData.CardCurrency)
	}
	if rechargeData.CardAvailableBalance != "10199" {
		t.Errorf("Expected card available balance 10199, got %s", rechargeData.CardAvailableBalance)
	}
	if rechargeData.CardStatus != CardStatusActive {
		t.Errorf("Expected card status %s, got %s", CardStatusActive, rechargeData.CardStatus)
	}
	if rechargeData.OrderStatus != "SUCCESS" {
		t.Errorf("Expected order status SUCCESS, got %s", rechargeData.OrderStatus)
	}
	if rechargeData.CompleteTime != "2026-01-26T15:39:02+08:00" {
		t.Errorf("Expected complete time 2026-01-26T15:39:02+08:00, got %s", rechargeData.CompleteTime)
	}
	if rechargeData.UpdateTime != "2026-01-26T15:39:02+08:00" {
		t.Errorf("Expected update time 2026-01-26T15:39:02+08:00, got %s", rechargeData.UpdateTime)
	}
}

func TestParseCardActivationCodeEvent(t *testing.T) {
	jsonPayload := `{
		"version": "V1.6.0",
		"event_name": "ISSUING",
		"event_type": "card.activation.code",
		"event_id": "8a78af1e-de83-43a5-b177-ecbc6a8a9fc6",
		"source_id": "52d51b68-9691-4448-ac5f-dcbbef9c62f8",
		"data": {
			"card_id": "52d51b68-9691-4448-ac5f-dcbbef9c62f8",
			"card_number": "12345678****0000",
			"activation_code": "12341234"
		}
	}`

	var event Event
	err := json.Unmarshal([]byte(jsonPayload), &event)
	if err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Verify event envelope
	if event.Version != "V1.6.0" {
		t.Errorf("Expected version V1.6.0, got %s", event.Version)
	}
	if event.EventName != EventNameIssuing {
		t.Errorf("Expected event name %s, got %s", EventNameIssuing, event.EventName)
	}
	if event.EventType != EventTypeCardActivationCode {
		t.Errorf("Expected event type %s, got %s", EventTypeCardActivationCode, event.EventType)
	}

	// Verify helper methods
	if !event.IsIssuingEvent() {
		t.Error("Expected IsIssuingEvent to return true")
	}
	if !event.IsCardEvent() {
		t.Error("Expected IsCardEvent to return true")
	}
	if !event.IsCardActivationCodeEvent() {
		t.Error("Expected IsCardActivationCodeEvent to return true")
	}

	// Parse card activation code data
	codeData, err := event.ParseCardActivationCodeData()
	if err != nil {
		t.Fatalf("Failed to parse card activation code data: %v", err)
	}

	// Verify card activation code data fields
	if codeData.CardID != "52d51b68-9691-4448-ac5f-dcbbef9c62f8" {
		t.Errorf("Expected card ID 52d51b68-9691-4448-ac5f-dcbbef9c62f8, got %s", codeData.CardID)
	}
	if codeData.CardNumber != "12345678****0000" {
		t.Errorf("Expected card number 12345678****0000, got %s", codeData.CardNumber)
	}
	if codeData.ActivationCode != "12341234" {
		t.Errorf("Expected activation code 12341234, got %s", codeData.ActivationCode)
	}
}

func TestCardEventTypeChecks(t *testing.T) {
	testCases := []struct {
		eventType   string
		isCardEvent bool
	}{
		{EventTypeCardCreateSucceeded, true},
		{EventTypeCardCreateFailed, true},
		{EventTypeCardUpdateSucceeded, true},
		{EventTypeCardUpdateFailed, true},
		{EventTypeCardRechargeSucceeded, true},
		{EventTypeCardRechargeFailed, true},
		{EventTypeCardActivationCode, true},
		{EventTypeCardActivated, true},
		{EventTypeCardSuspended, true},
		{EventTypeCardClosed, true},
		{EventTypeCardStatusUpdateSucceeded, true},
		{EventTypeCardStatusUpdateFailed, true},
		{EventTypePaymentIntentCreated, false},
		{EventTypeAccountCreate, false},
	}

	for _, tc := range testCases {
		event := Event{EventType: tc.eventType}
		if event.IsCardEvent() != tc.isCardEvent {
			t.Errorf("Expected IsCardEvent for %s to be %v", tc.eventType, tc.isCardEvent)
		}
	}
}

func TestParseCardDataError(t *testing.T) {
	// Test with non-card event type
	event := Event{
		EventType: EventTypePaymentIntentCreated,
		Data:      []byte(`{}`),
	}

	_, err := event.ParseCardData()
	if err == nil {
		t.Error("Expected error when parsing non-card event as card data")
	}
}
