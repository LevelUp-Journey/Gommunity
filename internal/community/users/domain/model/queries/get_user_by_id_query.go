package queries

import (
	"Gommunity/internal/community/users/domain/model/valueobjects"
	"errors"
)

type GetUserByIDQuery struct {
	userID valueobjects.UserID
}

func NewGetUserByIDQuery(userID valueobjects.UserID) (GetUserByIDQuery, error) {
	if userID.IsZero() {
		return GetUserByIDQuery{}, errors.New("userID cannot be empty")
	}

	return GetUserByIDQuery{userID: userID}, nil
}

func (q GetUserByIDQuery) UserID() valueobjects.UserID {
	return q.userID
}
