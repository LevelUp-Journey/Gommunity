package messaging

import (
	"context"
	"encoding/json"
	"log"

	"Gommunity/platform/users/application/eventhandlers"
	"Gommunity/platform/users/domain/model/events"
)

const (
	TopicCommunityRegistration = "community.registration"
	TopicProfileUpdated        = "community.profile.updated"
)

type KafkaEventConsumer struct {
	registrationHandler  *eventhandlers.UserRegistrationHandler
	profileUpdateHandler *eventhandlers.ProfileUpdatedHandler
}

func NewKafkaEventConsumer(
	registrationHandler *eventhandlers.UserRegistrationHandler,
	profileUpdateHandler *eventhandlers.ProfileUpdatedHandler,
) *KafkaEventConsumer {
	return &KafkaEventConsumer{
		registrationHandler:  registrationHandler,
		profileUpdateHandler: profileUpdateHandler,
	}
}

// HandleMessage routes messages to appropriate handlers based on topic
func (kec *KafkaEventConsumer) HandleMessage(topic string, message []byte) error {
	ctx := context.Background()
	log.Printf("Handling message from topic: %s", topic)

	switch topic {
	case TopicCommunityRegistration:
		return kec.handleRegistrationEvent(ctx, message)
	case TopicProfileUpdated:
		return kec.handleProfileUpdatedEvent(ctx, message)
	default:
		log.Printf("Unknown topic: %s", topic)
		return nil
	}
}

func (kec *KafkaEventConsumer) handleRegistrationEvent(ctx context.Context, message []byte) error {
	var event events.CommunityRegistrationEvent
	if err := json.Unmarshal(message, &event); err != nil {
		log.Printf("Error unmarshalling registration event: %v", err)
		return err
	}

	log.Printf("Processing registration event: UserID=%s, Username=%s", event.UserID, event.Username)
	return kec.registrationHandler.Handle(ctx, event)
}

func (kec *KafkaEventConsumer) handleProfileUpdatedEvent(ctx context.Context, message []byte) error {
	var event events.ProfileUpdatedEvent
	if err := json.Unmarshal(message, &event); err != nil {
		log.Printf("Error unmarshalling profile updated event: %v", err)
		return err
	}

	log.Printf("Processing profile updated event: UserID=%s, Username=%s", event.UserID, event.Username)
	return kec.profileUpdateHandler.Handle(ctx, event)
}
