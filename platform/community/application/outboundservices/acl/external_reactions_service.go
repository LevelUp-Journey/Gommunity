package acl

import (
	"context"

	reactions_repos "Gommunity/platform/reactions/domain/repositories"
	reactions_vo "Gommunity/platform/reactions/domain/model/valueobjects"
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
	if len(postIDs) == 0 {
		return nil
	}

	var reactionsPostIDs []reactions_vo.PostID
	for _, id := range postIDs {
		converted, err := reactions_vo.NewPostID(id.Value())
		if err != nil {
			return err
		}
		reactionsPostIDs = append(reactionsPostIDs, converted)
	}

	return s.reactionRepository.DeleteByPostIDs(ctx, reactionsPostIDs)
}
