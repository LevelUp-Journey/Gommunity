package commands

import (
	"errors"

	"Gommunity/internal/community/communities/domain/model/valueobjects"
)

type UpdateCommunityPrivacyCommand struct {
	communityID valueobjects.CommunityID
	isPrivate   bool
}

func NewUpdateCommunityPrivacyCommand(
	communityID valueobjects.CommunityID,
	isPrivate bool,
) (UpdateCommunityPrivacyCommand, error) {
	if communityID.IsZero() {
		return UpdateCommunityPrivacyCommand{}, errors.New("communityID cannot be empty")
	}

	return UpdateCommunityPrivacyCommand{
		communityID: communityID,
		isPrivate:   isPrivate,
	}, nil
}

func (c UpdateCommunityPrivacyCommand) CommunityID() valueobjects.CommunityID {
	return c.communityID
}

func (c UpdateCommunityPrivacyCommand) IsPrivate() bool {
	return c.isPrivate
}
