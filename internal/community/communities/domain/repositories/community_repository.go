package repositories

import (
	"context"

	"Gommunity/internal/community/communities/domain/model/entities"
	"Gommunity/internal/community/communities/domain/model/valueobjects"
)

type CommunityRepository interface {
	Save(ctx context.Context, community *entities.Community) error
	Update(ctx context.Context, community *entities.Community) error
	FindByID(ctx context.Context, communityID valueobjects.CommunityID) (*entities.Community, error)
	FindByOwnerID(ctx context.Context, ownerID valueobjects.OwnerID) ([]*entities.Community, error)
	FindAll(ctx context.Context) ([]*entities.Community, error)
	Delete(ctx context.Context, communityID valueobjects.CommunityID) error
	ExistsByID(ctx context.Context, communityID valueobjects.CommunityID) (bool, error)
}
