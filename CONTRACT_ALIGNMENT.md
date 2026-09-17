# API contract alignment

Contract reference: [OpenAPI revision 1feb1d2](https://github.com/uqpay/uqpay-docs/tree/1feb1d26d032c53b79ab44a7d48e88c9a91d397d/docs).

## PIN management

`Cards.ResetPIN` calls `/v1/issuing/cards/pin`. Despite historical reset naming in the SDKs, omitting `type` means initial `SET`. Explicitly pass `RESET` to reset without the current PIN, or `UPDATE` with `old_pin` to verify the current PIN before changing it. Both PIN values must contain six digits. `old_pin` is prohibited for `SET` and `RESET`. Validation and supported operations remain server-authoritative.

A response with `request_status=SUCCESS` and `order_status=PROCESSING` means the request was accepted. Use `card_order_id` with the existing card order retrieval operation to obtain `SUCCESS` or `FAILED`; a failed PIN order may include `failure_code`. PIN orders do not contain amount or card currency.

The legacy `/v1/issuing/cards/manage/pin` operation retains its existing behavior. Its `RESET` means changing with the old PIN; when migrating to the new endpoint use `UPDATE`, not `RESET`. Do not send legacy four-digit PINs to the new endpoint.

## RFI answers

Send the full `rfi_id`, including its prefix. A `TEXT` answer needs non-empty `text`; an `ATTACHMENT` answer contains uploaded file IDs in `attachments`.

RFI response answers now use a separate response model. Their `attachments` contain file detail objects (`file_type`, `file_name`, `size`, `url`), not upload file IDs. Request attachments remain string file IDs. Code constructing or accessing typed response answers must use the response model.

## Simulated deposits

Supply `account_id` explicitly, together with `amount`, `currency` and `sender_swift_code`. The recipient must be active and verified. The SDK does not infer the recipient from credentials or request headers. This operation is Sandbox-only.

## Transaction detail

`settlement_status` is available on transaction detail, and may be absent from list items. Values are `UNKNOWN`, `UNSETTLED`, `SETTLED` and `NOT_APPLICABLE`. `SETTLED` means clearing has been recorded, including partial clearing; it does not confirm full settlement. `UNKNOWN` does not mean unsettled.

## Beneficiary checks

Provide at least one non-empty `account_number` or `iban`. Both are accepted; `account_number` takes precedence. For LOCAL currencies without a default clearing system (including CNH), supply both `bank_country_code` and `clearing_system`. The default-route currencies are USD, GBP, EUR, SGD, CAD, AUD, HKD, MYR, IDR, PHP and INR. For SWIFT, use `clearing_system=SWIFT`. Validation is server-authoritative.

## Card art updates

Card updates accept `card_art_id` and `name_on_card`. Card art changes apply to virtual and physical cards; both `card_status` and `processing_status` must be `ACTIVE`. An accepted `PROCESSING` response is asynchronous: use `card_order_id` to check the final result.

## Typed webhook and deposit fields

Payment method parsing covers all 26 contract types, card-present details, card name/number/network and static QR fields. The existing `AlipayDetails` type is retained for non-card details. Payment intent events expose `NextAction`; issuing transaction events expose open-string `WalletType`, and `Cardholder` exposes `Reason`. Deposit responses expose deposit method and sender classifications.

String amounts remain strings, including negative/high-precision values. Existing scalar string fields accept JSON null as an empty Go string; use the original `Event.Data` when absence, null and empty string must be distinguished. Optional timestamp pointers and object pointers retain their existing semantics.

## Boundary and response regression coverage

Offline fixtures exercise all three KYC entry points (create cardholder, update cardholder and create card with inline cardholder fields). They cover all six proof providers, reference lengths 9/10/64/65 and birth dates corresponding to ages 17/18/79/80 on 2026-09-17. These tests prove that the SDK preserves the submitted values; server acceptance, current configuration and environment rollout still require Sandbox verification. Inline cardholder fields include email, first/last name and country code.

Card list pagination preserves page sizes 1, 10 and 100. Banking list page sizes are 1–100. Currency conversions require a fresh quote per operation: `FUNDS_ARRIVED` is intermediate and `TRADE_SETTLED` is the successful terminal state.

Payment intent retrieval supports per-request `x-on-behalf-of` without requiring a caller-supplied idempotency key. Existing automatic GET header behavior remains language-specific; the contract no longer requires the header but does not prohibit it. POST idempotency and retry behavior remain unchanged.

Card creation orders use `CREATE_CARD`. Issuing transfer REST status uses uppercase `PENDING`/`FAILED`/`COMPLETED`; transfer webhook status is a distinct lowercase field. Wallet values remain open strings, including empty and unknown values. The SDK does not infer transaction context from them.

Account list/detail responses expose current company and individual fields. Optional response maps avoid applying company-create requiredness to list summaries. Card list JSON-text metadata is decoded into the existing `FlexibleStringMap`; detail object/null metadata remains supported.

**Migration:** simulated authorization response amounts (`CardAvailableBalance`, `BillingAmount`, `TransactionAmount`) now use `common.FlexibleString` instead of `float64`. Use `.String()` or an exact decimal library; request amounts remain numeric. Both legacy JSON numbers and current decimal strings decode without a float conversion.

`NetworkProtectionFeeData` and `IssuingTransferStatusChangedData` expose string amounts. Representative webhook models accept null or array `other_documents`. `CardDetails.IssuerCountryCode` is optional and is based on Sandbox commit `5450a0a9` (gateway head `fb887a5d9261ed57122ac2583c02d9651595fa3e`); production support is not assumed.

## Acquiring GET headers

Offline route-specific checks cover D189–D196: balance list/detail, bank account list/detail, payout list/detail, settlements list and payment intent detail. Each call is exercised with and without `x-on-behalf-of`, while preserving the configured `x-client-id`. Callers need not supply an idempotency key for these GET requests; existing automatic GET headers remain supported. Existing POST idempotency and retry tests remain part of verification.

Balance and payout `Get`/`List`, and `Reports.ListSettlements`, now accept optional `*common.RequestOptions`. Existing ordinary calls remain valid. Code assigning these methods to exact function types or implementing interfaces may need to update signatures for the variadic options. Bank account reads now preserve the configured client ID when options are omitted or only specify delegation.

## Banking balance precision

Offline D122–D129 fixtures cover `available_balance`, `frozen_balance`, `margin_balance` and `prepaid_balance` through both list and detail reads. Each field receives zero with trailing decimal places, a positive amount, a negative amount, large positive/negative amounts, and a long decimal string. Distinct values across fields detect accidental swaps; comparisons retain strings and precision. Long decimal fixtures test client robustness, not server-supported currency precision. CLI verification covers JSON output.

## Acquiring nullable responses

Offline fixtures distinguish missing fields, explicit null, empty strings/objects and populated controls. REST checks cover payment attempts, refunds and payouts. Empty REST event times remain strings. Webhook fixtures exercise intent, attempt, refund, payout and chargeback alert families; empty timestamps or incomplete objects are robustness probes, not claims of server-valid payloads.

Signed Webhook fixtures preserve raw data and reject a payload whose bytes change after signing. This is offline verification, not evidence of event delivery from Sandbox.

`PaymentAttempt` exposes `AdviceCode` and `AuthenticationData`. Scalar response strings merge missing/null with an empty value; optional timestamp pointers merge missing/null but retain empty/non-empty strings. Use `Event.Data` to distinguish raw presence and null. Acquiring payout and chargeback events use this raw interface; no new typed helper is claimed.
