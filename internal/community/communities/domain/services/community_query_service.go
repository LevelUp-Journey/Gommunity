package services

import (
	"context"

	"Gommunity/internal/community/communities/domain/model/entities"
	"Gommunity/internal/community/communities/domain/model/queries"
)

type CommunityQueryService interface {
	HandleGetByID(ctx context.Context, query queries.GetCommunityByIDQuery) (*entities.Community, error)
	HandleGetByOwner(ctx context.Context, query queries.GetCommunitiesByOwnerQuery) ([]*entities.Community, error)
	HandleGetAll(ctx context.Context, query queries.GetAllCommunitiesQuery) ([]*entities.Community, error)
}
