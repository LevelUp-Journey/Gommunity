package acl

import (
	"context"

	"Gommunity/platform/posts/domain/model/valueobjects"
	users_acl "Gommunity/platform/users/interfaces/acl"
)

// ExternalUsersService provides access to the Users bounded context.
type ExternalUsersService struct {
	usersFacade users_acl.UsersFacade
}

// NewExternalUsersService builds a new ExternalUsersService.
func NewExternalUsersService(usersFacade users_acl.UsersFacade) *ExternalUsersService {
	return &ExternalUsersService{
		usersFacade: usersFacade,
	}
}

// ValidateUserExists ensures the author exists in Users BC.
func (s *ExternalUsersService) ValidateUserExists(ctx context.Context, userID valueobjects.AuthorID) (bool, error) {
	return s.usersFacade.ValidateUserExists(ctx, userID.Value())
}
