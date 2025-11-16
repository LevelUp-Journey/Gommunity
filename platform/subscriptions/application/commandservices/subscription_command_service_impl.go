package commandservices

import (
	"context"
	"errors"
	"fmt"

	"Gommunity/platform/subscriptions/application/outboundservices/acl"
	"Gommunity/platform/subscriptions/domain/model/commands"
	"Gommunity/platform/subscriptions/domain/model/entities"
	"Gommunity/platform/subscriptions/domain/model/valueobjects"
	"Gommunity/platform/subscriptions/domain/repositories"
	"Gommunity/platform/subscriptions/domain/services"
)

type subscriptionCommandServiceImpl struct {
	subscriptionRepo           repositories.SubscriptionRepository
	externalUsersService       *acl.ExternalUsersService
	externalCommunitiesService *acl.ExternalCommunitiesService
}

// NewSubscriptionCommandService creates a new SubscriptionCommandService implementation
func NewSubscriptionCommandService(
	subscriptionRepo repositories.SubscriptionRepository,
	externalUsersService *acl.ExternalUsersService,
	externalCommunitiesService *acl.ExternalCommunitiesService,
) services.SubscriptionCommandService {
	return &subscriptionCommandServiceImpl{
		subscriptionRepo:           subscriptionRepo,
		externalUsersService:       externalUsersService,
		externalCommunitiesService: externalCommunitiesService,
	}
}

// Handle processes a SubscribeUserCommand to add a user to a community with a role
func (s *subscriptionCommandServiceImpl) Handle(ctx context.Context, cmd commands.SubscribeUserCommand) (*valueobjects.SubscriptionID, error) {
	// Step 1: Validate that the community exists
	communityExists, err := s.externalCommunitiesService.ValidateCommunityExists(ctx, cmd.CommunityID())
	if err != nil {
		return nil, fmt.Errorf("failed to validate community existence: %w", err)
	}
	if !communityExists {
		return nil, errors.New("community not found")
	}

	// Step 2: Validate that the user to be subscribed exists
	userExists, err := s.externalUsersService.ValidateUserExists(ctx, cmd.UserID())
	if err != nil {
		return nil, fmt.Errorf("failed to validate user existence: %w", err)
	}
	if !userExists {
		return nil, errors.New("user not found")
	}

	// Step 3: Validate that the requesting user exists
	requestedByExists, err := s.externalUsersService.ValidateUserExists(ctx, cmd.RequestedBy())
	if err != nil {
		return nil, fmt.Errorf("failed to validate requesting user existence: %w", err)
	}
	if !requestedByExists {
		return nil, errors.New("requesting user not found")
	}

	// Step 4: Validate that the role exists in the Users BC
	roleExists, err := s.externalUsersService.ValidateRoleExists(ctx, cmd.Role().Value())
	if err != nil {
		return nil, fmt.Errorf("failed to validate role existence: %w", err)
	}
	if !roleExists {
		return nil, fmt.Errorf("role '%s' not found in Users BC", cmd.Role().Value())
	}

	// Step 5: Check if this is a subscription to a private or public community
	isPrivate, err := s.externalCommunitiesService.IsCommunityPrivate(ctx, cmd.CommunityID())
	if err != nil {
		return nil, fmt.Errorf("failed to check community privacy: %w", err)
	}

	// Step 6: Determine the actual role to assign based on business rules
	var actualRole valueobjects.CommunityRole

	if cmd.IsSelfSubscription() {
		// IMPORTANT: When users subscribe themselves (follow a community), they ALWAYS get the 'member' role
		// regardless of what role they requested. This prevents privilege escalation.
		actualRole = valueobjects.MemberRole
	} else {
		// When an admin/teacher adds another user, they can specify the role
		actualRole = cmd.Role()
	}

	// Step 7: Apply business rules based on community privacy
	if isPrivate {
		// In private communities, only owner/admin can add users
		if !cmd.IsSelfSubscription() {
			// Verify the requester has permission to add users
			requesterSubscription, err := s.subscriptionRepo.FindByUserAndCommunity(ctx, cmd.RequestedBy(), cmd.CommunityID())
			if err != nil {
				return nil, fmt.Errorf("failed to check requester permissions: %w", err)
			}

			// Check if requester is the community owner
			// Get requester's profile ID to compare with owner ID
			requesterProfileID, err := s.externalUsersService.GetProfileIDByUserID(ctx, cmd.RequestedBy())
			if err != nil {
				return nil, fmt.Errorf("failed to get requester profile ID: %w", err)
			}

			isOwner, err := s.externalCommunitiesService.ValidateUserIsOwner(ctx, cmd.CommunityID(), requesterProfileID)
			if err != nil {
				return nil, fmt.Errorf("failed to validate owner status: %w", err)
			}

			// Requester must be owner OR have admin role in the community
			hasPermission := isOwner
			if requesterSubscription != nil && requesterSubscription.HasAdminOrOwnerRole() {
				hasPermission = true
			}

			if !hasPermission {
				return nil, errors.New("only community owner or admins can add users to private communities")
			}
		}
	} else {
		// In public communities, users can only subscribe themselves
		if !cmd.IsSelfSubscription() {
			return nil, errors.New("users can only subscribe themselves to public communities")
		}
	}

	// Step 8: Check if the user is already subscribed to this community
	alreadySubscribed, err := s.subscriptionRepo.ExistsByUserAndCommunity(ctx, cmd.UserID(), cmd.CommunityID())
	if err != nil {
		return nil, fmt.Errorf("failed to check existing subscription: %w", err)
	}
	if alreadySubscribed {
		return nil, errors.New("user is already subscribed to this community")
	}

	// Step 9: Create the subscription entity with the actual role (not the requested role)
	subscription, err := entities.NewSubscription(
		cmd.UserID(),
		cmd.CommunityID(),
		actualRole, // Use the actual role determined by business rules
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription entity: %w", err)
	}

	// Step 10: Persist the subscription
	if err := s.subscriptionRepo.Save(ctx, subscription); err != nil {
		return nil, fmt.Errorf("failed to save subscription: %w", err)
	}

	subscriptionID := subscription.SubscriptionID()
	return &subscriptionID, nil
}

// HandleUnsubscribe processes an UnsubscribeUserCommand to remove a user from a community
func (s *subscriptionCommandServiceImpl) HandleUnsubscribe(ctx context.Context, cmd commands.UnsubscribeUserCommand) error {
	// Step 1: Validate that the community exists
	communityExists, err := s.externalCommunitiesService.ValidateCommunityExists(ctx, cmd.CommunityID())
	if err != nil {
		return fmt.Errorf("failed to validate community existence: %w", err)
	}
	if !communityExists {
		return errors.New("community not found")
	}

	// Step 2: Validate that the subscription exists
	subscription, err := s.subscriptionRepo.FindByUserAndCommunity(ctx, cmd.UserID(), cmd.CommunityID())
	if err != nil {
		return fmt.Errorf("failed to find subscription: %w", err)
	}
	if subscription == nil {
		return errors.New("subscription not found")
	}

	// Step 3: Check if the requester has permission to remove this subscription
	if !cmd.IsSelfUnsubscription() {
		// If it's not self-unsubscription, verify the requester has permission
		requesterSubscription, err := s.subscriptionRepo.FindByUserAndCommunity(ctx, cmd.RequestedBy(), cmd.CommunityID())
		if err != nil {
			return fmt.Errorf("failed to check requester permissions: %w", err)
		}

		// Check if requester is the community owner
		// Get requester's profile ID to compare with owner ID
		requesterProfileID, err := s.externalUsersService.GetProfileIDByUserID(ctx, cmd.RequestedBy())
		if err != nil {
			return fmt.Errorf("failed to get requester profile ID: %w", err)
		}

		isOwner, err := s.externalCommunitiesService.ValidateUserIsOwner(ctx, cmd.CommunityID(), requesterProfileID)
		if err != nil {
			return fmt.Errorf("failed to validate owner status: %w", err)
		}

		// Requester must be owner OR have admin role in the community
		hasPermission := isOwner
		if requesterSubscription != nil && requesterSubscription.HasAdminOrOwnerRole() {
			hasPermission = true
		}

		if !hasPermission {
			return errors.New("only community owner, admins, or the user themselves can remove subscriptions")
		}
	}

	// Step 4: Prevent owner from unsubscribing (owner should always be subscribed)
	// Get user's profile ID to compare with owner ID
	userProfileID, err := s.externalUsersService.GetProfileIDByUserID(ctx, cmd.UserID())
	if err != nil {
		return fmt.Errorf("failed to get user profile ID: %w", err)
	}

	isOwner, err := s.externalCommunitiesService.ValidateUserIsOwner(ctx, cmd.CommunityID(), userProfileID)
	if err != nil {
		return fmt.Errorf("failed to validate owner status: %w", err)
	}
	if isOwner {
		return errors.New("community owner cannot unsubscribe from their own community")
	}

	// Step 5: Delete the subscription (this also removes the community-specific role)
	if err := s.subscriptionRepo.DeleteByUserAndCommunity(ctx, cmd.UserID(), cmd.CommunityID()); err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}
