package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/streadway/amqp"
	"github.com/wesleybruno/desafio-clean-arch/configs"
	"github.com/wesleybruno/desafio-clean-arch/internal/event"
	"github.com/wesleybruno/desafio-clean-arch/internal/event/handler"
	"github.com/wesleybruno/desafio-clean-arch/internal/infra/database"
	"github.com/wesleybruno/desafio-clean-arch/internal/infra/graph"
	"github.com/wesleybruno/desafio-clean-arch/internal/infra/grpc/pb"
	"github.com/wesleybruno/desafio-clean-arch/internal/infra/grpc/service"
	"github.com/wesleybruno/desafio-clean-arch/internal/infra/web"
	server "github.com/wesleybruno/desafio-clean-arch/internal/infra/web/webserver"
	"github.com/wesleybruno/desafio-clean-arch/pkg/events"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.ActionEventHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	// listOrderUsecase := NewListOrderUseCase(db, eventDispatcher)

	webserver := server.NewWebServer(configs.WebServerPort)

	orderRepository := database.NewOrderRepository(db)
	orderCreatedEvent := event.NewOrderCreatedActionEvent()

	createOrderHandler := web.NewWebCreateOrderHandler(eventDispatcher, orderRepository, orderCreatedEvent)
	listOrderHandler := web.NewWebListOrderHandler(orderRepository)

	webHandler := server.NewHandlerMethod(
		"/order",
		"POST",
		createOrderHandler.Create,
	)

	listOrderwebHandler := server.NewHandlerMethod(
		"/order",
		"GET",
		listOrderHandler.List,
	)

	webserver.AddHandler(*webHandler)
	webserver.AddHandler(*listOrderwebHandler)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
