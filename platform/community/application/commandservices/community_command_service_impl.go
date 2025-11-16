package commandservices

import (
	"context"
	"errors"
	"log"

	"Gommunity/platform/community/domain/model/commands"
	"Gommunity/platform/community/domain/model/entities"
	"Gommunity/platform/community/domain/model/valueobjects"
	"Gommunity/platform/community/domain/repositories"
	"Gommunity/platform/community/domain/services"
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
		cmd.IconURL(),
		cmd.BannerURL(),
		cmd.IsPrivate(),
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

func (s *communityCommandServiceImpl) HandleUpdatePrivacy(ctx context.Context, cmd commands.UpdateCommunityPrivacyCommand) error {
	log.Printf("Updating privacy for community: %s", cmd.CommunityID().Value())

	// Find community
	community, err := s.communityRepo.FindByID(ctx, cmd.CommunityID())
	if err != nil {
		log.Printf("Error finding community: %v", err)
		return err
	}

	if community == nil {
		return errors.New("community not found")
	}

	// Update privacy status
	community.UpdatePrivacy(cmd.IsPrivate())

	// Save updated community
	if err := s.communityRepo.Update(ctx, community); err != nil {
		log.Printf("Error updating community privacy: %v", err)
		return err
	}

	log.Printf("Community privacy updated successfully: %s, isPrivate: %v", cmd.CommunityID().Value(), cmd.IsPrivate())
	return nil
}

func (s *communityCommandServiceImpl) HandleUpdateInfo(ctx context.Context, cmd commands.UpdateCommunityInfoCommand) error {
	log.Printf("Updating info for community: %s", cmd.CommunityID().Value())

	// Find community
	community, err := s.communityRepo.FindByID(ctx, cmd.CommunityID())
	if err != nil {
		log.Printf("Error finding community: %v", err)
		return err
	}

	if community == nil {
		return errors.New("community not found")
	}

	// Update community info
	community.UpdateInfo(cmd.Name(), cmd.Description(), cmd.IconURL(), cmd.BannerURL())

	// Save updated community
	if err := s.communityRepo.Update(ctx, community); err != nil {
		log.Printf("Error updating community info: %v", err)
		return err
	}

	log.Printf("Community info updated successfully: %s", cmd.CommunityID().Value())
	return nil
}
