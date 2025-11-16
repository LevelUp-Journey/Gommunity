package valueobjects

import (
	"errors"
)

type CommunityID struct {
	value string `json:"value"`
}

func NewCommunityID(value string) (CommunityID, error) {
	if value == "" {
		return CommunityID{}, errors.New("community ID cannot be empty")
	}
	return CommunityID{value: value}, nil
}

func (c CommunityID) Value() string {
	return c.value
}

func (c CommunityID) String() string {
	return c.value
}

func (c CommunityID) IsEmpty() bool {
	return c.value == ""
}

func (c CommunityID) Equals(other CommunityID) bool {
	return c.value == other.value
}
