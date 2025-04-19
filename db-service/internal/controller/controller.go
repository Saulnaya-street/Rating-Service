package controller

import (
	"awesomeProject/db-service/internal/delivery/handler"
	"awesomeProject/db-service/internal/kafka"
	"awesomeProject/db-service/internal/repository"
	"awesomeProject/db-service/internal/service"
	"awesomeProject/db-service/pkg"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type ControllerOptions struct {
	DB          *pgxpool.Pool
	KafkaClient kafka.IKafkaClient
}

type Controller struct {
	db            *pgxpool.Pool
	kafkaClient   kafka.IKafkaClient
	ratingService service.IRatingService
	server        *pkg.Server
	ratingHandler handler.IRatingHandler
	eventConsumer *kafka.EventConsumer
}

func NewController(opts ControllerOptions) *Controller {
	ratingRepo := repository.NewRatingRepository(opts.DB)

	ratingService := service.NewRatingService(ratingRepo)

	eventConsumer := kafka.NewEventConsumer(opts.KafkaClient, ratingService)

	ratingHandler := handler.NewRatingHandler(ratingService)

	server := pkg.NewServer()

	deliveryRouter := handler.NewRouter(ratingHandler)
	deliveryRouter.RegisterRoutes(server.GetRouter())

	return &Controller{
		db:            opts.DB,
		kafkaClient:   opts.KafkaClient,
		ratingService: ratingService,
		server:        server,
		ratingHandler: ratingHandler,
		eventConsumer: eventConsumer,
	}
}

func (c *Controller) StartEventConsumer() {
	c.eventConsumer.StartConsuming()
}

func (c *Controller) GetServer() *pkg.Server {
	return c.server
}

func (c *Controller) CloseConnections() {
	if c.kafkaClient != nil {
		if err := c.kafkaClient.Close(); err != nil {
			log.Printf("Error closing connection to Kafka: %v", err)
		}
	}

	if c.db != nil {
		c.db.Close()
	}

	log.Println("All connections are closed")
}
