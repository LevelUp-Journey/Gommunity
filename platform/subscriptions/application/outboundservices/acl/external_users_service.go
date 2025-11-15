package acl

import (
	"context"
	"errors"

	users_acl "Gommunity/platform/users/interfaces/acl"
	"Gommunity/platform/subscriptions/domain/model/valueobjects"
)

// ExternalUsersService provides ACL implementation for accessing User bounded context
type ExternalUsersService struct {
	usersFacade users_acl.UsersFacade
}

// NewExternalUsersService creates a new ExternalUsersService
func NewExternalUsersService(usersFacade users_acl.UsersFacade) *ExternalUsersService {
	return &ExternalUsersService{
		usersFacade: usersFacade,
	}
}

// FetchUserIDByUsername retrieves a user ID by username from the Users BC
func (s *ExternalUsersService) FetchUserIDByUsername(ctx context.Context, username string) (*valueobjects.UserID, error) {
	userIDValue, err := s.usersFacade.FindUserIDByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if userIDValue == "" {
		return nil, errors.New("user not found")
	}

	userID, err := valueobjects.NewUserID(userIDValue)
	if err != nil {
		return nil, err
	}

	return &userID, nil
}

// ValidateUserExists checks if a user exists in the Users BC
func (s *ExternalUsersService) ValidateUserExists(ctx context.Context, userID valueobjects.UserID) (bool, error) {
	return s.usersFacade.ValidateUserExists(ctx, userID.Value())
}

// ValidateRoleExists checks if a role exists in the Users BC
func (s *ExternalUsersService) ValidateRoleExists(ctx context.Context, roleName string) (bool, error) {
	return s.usersFacade.ValidateRoleExists(ctx, roleName)
}

// GetProfileIDByUserID retrieves a user's profile ID (UUID) by user ID
func (s *ExternalUsersService) GetProfileIDByUserID(ctx context.Context, userID valueobjects.UserID) (string, error) {
	return s.usersFacade.GetProfileIDByUserID(ctx, userID.Value())
}
