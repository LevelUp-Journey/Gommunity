package services

import (
	"context"

	"Gommunity/platform/users/domain/model/entities"
	"Gommunity/platform/users/domain/model/queries"
)

type UserQueryService interface {
	HandleGetByID(ctx context.Context, query queries.GetUserByIDQuery) (*entities.User, error)
	HandleGetByUsername(ctx context.Context, query queries.GetUserByUsernameQuery) (*entities.User, error)
}
