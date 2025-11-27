package acl

import (
	"context"

	subscriptions_acl "Gommunity/platform/subscriptions/interfaces/acl"
)

// ExternalSubscriptionsService provides ACL access to subscriptions context
type ExternalSubscriptionsService struct {
	subscriptionsFacade subscriptions_acl.SubscriptionsFacade
}

func NewExternalSubscriptionsService(subscriptionsFacade subscriptions_acl.SubscriptionsFacade) *ExternalSubscriptionsService {
	return &ExternalSubscriptionsService{
		subscriptionsFacade: subscriptionsFacade,
	}
}

// GetUserCommunities retrieves all community IDs for a given user
func (s *ExternalSubscriptionsService) GetUserCommunities(ctx context.Context, userID string) ([]string, error) {
	return s.subscriptionsFacade.GetUserCommunityIDs(ctx, userID)
}
