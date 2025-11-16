package acl

import "context"

// UsersFacade provides an anti-corruption layer for accessing User bounded context functionality
// This facade exposes only the necessary operations needed by other bounded contexts
type UsersFacade interface {
	// FindUserIDByUsername retrieves a user ID by username (returns UUID string)
	FindUserIDByUsername(ctx context.Context, username string) (string, error)

	// ValidateUserExists checks if a user exists by ID (UUID string)
	ValidateUserExists(ctx context.Context, userID string) (bool, error)

	// ValidateRoleExists checks if a role exists by name
	ValidateRoleExists(ctx context.Context, roleName string) (bool, error)

	// GetUserRoleInCommunity retrieves the user's role in a specific community
	// Returns empty string if user has no role in that community
	GetUserRoleInCommunity(ctx context.Context, userID string, communityID string) (string, error)

	// GetProfileIDByUserID retrieves a user's profile ID (UUID) by user ID (UUID string)
	GetProfileIDByUserID(ctx context.Context, userID string) (string, error)
}
