package eventhandlers

import (
	"context"
	"log"

	"Gommunity/platform/users/domain/model/events"
	"Gommunity/platform/users/domain/model/valueobjects"
	"Gommunity/platform/users/domain/repositories"
)

type ProfileUpdatedHandler struct {
	userRepository repositories.UserRepository
}

func NewProfileUpdatedHandler(userRepository repositories.UserRepository) *ProfileUpdatedHandler {
	return &ProfileUpdatedHandler{
		userRepository: userRepository,
	}
}

// Handle processes the ProfileUpdatedEvent
func (h *ProfileUpdatedHandler) Handle(ctx context.Context, event events.ProfileUpdatedEvent) error {
	log.Printf("Processing profile updated event for user: %s", event.UserID)

	// Create value objects
	userID, err := valueobjects.NewUserID(event.UserID)
	if err != nil {
		log.Printf("Error creating UserID: %v", err)
		return err
	}

	username, err := valueobjects.NewUsername(event.Username)
	if err != nil {
		log.Printf("Error creating Username: %v", err)
		return err
	}

	// Find existing user
	user, err := h.userRepository.FindByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return err
	}

	if user == nil {
		log.Printf("User not found: %s", event.UserID)
		return nil // Skip if user doesn't exist
	}

	// Update user profile
	user.UpdateProfile(username, event.ProfileURL, nil)

	// Save updated user
	if err := h.userRepository.Update(ctx, user); err != nil {
		log.Printf("Error updating user: %v", err)
		return err
	}

	log.Printf("User profile updated successfully: %s", event.UserID)
	return nil
}
