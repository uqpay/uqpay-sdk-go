package issuing

import (
	"encoding/json"
	"testing"
)

func TestCardMetadataJSONText(t *testing.T) {
	var card RetrieveCardResponse
	if err := json.Unmarshal([]byte(`{"card_limit":"12345678901234567890.12345678","metadata":"{\"reference\":\"0001\"}","risk_controls":null}`), &card); err != nil {
		t.Fatal(err)
	}
	if card.Metadata["reference"] != "0001" {
		t.Fatalf("list metadata discarded: %#v", card.Metadata)
	}
	if string(card.CardLimit) != "12345678901234567890.12345678" {
		t.Fatal(card.CardLimit)
	}
}

func TestTransactionAndTransferStrings(t *testing.T) {
	var tx Transaction
	if err := json.Unmarshal([]byte(`{"transaction_amount":"12345678901234567890.12345678","billing_amount":"0.00","transaction_fee":"0.01","card_available_balance":"-0.01","original_transaction_id":"","short_transaction_id":"not-a-uuid","wallet_type":"FUTURE_WALLET","merchant_data":{}}`), &tx); err != nil {
		t.Fatal(err)
	}
	if tx.TransactionAmount != "12345678901234567890.12345678" || tx.BillingAmount != "0.00" || tx.CardAvailableBalance != "-0.01" || tx.OriginalTransactionID != "" || *tx.WalletType != "FUTURE_WALLET" {
		t.Fatal(tx)
	}
	var transfer Transfer
	if err := json.Unmarshal([]byte(`{"amount":"0.01","fee_amount":"0.00","creator_id":"","reference_id":"transfer-ref","transfer_status":"COMPLETED"}`), &transfer); err != nil {
		t.Fatal(err)
	}
	if transfer.Amount != "0.01" || transfer.FeeAmount != "0.00" || transfer.TransferStatus != "COMPLETED" {
		t.Fatal(transfer)
	}
}
