package repositories

import (
	"context"

	"Gommunity/platform/subscriptions/domain/model/entities"
	"Gommunity/platform/subscriptions/domain/model/valueobjects"
)

// SubscriptionRepository defines the contract for subscription persistence operations
type SubscriptionRepository interface {
	// Save persists a subscription
	Save(ctx context.Context, subscription *entities.Subscription) error

	// FindByID retrieves a subscription by its ID
	FindByID(ctx context.Context, id valueobjects.SubscriptionID) (*entities.Subscription, error)

	// FindByUserAndCommunity retrieves a subscription by user ID and community ID
	FindByUserAndCommunity(ctx context.Context, userID valueobjects.UserID, communityID valueobjects.CommunityID) (*entities.Subscription, error)

	// FindAllByCommunityID retrieves all subscriptions for a specific community
	FindAllByCommunityID(ctx context.Context, communityID valueobjects.CommunityID, limit, offset *int) ([]*entities.Subscription, error)

	// FindAllByUserID retrieves all subscriptions for a specific user
	FindAllByUserID(ctx context.Context, userID valueobjects.UserID) ([]*entities.Subscription, error)

	// CountByCommunityID returns the total number of subscriptions for a community
	CountByCommunityID(ctx context.Context, communityID valueobjects.CommunityID) (int64, error)

	// ExistsByUserAndCommunity checks if a subscription exists for a user in a community
	ExistsByUserAndCommunity(ctx context.Context, userID valueobjects.UserID, communityID valueobjects.CommunityID) (bool, error)

	// Delete removes a subscription
	Delete(ctx context.Context, id valueobjects.SubscriptionID) error

	// DeleteByUserAndCommunity removes a subscription by user and community
	DeleteByUserAndCommunity(ctx context.Context, userID valueobjects.UserID, communityID valueobjects.CommunityID) error

	// DeleteByCommunity removes all subscriptions for a given community
	DeleteByCommunity(ctx context.Context, communityID valueobjects.CommunityID) error
}
