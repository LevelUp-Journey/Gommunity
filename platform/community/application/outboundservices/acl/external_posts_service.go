package acl

import (
	"context"

	posts_repos "Gommunity/platform/posts/domain/repositories"
	posts_vo "Gommunity/platform/posts/domain/model/valueobjects"
	community_vo "Gommunity/platform/community/domain/model/valueobjects"
)

// ExternalPostsService provides minimal operations against Posts BC needed for cascades.
type ExternalPostsService struct {
	postRepository posts_repos.PostRepository
}

func NewExternalPostsService(postRepository posts_repos.PostRepository) *ExternalPostsService {
	return &ExternalPostsService{
		postRepository: postRepository,
	}
}

// GetPostIDsByCommunity returns post IDs for the community.
func (s *ExternalPostsService) GetPostIDsByCommunity(ctx context.Context, communityID community_vo.CommunityID) ([]posts_vo.PostID, error) {
	postCommunityID, err := posts_vo.NewCommunityID(communityID.Value())
	if err != nil {
		return nil, err
	}
	return s.postRepository.FindPostIDsByCommunity(ctx, postCommunityID)
}

// DeletePostsByCommunity deletes all posts for the community.
func (s *ExternalPostsService) DeletePostsByCommunity(ctx context.Context, communityID community_vo.CommunityID) error {
	postCommunityID, err := posts_vo.NewCommunityID(communityID.Value())
	if err != nil {
		return err
	}
	return s.postRepository.DeleteByCommunity(ctx, postCommunityID)
}
