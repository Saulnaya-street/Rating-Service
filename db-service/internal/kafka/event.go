package kafka

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	BookCreated EventType = "book.created"
	BookUpdated EventType = "book.updated"
	BookDeleted EventType = "book.deleted"
)

// Event базовая структура событий
type Event struct {
	ID        string      `json:"id"`
	Type      EventType   `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

type BookEvent struct {
	BookID uuid.UUID `json:"book_id"`
	Name   string    `json:"name,omitempty"`
	Author string    `json:"author,omitempty"`
	Year   int       `json:"year,omitempty"`
	Genre  string    `json:"genre,omitempty"`
}

func DeserializeEvent(data []byte) (*Event, error) {
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("error deserializing event: %w", err)
	}
	return &event, nil
}

func DeserializeBookPayload(event *Event) (*BookEvent, error) {
	payloadJSON, err := json.Marshal(event.Payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing payload: %w", err)
	}

	var bookPayload BookEvent
	if err := json.Unmarshal(payloadJSON, &bookPayload); err != nil {
		return nil, fmt.Errorf("error deserializing book payload: %w", err)
	}

	return &bookPayload, nil
}
