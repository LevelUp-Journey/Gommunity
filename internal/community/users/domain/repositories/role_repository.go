package repositories

import (
	"context"

	"Gommunity/internal/community/users/domain/model/entities"
	"Gommunity/internal/community/users/domain/model/valueobjects"
)

type RoleRepository interface {
	Save(ctx context.Context, role *entities.Role) error
	FindByID(ctx context.Context, roleID valueobjects.RoleID) (*entities.Role, error)
	FindByName(ctx context.Context, name string) (*entities.Role, error)
	FindAll(ctx context.Context) ([]*entities.Role, error)
}
