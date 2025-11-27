package acl

import "context"

// CommunitiesFacade provides an anti-corruption layer for accessing Community bounded context functionality
// This facade exposes only the necessary operations needed by other bounded contexts
type CommunitiesFacade interface {
	// ValidateCommunityExists checks if a community exists by ID
	ValidateCommunityExists(ctx context.Context, communityID string) (bool, error)

	// IsCommunityPrivate checks if a community is private
	IsCommunityPrivate(ctx context.Context, communityID string) (bool, error)

	// GetCommunityOwnerID retrieves the owner ID of a community (as string UUID)
	GetCommunityOwnerID(ctx context.Context, communityID string) (string, error)

	// ValidateUserIsOwner checks if a user is the owner of a community
	ValidateUserIsOwner(ctx context.Context, communityID string, ownerID string) (bool, error)
}
