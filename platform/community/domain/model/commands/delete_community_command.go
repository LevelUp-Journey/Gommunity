package commands

import (
	"errors"

	"Gommunity/platform/community/domain/model/valueobjects"
)

type DeleteCommunityCommand struct {
	communityID valueobjects.CommunityID
	ownerID     valueobjects.OwnerID
}

func NewDeleteCommunityCommand(
	communityID valueobjects.CommunityID,
	ownerID valueobjects.OwnerID,
) (DeleteCommunityCommand, error) {
	if communityID.IsZero() {
		return DeleteCommunityCommand{}, errors.New("communityID cannot be empty")
	}

	if ownerID.IsZero() {
		return DeleteCommunityCommand{}, errors.New("ownerID cannot be empty")
	}

	return DeleteCommunityCommand{
		communityID: communityID,
		ownerID:     ownerID,
	}, nil
}

func (c DeleteCommunityCommand) CommunityID() valueobjects.CommunityID {
	return c.communityID
}

func (c DeleteCommunityCommand) OwnerID() valueobjects.OwnerID {
	return c.ownerID
}
