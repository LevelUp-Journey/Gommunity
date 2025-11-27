package acl

import (
	"context"
	"fmt"

	community_vo "Gommunity/platform/community/domain/model/valueobjects"
	subscription_commands "Gommunity/platform/subscriptions/domain/model/commands"
	subscription_vo "Gommunity/platform/subscriptions/domain/model/valueobjects"
	subscription_services "Gommunity/platform/subscriptions/domain/services"
)

// ExternalSubscriptionsService provides access to Subscriptions BC operations
type ExternalSubscriptionsService struct {
	subscriptionCommandService subscription_services.SubscriptionCommandService
}

func NewExternalSubscriptionsService(
	subscriptionCommandService subscription_services.SubscriptionCommandService,
) *ExternalSubscriptionsService {
	return &ExternalSubscriptionsService{
		subscriptionCommandService: subscriptionCommandService,
	}
}

// CreateOwnerSubscription creates a subscription with 'owner' role when a community is created
func (s *ExternalSubscriptionsService) CreateOwnerSubscription(
	ctx context.Context,
	userID string,
	communityID community_vo.CommunityID,
) error {
	// Convert community ID to subscription's CommunityID value object
	subCommunityID, err := subscription_vo.NewCommunityID(communityID.Value())
	if err != nil {
		return fmt.Errorf("failed to create community ID: %w", err)
	}

	// Create UserID value object for subscription
	subUserID, err := subscription_vo.NewUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to create user ID: %w", err)
	}

	// Create subscribe command with 'owner' role
	// The user is subscribing themselves as owner when creating the community
	cmd, err := subscription_commands.NewSubscribeUserCommand(
		subUserID,
		subCommunityID,
		subscription_vo.OwnerRole, // Owner role for community creator
		subUserID,                 // Requested by themselves
	)
	if err != nil {
		return fmt.Errorf("failed to create subscribe command: %w", err)
	}

	// Execute subscription command
	_, err = s.subscriptionCommandService.Handle(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to create owner subscription: %w", err)
	}

	return nil
}

// DeleteSubscriptionsByCommunity removes all subscriptions for the given community
func (s *ExternalSubscriptionsService) DeleteSubscriptionsByCommunity(ctx context.Context, communityID community_vo.CommunityID) error {
	subCommunityID, err := subscription_vo.NewCommunityID(communityID.Value())
	if err != nil {
		return fmt.Errorf("failed to create community ID: %w", err)
	}

	return s.subscriptionCommandService.HandleDeleteByCommunity(ctx, subCommunityID)
}
