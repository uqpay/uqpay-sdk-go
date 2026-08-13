# Changelog

All notable changes to the UQPAY Go SDK are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Versions follow UQPAY's shared `MAJOR.MINOR` and repository-specific `PATCH`
policy, which is based on [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Virtual Account Application list and retrieve clients with precise summary,
  full application, result error, bank-detail, and clearing-system models.
- Typed `virtual.account.create`, `virtual.account.update`, and
  `virtual.account.closed` parsing for webhook mappings on `V1.5.1`, `V1.5.2`,
  and `V1.6.0`.

### Changed

- Webhook freshness validation accepts Webhook Hub's Unix-millisecond
  `x-wk-timestamp` while retaining Unix-second compatibility and signing the
  unmodified header value.
- Create Virtual Account now requires a country and one currency, accepts an
  optional `LOCAL` or `SWIFT` method and nickname, and returns the asynchronous
  application response. Existing request options continue to forward
  `IdempotencyKey` as `x-idempotency-key` and `OnBehalfOf` as
  `x-on-behalf-of`.

## [1.2.1]

This Go-only release is an explicitly approved, one-time parity-repair exception.
It restores a capability already present in the Node.js, Python, and Java SDKs but
omitted from the Go SDK's `1.2.0` bootstrap alignment. It does not establish a
general policy of adding public SDK capabilities in PATCH releases.

### Added

- PGP-encrypted real-time authorization decisions through
  `client.Issuing.AuthDecision`.
- RSA-2048 key generation, armored strings and key-file configuration, and
  passphrase-protected customer private keys.
- Typed authorization transactions, string-preserving monetary values, and
  encrypted approve/decline responses.
- Configurable decision timeouts, request-size limits, and fail-closed HTTP error
  handling so UQPAY can apply the configured timeout action.

### Compatibility and security

- Go 1.19 source compatibility is retained.
- The authorization-decision path accepts only RSA keys of at least 2048 bits.
- `govulncheck` reports conditional findings in the indirect CIRCL dependency for
  P-384 (`GO-2026-4550`) and FourQ (`GO-2025-3754`). Those algorithms are not
  accepted by this RSA-only path; CIRCL versions containing the upstream fixes
  require Go 1.22 or later.

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
