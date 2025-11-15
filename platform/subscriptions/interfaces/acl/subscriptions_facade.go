package acl

import "context"

// SubscriptionsFacade exposes subscription-specific operations to other bounded contexts.
type SubscriptionsFacade interface {
	// GetUserRoleInCommunity returns the role name (member, admin, owner) for a user in a community.
	// Returns empty string when no subscription exists.
	GetUserRoleInCommunity(ctx context.Context, userID string, communityID string) (string, error)

	// IsUserSubscribed indicates whether the user belongs to the community.
	IsUserSubscribed(ctx context.Context, userID string, communityID string) (bool, error)
}
