package queries

import (
	"Gommunity/platform/feed/domain/model/valueobjects"
	"errors"
)

type GetUserFeedQuery struct {
	userID valueobjects.UserID
	limit  *int
	offset *int
}

func NewGetUserFeedQuery(userID valueobjects.UserID) (GetUserFeedQuery, error) {
	if userID.IsEmpty() {
		return GetUserFeedQuery{}, errors.New("user ID cannot be empty")
	}
	return GetUserFeedQuery{userID: userID}, nil
}

func (q GetUserFeedQuery) WithPagination(limit, offset int) GetUserFeedQuery {
	q.limit = &limit
	q.offset = &offset
	return q
}

func (q GetUserFeedQuery) UserID() valueobjects.UserID {
	return q.userID
}

func (q GetUserFeedQuery) Limit() *int {
	return q.limit
}

func (q GetUserFeedQuery) Offset() *int {
	return q.offset
}
