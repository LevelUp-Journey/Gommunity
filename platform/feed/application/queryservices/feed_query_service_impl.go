package queryservices

import (
	"context"
	"log"

	"Gommunity/platform/feed/application/outboundservices/acl"
	"Gommunity/platform/feed/domain/model/entities"
	"Gommunity/platform/feed/domain/model/queries"
	"Gommunity/platform/feed/domain/services"
)

type feedQueryServiceImpl struct {
	subscriptionsService *acl.ExternalSubscriptionsService
	postsService         *acl.ExternalPostsService
}

func NewFeedQueryService(
	subscriptionsService *acl.ExternalSubscriptionsService,
	postsService *acl.ExternalPostsService,
) services.FeedQueryService {
	return &feedQueryServiceImpl{
		subscriptionsService: subscriptionsService,
		postsService:         postsService,
	}
}

func (s *feedQueryServiceImpl) Handle(ctx context.Context, query queries.GetUserFeedQuery) ([]*entities.FeedItem, error) {
	log.Printf("Getting feed for user: %s", query.UserID().Value())

	// Step 1: Get all communities the user is subscribed to
	communityIDs, err := s.subscriptionsService.GetUserCommunities(ctx, query.UserID().Value())
	if err != nil {
		log.Printf("Error getting user communities: %v", err)
		return nil, err
	}

	log.Printf("User is subscribed to %d communities", len(communityIDs))

	// Step 2: If user is not subscribed to any community, return empty feed
	if len(communityIDs) == 0 {
		return []*entities.FeedItem{}, nil
	}

	// Step 3: Get announcements from those communities
	feedItems, err := s.postsService.GetAnnouncementsForCommunities(ctx, communityIDs, query.Limit(), query.Offset())
	if err != nil {
		log.Printf("Error getting announcements: %v", err)
		return nil, err
	}

	log.Printf("Found %d feed items for user", len(feedItems))
	return feedItems, nil
}
