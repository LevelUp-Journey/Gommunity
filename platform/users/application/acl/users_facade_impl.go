package acl

import (
	"context"
	"errors"

	"Gommunity/platform/users/domain/model/valueobjects"
	"Gommunity/platform/users/domain/repositories"
	"Gommunity/platform/users/interfaces/acl"
)

type usersFacadeImpl struct {
	userRepository repositories.UserRepository
}

// NewUsersFacade creates a new UsersFacade implementation
func NewUsersFacade(
	userRepository repositories.UserRepository,
) acl.UsersFacade {
	return &usersFacadeImpl{
		userRepository: userRepository,
	}
}

// FindUserIDByUsername retrieves a user ID by username (returns UUID string)
func (f *usersFacadeImpl) FindUserIDByUsername(ctx context.Context, username string) (string, error) {
	usernameVO, err := valueobjects.NewUsername(username)
	if err != nil {
		return "", err
	}

	user, err := f.userRepository.FindByUsername(ctx, usernameVO)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("user not found")
	}

	return user.UserID().Value(), nil
}

// ValidateUserExists checks if a user exists by ID (UUID string)
func (f *usersFacadeImpl) ValidateUserExists(ctx context.Context, userID string) (bool, error) {
	userIDVO, err := valueobjects.NewUserID(userID)
	if err != nil {
		return false, err
	}

	return f.userRepository.ExistsByUserID(ctx, userIDVO)
}

// ValidateRoleExists checks if a role exists by name
// Note: Users BC no longer manages roles. Roles are managed per-community in Subscriptions BC.
// This method always returns true as role validation is not needed at this level.
func (f *usersFacadeImpl) ValidateRoleExists(ctx context.Context, roleName string) (bool, error) {
	// Users BC doesn't manage roles anymore
	// Role validation happens in Subscriptions BC for community-specific roles
	// IAM roles (STUDENT, TEACHER, ADMIN) come from JWT and don't need validation
	return true, nil
}

// GetUserRoleInCommunity retrieves the user's role in a specific community
// Note: In the Users BC, we don't store community-specific roles
// This should be handled by the Subscriptions BC
// This method returns empty string as Users BC doesn't manage community roles
func (f *usersFacadeImpl) GetUserRoleInCommunity(ctx context.Context, userID string, communityID string) (string, error) {
	// Users BC doesn't manage community-specific roles
	// This is managed by the Subscriptions BC
	// Return empty string to indicate no role at this level
	return "", nil
}

// GetProfileIDByUserID retrieves a user's profile ID (UUID) by user ID (UUID string)
func (f *usersFacadeImpl) GetProfileIDByUserID(ctx context.Context, userID string) (string, error) {
	userIDVO, err := valueobjects.NewUserID(userID)
	if err != nil {
		return "", err
	}

	user, err := f.userRepository.FindByUserID(ctx, userIDVO)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("user not found")
	}

	return user.ProfileID().Value(), nil
}
