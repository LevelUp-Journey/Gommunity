package services

import (
	"context"

	"Gommunity/platform/posts/domain/model/commands"
	"Gommunity/platform/posts/domain/model/valueobjects"
)

// PostCommandService defines command handling behavior for posts.
type PostCommandService interface {
	HandlePublish(ctx context.Context, cmd commands.CreatePostCommand) (*valueobjects.PostID, error)
	HandleDelete(ctx context.Context, cmd commands.DeletePostCommand) error
}
