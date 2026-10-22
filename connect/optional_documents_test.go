package connect

import (
	"encoding/json"
	"testing"
)

func TestCompanyOptionalDocumentsOmittedAndPreserved(t *testing.T) {
	for _, docs := range [][]string{nil, {"file-synthetic"}} {
		body, err := json.Marshal(SubAccountOwnershipDetails{Representatives: []SubAccountRepresentative{{FaceDocs: docs}}, ShareholderDocs: docs})
		if err != nil {
			t.Fatal(err)
		}
		var got map[string]json.RawMessage
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatal(err)
		}
		_, present := got["shareholder_docs"]
		if present != (docs != nil) {
			t.Fatalf("unexpected shareholder_docs: %s", body)
		}
		var reps []map[string]json.RawMessage
		if err := json.Unmarshal(got["representatives"], &reps); err != nil {
			t.Fatal(err)
		}
		_, present = reps[0]["face_docs"]
		if present != (docs != nil) {
			t.Fatalf("unexpected face_docs: %s", body)
		}
		if docs != nil && string(got["shareholder_docs"]) != `["file-synthetic"]` {
			t.Fatalf("documents changed: %s", body)
		}
	}
}
