package web

import (
	"encoding/json"
	"net/http"

	"github.com/wesleybruno/desafio-clean-arch/internal/entity"
	"github.com/wesleybruno/desafio-clean-arch/internal/usecase"
	"github.com/wesleybruno/desafio-clean-arch/pkg/events"
)

type WebListOrderHandler struct {
	EventDispatcher events.EventDispatcherInterface
	OrderRepository entity.OrderRepositoryInterface
}

func NewWebListOrderHandler(
	EventDispatcher events.EventDispatcherInterface,
	OrderRepository entity.OrderRepositoryInterface,

) *WebListOrderHandler {
	return &WebListOrderHandler{
		EventDispatcher: EventDispatcher,
		OrderRepository: OrderRepository,
	}
}

func (h *WebListOrderHandler) List(w http.ResponseWriter, r *http.Request) {

	listOrderUseCase := usecase.NewListOrderUseCase(h.OrderRepository, h.EventDispatcher)
	output, err := listOrderUseCase.Execute()
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
