package commands

import (
	"errors"

	"Gommunity/internal/community/communities/domain/model/valueobjects"
)

type CreateCommunityCommand struct {
	ownerID     valueobjects.OwnerID
	name        valueobjects.CommunityName
	description valueobjects.Description
	iconURL     *string
	bannerURL   *string
}

func NewCreateCommunityCommand(
	ownerID valueobjects.OwnerID,
	name valueobjects.CommunityName,
	description valueobjects.Description,
	iconURL *string,
	bannerURL *string,
) (CreateCommunityCommand, error) {
	if ownerID.IsZero() {
		return CreateCommunityCommand{}, errors.New("ownerID cannot be empty")
	}

	if name.IsZero() {
		return CreateCommunityCommand{}, errors.New("name cannot be empty")
	}

	if description.IsZero() {
		return CreateCommunityCommand{}, errors.New("description cannot be empty")
	}

	return CreateCommunityCommand{
		ownerID:     ownerID,
		name:        name,
		description: description,
		iconURL:     iconURL,
		bannerURL:   bannerURL,
	}, nil
}

func (c CreateCommunityCommand) OwnerID() valueobjects.OwnerID {
	return c.ownerID
}

func (c CreateCommunityCommand) Name() valueobjects.CommunityName {
	return c.name
}

func (c CreateCommunityCommand) Description() valueobjects.Description {
	return c.description
}

func (c CreateCommunityCommand) IconURL() *string {
	return c.iconURL
}

func (c CreateCommunityCommand) BannerURL() *string {
	return c.bannerURL
}
