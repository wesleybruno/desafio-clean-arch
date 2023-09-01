//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"
	"github.com/google/wire"
	"github.com/wesleybruno/desafio-clean-arch/internal/entity"
	"github.com/wesleybruno/desafio-clean-arch/internal/event"
	"github.com/wesleybruno/desafio-clean-arch/internal/infra/database"
	"github.com/wesleybruno/desafio-clean-arch/internal/usecase"
	"github.com/wesleybruno/desafio-clean-arch/pkg/events"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setActionEvent = wire.NewSet(
	event.NewOrderCreatedActionEvent,
	wire.Bind(new(events.EventInterface), new(*event.ActionEvent)),
)

func NewListOrderUseCase(db *sql.DB) *usecase.ListOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		usecase.NewListOrderUseCase,
	)
	return &usecase.ListOrderUseCase{}
}

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setActionEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}
