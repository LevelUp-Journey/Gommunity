package queryservices

import (
	"context"
	"fmt"

	"Gommunity/platform/reactions/domain/model/entities"
	"Gommunity/platform/reactions/domain/model/queries"
	"Gommunity/platform/reactions/domain/repositories"
	"Gommunity/platform/reactions/domain/services"
)

type reactionQueryServiceImpl struct {
	reactionRepository repositories.ReactionRepository
}

// NewReactionQueryService constructs the reactions query service implementation.
func NewReactionQueryService(
	reactionRepository repositories.ReactionRepository,
) services.ReactionQueryService {
	return &reactionQueryServiceImpl{
		reactionRepository: reactionRepository,
	}
}

// HandleGetByPost retrieves all reactions for a post.
func (s *reactionQueryServiceImpl) HandleGetByPost(ctx context.Context, query queries.GetReactionsByPostQuery) ([]*entities.Reaction, error) {
	reactions, err := s.reactionRepository.FindByPost(ctx, query.PostID())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve reactions: %w", err)
	}
	return reactions, nil
}

// HandleGetCountByPost retrieves reaction count summary for a post.
func (s *reactionQueryServiceImpl) HandleGetCountByPost(ctx context.Context, query queries.GetReactionCountByPostQuery) (*services.ReactionSummary, error) {
	counts, err := s.reactionRepository.CountByPost(ctx, query.PostID())
	if err != nil {
		return nil, fmt.Errorf("failed to count reactions: %w", err)
	}

	totalCount := 0
	for _, count := range counts {
		totalCount += count
	}

	return &services.ReactionSummary{
		TotalCount: totalCount,
		Counts:     counts,
	}, nil
}

// HandleGetUserReactionOnPost retrieves a user's reaction on a specific post.
func (s *reactionQueryServiceImpl) HandleGetUserReactionOnPost(ctx context.Context, query queries.GetUserReactionOnPostQuery) (*entities.Reaction, error) {
	reaction, err := s.reactionRepository.FindByPostAndUser(ctx, query.PostID(), query.UserID())
	if err != nil {
		return nil, fmt.Errorf("failed to find user reaction: %w", err)
	}
	return reaction, nil
}
