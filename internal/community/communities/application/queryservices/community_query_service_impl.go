package queryservices

import (
	"context"

	"Gommunity/internal/community/communities/domain/model/entities"
	"Gommunity/internal/community/communities/domain/model/queries"
	"Gommunity/internal/community/communities/domain/repositories"
	"Gommunity/internal/community/communities/domain/services"
)

type communityQueryServiceImpl struct {
	communityRepo repositories.CommunityRepository
}

func NewCommunityQueryService(communityRepo repositories.CommunityRepository) services.CommunityQueryService {
	return &communityQueryServiceImpl{
		communityRepo: communityRepo,
	}
}

func (s *communityQueryServiceImpl) HandleGetByID(ctx context.Context, query queries.GetCommunityByIDQuery) (*entities.Community, error) {
	return s.communityRepo.FindByID(ctx, query.CommunityID())
}

func (s *communityQueryServiceImpl) HandleGetByOwner(ctx context.Context, query queries.GetCommunitiesByOwnerQuery) ([]*entities.Community, error) {
	return s.communityRepo.FindByOwnerID(ctx, query.OwnerID())
}

func (s *communityQueryServiceImpl) HandleGetAll(ctx context.Context, query queries.GetAllCommunitiesQuery) ([]*entities.Community, error) {
	// TODO: Implement pagination using query.Limit() and query.Offset()
	return s.communityRepo.FindAll(ctx)
}
