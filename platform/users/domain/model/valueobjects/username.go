package valueobjects

import (
	"errors"
	"regexp"
)

type Username struct {
	value string `json:"value" bson:"username"`
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)

func NewUsername(value string) (Username, error) {
	if value == "" {
		return Username{}, errors.New("username cannot be empty")
	}

	if !usernameRegex.MatchString(value) {
		return Username{}, errors.New("username must be 3-30 characters and contain only letters, numbers, and underscores")
	}

	return Username{value: value}, nil
}

func (u Username) Value() string {
	return u.value
}

func (u Username) String() string {
	return u.value
}

func (u Username) IsZero() bool {
	return u.value == ""
}
