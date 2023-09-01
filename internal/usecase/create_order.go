package usecase

import (
	"github.com/wesleybruno/desafio-clean-arch/internal/entity"
	"github.com/wesleybruno/desafio-clean-arch/pkg/events"
)

type OrderInputDTO struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type CreateOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	ActionEvent     events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewCreateOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	ActionEvent events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: OrderRepository,
		ActionEvent:     ActionEvent,
		EventDispatcher: EventDispatcher,
	}
}

func (c *CreateOrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {
	order := entity.Order{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}
	order.CalculateFinalPrice()
	if err := c.OrderRepository.Save(&order); err != nil {
		return OrderOutputDTO{}, err
	}

	dto := OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.Price + order.Tax,
	}

	c.ActionEvent.SetPayload(dto)
	c.EventDispatcher.Dispatch(c.ActionEvent)

	return dto, nil
}
