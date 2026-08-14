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

// ParseVirtualAccountApplicationData supports application events for webhook
// versions V1.5.1, V1.5.2, and V1.6.0. SourceID equals ApplicationID; callers
// should order changes by ApplicationID and PublicVersion.
func (e *Event) ParseVirtualAccountApplicationData() (*banking.VirtualAccountApplication, error) {
	if e.EventType != EventTypeVirtualAccountCreate && e.EventType != EventTypeVirtualAccountUpdate && e.EventType != EventTypeVirtualAccountClosed {
		return nil, fmt.Errorf("event type %s is not a virtual account application event", e.EventType)
	}
	var data banking.VirtualAccountApplication
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to parse virtual account application data: %w", err)
	}
	return &data, nil
}
