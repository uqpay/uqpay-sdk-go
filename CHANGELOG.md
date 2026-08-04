# Changelog

All notable changes to the UQPAY Go SDK are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0]

This bootstrap alignment release establishes the shared stable `1.2` capability
baseline used by all five UQPAY customer SDKs. It covers all 98 callable operations
in the current business API contract; Ramp remains outside the SDK product scope.

### Added

- Connect RFI list, retrieve, and answer clients.
- Issuing card limit, risk, PIN, ART, merchant-brand, and unsolicited-refund
  release operations.
- Payment terminal registration and PIN-key operations.
- Simulator clients for deposit and issuing test flows.
- Webhook signature verification using
  `HMAC-SHA512(secret, rawPayload + timestamp)`.

### Compatibility

- Go 1.19 source compatibility and the existing
  `github.com/uqpay/uqpay-sdk-go` module/import path are retained.
- Existing public APIs are not removed or renamed.
