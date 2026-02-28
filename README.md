# UQPAY Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/uqpay/uqpay-sdk-go.svg)](https://pkg.go.dev/github.com/uqpay/uqpay-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/uqpay/uqpay-sdk-go)](https://goreportcard.com/report/github.com/uqpay/uqpay-sdk-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Official Go SDK for UQPAY - A comprehensive payment and card issuing platform.

## Features

- 🚀 **Easy Integration** - Simple and intuitive API
- 💳 **Card Issuing** - Create and manage virtual/physical cards
- 👤 **Cardholder Management** - Full cardholder lifecycle management
- 💰 **Card Operations** - Recharge, withdraw, freeze, and manage card status
- 🏦 **Banking** - Balances, transfers, deposits, payouts, beneficiaries, virtual accounts, conversions, and exchange rates
- 📊 **Transaction Tracking** - Real-time transaction monitoring
- 🔒 **Secure** - Built-in OAuth2 authentication with automatic token management
- ⚡ **Idempotency** - Automatic idempotency key generation for safe retries
- 🌍 **Multi-Environment** - Support for Sandbox and Production environments

## Installation

```bash
go get github.com/uqpay/uqpay-sdk-go@latest
```

**Requirements**: Go 1.19 or higher

## Quick Start

### Initialize the SDK

```go
package main

import (
    "context"
    "log"

    "github.com/uqpay/uqpay-sdk-go"
    "github.com/uqpay/uqpay-sdk-go/configuration"
    "github.com/uqpay/uqpay-sdk-go/issuing"
)

func main() {
    // Create client with Sandbox environment
    client, err := uqpay.NewClient(
        "your-client-id",
        "your-api-key",
        configuration.Sandbox(),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Your code here...
}
```

### Create a Cardholder

```go
cardholder, err := client.Issuing.Cardholders.Create(ctx, &issuing.CreateCardholderRequest{
    Email:       "user@example.com",
    PhoneNumber: "1234567890",
    FirstName:   "John",
    LastName:    "Doe",
    CountryCode: "US",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Cardholder ID: %s\n", cardholder.CardholderID)
```

### List Card Products

```go
products, err := client.Issuing.Products.List(ctx, &issuing.ListProductsRequest{
    PageSize:   10,
    PageNumber: 1,
})
if err != nil {
    log.Fatal(err)
}

for _, product := range products.Data {
    fmt.Printf("Product: %s (%s)\n", product.ProductID, product.CardScheme)
}
```

### Create a Card

```go
card, err := client.Issuing.Cards.Create(ctx, &issuing.CreateCardRequest{
    CardCurrency:  "USD",
    CardholderID:  cardholder.CardholderID,
    CardProductID: "product-id",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Card ID: %s (Status: %s)\n", card.CardID, card.CardStatus)
```

### Get Secure Card Details

```go
secureInfo, err := client.Issuing.Cards.GetSecure(ctx, card.CardID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Card Number: %s\n", secureInfo.CardNumber)
fmt.Printf("CVV: %s\n", secureInfo.CVV)
fmt.Printf("Expiry: %s\n", secureInfo.ExpiryDate)
```

### Recharge a Card

```go
order, err := client.Issuing.Cards.Recharge(ctx, card.CardID, &issuing.CardOrderRequest{
    Amount: 100.50,
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Recharge Order ID: %s\n", order.OrderID)
```

### Withdraw from a Card

```go
order, err := client.Issuing.Cards.Withdraw(ctx, card.CardID, &issuing.CardOrderRequest{
    Amount: 50.00,
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Withdraw Order ID: %s (Status: %s)\n", order.OrderID, order.Status)
```

### Update Card Status

```go
// Freeze a card
err = client.Issuing.Cards.UpdateStatus(ctx, card.CardID, &issuing.UpdateCardStatusRequest{
    CardStatus: "FROZEN",
})
if err != nil {
    log.Fatal(err)
}

// Unfreeze a card
err = client.Issuing.Cards.UpdateStatus(ctx, card.CardID, &issuing.UpdateCardStatusRequest{
    CardStatus: "ACTIVE",
})
```

### List Transactions

```go
transactions, err := client.Issuing.Transactions.List(ctx, &issuing.ListTransactionsRequest{
    PageSize:   10,
    PageNumber: 1,
    CardID:     card.CardID,
})
if err != nil {
    log.Fatal(err)
}

for _, txn := range transactions.Data {
    fmt.Printf("Transaction: %s - %s %s\n",
        txn.TransactionID,
        txn.TransactionAmount,
        txn.TransactionCurrency,
    )
}
```

## Configuration

### Environment Configuration

```go
// Sandbox (for testing)
client, err := uqpay.NewClient(clientID, apiKey, configuration.Sandbox())

// Production
client, err := uqpay.NewClient(clientID, apiKey, configuration.Production())

// Custom environment
client, err := uqpay.NewClient(clientID, apiKey, &configuration.Config{
    BaseURL: "https://custom-api.example.com/api",
})
```

### Environment Variables

For testing, you can use environment variables:

```bash
export UQPAY_CLIENT_ID="your-client-id"
export UQPAY_API_KEY="your-api-key"
```

Or create a `.env` file (see `.env.example`):

```bash
cp .env.example .env
# Edit .env with your credentials
```

## API Coverage

### Banking API

> 详细使用文档: [docs/banking-usage.md](docs/banking-usage.md)

| Resource | Operations |
|----------|------------|
| **Balances** | Get, List, ListTransactions |
| **Transfers** | Create, List, Get |
| **Deposits** | List, Get |
| **Beneficiaries** | Create, List, Get, Update, Delete, ListPaymentMethods, Check |
| **Payouts** | Create, List, Get |
| **Virtual Accounts** | Create, List |
| **Conversions** | CreateQuote, Create, List, Get, ListConversionDates |
| **Exchange Rates** | List |

### Issuing API

| Resource | Operations |
|----------|------------|
| **Cardholders** | Create, Get, List |
| **Cards** | Create, Get, GetSecure, List, Recharge, Withdraw, UpdateStatus |
| **Transactions** | Get, List |
| **Products** | List |

## Error Handling

The SDK returns detailed error information:

```go
card, err := client.Issuing.Cards.Get(ctx, cardID)
if err != nil {
    // Error includes HTTP status code and API error details
    log.Printf("Error: %v\n", err)
    return
}
```

Example error format:
```
failed to get card: 404: card_not_found: Card not found (HTTP 404)
```

## Features

### Automatic OAuth2 Token Management

The SDK automatically handles OAuth2 authentication:
- Fetches access tokens using client credentials
- Caches tokens until expiration
- Automatically refreshes expired tokens
- Thread-safe token management

### Automatic Idempotency Keys

Every API request automatically includes a unique idempotency key to ensure safe retries and prevent duplicate operations.

### Type Safety

All API requests and responses are strongly typed with proper Go structs:

```go
type Card struct {
    CardID           string `json:"card_id"`
    CardNumber       string `json:"card_number"`
    CardCurrency     string `json:"card_currency"`
    CardholderID     string `json:"cardholder_id"`
    CardProductID    string `json:"card_product_id"`
    CardStatus       string `json:"card_status"`
    AvailableBalance string `json:"available_balance"`
    CreateTime       string `json:"create_time"`
}
```

## Testing

### Run Tests

```bash
# Set credentials
export UQPAY_CLIENT_ID="your-client-id"
export UQPAY_API_KEY="your-api-key"

# Run all tests
go test -v ./test/...

# Run specific test
go test -v ./test -run TestCardholders

# Skip integration tests (for CI)
export SKIP_INTEGRATION_TESTS=true
go test -v ./...
```

### Test Coverage

The SDK includes comprehensive integration tests covering:
- Cardholder creation and retrieval
- Card product listing
- Card creation and management
- Secure card information retrieval
- Card recharge operations
- Card withdraw operations
- Card status updates
- Transaction listing and retrieval

## Development

### Project Structure

```
uqpay-sdk-go/
├── auth/              # OAuth2 authentication
├── banking/           # Banking API client
│   ├── balances.go
│   ├── beneficiaries.go
│   ├── conversion.go
│   ├── deposits.go
│   ├── exchange_rates.go
│   ├── payouts.go
│   ├── transfers.go
│   └── virtual_accounts.go
├── common/            # Shared API client
├── configuration/     # Environment configuration
├── issuing/           # Issuing API client
│   ├── cardholders.go
│   ├── cards.go
│   ├── transactions.go
│   └── products.go
├── docs/              # Documentation
├── test/              # Integration tests
└── version.go         # SDK version
```

### Build

```bash
# Build all packages
go build ./...

# Run code formatting
gofmt -w .

# Run linter
go vet ./...
```

## Versioning

This SDK follows [Semantic Versioning](https://semver.org/).

Current version: `1.0.3`

To install a specific version:

```bash
go get github.com/uqpay/uqpay-sdk-go@v1.0.3
```

View all releases: [GitHub Releases](https://github.com/uqpay/uqpay-sdk-go/releases)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

- 📧 Email: support@uqpay.com
- 📚 Documentation: https://docs.uqpay.com
- 🐛 Issues: https://github.com/uqpay/uqpay-sdk-go/issues

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Changelog

### v1.0.3 (Latest)
- Initial stable release
- Full Issuing API support
- Automatic OAuth2 token management
- Comprehensive test coverage
- Production ready

---

Made with ❤️ by UQPAY
