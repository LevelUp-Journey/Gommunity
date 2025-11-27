package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

type CommunityID struct {
	value string `json:"value" bson:"community_id"`
}

func NewCommunityID(value string) (CommunityID, error) {
	if value == "" {
		return CommunityID{}, errors.New("community ID cannot be empty")
	}

	// Validate UUID format
	if _, err := uuid.Parse(value); err != nil {
		return CommunityID{}, errors.New("community ID must be a valid UUID")
	}

	return CommunityID{value: value}, nil
}

func (c CommunityID) Value() string {
	return c.value
}

func (c CommunityID) String() string {
	return c.value
}

func (c CommunityID) IsZero() bool {
	return c.value == ""
}
