package services

import (
	"context"

	"Gommunity/platform/reactions/domain/model/entities"
	"Gommunity/platform/reactions/domain/model/queries"
)

// ReactionSummary represents aggregated reaction counts by type.
type ReactionSummary struct {
	TotalCount int
	Counts     map[string]int // key: reaction type, value: count
}

// ReactionQueryService defines query operations for reactions.
type ReactionQueryService interface {
	HandleGetByPost(ctx context.Context, query queries.GetReactionsByPostQuery) ([]*entities.Reaction, error)
	HandleGetCountByPost(ctx context.Context, query queries.GetReactionCountByPostQuery) (*ReactionSummary, error)
	HandleGetUserReactionOnPost(ctx context.Context, query queries.GetUserReactionOnPostQuery) (*entities.Reaction, error)
}
