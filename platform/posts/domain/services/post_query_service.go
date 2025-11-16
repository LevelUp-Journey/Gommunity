package services

import (
	"context"

	"Gommunity/platform/posts/domain/model/entities"
	"Gommunity/platform/posts/domain/model/queries"
)

// PostQueryService defines query handling behavior for posts.
type PostQueryService interface {
	HandleGetByID(ctx context.Context, query queries.GetPostByIDQuery) (*entities.Post, error)
	HandleGetByCommunity(ctx context.Context, query queries.GetPostsByCommunityQuery) ([]*entities.Post, error)
}
