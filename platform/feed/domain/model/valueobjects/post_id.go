package valueobjects

import (
	"errors"
)

type PostID struct {
	value string `json:"value"`
}

func NewPostID(value string) (PostID, error) {
	if value == "" {
		return PostID{}, errors.New("post ID cannot be empty")
	}
	return PostID{value: value}, nil
}

func (p PostID) Value() string {
	return p.value
}

func (p PostID) String() string {
	return p.value
}

func (p PostID) IsEmpty() bool {
	return p.value == ""
}

func (p PostID) Equals(other PostID) bool {
	return p.value == other.value
}
