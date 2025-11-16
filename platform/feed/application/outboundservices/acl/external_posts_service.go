package acl

import (
	"Gommunity/platform/feed/domain/model/entities"
	"Gommunity/platform/feed/domain/model/valueobjects"
	"context"

	posts_acl "Gommunity/platform/posts/interfaces/acl"
)

// ExternalPostsService provides ACL access to posts context
type ExternalPostsService struct {
	postsFacade posts_acl.PostsFacade
}

func NewExternalPostsService(postsFacade posts_acl.PostsFacade) *ExternalPostsService {
	return &ExternalPostsService{
		postsFacade: postsFacade,
	}
}

// GetAnnouncementsForCommunities retrieves announcements from multiple communities
func (s *ExternalPostsService) GetAnnouncementsForCommunities(ctx context.Context, communityIDs []string, limit, offset *int) ([]*entities.FeedItem, error) {
	postsData, err := s.postsFacade.GetAnnouncementsByCommunities(ctx, communityIDs, limit, offset)
	if err != nil {
		return nil, err
	}

	feedItems := make([]*entities.FeedItem, len(postsData))
	for i, postData := range postsData {
		postID, _ := valueobjects.NewPostID(postData.PostID)
		communityID, _ := valueobjects.NewCommunityID(postData.CommunityID)

		feedItems[i] = entities.NewFeedItem(
			postID,
			communityID,
			postData.AuthorID,
			postData.Content,
			postData.MessageType,
			postData.CreatedAt,
			postData.UpdatedAt,
		)
	}

	return feedItems, nil
}
