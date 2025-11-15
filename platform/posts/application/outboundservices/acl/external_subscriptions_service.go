package acl

import (
	"context"

	"Gommunity/platform/posts/domain/model/valueobjects"
	subscriptions_acl "Gommunity/platform/subscriptions/interfaces/acl"
)

// ExternalSubscriptionsService provides read access to subscription data.
type ExternalSubscriptionsService struct {
	subscriptionsFacade subscriptions_acl.SubscriptionsFacade
}

// NewExternalSubscriptionsService builds a new ExternalSubscriptionsService.
func NewExternalSubscriptionsService(subscriptionsFacade subscriptions_acl.SubscriptionsFacade) *ExternalSubscriptionsService {
	return &ExternalSubscriptionsService{
		subscriptionsFacade: subscriptionsFacade,
	}
}

// GetUserRole retrieves the user's community role (member/admin/owner).
func (s *ExternalSubscriptionsService) GetUserRole(ctx context.Context, userID valueobjects.AuthorID, communityID valueobjects.CommunityID) (*valueobjects.CommunityRole, error) {
	roleValue, err := s.subscriptionsFacade.GetUserRoleInCommunity(ctx, userID.Value(), communityID.Value())
	if err != nil {
		return nil, err
	}

	if roleValue == "" {
		return nil, nil
	}

	role, err := valueobjects.NewCommunityRole(roleValue)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// IsUserSubscribed checks whether the user belongs to the community.
func (s *ExternalSubscriptionsService) IsUserSubscribed(ctx context.Context, userID valueobjects.AuthorID, communityID valueobjects.CommunityID) (bool, error) {
	return s.subscriptionsFacade.IsUserSubscribed(ctx, userID.Value(), communityID.Value())
}
