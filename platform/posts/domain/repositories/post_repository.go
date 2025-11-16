package repositories

import (
	"context"

	"Gommunity/platform/posts/domain/model/entities"
	"Gommunity/platform/posts/domain/model/valueobjects"
)

// PostRepository defines persistence operations for post aggregates.
type PostRepository interface {
	Save(ctx context.Context, post *entities.Post) error
	FindByID(ctx context.Context, postID valueobjects.PostID) (*entities.Post, error)
	FindByCommunity(ctx context.Context, communityID valueobjects.CommunityID, limit, offset *int) ([]*entities.Post, error)
	FindByAuthorAndCommunity(ctx context.Context, authorID valueobjects.AuthorID, communityID valueobjects.CommunityID) (*entities.Post, error)
	FindByCommunities(ctx context.Context, communityIDs []valueobjects.CommunityID, postType *valueobjects.PostType, limit, offset *int) ([]*entities.Post, error)
	Delete(ctx context.Context, postID valueobjects.PostID) error
}
