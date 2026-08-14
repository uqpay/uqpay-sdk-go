package webhook

import (
	"encoding/json"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/v2/banking"
)

const (
	EventNameVirtual              = "VIRTUAL"
	EventTypeVirtualAccountCreate = "virtual.account.create"
	EventTypeVirtualAccountUpdate = "virtual.account.update"
	EventTypeVirtualAccountClosed = "virtual.account.closed"
)

var virtualAccountApplicationVersions = map[string]struct{}{
	"V1.5.1": {},
	"V1.5.2": {},
	"V1.6.0": {},
}

// VirtualAccountApplicationData is the webhook-specific VA application shape
// added in the current SDK scope.
// AccountID is the UUID of the account that owns the application. DirectID is a
// string: "0" for a main account, or the connected account's main account ID.
// REST public SDK types remain unchanged pending a confirmed Developer Docs
// contract for these fields.
type VirtualAccountApplicationData struct {
	banking.VirtualAccountApplication
	AccountID string `json:"account_id"`
	DirectID  string `json:"direct_id"`
}

// ParseVirtualAccountApplicationData supports application events for webhook
// versions V1.5.1, V1.5.2, and V1.6.0. SourceID equals ApplicationID; callers
// should order changes by ApplicationID and PublicVersion.
func (e *Event) ParseVirtualAccountApplicationData() (*VirtualAccountApplicationData, error) {
	if e.EventType != EventTypeVirtualAccountCreate && e.EventType != EventTypeVirtualAccountUpdate && e.EventType != EventTypeVirtualAccountClosed {
		return nil, fmt.Errorf("event type %s is not a virtual account application event", e.EventType)
	}
	if _, ok := virtualAccountApplicationVersions[e.Version]; !ok {
		return nil, fmt.Errorf("webhook version %s does not use the virtual account application contract", e.Version)
	}
	var data VirtualAccountApplicationData
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to parse virtual account application data: %w", err)
	}
	if data.AccountID == "" || data.DirectID == "" {
		return nil, fmt.Errorf("virtual account application data requires account_id and direct_id")
	}
	if data.ApplicationID == "" || e.SourceID != data.ApplicationID {
		return nil, fmt.Errorf("virtual account application source_id must equal application_id")
	}
	return &data, nil
}
