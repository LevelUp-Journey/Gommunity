package queries

import (
	"errors"

	"Gommunity/platform/subscriptions/domain/model/valueobjects"
)

// GetAllSubscriptionsByCommunityQuery represents a request to get all subscriptions for a specific community
type GetAllSubscriptionsByCommunityQuery struct {
	communityID valueobjects.CommunityID
	limit       *int
	offset      *int
}

func NewGetAllSubscriptionsByCommunityQuery(communityID valueobjects.CommunityID) (GetAllSubscriptionsByCommunityQuery, error) {
	if communityID.IsZero() {
		return GetAllSubscriptionsByCommunityQuery{}, errors.New("community ID cannot be empty")
	}

	return GetAllSubscriptionsByCommunityQuery{
		communityID: communityID,
	}, nil
}

func (q GetAllSubscriptionsByCommunityQuery) WithPagination(limit, offset int) GetAllSubscriptionsByCommunityQuery {
	q.limit = &limit
	q.offset = &offset
	return q
}

func (q GetAllSubscriptionsByCommunityQuery) CommunityID() valueobjects.CommunityID {
	return q.communityID
}

func (q GetAllSubscriptionsByCommunityQuery) Limit() *int {
	return q.limit
}

func (q GetAllSubscriptionsByCommunityQuery) Offset() *int {
	return q.offset
}
