package eventhandlers

import (
	"context"
	"log"

	"Gommunity/platform/users/domain/model/entities"
	"Gommunity/platform/users/domain/model/events"
	"Gommunity/platform/users/domain/model/valueobjects"
	"Gommunity/platform/users/domain/repositories"
)

type UserRegistrationHandler struct {
	userRepository repositories.UserRepository
}

func NewUserRegistrationHandler(userRepository repositories.UserRepository) *UserRegistrationHandler {
	return &UserRegistrationHandler{
		userRepository: userRepository,
	}
}

// Handle processes the CommunityRegistrationEvent
func (h *UserRegistrationHandler) Handle(ctx context.Context, event events.CommunityRegistrationEvent) error {
	log.Printf("Processing community registration event for user: %s", event.UserID)

	// Create value objects
	userID, err := valueobjects.NewUserID(event.UserID)
	if err != nil {
		log.Printf("Error creating UserID: %v", err)
		return err
	}

	profileID, err := valueobjects.NewProfileID(event.ProfileID)
	if err != nil {
		log.Printf("Error creating ProfileID: %v", err)
		return err
	}

	username, err := valueobjects.NewUsername(event.Username)
	if err != nil {
		log.Printf("Error creating Username: %v", err)
		return err
	}

	// Check if user already exists
	exists, err := h.userRepository.ExistsByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		return err
	}

	if exists {
		log.Printf("User already exists: %s", event.UserID)
		return nil // Skip if already exists
	}

	// Create new user entity
	user, err := entities.NewUser(userID, profileID, username, event.ProfileURL)
	if err != nil {
		log.Printf("Error creating user entity: %v", err)
		return err
	}

	// Save user
	if err := h.userRepository.Save(ctx, user); err != nil {
		log.Printf("Error saving user: %v", err)
		return err
	}

	log.Printf("User registered successfully: %s", event.UserID)
	return nil
}
