package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

// AuthorID represents the identifier of the user who interacts with posts.
type AuthorID struct {
	value string `json:"value" bson:"author_id"`
}

// NewAuthorID validates and creates a new AuthorID.
func NewAuthorID(value string) (AuthorID, error) {
	if value == "" {
		return AuthorID{}, errors.New("author ID cannot be empty")
	}
	if _, err := uuid.Parse(value); err != nil {
		return AuthorID{}, errors.New("author ID must be a valid UUID")
	}
	return AuthorID{value: value}, nil
}

// Value returns the identifier value.
func (a AuthorID) Value() string {
	return a.value
}

// String returns the identifier as string.
func (a AuthorID) String() string {
	return a.value
}

// IsZero indicates whether the identifier is set.
func (a AuthorID) IsZero() bool {
	return a.value == ""
}

// Equals compares two author identifiers.
func (a AuthorID) Equals(other AuthorID) bool {
	return a.value == other.value
}
