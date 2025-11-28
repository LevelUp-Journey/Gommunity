package services

import (
	"context"

	"Gommunity/platform/subscriptions/domain/model/entities"
	"Gommunity/platform/subscriptions/domain/model/queries"
)

// SubscriptionQueryService defines the contract for subscription query operations
type SubscriptionQueryService interface {
	// Handle processes a GetSubscriptionByUserAndCommunityQuery to retrieve a specific subscription
	Handle(ctx context.Context, query queries.GetSubscriptionByUserAndCommunityQuery) (*entities.Subscription, error)

	// HandleCount processes a GetSubscriptionCountByCommunityQuery to get total subscriptions for a community
	HandleCount(ctx context.Context, query queries.GetSubscriptionCountByCommunityQuery) (int64, error)

	// HandleAll processes a GetAllSubscriptionsByCommunityQuery to retrieve all subscriptions for a community
	HandleAll(ctx context.Context, query queries.GetAllSubscriptionsByCommunityQuery) ([]*entities.Subscription, error)

	// HandleGetCommunitiesByUser processes a GetCommunitiesByUserQuery to retrieve all communities a user is subscribed to
	HandleGetCommunitiesByUser(ctx context.Context, query queries.GetCommunitiesByUserQuery) ([]*entities.Subscription, error)
}
