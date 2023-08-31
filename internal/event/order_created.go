package event

import "time"

type ActionEvent struct {
	Name    string
	Payload interface{}
}

func NewOrderCreated() *ActionEvent {
	return &ActionEvent{
		Name: "OrderCreated",
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
