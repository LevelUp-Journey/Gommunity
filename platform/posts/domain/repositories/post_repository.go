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
	FindByCommunities(ctx context.Context, communityIDs []valueobjects.CommunityID, limit, offset *int) ([]*entities.Post, error)
	Delete(ctx context.Context, postID valueobjects.PostID) error

	// FindPostIDsByCommunity returns only post IDs for a community (for cascade deletion)
	FindPostIDsByCommunity(ctx context.Context, communityID valueobjects.CommunityID) ([]valueobjects.PostID, error)

	// DeleteByCommunity removes all posts for a community
	DeleteByCommunity(ctx context.Context, communityID valueobjects.CommunityID) error
}
