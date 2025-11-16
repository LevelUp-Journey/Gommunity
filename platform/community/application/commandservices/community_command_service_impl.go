package commandservices

import (
	"context"
	"errors"
	"log"

	"Gommunity/platform/community/application/outboundservices/acl"
	"Gommunity/platform/community/domain/model/commands"
	"Gommunity/platform/community/domain/model/entities"
	"Gommunity/platform/community/domain/model/valueobjects"
	"Gommunity/platform/community/domain/repositories"
	"Gommunity/platform/community/domain/services"
)

type communityCommandServiceImpl struct {
	communityRepo                repositories.CommunityRepository
	externalSubscriptionsService *acl.ExternalSubscriptionsService
}

func NewCommunityCommandService(
	communityRepo repositories.CommunityRepository,
	externalSubscriptionsService *acl.ExternalSubscriptionsService,
) services.CommunityCommandService {
	return &communityCommandServiceImpl{
		communityRepo:                communityRepo,
		externalSubscriptionsService: externalSubscriptionsService,
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

	// Automatically create subscription with 'owner' role for the community creator
	// This replaces the old TODO about publishing events to update user roles
	if err := s.externalSubscriptionsService.CreateOwnerSubscription(ctx, cmd.OwnerID().Value(), community.CommunityID()); err != nil {
		log.Printf("Warning: Failed to create owner subscription for community %s: %v", community.CommunityID().Value(), err)
		// Note: We don't fail the community creation if subscription fails
		// The community is already created, we just log the error
	} else {
		log.Printf("Owner subscription created successfully for user %s in community %s", cmd.OwnerID().Value(), community.CommunityID().Value())
	}

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
