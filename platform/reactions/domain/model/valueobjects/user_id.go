package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

// UserID represents a reference to a user from the users bounded context.
type UserID struct {
	value string `json:"value" bson:"user_id"`
}

// NewUserID creates and validates a UserID.
func NewUserID(value string) (UserID, error) {
	if value == "" {
		return UserID{}, errors.New("user ID cannot be empty")
	}
	if _, err := uuid.Parse(value); err != nil {
		return UserID{}, errors.New("user ID must be a valid UUID")
	}
	return UserID{value: value}, nil
}

// Value returns the string value of the UserID.
func (u UserID) Value() string {
	return u.value
}

// String returns the string representation of the UserID.
func (u UserID) String() string {
	return u.value
}

// IsZero checks if the UserID is unset.
func (u UserID) IsZero() bool {
	return u.value == ""
}
