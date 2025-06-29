// internal/event/event/comboname_created_event.go
package event

import (
	"catalog/pkg/events"
	"time"
)

type ComboNameCreated struct {
	Name     string
	ID       int64
	DateTime time.Time
}

func NewComboNameCreatedEvent() events.EventInterface {
	return events.NewBaseEvent("ComboNameCreated", nil)
}
