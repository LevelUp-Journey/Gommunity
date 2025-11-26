package services

import (
	"context"

	"Gommunity/platform/subscriptions/domain/model/commands"
	"Gommunity/platform/subscriptions/domain/model/valueobjects"
)

// SubscriptionCommandService defines the contract for subscription command operations
type SubscriptionCommandService interface {
	// Handle processes a SubscribeUserCommand to add a user to a community with a role
	Handle(ctx context.Context, cmd commands.SubscribeUserCommand) (*valueobjects.SubscriptionID, error)

	// HandleUnsubscribe processes an UnsubscribeUserCommand to remove a user from a community
	HandleUnsubscribe(ctx context.Context, cmd commands.UnsubscribeUserCommand) error

	// HandleDeleteByCommunity removes all subscriptions linked to a community
	HandleDeleteByCommunity(ctx context.Context, communityID valueobjects.CommunityID) error
}
