package queries

import (
	"errors"

	"Gommunity/platform/subscriptions/domain/model/valueobjects"
)

// GetCommunitiesByUserQuery represents a request to get all communities a user is subscribed to
type GetCommunitiesByUserQuery struct {
	userID valueobjects.UserID
}

func NewGetCommunitiesByUserQuery(
	userID valueobjects.UserID,
) (GetCommunitiesByUserQuery, error) {
	if userID.IsZero() {
		return GetCommunitiesByUserQuery{}, errors.New("user ID cannot be zero")
	}

	return GetCommunitiesByUserQuery{
		userID: userID,
	}, nil
}

func (q GetCommunitiesByUserQuery) UserID() valueobjects.UserID {
	return q.userID
}
