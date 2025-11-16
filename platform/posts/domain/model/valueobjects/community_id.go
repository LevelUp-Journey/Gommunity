package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

// CommunityID represents the identifier of a community within the posts context.
type CommunityID struct {
	value string `json:"value" bson:"community_id"`
}

// NewCommunityID validates and creates a CommunityID.
func NewCommunityID(value string) (CommunityID, error) {
	if value == "" {
		return CommunityID{}, errors.New("community ID cannot be empty")
	}
	if _, err := uuid.Parse(value); err != nil {
		return CommunityID{}, errors.New("community ID must be a valid UUID")
	}
	return CommunityID{value: value}, nil
}

// Value returns the identifier.
func (c CommunityID) Value() string {
	return c.value
}

// String returns the identifier as string.
func (c CommunityID) String() string {
	return c.value
}

// IsZero indicates if the identifier is empty.
func (c CommunityID) IsZero() bool {
	return c.value == ""
}
