package web

import (
	"encoding/json"
	"net/http"

	"github.com/wesleybruno/desafio-clean-arch/internal/entity"
	"github.com/wesleybruno/desafio-clean-arch/internal/usecase"
	"github.com/wesleybruno/desafio-clean-arch/pkg/events"
)

type WebCreateOrderHandler struct {
	EventDispatcher events.EventDispatcherInterface
	OrderRepository entity.OrderRepositoryInterface
	ActionEvent     events.EventInterface
}

func NewWebCreateOrderHandler(
	EventDispatcher events.EventDispatcherInterface,
	OrderRepository entity.OrderRepositoryInterface,
	ActionEvent events.EventInterface,
) *WebCreateOrderHandler {
	return &WebCreateOrderHandler{
		EventDispatcher: EventDispatcher,
		OrderRepository: OrderRepository,
		ActionEvent:     ActionEvent,
	}
}

func (h *WebCreateOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository, h.ActionEvent, h.EventDispatcher)
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
