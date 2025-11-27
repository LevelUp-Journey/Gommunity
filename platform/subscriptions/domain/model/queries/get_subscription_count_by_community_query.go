package queries

import (
	"errors"

	"Gommunity/platform/subscriptions/domain/model/valueobjects"
)

// GetSubscriptionCountByCommunityQuery represents a request to get the total number of subscriptions for a community
type GetSubscriptionCountByCommunityQuery struct {
	communityID valueobjects.CommunityID
}

func NewGetSubscriptionCountByCommunityQuery(communityID valueobjects.CommunityID) (GetSubscriptionCountByCommunityQuery, error) {
	if communityID.IsZero() {
		return GetSubscriptionCountByCommunityQuery{}, errors.New("community ID cannot be empty")
	}

	return GetSubscriptionCountByCommunityQuery{
		communityID: communityID,
	}, nil
}

func (q GetSubscriptionCountByCommunityQuery) CommunityID() valueobjects.CommunityID {
	return q.communityID
}
