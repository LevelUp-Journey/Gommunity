package queryservices

import (
	"context"

	"Gommunity/internal/community/users/domain/model/entities"
	"Gommunity/internal/community/users/domain/model/queries"
	"Gommunity/internal/community/users/domain/repositories"
	"Gommunity/internal/community/users/domain/services"
)

type userQueryServiceImpl struct {
	userRepository repositories.UserRepository
}

func NewUserQueryService(userRepository repositories.UserRepository) services.UserQueryService {
	return &userQueryServiceImpl{
		userRepository: userRepository,
	}
}

func (s *userQueryServiceImpl) HandleGetByID(ctx context.Context, query queries.GetUserByIDQuery) (*entities.User, error) {
	return s.userRepository.FindByUserID(ctx, query.UserID())
}

func (s *userQueryServiceImpl) HandleGetByUsername(ctx context.Context, query queries.GetUserByUsernameQuery) (*entities.User, error) {
	return s.userRepository.FindByUsername(ctx, query.Username())
}
