package services

import (
	"context"

	"Gommunity/platform/reactions/domain/model/commands"
	"Gommunity/platform/reactions/domain/model/valueobjects"
)

// ReactionCommandService defines command operations for reactions.
type ReactionCommandService interface {
	HandleAdd(ctx context.Context, cmd commands.AddReactionCommand) (*valueobjects.ReactionID, error)
	HandleRemove(ctx context.Context, cmd commands.RemoveReactionCommand) error
}
