package kafka

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type KafkaConfig struct {
	BootstrapServers string
	GroupID          string
	Topics           []string
	SecurityProtocol string
	SASLMechanism    string
	SASLUsername     string
	SASLPassword     string
}

type KafkaConsumer struct {
	readers []*kafka.Reader
}

// NewKafkaConsumer creates a new Kafka consumer with Azure Event Hub support
func NewKafkaConsumer(config KafkaConfig) *KafkaConsumer {
	var readers []*kafka.Reader

	// Configure SASL mechanism for Azure Event Hub
	var mechanism sasl.Mechanism
	var dialer *kafka.Dialer

	if config.SecurityProtocol == "SASL_SSL" {
		log.Println("========================================")
		log.Println("Configuring Kafka consumer for Azure Event Hub...")
		log.Printf("Bootstrap Servers: %s", config.BootstrapServers)
		log.Printf("SASL Mechanism: %s", config.SASLMechanism)
		log.Printf("SASL Username: %s", config.SASLUsername)

		// Validate configuration for Azure Event Hub
		validationErrors := []string{}

		// Mask password for security but show if it's set
		if config.SASLPassword != "" {
			log.Printf("SASL Password: [SET - length: %d characters]", len(config.SASLPassword))
			// Check if password starts with "Endpoint=sb://" (valid Azure Event Hub connection string)
			if len(config.SASLPassword) > 12 && config.SASLPassword[:12] == "Endpoint=sb:" {
				log.Println("✓ Password appears to be a valid Azure Event Hub connection string")
			} else {
				validationErrors = append(validationErrors, "Password does NOT appear to be a valid Azure Event Hub connection string (should start with: Endpoint=sb://)")
			}
		} else {
			validationErrors = append(validationErrors, "SASL Password is EMPTY!")
		}

		// Validate username
		if config.SASLUsername != "$ConnectionString" {
			validationErrors = append(validationErrors, "SASL Username must be exactly '$ConnectionString' for Azure Event Hub")
		}

		// Validate bootstrap servers
		if config.BootstrapServers == "" || config.BootstrapServers == "localhost:9092" {
			validationErrors = append(validationErrors, "Bootstrap servers not properly configured for Azure Event Hub")
		}

		// Display validation results
		if len(validationErrors) > 0 {
			log.Println("⚠️  CONFIGURATION VALIDATION ERRORS:")
			for i, err := range validationErrors {
				log.Printf("   %d. %s", i+1, err)
			}
			log.Println("⚠️  Azure Event Hub connection will likely FAIL with these errors!")
			log.Println("⚠️  Please check your .env file configuration.")
		} else {
			log.Println("✓ All Azure Event Hub configuration validations passed")
		}

		log.Println("Required Event Hubs (topics) in Azure:")
		for i, topic := range config.Topics {
			log.Printf("   %d. %s", i+1, topic)
		}
		log.Println("⚠️  Make sure these Event Hubs exist in your Azure Event Hubs Namespace!")
		log.Println("========================================")

		if config.SASLMechanism == "PLAIN" {
			mechanism = plain.Mechanism{
				Username: config.SASLUsername,
				Password: config.SASLPassword,
			}
		}

		// Configure TLS for Azure Event Hub
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}

		dialer = &kafka.Dialer{
			Timeout:       30 * time.Second, // Increased timeout for Azure
			DualStack:     true,
			TLS:           tlsConfig,
			SASLMechanism: mechanism,
		}
	} else {
		// Default dialer without TLS/SASL
		dialer = &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
		}
	}

	// Create a reader for each topic
	for _, topic := range config.Topics {
		readerConfig := kafka.ReaderConfig{
			Brokers:        []string{config.BootstrapServers},
			GroupID:        config.GroupID,
			Topic:          topic,
			MinBytes:       10e3, // 10KB
			MaxBytes:       10e6, // 10MB
			CommitInterval: time.Second,
			StartOffset:    kafka.LastOffset,
			Dialer:         dialer,
		}

		reader := kafka.NewReader(readerConfig)
		readers = append(readers, reader)
	}

	log.Printf("Kafka consumer created for topics: %v with group ID: %s (Security: %s)",
		config.Topics, config.GroupID, config.SecurityProtocol)

	return &KafkaConsumer{
		readers: readers,
	}
}

// ConsumeMessages starts consuming messages from Kafka with retry logic
func (kc *KafkaConsumer) ConsumeMessages(ctx context.Context, handler func(topic string, message []byte) error) error {
	log.Println("Starting Kafka message consumption...")

	// Start a goroutine for each reader with retry logic
	for _, reader := range kc.readers {
		go func(r *kafka.Reader) {
			retryDelay := 10 * time.Second
			maxRetryDelay := 5 * time.Minute
			consecutiveErrors := 0

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

						consecutiveErrors++

						// Enhanced error logging for Azure Event Hub
						log.Printf("❌ Error fetching message (attempt #%d): %v", consecutiveErrors, err)

						// Provide helpful troubleshooting tips
						if consecutiveErrors == 1 {
							log.Println("💡 Troubleshooting tips for Azure Event Hub:")
							log.Println("   1. Verify the Event Hub (topic) exists in Azure Portal")
							log.Println("   2. Check that your Connection String is correct and not expired")
							log.Println("   3. Ensure your Shared Access Policy has 'Listen' permission")
							log.Println("   4. Confirm the username is exactly: $ConnectionString")
							log.Println("   5. Verify your Consumer Group exists (default: $Default)")
						}

						// Exponential backoff with jitter
						if consecutiveErrors > 3 {
							if retryDelay < maxRetryDelay {
								retryDelay = retryDelay * 2
								if retryDelay > maxRetryDelay {
									retryDelay = maxRetryDelay
								}
							}
							log.Printf("⏳ Waiting %v before retry...", retryDelay)
							time.Sleep(retryDelay)
						} else {
							time.Sleep(5 * time.Second)
						}
						continue
					}

					// Reset error counter on successful fetch
					if consecutiveErrors > 0 {
						log.Printf("✓ Successfully reconnected to Azure Event Hub after %d errors", consecutiveErrors)
						consecutiveErrors = 0
						retryDelay = 10 * time.Second
					}

					log.Printf("📨 Received message from topic %s: partition=%d offset=%d",
						msg.Topic, msg.Partition, msg.Offset)

					if err := handler(msg.Topic, msg.Value); err != nil {
						log.Printf("❌ Error handling message: %v", err)
						// Continue processing other messages even if one fails
						continue
					}

					// Commit the message after successful processing
					if err := r.CommitMessages(ctx, msg); err != nil {
						log.Printf("⚠️  Error committing message: %v", err)
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
