package connect

import (
	"encoding/json"
	"testing"
)

func TestRFITextAndResponseAttachments(t *testing.T) {
	b, err := json.Marshal(RFIAnswerItem{Key: "note", Type: "TEXT", Text: "source of funds"})
	if err != nil {
		t.Fatal(err)
	}
	var wire map[string]interface{}
	_ = json.Unmarshal(b, &wire)
	if wire["text"] != "source of funds" {
		t.Fatalf("lost text: %s", b)
	}
	if _, ok := wire["attachments"]; ok {
		t.Fatalf("text emits attachments: %s", b)
	}
	var rfi RFI
	err = json.Unmarshal([]byte(`{"rfi_id":"ACTREQ-test","request":[{"answer":{"key":"document","type":"ATTACHMENT","attachments":[{"file_type":"pdf","file_name":"proof.pdf","size":42,"url":"https://example.test/proof"}]}},{"answer":{"type":"TEXT","text":"source of funds"}}]}`), &rfi)
	if err != nil {
		t.Fatal(err)
	}
	if rfi.Request[0].Answer.Attachments[0].FileName != "proof.pdf" || rfi.Request[1].Answer.Text != "source of funds" {
		t.Fatalf("lost RFI answers: %+v", rfi)
	}
}
