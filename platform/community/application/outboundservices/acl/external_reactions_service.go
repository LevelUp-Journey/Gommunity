package acl

import (
	"context"

	reactions_repos "Gommunity/platform/reactions/domain/repositories"
	posts_vo "Gommunity/platform/posts/domain/model/valueobjects"
)

// ExternalReactionsService provides deletion support for reactions tied to posts.
type ExternalReactionsService struct {
	reactionRepository reactions_repos.ReactionRepository
}

func NewExternalReactionsService(reactionRepository reactions_repos.ReactionRepository) *ExternalReactionsService {
	return &ExternalReactionsService{
		reactionRepository: reactionRepository,
	}
}

// DeleteReactionsByPostIDs deletes reactions for the provided posts.
func (s *ExternalReactionsService) DeleteReactionsByPostIDs(ctx context.Context, postIDs []posts_vo.PostID) error {
	return s.reactionRepository.DeleteByPostIDs(ctx, postIDs)
}
