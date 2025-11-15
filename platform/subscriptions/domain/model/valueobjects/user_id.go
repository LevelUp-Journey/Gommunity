package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

// UserID represents a User ID from the Users bounded context (UUID string)
type UserID struct {
	value string `json:"value" bson:"user_id"`
}

func NewUserID(value string) (UserID, error) {
	if value == "" {
		return UserID{}, errors.New("user ID cannot be empty")
	}

	// Validate UUID format
	if _, err := uuid.Parse(value); err != nil {
		return UserID{}, errors.New("user ID must be a valid UUID")
	}

	return UserID{value: value}, nil
}

func (u UserID) Value() string {
	return u.value
}

func (u UserID) String() string {
	return u.value
}

func (u UserID) IsZero() bool {
	return u.value == ""
}

func (u UserID) Equals(other UserID) bool {
	return u.value == other.value
}
