package commandservices

import (
	"context"
	"errors"
	"log"

	"Gommunity/internal/community/communities/domain/model/commands"
	"Gommunity/internal/community/communities/domain/model/entities"
	"Gommunity/internal/community/communities/domain/model/valueobjects"
	"Gommunity/internal/community/communities/domain/repositories"
	"Gommunity/internal/community/communities/domain/services"
)

type communityCommandServiceImpl struct {
	communityRepo repositories.CommunityRepository
	// TODO: Add ACL service to update user role when we integrate with users context
}

func NewCommunityCommandService(
	communityRepo repositories.CommunityRepository,
) services.CommunityCommandService {
	return &communityCommandServiceImpl{
		communityRepo: communityRepo,
	}
}

func (s *communityCommandServiceImpl) HandleCreate(ctx context.Context, cmd commands.CreateCommunityCommand) (*valueobjects.CommunityID, error) {
	log.Printf("Creating community for owner: %s", cmd.OwnerID().Value())

	// Create community entity
	community, err := entities.NewCommunity(
		cmd.OwnerID(),
		cmd.Name(),
		cmd.Description(),
	)
	if err != nil {
		log.Printf("Error creating community entity: %v", err)
		return nil, err
	}

	// Save community
	if err := s.communityRepo.Save(ctx, community); err != nil {
		log.Printf("Error saving community: %v", err)
		return nil, err
	}

	// TODO: Publish CommunityCreatedEvent to update user role from ROLE_TEACHER to ROLE_OWNER
	// This will be done via Kafka event publishing
	// Event should contain: communityID, ownerID, name

	log.Printf("Community created successfully: %s", community.CommunityID().Value())

	communityID := community.CommunityID()
	return &communityID, nil
}

func (s *communityCommandServiceImpl) HandleDelete(ctx context.Context, cmd commands.DeleteCommunityCommand) error {
	log.Printf("Deleting community: %s by owner: %s", cmd.CommunityID().Value(), cmd.OwnerID().Value())

	// Find community
	community, err := s.communityRepo.FindByID(ctx, cmd.CommunityID())
	if err != nil {
		log.Printf("Error finding community: %v", err)
		return err
	}

	if community == nil {
		return errors.New("community not found")
	}

	// Verify that the user is the owner
	if !community.IsOwner(cmd.OwnerID().Value()) {
		log.Printf("User %s is not the owner of community %s", cmd.OwnerID().Value(), cmd.CommunityID().Value())
		return errors.New("only the owner can delete the community")
	}

	// Delete community
	if err := s.communityRepo.Delete(ctx, cmd.CommunityID()); err != nil {
		log.Printf("Error deleting community: %v", err)
		return err
	}

	log.Printf("Community deleted successfully: %s", cmd.CommunityID().Value())
	return nil
}
