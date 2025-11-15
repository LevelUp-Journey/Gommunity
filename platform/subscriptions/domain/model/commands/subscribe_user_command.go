package commands

import (
	"errors"

	"Gommunity/platform/subscriptions/domain/model/valueobjects"
)

// SubscribeUserCommand represents the intention to subscribe a user to a community
// In public communities: user subscribes themselves
// In private communities: teacher/admin adds user and assigns role
type SubscribeUserCommand struct {
	userID      valueobjects.UserID
	communityID valueobjects.CommunityID
	role        valueobjects.CommunityRole
	requestedBy valueobjects.UserID // Who is making the request (can be same as userID for self-subscription)
}

func NewSubscribeUserCommand(
	userID valueobjects.UserID,
	communityID valueobjects.CommunityID,
	role valueobjects.CommunityRole,
	requestedBy valueobjects.UserID,
) (SubscribeUserCommand, error) {
	if userID.IsZero() {
		return SubscribeUserCommand{}, errors.New("user ID cannot be zero")
	}
	if communityID.IsZero() {
		return SubscribeUserCommand{}, errors.New("community ID cannot be empty")
	}
	if role.IsZero() {
		return SubscribeUserCommand{}, errors.New("role cannot be empty")
	}
	if requestedBy.IsZero() {
		return SubscribeUserCommand{}, errors.New("requestedBy ID cannot be zero")
	}

	return SubscribeUserCommand{
		userID:      userID,
		communityID: communityID,
		role:        role,
		requestedBy: requestedBy,
	}, nil
}

func (c SubscribeUserCommand) UserID() valueobjects.UserID {
	return c.userID
}

func (c SubscribeUserCommand) CommunityID() valueobjects.CommunityID {
	return c.communityID
}

func (c SubscribeUserCommand) Role() valueobjects.CommunityRole {
	return c.role
}

func (c SubscribeUserCommand) RequestedBy() valueobjects.UserID {
	return c.requestedBy
}

func (c SubscribeUserCommand) IsSelfSubscription() bool {
	return c.userID.Equals(c.requestedBy)
}
