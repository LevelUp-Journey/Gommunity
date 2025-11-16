package acl

import (
	"context"

	communities_acl "Gommunity/platform/community/interfaces/acl"
	"Gommunity/platform/subscriptions/domain/model/valueobjects"
)

// ExternalCommunitiesService provides ACL implementation for accessing Community bounded context
type ExternalCommunitiesService struct {
	communitiesFacade communities_acl.CommunitiesFacade
}

// NewExternalCommunitiesService creates a new ExternalCommunitiesService
func NewExternalCommunitiesService(communitiesFacade communities_acl.CommunitiesFacade) *ExternalCommunitiesService {
	return &ExternalCommunitiesService{
		communitiesFacade: communitiesFacade,
	}
}

// ValidateCommunityExists checks if a community exists in the Communities BC
func (s *ExternalCommunitiesService) ValidateCommunityExists(ctx context.Context, communityID valueobjects.CommunityID) (bool, error) {
	return s.communitiesFacade.ValidateCommunityExists(ctx, communityID.Value())
}

// IsCommunityPrivate checks if a community is private
func (s *ExternalCommunitiesService) IsCommunityPrivate(ctx context.Context, communityID valueobjects.CommunityID) (bool, error) {
	return s.communitiesFacade.IsCommunityPrivate(ctx, communityID.Value())
}

// GetCommunityOwnerID retrieves the owner ID of a community (as string UUID)
func (s *ExternalCommunitiesService) GetCommunityOwnerID(ctx context.Context, communityID valueobjects.CommunityID) (string, error) {
	return s.communitiesFacade.GetCommunityOwnerID(ctx, communityID.Value())
}

// ValidateUserIsOwner checks if a user is the owner of a community
func (s *ExternalCommunitiesService) ValidateUserIsOwner(ctx context.Context, communityID valueobjects.CommunityID, ownerID string) (bool, error) {
	return s.communitiesFacade.ValidateUserIsOwner(ctx, communityID.Value(), ownerID)
}
