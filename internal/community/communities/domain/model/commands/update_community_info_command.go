package commands

import (
	"errors"

	"Gommunity/internal/community/communities/domain/model/valueobjects"
)

type UpdateCommunityInfoCommand struct {
	communityID valueobjects.CommunityID
	name        valueobjects.CommunityName
	description valueobjects.Description
	iconURL     *string
	bannerURL   *string
}

func NewUpdateCommunityInfoCommand(
	communityID valueobjects.CommunityID,
	name valueobjects.CommunityName,
	description valueobjects.Description,
	iconURL *string,
	bannerURL *string,
) (UpdateCommunityInfoCommand, error) {
	if communityID.IsZero() {
		return UpdateCommunityInfoCommand{}, errors.New("communityID cannot be empty")
	}

	if name.IsZero() {
		return UpdateCommunityInfoCommand{}, errors.New("name cannot be empty")
	}

	if description.IsZero() {
		return UpdateCommunityInfoCommand{}, errors.New("description cannot be empty")
	}

	return UpdateCommunityInfoCommand{
		communityID: communityID,
		name:        name,
		description: description,
		iconURL:     iconURL,
		bannerURL:   bannerURL,
	}, nil
}

func (c UpdateCommunityInfoCommand) CommunityID() valueobjects.CommunityID {
	return c.communityID
}

func (c UpdateCommunityInfoCommand) Name() valueobjects.CommunityName {
	return c.name
}

func (c UpdateCommunityInfoCommand) Description() valueobjects.Description {
	return c.description
}

func (c UpdateCommunityInfoCommand) IconURL() *string {
	return c.iconURL
}

func (c UpdateCommunityInfoCommand) BannerURL() *string {
	return c.bannerURL
}
