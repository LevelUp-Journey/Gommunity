package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	BootstrapServers string
	GroupID          string
	Topics           []string
}

type KafkaConsumer struct {
	readers []*kafka.Reader
}

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer(config KafkaConfig) *KafkaConsumer {
	var readers []*kafka.Reader

	// Create a reader for each topic
	for _, topic := range config.Topics {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{config.BootstrapServers},
			GroupID:        config.GroupID,
			Topic:          topic,
			MinBytes:       10e3, // 10KB
			MaxBytes:       10e6, // 10MB
			CommitInterval: time.Second,
			StartOffset:    kafka.LastOffset,
		})
		readers = append(readers, reader)
	}

	log.Printf("Kafka consumer created for topics: %v with group ID: %s", config.Topics, config.GroupID)

	return &KafkaConsumer{
		readers: readers,
	}
}

// ConsumeMessages starts consuming messages from Kafka
func (kc *KafkaConsumer) ConsumeMessages(ctx context.Context, handler func(topic string, message []byte) error) error {
	log.Println("Starting Kafka message consumption...")

	// Start a goroutine for each reader
	for _, reader := range kc.readers {
		go func(r *kafka.Reader) {
			for {
				select {
				case <-ctx.Done():
					log.Println("Stopping Kafka consumer for topic...")
					return
				default:
					msg, err := r.FetchMessage(ctx)
					if err != nil {
						if ctx.Err() != nil {
							return
						}
						log.Printf("Error fetching message: %v", err)
						continue
					}

					log.Printf("Received message from topic %s: partition=%d offset=%d",
						msg.Topic, msg.Partition, msg.Offset)

					if err := handler(msg.Topic, msg.Value); err != nil {
						log.Printf("Error handling message: %v", err)
						// Continue processing other messages even if one fails
						continue
					}

					// Commit the message after successful processing
					if err := r.CommitMessages(ctx, msg); err != nil {
						log.Printf("Error committing message: %v", err)
					}
				}
			}
		}(reader)
	}

	// Wait for context cancellation
	<-ctx.Done()
	log.Println("Stopping Kafka consumer...")
	return kc.Close()
}

// Close closes the Kafka consumer
func (kc *KafkaConsumer) Close() error {
	for _, reader := range kc.readers {
		if err := reader.Close(); err != nil {
			log.Printf("Error closing Kafka reader: %v", err)
		}
	}
	log.Println("Kafka consumer closed")
	return nil
}
