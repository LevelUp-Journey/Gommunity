package repositories

import (
	"context"

	"Gommunity/platform/reactions/domain/model/entities"
	"Gommunity/platform/reactions/domain/model/valueobjects"
)

// ReactionRepository defines persistence operations for reaction aggregates.
type ReactionRepository interface {
	Save(ctx context.Context, reaction *entities.Reaction) error
	Update(ctx context.Context, reaction *entities.Reaction) error
	FindByID(ctx context.Context, reactionID valueobjects.ReactionID) (*entities.Reaction, error)
	FindByPostAndUser(ctx context.Context, postID valueobjects.PostID, userID valueobjects.UserID) (*entities.Reaction, error)
	FindByPost(ctx context.Context, postID valueobjects.PostID) ([]*entities.Reaction, error)
	CountByPost(ctx context.Context, postID valueobjects.PostID) (map[string]int, error)
	Delete(ctx context.Context, reactionID valueobjects.ReactionID) error
	DeleteByPostAndUser(ctx context.Context, postID valueobjects.PostID, userID valueobjects.UserID) error
}
