package usecase

import (
	"github.com/wesleybruno/desafio-clean-arch/internal/entity"
)

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *ListOrderUseCase) Execute() ([]OrderOutputDTO, error) {

	orders, err := c.OrderRepository.ListAll()
	if err != nil {
		return nil, err
	}

	var ordersDTO []OrderOutputDTO

	for _, o := range orders {

		dto := OrderOutputDTO{
			ID:         o.ID,
			Price:      o.Price,
			Tax:        o.Tax,
			FinalPrice: o.FinalPrice,
		}

		ordersDTO = append(ordersDTO, dto)

	}

	return ordersDTO, nil
}
