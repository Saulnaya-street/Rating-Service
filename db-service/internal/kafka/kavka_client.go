package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type IKafkaClient interface {
	Close() error
	Subscribe(handler func(ctx context.Context, msg []byte) error)
}

type KafkaClient struct {
	reader *kafka.Reader
	config KafkaConfig
}

func NewKafkaClient(cfg KafkaConfig) (IKafkaClient, error) {
	client := &KafkaClient{
		config: cfg,
	}

	client.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          cfg.Topic,
		GroupID:        cfg.GroupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		MaxWait:        time.Second,
		StartOffset:    kafka.FirstOffset,
		CommitInterval: time.Second,
	})

	log.Printf("Configured Kafka consumer for topic %s with group %s", cfg.Topic, cfg.GroupID)

	// Проверка соединения
	log.Printf("Checking connection to Kafka broker %s...", cfg.Brokers[0])
	conn, err := kafka.DialLeader(context.Background(), "tcp", cfg.Brokers[0], cfg.Topic, 0)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Kafka: %w", err)
	}
	conn.Close()
	log.Printf("Connection to Kafka established successfully")

	return client, nil
}

func (k *KafkaClient) Close() error {
	if k.reader != nil {
		log.Println("Kafka Consumer Shutdown...")
		err := k.reader.Close()
		if err != nil {
			return fmt.Errorf("consumer closing error: %w", err)
		}
		log.Println("Kafka consumer successfully closed")
	}

	return nil
}

func (k *KafkaClient) Subscribe(handler func(ctx context.Context, msg []byte) error) {
	go func() {
		for {
			ctx := context.Background()
			msg, err := k.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message from Kafka: %v", err)
				time.Sleep(time.Second)
				continue
			}

			log.Printf("Received message from topic %s with key: %s", msg.Topic, string(msg.Key))

			if err := handler(ctx, msg.Value); err != nil {
				log.Printf("Error processing message: %v", err)
				continue
			}

		}
	}()
}
