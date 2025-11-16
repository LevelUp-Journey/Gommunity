package acl

import (
	"context"
	"errors"

	"Gommunity/platform/community/domain/model/valueobjects"
	"Gommunity/platform/community/domain/repositories"
	"Gommunity/platform/community/interfaces/acl"
)

type communitiesFacadeImpl struct {
	communityRepository repositories.CommunityRepository
}

// NewCommunitiesFacade creates a new CommunitiesFacade implementation
func NewCommunitiesFacade(
	communityRepository repositories.CommunityRepository,
) acl.CommunitiesFacade {
	return &communitiesFacadeImpl{
		communityRepository: communityRepository,
	}
}

// ValidateCommunityExists checks if a community exists by ID
func (f *communitiesFacadeImpl) ValidateCommunityExists(ctx context.Context, communityID string) (bool, error) {
	communityIDVO, err := valueobjects.NewCommunityID(communityID)
	if err != nil {
		return false, err
	}

	return f.communityRepository.ExistsByID(ctx, communityIDVO)
}

// IsCommunityPrivate checks if a community is private
func (f *communitiesFacadeImpl) IsCommunityPrivate(ctx context.Context, communityID string) (bool, error) {
	communityIDVO, err := valueobjects.NewCommunityID(communityID)
	if err != nil {
		return false, err
	}

	community, err := f.communityRepository.FindByID(ctx, communityIDVO)
	if err != nil {
		return false, err
	}

	if community == nil {
		return false, errors.New("community not found")
	}

	return community.IsPrivate(), nil
}

// GetCommunityOwnerID retrieves the owner ID of a community (as string UUID)
func (f *communitiesFacadeImpl) GetCommunityOwnerID(ctx context.Context, communityID string) (string, error) {
	communityIDVO, err := valueobjects.NewCommunityID(communityID)
	if err != nil {
		return "", err
	}

	community, err := f.communityRepository.FindByID(ctx, communityIDVO)
	if err != nil {
		return "", err
	}

	if community == nil {
		return "", errors.New("community not found")
	}

	return community.OwnerID().Value(), nil
}

// ValidateUserIsOwner checks if a user is the owner of a community
func (f *communitiesFacadeImpl) ValidateUserIsOwner(ctx context.Context, communityID string, ownerID string) (bool, error) {
	actualOwnerID, err := f.GetCommunityOwnerID(ctx, communityID)
	if err != nil {
		return false, err
	}

	return actualOwnerID == ownerID, nil
}
