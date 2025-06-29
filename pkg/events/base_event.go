package events

import "time"

type BaseEvent struct {
	Name     string
	DateTime time.Time
	Payload  interface{}
}

func NewBaseEvent(name string, payload interface{}) EventInterface {
	return &BaseEvent{
		Name:     name,
		DateTime: time.Now(),
		Payload:  payload,
	}
}

func (e *BaseEvent) GetName() string {
	return e.Name
}

func (e *BaseEvent) GetDateTime() time.Time {
	return e.DateTime
}

func (e *BaseEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *BaseEvent) SetPayload(payload interface{}) {
	e.Payload = payload
}
