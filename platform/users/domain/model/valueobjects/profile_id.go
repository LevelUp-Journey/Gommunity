package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

type ProfileID struct {
	value string `json:"value" bson:"profile_id"`
}

func NewProfileID(value string) (ProfileID, error) {
	if value == "" {
		return ProfileID{}, errors.New("profile ID cannot be empty")
	}

	// Validate UUID format
	if _, err := uuid.Parse(value); err != nil {
		return ProfileID{}, errors.New("profile ID must be a valid UUID")
	}

	return ProfileID{value: value}, nil
}

func (p ProfileID) Value() string {
	return p.value
}

func (p ProfileID) String() string {
	return p.value
}

func (p ProfileID) IsZero() bool {
	return p.value == ""
}
