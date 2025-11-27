package acl


import (
	"context"
	"fmt"

	"Gommunity/platform/reactions/domain/model/valueobjects"
	posts_acl "Gommunity/platform/posts/interfaces/acl"
)

// ExternalPostsService validates posts from the posts bounded context.
type ExternalPostsService struct {
	postsFacade posts_acl.PostsFacade
}

// NewExternalPostsService constructs the external posts service.
func NewExternalPostsService(postsFacade posts_acl.PostsFacade) *ExternalPostsService {
	return &ExternalPostsService{
		postsFacade: postsFacade,
	}
}

// ValidatePostExists checks if a post exists.
func (s *ExternalPostsService) ValidatePostExists(ctx context.Context, postID valueobjects.PostID) (bool, error) {
	exists, err := s.postsFacade.PostExists(ctx, postID.Value())
	if err != nil {
		return false, fmt.Errorf("failed to validate post existence: %w", err)
	}
	return exists, nil
}
