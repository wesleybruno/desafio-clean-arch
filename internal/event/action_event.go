package event

import "time"

type ActionEvent struct {
	Name    string
	Payload interface{}
}

func NewOrderCreatedActionEvent() *ActionEvent {
	return &ActionEvent{
		Name: "OrderCreated",
	}
}

func NewActionEvent(event string) *ActionEvent {
	return &ActionEvent{
		Name: event,
	}
}

func (e *ActionEvent) GetName() string {
	return e.Name
}

func (e *ActionEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *ActionEvent) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *ActionEvent) GetDateTime() time.Time {
	return time.Now()
}
