package repositories

import (
	"context"

	"Gommunity/platform/users/domain/model/entities"
	"Gommunity/platform/users/domain/model/valueobjects"
)

type UserRepository interface {
	Save(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	FindByUserID(ctx context.Context, userID valueobjects.UserID) (*entities.User, error)
	FindByProfileID(ctx context.Context, profileID valueobjects.ProfileID) (*entities.User, error)
	FindByUsername(ctx context.Context, username valueobjects.Username) (*entities.User, error)
	ExistsByUserID(ctx context.Context, userID valueobjects.UserID) (bool, error)
	Delete(ctx context.Context, userID valueobjects.UserID) error
}
