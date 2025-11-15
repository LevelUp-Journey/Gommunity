package queryservices

import (
	"context"

	"Gommunity/platform/posts/domain/model/entities"
	"Gommunity/platform/posts/domain/model/queries"
	"Gommunity/platform/posts/domain/repositories"
	"Gommunity/platform/posts/domain/services"
)

type postQueryServiceImpl struct {
	postRepository repositories.PostRepository
}

// NewPostQueryService creates a query service implementation.
func NewPostQueryService(postRepository repositories.PostRepository) services.PostQueryService {
	return &postQueryServiceImpl{
		postRepository: postRepository,
	}
}

// HandleGetByID retrieves a post by identifier.
func (s *postQueryServiceImpl) HandleGetByID(ctx context.Context, query queries.GetPostByIDQuery) (*entities.Post, error) {
	return s.postRepository.FindByID(ctx, query.PostID())
}

// HandleGetByCommunity retrieves posts for a community.
func (s *postQueryServiceImpl) HandleGetByCommunity(ctx context.Context, query queries.GetPostsByCommunityQuery) ([]*entities.Post, error) {
	return s.postRepository.FindByCommunity(ctx, query.CommunityID(), query.Limit(), query.Offset())
}
