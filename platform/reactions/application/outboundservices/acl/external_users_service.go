package acl

import (
	"context"
	"fmt"

	"Gommunity/platform/reactions/domain/model/valueobjects"
	users_acl "Gommunity/platform/users/interfaces/acl"
)

// ExternalUsersService validates users from the users bounded context.
type ExternalUsersService struct {
	usersFacade users_acl.UsersFacade
}

// NewExternalUsersService constructs the external users service.
func NewExternalUsersService(usersFacade users_acl.UsersFacade) *ExternalUsersService {
	return &ExternalUsersService{
		usersFacade: usersFacade,
	}
}

// ValidateUserExists checks if a user exists.
func (s *ExternalUsersService) ValidateUserExists(ctx context.Context, userID valueobjects.UserID) (bool, error) {
	exists, err := s.usersFacade.ValidateUserExists(ctx, userID.Value())
	if err != nil {
		return false, fmt.Errorf("failed to validate user existence: %w", err)
	}
	return exists, nil
}
