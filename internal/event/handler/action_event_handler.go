package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/streadway/amqp"
	"github.com/wesleybruno/desafio-clean-arch/pkg/events"
)

type ActionEventHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewActionEventHandler(rabbitMQChannel *amqp.Channel) *ActionEventHandler {
	return &ActionEventHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *ActionEventHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("%v: %v", event.GetName(), event.GetPayload())
	jsonOutput, _ := json.Marshal(event.GetPayload())

	msgRabbitmq := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	h.RabbitMQChannel.Publish(
		"amq.direct", // exchange
		"",           // key name
		false,        // mandatory
		false,        // immediate
		msgRabbitmq,  // message to publish
	)
}
