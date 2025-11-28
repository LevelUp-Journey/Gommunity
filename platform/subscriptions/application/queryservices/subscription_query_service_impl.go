package queryservices

import (
	"context"
	"fmt"

	"Gommunity/platform/subscriptions/domain/model/entities"
	"Gommunity/platform/subscriptions/domain/model/queries"
	"Gommunity/platform/subscriptions/domain/repositories"
	"Gommunity/platform/subscriptions/domain/services"
)

type subscriptionQueryServiceImpl struct {
	subscriptionRepo repositories.SubscriptionRepository
}

// NewSubscriptionQueryService creates a new SubscriptionQueryService implementation
func NewSubscriptionQueryService(
	subscriptionRepo repositories.SubscriptionRepository,
) services.SubscriptionQueryService {
	return &subscriptionQueryServiceImpl{
		subscriptionRepo: subscriptionRepo,
	}
}

// Handle processes a GetSubscriptionByUserAndCommunityQuery to retrieve a specific subscription
func (s *subscriptionQueryServiceImpl) Handle(ctx context.Context, query queries.GetSubscriptionByUserAndCommunityQuery) (*entities.Subscription, error) {
	subscription, err := s.subscriptionRepo.FindByUserAndCommunity(ctx, query.UserID(), query.CommunityID())
	if err != nil {
		return nil, fmt.Errorf("failed to find subscription: %w", err)
	}

	return subscription, nil
}

// HandleCount processes a GetSubscriptionCountByCommunityQuery to get total subscriptions for a community
func (s *subscriptionQueryServiceImpl) HandleCount(ctx context.Context, query queries.GetSubscriptionCountByCommunityQuery) (int64, error) {
	count, err := s.subscriptionRepo.CountByCommunityID(ctx, query.CommunityID())
	if err != nil {
		return 0, fmt.Errorf("failed to count subscriptions: %w", err)
	}

	return count, nil
}

// HandleAll processes a GetAllSubscriptionsByCommunityQuery to retrieve all subscriptions for a community
func (s *subscriptionQueryServiceImpl) HandleAll(ctx context.Context, query queries.GetAllSubscriptionsByCommunityQuery) ([]*entities.Subscription, error) {
	subscriptions, err := s.subscriptionRepo.FindAllByCommunityID(ctx, query.CommunityID(), query.Limit(), query.Offset())
	if err != nil {
		return nil, fmt.Errorf("failed to find subscriptions: %w", err)
	}

	if subscriptions == nil {
		return []*entities.Subscription{}, nil
	}

	return subscriptions, nil
}

// HandleGetCommunitiesByUser processes a GetCommunitiesByUserQuery to retrieve all communities a user is subscribed to
func (s *subscriptionQueryServiceImpl) HandleGetCommunitiesByUser(ctx context.Context, query queries.GetCommunitiesByUserQuery) ([]*entities.Subscription, error) {
	subscriptions, err := s.subscriptionRepo.FindAllByUserID(ctx, query.UserID())
	if err != nil {
		return nil, fmt.Errorf("failed to find user subscriptions: %w", err)
	}

	if subscriptions == nil {
		return []*entities.Subscription{}, nil
	}

	return subscriptions, nil
}
