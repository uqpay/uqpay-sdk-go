# API contract alignment

Contract reference: [OpenAPI revision 8267056](https://github.com/uqpay/uqpay-docs/tree/8267056f7fefc1183b069e0387e5ebaecbd2ad27/docs).

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
