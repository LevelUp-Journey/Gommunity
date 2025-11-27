package valueobjects

import (
	"errors"
)

type FeedID struct {
	value string `json:"value"`
}

func NewFeedID(value string) (FeedID, error) {
	if value == "" {
		return FeedID{}, errors.New("feed ID cannot be empty")
	}
	return FeedID{value: value}, nil
}

func (f FeedID) Value() string {
	return f.value
}

func (f FeedID) String() string {
	return f.value
}

func (f FeedID) IsEmpty() bool {
	return f.value == ""
}

func (f FeedID) Equals(other FeedID) bool {
	return f.value == other.value
}
