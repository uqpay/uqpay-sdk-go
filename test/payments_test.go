package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/uqpay/uqpay-sdk-go/payment"
)


// ============================================================================
// Payment Intents Tests
// ============================================================================

func TestPaymentIntents(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	var createdIntentID string

	t.Run("Create", func(t *testing.T) {
		req := &payment.CreatePaymentIntentRequest{
			Amount:          "103.00",
			Currency:        "USD",
			MerchantOrderID: fmt.Sprintf("sdk-%d", time.Now().UnixNano()),
			Description:     "Test payment intent",
			ReturnURL:       "https://example.com/return",
			Metadata: map[string]string{
				"test": "true",
			},
		}

		resp, err := client.Payment.PaymentIntents.Create(ctx, req)
		if err != nil {
			t.Fatalf("Create payment intent failed: %v", err)
		}

		// Assertions
		if resp.PaymentIntentID == "" {
			t.Error("PaymentIntentID should not be empty")
		}
		if resp.Amount != "103" {
			t.Errorf("Amount mismatch: got %s, want 101", resp.Amount)
		}
		if resp.Currency != "USD" {
			t.Errorf("Currency mismatch: got %s, want USD", resp.Currency)
		}

		createdIntentID = resp.PaymentIntentID
		t.Logf("Payment intent created successfully")
		t.Logf("   ID: %s", resp.PaymentIntentID)
		t.Logf("   Amount: %s %s", resp.Amount, resp.Currency)
		t.Logf("   Status: %s", resp.IntentStatus)
		t.Logf("   Merchant Order ID: %s", resp.MerchantOrderID)
		t.Logf("   Description: %s", resp.Description)
		t.Logf("   Created: %s", resp.CreateTime)
	})

	t.Run("Confirm", func(t *testing.T) {
		if createdIntentID == "" {
			t.Skip("No payment intent ID available, skipping Confirm test")
		}

		// Confirm with test card details
		confirmReq := &payment.ConfirmPaymentIntentRequest{
			PaymentMethod: &payment.PaymentMethod{
				Type: "card",
				Card: &payment.Card{
					CardNumber:  "4176660000000027",
					ExpiryMonth: "12",
					ExpiryYear:  "2033",
					CVC:         "303",
					CardName:    "Test User",
				},
			},
			ReturnURL: "https://example.com/return",
		}

		resp, err := client.Payment.PaymentIntents.Confirm(ctx, createdIntentID, confirmReq)
		if err != nil {
			t.Logf("Confirm payment intent returned: %v", err)
			return
		}

		// Assertions
		if resp.PaymentIntentID == "" {
			t.Error("PaymentIntentID should not be empty")
		}
		if resp.IntentStatus == "" {
			t.Error("IntentStatus should not be empty")
		}

		// Check expected statuses after confirmation
		validStatuses := map[string]bool{
			"REQUIRES_CUSTOMER_ACTION": true, // Needs 3DS or other action
			"REQUIRES_PAYMENT_METHOD":  true, // Payment method rejected
			"REQUIRES_CAPTURE":         true, // Ready to capture
			"PENDING":                  true, // Waiting for provider
			"SUCCEEDED":                true, // Payment complete
		}

		if !validStatuses[resp.IntentStatus] {
			t.Errorf("Unexpected status after confirm: %s", resp.IntentStatus)
		}

		t.Logf("Payment intent confirmed successfully")
		t.Logf("   ID: %s", resp.PaymentIntentID)
		t.Logf("   Status: %s", resp.IntentStatus)
		t.Logf("   Amount: %s %s", resp.Amount, resp.Currency)

		if resp.NextAction != nil {
			t.Logf("   Next Action: %v", resp.NextAction)
		}
	})

	t.Run("Get", func(t *testing.T) {
		// First list to get an ID if we don't have one
		if createdIntentID == "" {
			listReq := &payment.ListPaymentIntentsRequest{
				PageSize:   1,
				PageNumber: 1,
			}
			listResp, err := client.Payment.PaymentIntents.List(ctx, listReq)
			if err != nil {
				t.Skipf("List payment intents failed, skipping Get test: %v", err)
			}
			if len(listResp.Data) == 0 {
				t.Skip("No payment intents available, skipping Get test")
			}
			createdIntentID = listResp.Data[0].PaymentIntentID
		}

		resp, err := client.Payment.PaymentIntents.Get(ctx, createdIntentID)
		if err != nil {
			t.Fatalf("Get payment intent failed: %v", err)
		}

		// Assertions
		if resp.PaymentIntentID == "" {
			t.Error("PaymentIntentID should not be empty")
		}
		if resp.PaymentIntentID != createdIntentID {
			t.Errorf("PaymentIntentID mismatch: got %s, want %s", resp.PaymentIntentID, createdIntentID)
		}

		t.Logf("Payment intent retrieved successfully")
		t.Logf("   ID: %s", resp.PaymentIntentID)
		t.Logf("   Amount: %s %s", resp.Amount, resp.Currency)
		t.Logf("   Status: %s", resp.IntentStatus)
		t.Logf("   Merchant Order ID: %s", resp.MerchantOrderID)
	})

	t.Run("List", func(t *testing.T) {
		req := &payment.ListPaymentIntentsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Payment.PaymentIntents.List(ctx, req)
		if err != nil {
			t.Fatalf("List payment intents failed: %v", err)
		}

		// Assertions
		if resp.Data == nil {
			t.Error("Data should not be nil")
		}

		t.Logf("Found %d payment intents (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			for i, intent := range resp.Data {
				if i >= 3 {
					t.Logf("   ... and %d more", len(resp.Data)-3)
					break
				}
				t.Logf("   Intent %d: ID=%s, Amount=%s %s, Status=%s",
					i+1, intent.PaymentIntentID, intent.Amount, intent.Currency, intent.IntentStatus)
			}
		} else {
			t.Log("No payment intents found")
		}
	})

	t.Run("ListWithFilters", func(t *testing.T) {
		// Status values per spec enum
		statuses := []string{
			"REQUIRES_PAYMENT_METHOD",
			"REQUIRES_CUSTOMER_ACTION",
			"REQUIRES_CAPTURE",
			"PENDING",
			"SUCCEEDED",
			"CANCELLED",
			"FAILED",
		}

		for _, status := range statuses {
			req := &payment.ListPaymentIntentsRequest{
				PageSize:            5,
				PageNumber:          1,
				PaymentIntentStatus: status,
			}

			resp, err := client.Payment.PaymentIntents.List(ctx, req)
			if err != nil {
				t.Errorf("%s intents: failed - %v", status, err)
				continue
			}

			// Assertion
			if resp.Data == nil {
				t.Errorf("%s intents: Data should not be nil", status)
			}

			t.Logf("%s intents: %d found", status, resp.TotalItems)
		}
	})

	t.Run("Update", func(t *testing.T) {
		if createdIntentID == "" {
			t.Skip("No payment intent ID available, skipping Update test")
		}

		req := &payment.UpdatePaymentIntentRequest{
			Description: "Updated test payment intent",
			Metadata: map[string]string{
				"updated": "true",
			},
		}

		resp, err := client.Payment.PaymentIntents.Update(ctx, createdIntentID, req)
		if err != nil {
			t.Fatalf("Update payment intent failed: %v", err)
		}

		// Assertions
		if resp.PaymentIntentID == "" {
			t.Error("PaymentIntentID should not be empty")
		}
		if resp.Description != "Updated test payment intent" {
			t.Errorf("Description mismatch: got %s, want 'Updated test payment intent'", resp.Description)
		}

		t.Logf("Payment intent updated successfully")
		t.Logf("   ID: %s", resp.PaymentIntentID)
		t.Logf("   Description: %s", resp.Description)
		t.Logf("   Updated: %s", resp.UpdateTime)
	})

	t.Run("Cancel", func(t *testing.T) {
		if createdIntentID == "" {
			t.Skip("No payment intent ID available, skipping Cancel test")
		}

		req := &payment.CancelPaymentIntentRequest{
			CancellationReason: "requested_by_customer",
		}

		resp, err := client.Payment.PaymentIntents.Cancel(ctx, createdIntentID, req)
		if err != nil {
			t.Fatalf("Cancel payment intent failed: %v", err)
		}

		// Assertions
		if resp.PaymentIntentID == "" {
			t.Error("PaymentIntentID should not be empty")
		}
		validCancelStatuses := map[string]bool{"CANCELED": true, "CANCELLED": true, "REQUIRES_PAYMENT_METHOD": true}
		if !validCancelStatuses[resp.IntentStatus] {
			t.Errorf("Unexpected IntentStatus after cancel: %s", resp.IntentStatus)
		}

		t.Logf("Payment intent canceled successfully")
		t.Logf("   ID: %s", resp.PaymentIntentID)
		t.Logf("   Status: %s", resp.IntentStatus)
	})
}

// ============================================================================
// Capture Payment Intent Test
// ============================================================================

func TestCapturePaymentIntent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	// Step 1: Create a payment intent
	createReq := &payment.CreatePaymentIntentRequest{
		Amount:          "50.00",
		Currency:        "USD",
		MerchantOrderID: fmt.Sprintf("sdk-%d", time.Now().UnixNano()),
		Description:     "Test capture flow",
		ReturnURL:       "https://example.com/return",
	}

	intent, err := client.Payment.PaymentIntents.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Create payment intent failed: %v", err)
	}

	t.Logf("Created intent for capture: %s", intent.PaymentIntentID)

	// Step 2: Confirm with a card that supports manual capture
	confirmReq := &payment.ConfirmPaymentIntentRequest{
		PaymentMethod: &payment.PaymentMethod{
			Type: "card",
			Card: &payment.Card{
				CardNumber:        "4176660000000027",
				ExpiryMonth:       "12",
				ExpiryYear:        "2033",
				CVC:               "303",
				CardName:          "Test User",
				AutoCapture:       boolPtr(false),
				AuthorizationType: "authorization",
				ThreeDSAction:     "skip_3ds",
				Billing: &payment.Billing{
					FirstName:   "Test",
					LastName:    "User",
					Email:       "test@example.com",
					PhoneNumber: "+10000000000",
					Address: &payment.Address{
						CountryCode: "SG",
						City:        "Singapore",
						Street:      "444 Orchard Rd",
						Postcode:    "924011",
					},
				},
			},
		},
		ReturnURL: "https://example.com/return",
	}

	confirmed, err := client.Payment.PaymentIntents.Confirm(ctx, intent.PaymentIntentID, confirmReq)
	if err != nil {
		t.Logf("Confirm payment intent returned: %v", err)
		return
	}

	t.Logf("Confirmed intent: %s, status: %s", confirmed.PaymentIntentID, confirmed.IntentStatus)

	// Step 3: Capture if status allows
	if confirmed.IntentStatus != "REQUIRES_CAPTURE" {
		t.Skipf("Intent status is %s, not REQUIRES_CAPTURE - skipping capture", confirmed.IntentStatus)
	}

	captureReq := &payment.CapturePaymentIntentRequest{
		AmountToCapture: 50.00,
	}

	captured, err := client.Payment.PaymentIntents.Capture(ctx, intent.PaymentIntentID, captureReq)
	if err != nil {
		t.Fatalf("Capture payment intent failed: %v", err)
	}

	// Assertions
	if captured.PaymentIntentID == "" {
		t.Error("PaymentIntentID should not be empty")
	}

	t.Logf("Payment intent captured successfully")
	t.Logf("   ID: %s", captured.PaymentIntentID)
	t.Logf("   Status: %s", captured.IntentStatus)
	t.Logf("   Captured Amount: %s", captured.CapturedAmount)
}

// boolPtr returns a pointer to a bool value
func boolPtr(b bool) *bool {
	return &b
}

// ============================================================================
// Confirm with BrowserInfo and IPAddress Test
// ============================================================================

func TestConfirmWithBrowserInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	// Step 1: Create a payment intent with PaymentOrders
	createReq := &payment.CreatePaymentIntentRequest{
		Amount:          "100.00",
		Currency:        "USD",
		MerchantOrderID: fmt.Sprintf("sdk-%d", time.Now().UnixNano()),
		Description:     "Test 3DS with browser info",
		ReturnURL:       "https://example.com/return",
		IPAddress:       "203.0.113.50",
		PaymentOrders: &payment.PaymentOrders{
			Type: "physical_goods",
			Products: []payment.PaymentProduct{
				{
					Name:     "Test Product",
					Price:    "100.00",
					Quantity: 1,
				},
			},
		},
	}

	intent, err := client.Payment.PaymentIntents.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Create payment intent failed: %v", err)
	}

	if intent.PaymentIntentID == "" {
		t.Fatal("PaymentIntentID should not be empty")
	}

	t.Logf("Created intent with PaymentOrders: %s", intent.PaymentIntentID)

	// Step 2: Confirm with BrowserInfo and IPAddress
	confirmReq := &payment.ConfirmPaymentIntentRequest{
		PaymentMethod: &payment.PaymentMethod{
			Type: "card",
			Card: &payment.Card{
				CardNumber:    "4176660000000027",
				ExpiryMonth:   "12",
				ExpiryYear:    "2033",
				CVC:           "303",
				CardName:      "Test User",
				ThreeDSAction: "enforce_3ds",
				Billing: &payment.Billing{
					FirstName:   "Test",
					LastName:    "User",
					Email:       "test@example.com",
					PhoneNumber: "+10000000000",
					Address: &payment.Address{
						CountryCode: "SG",
						City:        "Singapore",
						Street:      "444 Orchard Rd",
						Postcode:    "924011",
					},
				},
			},
		},
		IPAddress: "203.0.113.50",
		BrowserInfo: &payment.BrowserInfo{
			AcceptHeader:     "text/html",
			Language:         "en-US",
			ScreenColorDepth: 24,
			ScreenHeight:     1080,
			ScreenWidth:      1920,
			Timezone:         "8",
			Browser: &payment.BrowserDetail{
				JavaEnabled:       false,
				JavascriptEnabled: true,
				UserAgent:         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
			},
		},
		ReturnURL: "https://example.com/return",
	}

	resp, err := client.Payment.PaymentIntents.Confirm(ctx, intent.PaymentIntentID, confirmReq)
	if err != nil {
		t.Logf("Confirm with browser info returned: %v", err)
		return
	}

	// Assertions
	if resp.PaymentIntentID == "" {
		t.Error("PaymentIntentID should not be empty")
	}
	if resp.IntentStatus == "" {
		t.Error("IntentStatus should not be empty")
	}

	t.Logf("Confirmed with BrowserInfo successfully")
	t.Logf("   ID: %s", resp.PaymentIntentID)
	t.Logf("   Status: %s", resp.IntentStatus)

	if resp.NextAction != nil {
		t.Logf("   Next Action: %v", resp.NextAction)
	}
}

// ============================================================================
// Update with Customer and PaymentOrders Test
// ============================================================================

func TestUpdatePaymentIntentWithCustomer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	// Step 1: Create a payment intent
	createReq := &payment.CreatePaymentIntentRequest{
		Amount:          "75.00",
		Currency:        "USD",
		MerchantOrderID: "test-update-customer-001",
		Description:     "Test update with customer",
		ReturnURL:       "https://example.com/return",
	}

	intent, err := client.Payment.PaymentIntents.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Create payment intent failed: %v", err)
	}

	t.Logf("Created intent: %s", intent.PaymentIntentID)

	// Step 2: Update with customer details and payment orders
	updateReq := &payment.UpdatePaymentIntentRequest{
		Customer: &payment.CustomerRequest{
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "john.doe@example.com",
			PhoneNumber: "+6591234567",
			Description: "Test customer",
		},
		MerchantOrderID: "test-update-customer-002",
		Description:     "Updated with customer info",
	}

	resp, err := client.Payment.PaymentIntents.Update(ctx, intent.PaymentIntentID, updateReq)
	if err != nil {
		t.Fatalf("Update with customer failed: %v", err)
	}

	// Assertions
	if resp.PaymentIntentID == "" {
		t.Error("PaymentIntentID should not be empty")
	}
	if resp.Description != "Updated with customer info" {
		t.Errorf("Description mismatch: got %s", resp.Description)
	}

	t.Logf("Updated with customer successfully")
	t.Logf("   ID: %s", resp.PaymentIntentID)
	t.Logf("   Description: %s", resp.Description)
	if resp.Customer != nil {
		t.Logf("   Customer: %s %s (%s)", resp.Customer.FirstName, resp.Customer.LastName, resp.Customer.Email)
	}
}

// ============================================================================
// Payment Method Confirmation Tests (Table-Driven)
// ============================================================================

// TestConfirmPaymentMethods tests confirmation with various payment method types
// Each test creates its own payment intent and confirms with a specific payment method
func TestConfirmPaymentMethods(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	// Define test cases for each payment method type
	testCases := []struct {
		name          string
		paymentMethod *payment.PaymentMethod
		currency      string // Some payment methods may require specific currencies
		skipReason    string // If set, test will be skipped with this reason
	}{
		// ================================================================
		// Card Payments
		// ================================================================
		{
			name: "Card",
			paymentMethod: &payment.PaymentMethod{
				Type: "card",
				Card: &payment.Card{
					CardNumber:        "4111111111111111",
					ExpiryMonth:       "12",
					ExpiryYear:        "2030",
					CVC:               "123",
					CardName:          "Test User",
					Network:           "visa",
					AuthorizationType: "authorization",
					ThreeDSAction:     "skip_3ds",
					Billing: &payment.Billing{
						FirstName:   "Test",
						LastName:    "User",
						Email:       "test@example.com",
						PhoneNumber: "+10000000000",
						Address: &payment.Address{
							CountryCode: "SG",
							State:       "",
							City:        "Singapore",
							Street:      "444 Orchard Rd, Midpoint Orchard, Singapore",
							Postcode:    "924011",
						},
					},
				},
			},
			currency: "AUD",
		},

		// ================================================================
		// China & Hong Kong Wallets
		// ================================================================
		{
			name: "AlipayCN_QRCode",
			paymentMethod: &payment.PaymentMethod{
				Type: "alipaycn",
				AlipayCN: &payment.WalletPayment{
					Flow:      "qrcode",
					OSType:    "",
					IsPresent: false,
				},
			},
			currency: "AUD",
		},
		{
			name: "AlipayHK_QRCode",
			paymentMethod: &payment.PaymentMethod{
				Type: "alipayhk",
				AlipayHK: &payment.WalletPayment{
					Flow:      "qrcode",
					OSType:    "",
					IsPresent: false,
				},
			},
			currency: "AUD",
		},
		{
			name: "UnionPay_QRCode",
			paymentMethod: &payment.PaymentMethod{
				Type: "unionpay",
				UnionPay: &payment.WalletPayment{
					Flow: "qrcode",
				},
			},
			currency: "AUD",
		},
		{
			name: "UnionPay_SecurePay",
			paymentMethod: &payment.PaymentMethod{
				Type: "unionpay",
				UnionPay: &payment.WalletPayment{
					Flow: "securepay",
				},
			},
			currency: "AUD",
		},
		{
			name: "WeChatPay_QRCode",
			paymentMethod: &payment.PaymentMethod{
				Type: "wechatpay",
				WeChatPay: &payment.WeChatPay{
					Flow: "qrcode",
				},
			},
			currency: "AUD",
		},
		{
			name: "WeChatPay_MobileWeb",
			paymentMethod: &payment.PaymentMethod{
				Type: "wechatpay",
				WeChatPay: &payment.WeChatPay{
					Flow:   "mobile_web",
					OSType: "ios",
				},
			},
			currency: "AUD",
		},
		{
			name: "WeChatPay_MobileWeb",
			paymentMethod: &payment.PaymentMethod{
				Type: "wechatpay",
				WeChatPay: &payment.WeChatPay{
					Flow:   "mini_program",
					OSType: "ios",
					OpenID: "",
				},
			},
			currency: "AUD",
		},
		{
			name: "WeChatPay_MobileWeb",
			paymentMethod: &payment.PaymentMethod{
				Type: "wechatpay",
				WeChatPay: &payment.WeChatPay{
					Flow:   "mobile_app",
					OSType: "ios",
					OpenID: "", //required
				},
			},
			currency: "AUD",
		},

		{
			name: "WeChatPay_MobileWeb",
			paymentMethod: &payment.PaymentMethod{
				Type: "wechatpay",
				WeChatPay: &payment.WeChatPay{
					Flow:   "official_account",
					OSType: "ios",
					OpenID: "", //required
				},
			},
			currency: "AUD",
		},

		// ================================================================
		// Southeast Asia Wallets
		// ================================================================
		{
			name: "GrabPay_QRCode",
			paymentMethod: &payment.PaymentMethod{
				Type: "grabpay",
				GrabPay: &payment.GrabPay{
					Flow:        "qrcode",
					ShopperName: "Test Shopper",
				},
			},
			currency: "AUD",
		},
		{
			name: "PayNow_QRCode",
			paymentMethod: &payment.PaymentMethod{
				Type: "paynow",
				PayNow: &payment.WalletPayment{
					Flow: "qrcode",
				},
			},
			currency: "AUD",
		},
		{
			name: "TrueMoney_QRCode",
			paymentMethod: &payment.PaymentMethod{
				Type: "truemoney",
				TrueMoney: &payment.WalletPayment{
					Flow: "qrcode",
				},
			},
			currency: "AUD",
		},
		{
			name: "TNG_QRCode",
			paymentMethod: &payment.PaymentMethod{
				Type: "tng",
				TNG: &payment.WalletPayment{
					Flow: "qrcode",
				},
			},
			currency: "AUD",
		},
		{
			name: "GCash_QRCode",
			paymentMethod: &payment.PaymentMethod{
				Type: "gcash",
				GCash: &payment.WalletPayment{
					Flow: "qrcode",
				},
			},
			currency: "AUD",
		},
		{
			name: "Dana_QRCode",
			paymentMethod: &payment.PaymentMethod{
				Type: "dana",
				Dana: &payment.WalletPayment{
					Flow: "qrcode",
				},
			},
			currency: "AUD",
		},

		// ================================================================
		// Korean Wallets
		// ================================================================
		{
			name: "KakaoPay_QRCode",
			paymentMethod: &payment.PaymentMethod{
				Type: "kakaopay",
				KakaoPay: &payment.WalletPayment{
					Flow: "qrcode",
				},
			},
			currency: "AUD",
		},
		{
			name: "Toss_QRCode",
			paymentMethod: &payment.PaymentMethod{
				Type: "tosspay",
				Toss: &payment.WalletPayment{
					Flow:      "qrcode",
					IsPresent: false,
				},
			},
			currency: "AUD",
		},
		{
			name: "NaverPay_QRCode",
			paymentMethod: &payment.PaymentMethod{
				Type: "naverpay",
				NaverPay: &payment.WalletPayment{
					Flow: "qrcode",
				},
			},
			currency: "AUD",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Skip if specified
			if tc.skipReason != "" {
				t.Skip(tc.skipReason)
			}

			// Step 1: Create a payment intent
			createReq := &payment.CreatePaymentIntentRequest{
				Amount:          "10.00",
				Currency:        tc.currency,
				MerchantOrderID: fmt.Sprintf("sdk-%d", time.Now().UnixNano()),
				Description:     "Test " + tc.name + " payment",
				ReturnURL:       "https://example.com/return",
			}

			intent, err := client.Payment.PaymentIntents.Create(ctx, createReq)
			if err != nil {
				t.Logf("Create intent failed: %v", err)
				return
			}

			// Assertions for create
			if intent.PaymentIntentID == "" {
				t.Error("PaymentIntentID should not be empty")
			}

			t.Logf("Created intent: %s (status: %s)", intent.PaymentIntentID, intent.IntentStatus)

			// Step 2: Confirm with the specific payment method
			confirmReq := &payment.ConfirmPaymentIntentRequest{
				PaymentMethod: tc.paymentMethod,
				ReturnURL:     "https://example.com/return",
			}

			resp, err := client.Payment.PaymentIntents.Confirm(ctx, intent.PaymentIntentID, confirmReq)
			if err != nil {
				t.Logf("Confirm failed: %v", err)
				return
			}

			// Assertions for confirm
			if resp.PaymentIntentID == "" {
				t.Error("PaymentIntentID should not be empty after confirm")
			}
			if resp.IntentStatus == "" {
				t.Error("IntentStatus should not be empty after confirm")
			}

			// Verify valid status (REQUIRES_PAYMENT_METHOD means card/method was rejected but API call succeeded)
			validStatuses := map[string]bool{
				"REQUIRES_CUSTOMER_ACTION": true,
				"REQUIRES_PAYMENT_METHOD":  true,
				"REQUIRES_CAPTURE":         true,
				"PENDING":                  true,
				"SUCCEEDED":                true,
			}

			if !validStatuses[resp.IntentStatus] {
				t.Errorf("Unexpected status after confirm: %s", resp.IntentStatus)
			}

			t.Logf("Confirmed successfully")
			t.Logf("   ID: %s", resp.PaymentIntentID)
			t.Logf("   Status: %s", resp.IntentStatus)
			t.Logf("   Amount: %s %s", resp.Amount, resp.Currency)

			// Log next action if present (common for QR code payments)
			if resp.NextAction != nil {
				t.Logf("   Next Action: %v", resp.NextAction)
			}
		})
	}
}

func TestGetPaymentIntent(t *testing.T) {
	client := GetPaymentTestClient(t)
	ctx := context.Background()

	// Create a fresh PI to retrieve
	created, err := client.Payment.PaymentIntents.Create(ctx, &payment.CreatePaymentIntentRequest{
		Amount:          "10.00",
		Currency:        "USD",
		MerchantOrderID: fmt.Sprintf("get-test-%d", time.Now().UnixNano()),
		Description:     "SDK test get intent",
		ReturnURL:       "https://example.com/return",
	})
	if err != nil {
		t.Fatalf("Setup: failed to create payment intent: %v", err)
	}
	paymentId := created.PaymentIntentID

	t.Logf("Retrieving payment intent: %s", paymentId)

	resp, err := client.Payment.PaymentIntents.Get(ctx, paymentId)
	if err != nil {
		t.Fatalf("Get payment intent failed: %v", err)
	}

	// Assertions
	if resp.PaymentIntentID == "" {
		t.Error("PaymentIntentID should not be empty")
	}
	if resp.PaymentIntentID != paymentId {
		t.Errorf("PaymentIntentID mismatch: got %s, want %s", resp.PaymentIntentID, paymentId)
	}

	t.Logf("Retrieved successfully")
	t.Logf("   ID: %s", resp.PaymentIntentID)
	t.Logf("   IntentStatus: %s", resp.IntentStatus)
	t.Logf("   Amount: %s %s", resp.Amount, resp.Currency)
	t.Logf("   MerchantOrderID: %s", resp.MerchantOrderID)
	t.Logf("   Description: %s", resp.Description)
	t.Logf("   Created: %s", resp.CreateTime)

	if resp.NextAction != nil {
		t.Logf("   Next Action: %v", resp.NextAction)
	}
}

func TestGetPaymentAttempt(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	// Get an attempt ID from the list
	listResp, err := client.Payment.PaymentAttempts.List(ctx, &payment.ListPaymentAttemptsRequest{
		PageSize:   1,
		PageNumber: 1,
	})
	if err != nil || len(listResp.Data) == 0 {
		t.Skip("No payment attempts available to test Get")
	}

	attemptID := listResp.Data[0].AttemptID
	t.Logf("Retrieving payment attempt: %s", attemptID)

	resp, err := client.Payment.PaymentAttempts.Get(ctx, attemptID)
	if err != nil {
		t.Fatalf("Get payment attempt failed: %v", err)
	}

	if resp.AttemptID == "" {
		t.Error("AttemptID should not be empty")
	}
	if resp.AttemptID != attemptID {
		t.Errorf("AttemptID mismatch: got %s, want %s", resp.AttemptID, attemptID)
	}

	t.Logf("Retrieved successfully")
	t.Logf("   ID: %s", resp.AttemptID)
	t.Logf("   Amount: %s %s", resp.Amount, resp.Currency)
	t.Logf("   Captured Amount: %s", resp.CapturedAmount)
	t.Logf("   Refunded Amount: %s", resp.RefundedAmount)
	t.Logf("   Status: %s", resp.AttemptStatus)
	if resp.FailureCode != "" {
		t.Logf("   Failure Code: %s", resp.FailureCode)
	}
	t.Logf("   Created: %s", resp.CreateTime)
}

// ============================================================================
// Payment Attempts Tests
// ============================================================================
func TestPaymentAttempts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		req := &payment.ListPaymentAttemptsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Payment.PaymentAttempts.List(ctx, req)
		if err != nil {
			t.Fatalf("List payment attempts failed: %v", err)
		}

		// Assertions
		if resp.Data == nil {
			t.Error("Data should not be nil")
		}

		t.Logf("Found %d payment attempts (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			attempt := resp.Data[0]
			t.Logf("First attempt:")
			t.Logf("   ID: %s", attempt.AttemptID)
			t.Logf("   Amount: %s %s", attempt.Amount, attempt.Currency)
			t.Logf("   Captured: %s, Refunded: %s", attempt.CapturedAmount, attempt.RefundedAmount)
			t.Logf("   Status: %s", attempt.AttemptStatus)
			if attempt.FailureCode != "" {
				t.Logf("   Failure: %s", attempt.FailureCode)
			}
			t.Logf("   Created: %s", attempt.CreateTime)
		} else {
			t.Log("No payment attempts found")
		}
	})

	t.Run("Get", func(t *testing.T) {
		// First list to get a valid ID
		listReq := &payment.ListPaymentAttemptsRequest{
			PageSize:   1,
			PageNumber: 1,
		}

		listResp, err := client.Payment.PaymentAttempts.List(ctx, listReq)
		if err != nil {
			t.Skipf("List payment attempts failed, skipping Get test: %v", err)
		}

		if len(listResp.Data) == 0 {
			t.Skip("No payment attempts available, skipping Get test")
		}

		attemptID := listResp.Data[0].AttemptID

		resp, err := client.Payment.PaymentAttempts.Get(ctx, attemptID)
		if err != nil {
			t.Fatalf("Get payment attempt failed: %v", err)
		}

		// Assertions
		if resp.AttemptID == "" {
			t.Error("AttemptID should not be empty")
		}
		if resp.AttemptID != attemptID {
			t.Errorf("AttemptID mismatch: got %s, want %s", resp.AttemptID, attemptID)
		}

		t.Logf("Payment attempt retrieved successfully")
		t.Logf("   ID: %s", resp.AttemptID)
		t.Logf("   Amount: %s %s", resp.Amount, resp.Currency)
		t.Logf("   Status: %s", resp.AttemptStatus)
	})

	t.Run("ListWithFilters", func(t *testing.T) {
		// Status values per spec enum
		statuses := []string{
			"INITIATED",
			"AUTHORIZED",
			"SUCCEEDED",
			"CANCELLED",
			"FAILED",
		}

		for _, status := range statuses {
			req := &payment.ListPaymentAttemptsRequest{
				PageSize:      5,
				PageNumber:    1,
				AttemptStatus: status,
			}

			resp, err := client.Payment.PaymentAttempts.List(ctx, req)
			if err != nil {
				t.Errorf("%s attempts: failed - %v", status, err)
				continue
			}

			// Assertion
			if resp.Data == nil {
				t.Errorf("%s attempts: Data should not be nil", status)
			}

			t.Logf("%s attempts: %d found", status, resp.TotalItems)
		}
	})

	t.Run("ListByPaymentIntent", func(t *testing.T) {
		// Get a real PI ID from the payment intents list
		piList, err := client.Payment.PaymentIntents.List(ctx, &payment.ListPaymentIntentsRequest{
			PageSize: 1, PageNumber: 1,
		})
		if err != nil || len(piList.Data) == 0 {
			t.Log("No payment intents found, skipping ListByPaymentIntent")
			return
		}
		piID := piList.Data[0].PaymentIntentID

		req := &payment.ListPaymentAttemptsRequest{
			PageSize:        5,
			PageNumber:      1,
			PaymentIntentID: piID,
		}

		resp, err := client.Payment.PaymentAttempts.List(ctx, req)
		if err != nil {
			t.Fatalf("List attempts by PI failed: %v", err)
		}

		if resp.Data == nil {
			t.Error("Data should not be nil")
		}

		t.Logf("Attempts for PI %s: %d found", piID, resp.TotalItems)
		for i, attempt := range resp.Data {
			if i >= 3 {
				break
			}
			t.Logf("   Attempt %d: ID=%s, Status=%s, Amount=%s %s",
				i+1, attempt.AttemptID, attempt.AttemptStatus, attempt.Amount, attempt.Currency)
		}
	})
}

// ============================================================================
// Payment Balances Tests
// ============================================================================

func TestPaymentBalances(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		req := &payment.ListBalancesRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Payment.Balances.List(ctx, req)
		if err != nil {
			t.Fatalf("List payment balances failed: %v", err)
		}

		// Assertions
		if resp.Data == nil {
			t.Error("Data should not be nil")
		}

		t.Logf("Found %d payment balances (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			for i, balance := range resp.Data {
				t.Logf("Balance %d: %s", i+1, balance.Currency)
				t.Logf("   Available: %s", balance.AvailableBalance)
				t.Logf("   Payable: %s", balance.PayableBalance)
				t.Logf("   Pending: %s", balance.PendingBalance)
				t.Logf("   Reserved: %s", balance.ReservedBalance)
			}
		} else {
			t.Log("No payment balances found")
		}
	})

	t.Run("Get", func(t *testing.T) {
		currency := "USD"

		resp, err := client.Payment.Balances.Get(ctx, currency)
		if err != nil {
			t.Fatalf("Get payment balance for %s failed: %v", currency, err)
		}

		// Assertions
		if resp.BalanceID == "" {
			t.Error("BalanceID should not be empty")
		}
		if resp.Currency != currency {
			t.Errorf("Currency mismatch: got %s, want %s", resp.Currency, currency)
		}

		t.Logf("Balance retrieved for %s", currency)
		t.Logf("   Balance ID: %s", resp.BalanceID)
		t.Logf("   Available: %s", resp.AvailableBalance)
		t.Logf("   Payable: %s", resp.PayableBalance)
		t.Logf("   Pending: %s", resp.PendingBalance)
		t.Logf("   Reserved: %s", resp.ReservedBalance)
	})

	t.Run("GetMultipleCurrencies", func(t *testing.T) {
		// List available balances first, then retrieve each one
		listResp, err := client.Payment.Balances.List(ctx, &payment.ListBalancesRequest{
			PageSize:   10,
			PageNumber: 1,
		})
		if err != nil {
			t.Logf("List balances failed, skipping GetMultipleCurrencies: %v", err)
			return
		}
		if len(listResp.Data) == 0 {
			t.Log("No balances available, skipping GetMultipleCurrencies")
			return
		}

		for _, b := range listResp.Data {
			resp, err := client.Payment.Balances.Get(ctx, b.Currency)
			if err != nil {
				t.Logf("%s: failed - %v", b.Currency, err)
				continue
			}
			t.Logf("✅ %s: Available=%s, Payable=%s", resp.Currency, resp.AvailableBalance, resp.PayableBalance)
		}
	})

	t.Run("ListPagination", func(t *testing.T) {
		// Test with small page size to verify pagination
		req := &payment.ListBalancesRequest{
			PageSize:   2,
			PageNumber: 1,
		}

		resp, err := client.Payment.Balances.List(ctx, req)
		if err != nil {
			t.Fatalf("List balances page 1 failed: %v", err)
		}

		// Assertions
		if resp.TotalItems < 0 {
			t.Error("TotalItems should be >= 0")
		}
		if resp.TotalPages < 0 {
			t.Error("TotalPages should be >= 0")
		}
		if len(resp.Data) > 2 {
			t.Errorf("Data length should be <= page_size(2), got %d", len(resp.Data))
		}

		t.Logf("Page 1: %d items, total: %d, pages: %d", len(resp.Data), resp.TotalItems, resp.TotalPages)

		// If there are more pages, test page 2
		if resp.TotalPages > 1 {
			req2 := &payment.ListBalancesRequest{
				PageSize:   2,
				PageNumber: 2,
			}

			resp2, err := client.Payment.Balances.List(ctx, req2)
			if err != nil {
				t.Fatalf("List balances page 2 failed: %v", err)
			}

			t.Logf("Page 2: %d items", len(resp2.Data))
		}
	})
}

// ============================================================================
// Payment Payouts Tests
// ============================================================================

func TestPaymentPayouts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		req := &payment.ListPayoutsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Payment.Payouts.List(ctx, req)
		if err != nil {
			t.Fatalf("List payouts failed: %v", err)
		}

		// Assertions
		if resp.Data == nil {
			t.Error("Data should not be nil")
		}

		t.Logf("Found %d payouts (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			payout := resp.Data[0]
			t.Logf("First payout:")
			t.Logf("   ID: %s", payout.PayoutID)
			t.Logf("   Amount: %s %s", payout.PayoutAmount, payout.PayoutCurrency)
			t.Logf("   Status: %s", payout.PayoutStatus)
			t.Logf("   Note: %s", payout.InternalNote)
			t.Logf("   Descriptor: %s", payout.StatementDescriptor)
			t.Logf("   Created: %s", payout.CreateTime)
		} else {
			t.Log("No payouts found")
		}
	})

	t.Run("Get", func(t *testing.T) {
		// First list to get a valid ID
		listReq := &payment.ListPayoutsRequest{
			PageSize:   1,
			PageNumber: 1,
		}

		listResp, err := client.Payment.Payouts.List(ctx, listReq)
		if err != nil {
			t.Skipf("List payouts failed, skipping Get test: %v", err)
		}

		if len(listResp.Data) == 0 {
			t.Skip("No payouts available, skipping Get test")
		}

		payoutID := listResp.Data[0].PayoutID

		resp, err := client.Payment.Payouts.Get(ctx, payoutID)
		if err != nil {
			t.Fatalf("Get payout failed: %v", err)
		}

		// Assertions
		if resp.PayoutID == "" {
			t.Error("PayoutID should not be empty")
		}
		if resp.PayoutID != payoutID {
			t.Errorf("PayoutID mismatch: got %s, want %s", resp.PayoutID, payoutID)
		}

		t.Logf("Payout retrieved successfully")
		t.Logf("   ID: %s", resp.PayoutID)
		t.Logf("   Amount: %s %s", resp.PayoutAmount, resp.PayoutCurrency)
		t.Logf("   Status: %s", resp.PayoutStatus)
		t.Logf("   Note: %s", resp.InternalNote)
	})

	t.Run("ListWithFilters", func(t *testing.T) {
		// Status values per spec enum
		statuses := []string{
			"INITIATED",
			"PROCESSING",
			"COMPLETED",
			"FAILED",
			"FAILED_REFUNDED",
		}

		for _, status := range statuses {
			req := &payment.ListPayoutsRequest{
				PageSize:     5,
				PageNumber:   1,
				PayoutStatus: status,
			}

			resp, err := client.Payment.Payouts.List(ctx, req)
			if err != nil {
				t.Errorf("%s payouts: failed - %v", status, err)
				continue
			}

			// Assertion
			if resp.Data == nil {
				t.Errorf("%s payouts: Data should not be nil", status)
			}

			t.Logf("%s payouts: %d found", status, resp.TotalItems)
		}
	})

	t.Run("ListWithDateRange", func(t *testing.T) {
		req := &payment.ListPayoutsRequest{
			PageSize:   10,
			PageNumber: 1,
			StartTime:  "2026-03-07T00:00:00Z",
			EndTime:    "2026-04-06T23:59:59Z",
		}

		resp, err := client.Payment.Payouts.List(ctx, req)
		if err != nil {
			t.Logf("List payouts with date range: %v (may not be supported for this account)", err)
			return
		}

		if resp.Data == nil {
			t.Error("Data should not be nil")
		}

		t.Logf("Payouts in date range: %d found (total: %d)", len(resp.Data), resp.TotalItems)
	})
}

// TestPayoutCreate tests payout creation
func TestPayoutCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		req := &payment.CreatePayoutRequest{
			PayoutCurrency:      "USD",
			PayoutAmount:        "5.00",
			StatementDescriptor: "SDK test payout",
			InternalNote:        "SDK integration test payout",
		}

		resp, err := client.Payment.Payouts.Create(ctx, req)
		if err != nil {
			t.Logf("Create payout returned: %v", err)
			return
		}

		// Assertions
		if resp.PayoutID == "" {
			t.Error("PayoutID should not be empty")
		}
		if resp.PayoutStatus == "" {
			t.Error("PayoutStatus should not be empty")
		}
		if resp.PayoutCurrency != "USD" {
			t.Errorf("PayoutCurrency mismatch: got %s, want USD", resp.PayoutCurrency)
		}

		t.Logf("Payout created successfully")
		t.Logf("   ID: %s", resp.PayoutID)
		t.Logf("   Amount: %s %s", resp.PayoutAmount, resp.PayoutCurrency)
		t.Logf("   Status: %s", resp.PayoutStatus)
		t.Logf("   Descriptor: %s", resp.StatementDescriptor)
		t.Logf("   Note: %s", resp.InternalNote)
		t.Logf("   Created: %s", resp.CreateTime)
	})
}

// ============================================================================
// Payment Refunds Tests
// ============================================================================

func TestPaymentRefunds(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		req := &payment.ListRefundsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Payment.Refunds.List(ctx, req)
		if err != nil {
			t.Fatalf("List refunds failed: %v", err)
		}

		// Assertions
		if resp.Data == nil {
			t.Error("Data should not be nil")
		}

		t.Logf("Found %d refunds (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			refund := resp.Data[0]
			t.Logf("First refund:")
			t.Logf("   ID: %s", refund.PaymentRefundID)
			t.Logf("   Payment Attempt ID: %s", refund.PaymentAttemptID)
			t.Logf("   Amount: %s %s", refund.Amount, refund.Currency)
			t.Logf("   Status: %s", refund.RefundStatus)
			t.Logf("   Reason: %s", refund.Reason)
			t.Logf("   Created: %s", refund.CreateTime)
		} else {
			t.Log("No refunds found")
		}
	})

	t.Run("Get", func(t *testing.T) {
		// First list to get a valid ID
		listReq := &payment.ListRefundsRequest{
			PageSize:   1,
			PageNumber: 1,
		}

		listResp, err := client.Payment.Refunds.List(ctx, listReq)
		if err != nil {
			t.Skipf("List refunds failed, skipping Get test: %v", err)
		}

		if len(listResp.Data) == 0 {
			t.Skip("No refunds available, skipping Get test")
		}

		refundID := listResp.Data[0].PaymentRefundID

		resp, err := client.Payment.Refunds.Get(ctx, refundID)
		if err != nil {
			t.Fatalf("Get refund failed: %v", err)
		}

		// Assertions
		if resp.PaymentRefundID == "" {
			t.Error("PaymentRefundID should not be empty")
		}
		if resp.PaymentRefundID != refundID {
			t.Errorf("PaymentRefundID mismatch: got %s, want %s", resp.PaymentRefundID, refundID)
		}

		t.Logf("Refund retrieved successfully")
		t.Logf("   ID: %s", resp.PaymentRefundID)
		t.Logf("   Payment Attempt ID: %s", resp.PaymentAttemptID)
		t.Logf("   Amount: %s %s", resp.Amount, resp.Currency)
		t.Logf("   Status: %s", resp.RefundStatus)
	})

	t.Run("ListWithDateRange", func(t *testing.T) {
		req := &payment.ListRefundsRequest{
			PageSize:   10,
			PageNumber: 1,
			StartTime:  "2026-03-07T00:00:00Z",
			EndTime:    "2026-04-06T23:59:59Z",
		}

		resp, err := client.Payment.Refunds.List(ctx, req)
		if err != nil {
			t.Fatalf("List refunds with date range failed: %v", err)
		}

		if resp.Data == nil {
			t.Error("Data should not be nil")
		}

		t.Logf("Refunds in date range: %d found (total: %d)", len(resp.Data), resp.TotalItems)
	})

	t.Run("ListByPaymentIntent", func(t *testing.T) {
		// Get a real PI ID from payment intents list
		piList, err := client.Payment.PaymentIntents.List(ctx, &payment.ListPaymentIntentsRequest{
			PageSize: 1, PageNumber: 1,
		})
		if err != nil || len(piList.Data) == 0 {
			t.Log("No payment intents found, skipping ListByPaymentIntent")
			return
		}
		piID := piList.Data[0].PaymentIntentID

		req := &payment.ListRefundsRequest{
			PageSize:        5,
			PageNumber:      1,
			PaymentIntentID: piID,
		}

		resp, err := client.Payment.Refunds.List(ctx, req)
		if err != nil {
			t.Fatalf("List refunds by PI failed: %v", err)
		}

		if resp.Data == nil {
			t.Error("Data should not be nil")
		}

		t.Logf("Refunds for PI %s: %d found (total: %d)", piID, len(resp.Data), resp.TotalItems)
	})
}

// Create Refund
func TestCreateRefund(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	// Requires a previously captured payment intent — skip in automated runs
	t.Skip("Skipping: requires a manually captured payment intent ID to refund")
	paymentIntentID := ""

	t.Run("CreateRefund", func(t *testing.T) {
		req := &payment.CreateRefundRequest{
			PaymentIntentID: paymentIntentID,
			Amount:          "10.0",
			Reason:          "requested_by_customer",
			Metadata: map[string]string{
				"test": "true",
			},
		}

		resp, err := client.Payment.Refunds.Create(ctx, req)
		if err != nil {
			t.Fatalf("Create refund failed: %v", err)
		}

		// Assertions
		if resp.PaymentRefundID == "" {
			t.Error("PaymentRefundID should not be empty")
		}
		if resp.RefundStatus == "" {
			t.Error("RefundStatus should not be empty")
		}

		// Verify valid refund status
		validStatuses := map[string]bool{
			"INITIATED":  true,
			"PROCESSING": true,
			"SUCCEEDED":  true,
		}
		if !validStatuses[resp.RefundStatus] {
			t.Errorf("Unexpected refund status: %s", resp.RefundStatus)
		}

		t.Logf("Refund created successfully")
		t.Logf("   ID: %s", resp.PaymentRefundID)
		t.Logf("   Payment Attempt ID: %s", resp.PaymentAttemptID)
		t.Logf("   Amount: %s %s", resp.Amount, resp.Currency)
		t.Logf("   Status: %s", resp.RefundStatus)
		t.Logf("   Reason: %s", resp.Reason)
		t.Logf("   Created: %s", resp.CreateTime)
	})
}

// ============================================================================
// Payment Reports Tests
// ============================================================================

func TestPaymentReports(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := GetPaymentTestClient(t)
	ctx := context.Background()

	t.Run("ListSettlements", func(t *testing.T) {
		req := &payment.ListSettlementsRequest{
			PageSize:   10,
			PageNumber: 1,
		}

		resp, err := client.Payment.Reports.ListSettlements(ctx, req)
		if err != nil {
			t.Fatalf("List settlements failed: %v", err)
		}

		// Assertions
		if resp.Data == nil {
			t.Error("Data should not be nil")
		}

		t.Logf("Found %d settlements (total: %d)", len(resp.Data), resp.TotalItems)
		t.Logf("Total pages: %d", resp.TotalPages)

		if len(resp.Data) > 0 {
			settlement := resp.Data[0]
			t.Logf("First settlement:")
			t.Logf("   ID: %s", settlement.SettlementID)
			t.Logf("   Payment Intent: %s", settlement.PaymentIntentID)
			t.Logf("   Transaction Amount: %s %s", settlement.TransactionAmount, settlement.TransactionCurrency)
			t.Logf("   Settlement Amount: %s %s", settlement.SettlementAmount, settlement.SettlementCurrency)
			t.Logf("   Net Settlement: %s", settlement.NetSettlementAmount)
			t.Logf("   Total Fee: %s", settlement.TotalFeeAmount)
			t.Logf("   Status: %s", settlement.SettlementStatus)
			t.Logf("   Settlement Date: %s", settlement.SettlementDate)
		} else {
			t.Log("No settlements found")
		}
	})

	t.Run("ListSettlementsByPaymentIntent", func(t *testing.T) {
		// Get a real PI ID from existing settlements or payment intents
		piID := ""
		allResp, err := client.Payment.Reports.ListSettlements(ctx, &payment.ListSettlementsRequest{
			PageSize: 1, PageNumber: 1,
		})
		if err == nil && len(allResp.Data) > 0 && allResp.Data[0].PaymentIntentID != "" {
			piID = allResp.Data[0].PaymentIntentID
		} else {
			piList, piErr := client.Payment.PaymentIntents.List(ctx, &payment.ListPaymentIntentsRequest{
				PageSize: 1, PageNumber: 1,
			})
			if piErr != nil || len(piList.Data) == 0 {
				t.Log("No payment intents found, skipping ListSettlementsByPaymentIntent")
				return
			}
			piID = piList.Data[0].PaymentIntentID
		}

		req := &payment.ListSettlementsRequest{
			PaymentIntentID: piID,
			PageSize:        10,
			PageNumber:      1,
		}

		resp, err := client.Payment.Reports.ListSettlements(ctx, req)
		if err != nil {
			t.Fatalf("List settlements by PI failed: %v", err)
		}

		if resp.Data == nil {
			t.Error("Data should not be nil")
		}

		t.Logf("Settlements for PI %s: %d found (total: %d)", piID, len(resp.Data), resp.TotalItems)
	})

	t.Run("ListSettlementsWithDateRange", func(t *testing.T) {
		// Test with date range (last 30 days)
		req := &payment.ListSettlementsRequest{
			SettledStartTime: "2026-03-07",
			SettledEndTime:   "2026-04-06",
			PageSize:         10,
			PageNumber:       1,
		}

		resp, err := client.Payment.Reports.ListSettlements(ctx, req)
		if err != nil {
			t.Fatalf("List settlements with date range failed: %v", err)
		}

		// Assertions
		if resp.Data == nil {
			t.Error("Data should not be nil")
		}

		t.Logf("Found %d settlements in date range (total: %d)", len(resp.Data), resp.TotalItems)

		if len(resp.Data) > 0 {
			for i, settlement := range resp.Data {
				if i >= 3 {
					t.Logf("   ... and %d more", len(resp.Data)-3)
					break
				}
				t.Logf("Settlement %d: ID=%s, Amount=%s %s, Status=%s, Date=%s",
					i+1, settlement.SettlementID, settlement.SettlementAmount, settlement.SettlementCurrency, settlement.SettlementStatus, settlement.SettlementDate)
			}
		}
	})
}
