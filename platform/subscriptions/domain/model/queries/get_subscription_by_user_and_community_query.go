package queries

import (
	"errors"

	"Gommunity/platform/subscriptions/domain/model/valueobjects"
)

// GetSubscriptionByUserAndCommunityQuery represents a request to get a specific subscription by user and community
type GetSubscriptionByUserAndCommunityQuery struct {
	userID      valueobjects.UserID
	communityID valueobjects.CommunityID
}

func NewGetSubscriptionByUserAndCommunityQuery(
	userID valueobjects.UserID,
	communityID valueobjects.CommunityID,
) (GetSubscriptionByUserAndCommunityQuery, error) {
	if userID.IsZero() {
		return GetSubscriptionByUserAndCommunityQuery{}, errors.New("user ID cannot be zero")
	}
	if communityID.IsZero() {
		return GetSubscriptionByUserAndCommunityQuery{}, errors.New("community ID cannot be empty")
	}

	return GetSubscriptionByUserAndCommunityQuery{
		userID:      userID,
		communityID: communityID,
	}, nil
}

func (q GetSubscriptionByUserAndCommunityQuery) UserID() valueobjects.UserID {
	return q.userID
}

func (q GetSubscriptionByUserAndCommunityQuery) CommunityID() valueobjects.CommunityID {
	return q.communityID
}
