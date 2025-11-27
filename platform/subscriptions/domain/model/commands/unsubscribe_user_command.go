package commands

import (
	"errors"

	"Gommunity/platform/subscriptions/domain/model/valueobjects"
)

// UnsubscribeUserCommand represents the intention to remove a user's subscription from a community
// In public communities: user can unsubscribe themselves
// In private communities: user can unsubscribe themselves OR teacher/admin can remove them
type UnsubscribeUserCommand struct {
	userID      valueobjects.UserID
	communityID valueobjects.CommunityID
	requestedBy valueobjects.UserID // Who is making the request (can be same as userID for self-unsubscription)
}

func NewUnsubscribeUserCommand(
	userID valueobjects.UserID,
	communityID valueobjects.CommunityID,
	requestedBy valueobjects.UserID,
) (UnsubscribeUserCommand, error) {
	if userID.IsZero() {
		return UnsubscribeUserCommand{}, errors.New("user ID cannot be zero")
	}
	if communityID.IsZero() {
		return UnsubscribeUserCommand{}, errors.New("community ID cannot be empty")
	}
	if requestedBy.IsZero() {
		return UnsubscribeUserCommand{}, errors.New("requestedBy ID cannot be zero")
	}

	return UnsubscribeUserCommand{
		userID:      userID,
		communityID: communityID,
		requestedBy: requestedBy,
	}, nil
}

func (c UnsubscribeUserCommand) UserID() valueobjects.UserID {
	return c.userID
}

func (c UnsubscribeUserCommand) CommunityID() valueobjects.CommunityID {
	return c.communityID
}

func (c UnsubscribeUserCommand) RequestedBy() valueobjects.UserID {
	return c.requestedBy
}

func (c UnsubscribeUserCommand) IsSelfUnsubscription() bool {
	return c.userID.Equals(c.requestedBy)
}
