# UQPAY Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/uqpay/uqpay-sdk-go/v2.svg)](https://pkg.go.dev/github.com/uqpay/uqpay-sdk-go/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/uqpay/uqpay-sdk-go/v2)](https://goreportcard.com/report/github.com/uqpay/uqpay-sdk-go/v2)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Official Go SDK for UQPAY - A comprehensive payment and card issuing platform.

## Features

- 🚀 **Easy Integration** - Simple and intuitive API
- 💳 **Card Issuing** - Create and manage virtual/physical cards
- 👤 **Cardholder Management** - Full cardholder lifecycle management
- 💰 **Card Operations** - Recharge, withdraw, freeze, and manage card status
- 🏦 **Banking** - Balances, transfers, deposits, payouts, beneficiaries, virtual accounts, conversions, and exchange rates
- 📊 **Transaction Tracking** - Real-time transaction monitoring
- 🔒 **Secure** - Automatic UQPAY Access Token management
- 🔑 **Authorization Decisions** - PGP-encrypted real-time card transaction decisions
- ⚡ **Idempotency** - Automatic idempotency key generation for safe retries
- 🌍 **Multi-Environment** - Support for Sandbox and Production environments

## Installation

```bash
go get github.com/uqpay/uqpay-sdk-go/v2@latest
```

**Requirements**: Go 1.19 or higher

## Quick Start

### Initialize the SDK

```go
package main

import (
    "context"
    "log"

    "github.com/uqpay/uqpay-sdk-go/v2"
    "github.com/uqpay/uqpay-sdk-go/v2/configuration"
    "github.com/uqpay/uqpay-sdk-go/v2/issuing"
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

Store credentials in your environment and pass them explicitly when you create the client. The SDK does not load environment variables or `.env` files automatically.

```bash
export UQPAY_CLIENT_ID="your-client-id"
export UQPAY_API_KEY="your-api-key"
```

```go
client, err := uqpay.NewClient(
    os.Getenv("UQPAY_CLIENT_ID"),
    os.Getenv("UQPAY_API_KEY"),
    configuration.Sandbox(),
)
```

The repository's integration-test helper can load `.env` for local SDK testing. That test-only behavior is not part of the SDK runtime.

## API Coverage

The root client exposes these product namespaces and resource clients. See the [UQPAY API documentation](https://developers.uqpay.com/) and the exported Go types for operation-level request and response details.

| Namespace | Product | Resources |
|-----------|---------|-----------|
| `client.Connect` | Account Center | Accounts and RFIs |
| `client.Banking` | Global Account | Balances, transfers, deposits, beneficiaries, payouts, virtual accounts, conversions, and exchange rates |
| `client.Payment` | Global Acquiring | Payment intents, attempts, refunds, reports, balances, payouts, bank accounts, and terminals |
| `client.Issuing` | Card Issuance | Cards, cardholders, transactions, products, balances, transfers, reports, download center, merchant brands, and authorization decisions |
| `client.Supporting` | Supporting Services | File upload and download links |
| `client.Simulator` | Sandbox | Issuing authorization/reversal and deposit simulation |

Stablecoin Account (Ramp) is not included in the current SDK product scope.

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

## Authorization Decision (PGP)

Authorization decisions let your endpoint approve or decline card transactions in
real time. Before implementing the handler, review the UQPAY authorization decision
[workflow, enablement requirements, response codes, and configured timeout window](https://developers.uqpay.com/card-issuance/v1.6/guide/authorization-decisions).

Generate the required RSA 2048 key pair:

```go
keys, err := authdecision.GenerateKeyPair("Acme Corp", "issuing@acme.example")
if err != nil {
    log.Fatal(err)
}
// Exchange keys.PublicKey with UQPAY and store keys.PrivateKey securely.
```

Configure the SDK and register an HTTP handler. `DecisionTimeout` must be shorter
than the timeout configured with UQPAY; the default UQPAY window is two seconds.

```go
err := client.Issuing.AuthDecision.Configure(authdecision.Config{
    PrivateKey:     os.Getenv("AUTH_DECISION_PRIVATE_KEY"),
    UQPayPublicKey: os.Getenv("UQPAY_PGP_PUBLIC_KEY"),
    Passphrase:     os.Getenv("PGP_PASSPHRASE"),
})
if err != nil {
    log.Fatal(err)
}

handler, err := client.Issuing.AuthDecision.Handler(authdecision.HandlerOptions{
    DecisionTimeout: 1500 * time.Millisecond,
    Decide: func(ctx context.Context, tx authdecision.Transaction) (authdecision.Result, error) {
        if tx.BillingAmount == "10000.00" {
            return authdecision.Result{ResponseCode: "51"}, nil
        }
        return authdecision.Result{
            ResponseCode:       "00",
            PartnerReferenceID: "ref-001",
        }, nil
    },
    OnError: func(err error) {
        log.Printf("authorization decision failed: %v", err)
    },
})
if err != nil {
    log.Fatal(err)
}

http.Handle("/auth-decision", handler)
```

The SDK accepts monetary fields encoded as either JSON strings or numbers and
exposes them as strings to preserve decimal precision. On processing errors the
HTTP handler aborts the response so UQPAY applies the configured timeout action.

## Features

### Automatic Access Token Management

The SDK automatically manages UQPAY Access Tokens:
- Retrieves an Access Token using the Client ID and API key
- Caches the Token until it nears expiry
- Retrieves a new Token before the current one expires
- Thread-safe token management

### Automatic Idempotency Keys

Every API request automatically includes a unique idempotency key to ensure safe retries and prevent duplicate operations.

For Create Virtual Account, supply a stable key explicitly and reuse it only for
the same normalized application request:

```go
application, err := client.Banking.VirtualAccounts.Create(ctx,
    &banking.CreateVirtualAccountRequest{
        Country: "BH", Currency: "USD", PaymentMethod: "SWIFT",
        Nickname: "USD collections",
    },
    &common.RequestOptions{
        IdempotencyKey: "merchant-va-application-42",
        // OnBehalfOf: "connected-account-id",
    },
)

// Applications are distinct from issued Virtual Accounts.
applications, err := client.Banking.VirtualAccountApplications.List(ctx,
    &banking.ListVirtualAccountApplicationsRequest{
        PageNumber: 1, PageSize: 50, Status: "SUBMITTED",
        Country: "BH", Currency: "USD",
    },
)
latest, err := client.Banking.VirtualAccountApplications.Retrieve(
    ctx, application.Data.ApplicationID,
)
issued, err := client.Banking.VirtualAccounts.List(ctx,
    &banking.ListVirtualAccountsRequest{PageNumber: 1, PageSize: 50},
)
```

Create returns HTTP 200 with an accepted application; it does not mean bank
details are ready. Synchronous 400 errors create no application. Asynchronous
method failures are returned in `Results[].Error`.

The webhook parser accepts the application mapping used for
`V1.5.1`, `V1.5.2`, and `V1.6.0`. Handle `virtual.account.create`,
`virtual.account.update`, and `virtual.account.closed`; `SourceID` equals
`ApplicationID`. `AccountID` is the UUID of the account that owns the application.
`DirectID` is a string: `"0"` for a main account, or the connected account's main
account ID. The same required fields are typed on Create and Retrieve application
details and on every List application summary.
Deduplicate by event ID and apply only a higher `PublicVersion`.
Every returned bank detail has `CloseReason`; non-closed records and closed
records without a recorded reason use the empty string.

`ParseVirtualAccountApplicationData` enforces the supported version, both
account-context fields, and the source/application ID match. Unknown older
Virtual Account events remain available through the generic `Event.Data` raw JSON
and are not reclassified as application events.

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
├── auth/              # Access Token management
├── authdecision/      # PGP authorization decision handling
├── banking/           # Banking API client
│   ├── balances.go
│   ├── beneficiaries.go
│   ├── conversion.go
│   ├── deposits.go
│   ├── exchange_rates.go
│   ├── payouts.go
│   ├── transfers.go
│   └── virtual_accounts.go
├── common/            # Shared HTTP transport and request options
├── configuration/     # Sandbox and production environments
├── connect/           # Account Center API client
├── issuing/           # Card Issuance API client
├── payment/           # Global Acquiring API client
├── simulator/         # Sandbox transaction simulator
├── supporting/        # File upload and download links
├── test/              # Integration tests
├── webhook/           # Webhook signature verification
├── uqpay.go            # Root client
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

Install the latest compatible version:

```bash
go get github.com/uqpay/uqpay-sdk-go/v2@latest
```

View all releases: [GitHub Releases](https://github.com/uqpay/uqpay-sdk-go/releases)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

- 📧 Email: support@uqpay.com
- 📚 Documentation: https://developers.uqpay.com
- 🐛 Issues: https://github.com/uqpay/uqpay-sdk-go/issues

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version-specific changes and [GitHub Releases](https://github.com/uqpay/uqpay-sdk-go/releases) for published releases.

---

Made with ❤️ by UQPAY
