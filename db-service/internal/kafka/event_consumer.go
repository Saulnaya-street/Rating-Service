package kafka

import (
	"awesomeProject/db-service/internal/service"
	"context"
	"fmt"
	"log"
)

type EventConsumer struct {
	kafkaClient   IKafkaClient
	ratingService service.IRatingService
}

func NewEventConsumer(client IKafkaClient, ratingService service.IRatingService) *EventConsumer {
	return &EventConsumer{
		kafkaClient:   client,
		ratingService: ratingService,
	}
}

func (ec *EventConsumer) StartConsuming() {
	ec.kafkaClient.Subscribe(ec.processMessage)
	log.Println("Event consumer started")
}

func (ec *EventConsumer) processMessage(ctx context.Context, msgData []byte) error {
	event, err := DeserializeEvent(msgData)
	if err != nil {
		return fmt.Errorf("failed to deserialize event: %w", err)
	}

	log.Printf("Processing event: %s, id: %s", event.Type, event.ID)

	switch event.Type {
	case BookCreated:
		return ec.handleBookCreatedEvent(ctx, event)
	case BookUpdated:

		log.Printf("Book update event received, no action required")
		return nil
	case BookDeleted:

		log.Printf("Book delete event received, no action required")
		return nil
	default:
		log.Printf("Unknown event type: %s", event.Type)
		return nil
	}
}

func (ec *EventConsumer) handleBookCreatedEvent(ctx context.Context, event *Event) error {
	bookEvent, err := DeserializeBookPayload(event)
	if err != nil {
		return fmt.Errorf("failed to deserialize book created payload: %w", err)
	}

	log.Printf("Book created: %s", bookEvent.BookID)

	err = ec.ratingService.HandleNewBook(ctx, bookEvent.BookID)
	if err != nil {
		return fmt.Errorf("failed to handle new book event: %w", err)
	}

	log.Printf("Successfully processed book created event for book: %s", bookEvent.BookID)
	return nil
}
