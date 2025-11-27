package acl

import (
	"context"

	communities_acl "Gommunity/platform/community/interfaces/acl"
	"Gommunity/platform/posts/domain/model/valueobjects"
)

// ExternalCommunitiesService provides access to the Communities bounded context.
type ExternalCommunitiesService struct {
	communitiesFacade communities_acl.CommunitiesFacade
}

// NewExternalCommunitiesService builds a new ExternalCommunitiesService.
func NewExternalCommunitiesService(communitiesFacade communities_acl.CommunitiesFacade) *ExternalCommunitiesService {
	return &ExternalCommunitiesService{
		communitiesFacade: communitiesFacade,
	}
}

// ValidateCommunityExists checks whether the community exists.
func (s *ExternalCommunitiesService) ValidateCommunityExists(ctx context.Context, communityID valueobjects.CommunityID) (bool, error) {
	return s.communitiesFacade.ValidateCommunityExists(ctx, communityID.Value())
}

// ValidateUserIsOwner verifies whether the provided author owns the community.
func (s *ExternalCommunitiesService) ValidateUserIsOwner(ctx context.Context, communityID valueobjects.CommunityID, authorID valueobjects.AuthorID) (bool, error) {
	return s.communitiesFacade.ValidateUserIsOwner(ctx, communityID.Value(), authorID.Value())
}
