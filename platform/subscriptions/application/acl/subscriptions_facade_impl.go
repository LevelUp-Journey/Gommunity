package acl

import (
	"context"

	"Gommunity/platform/subscriptions/domain/model/valueobjects"
	"Gommunity/platform/subscriptions/domain/repositories"
	"Gommunity/platform/subscriptions/interfaces/acl"
)

type subscriptionsFacadeImpl struct {
	subscriptionRepository repositories.SubscriptionRepository
}

// NewSubscriptionsFacade creates a new SubscriptionsFacade implementation.
func NewSubscriptionsFacade(
	subscriptionRepository repositories.SubscriptionRepository,
) acl.SubscriptionsFacade {
	return &subscriptionsFacadeImpl{
		subscriptionRepository: subscriptionRepository,
	}
}

// GetUserRoleInCommunity returns the role granted to a user within a community.
func (f *subscriptionsFacadeImpl) GetUserRoleInCommunity(ctx context.Context, userID string, communityID string) (string, error) {
	userIDVO, err := valueobjects.NewUserID(userID)
	if err != nil {
		return "", err
	}

	communityIDVO, err := valueobjects.NewCommunityID(communityID)
	if err != nil {
		return "", err
	}

	subscription, err := f.subscriptionRepository.FindByUserAndCommunity(ctx, userIDVO, communityIDVO)
	if err != nil {
		return "", err
	}

	if subscription == nil {
		return "", nil
	}

	return subscription.Role().Value(), nil
}

// IsUserSubscribed checks whether a user belongs to a community.
func (f *subscriptionsFacadeImpl) IsUserSubscribed(ctx context.Context, userID string, communityID string) (bool, error) {
	userIDVO, err := valueobjects.NewUserID(userID)
	if err != nil {
		return false, err
	}

	communityIDVO, err := valueobjects.NewCommunityID(communityID)
	if err != nil {
		return false, err
	}

	return f.subscriptionRepository.ExistsByUserAndCommunity(ctx, userIDVO, communityIDVO)
}
