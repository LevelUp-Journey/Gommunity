package queries

import (
	"Gommunity/platform/users/domain/model/valueobjects"
	"errors"
)

type GetUserByUsernameQuery struct {
	username valueobjects.Username
}

func NewGetUserByUsernameQuery(username valueobjects.Username) (GetUserByUsernameQuery, error) {
	if username.IsZero() {
		return GetUserByUsernameQuery{}, errors.New("username cannot be empty")
	}

	return GetUserByUsernameQuery{username: username}, nil
}

func (q GetUserByUsernameQuery) Username() valueobjects.Username {
	return q.username
}
