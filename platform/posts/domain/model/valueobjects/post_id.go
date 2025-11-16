package valueobjects

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PostID represents the identifier for a post aggregate.
type PostID struct {
	value string `json:"value" bson:"post_id"`
}

// NewPostID creates a new PostID from an existing value.
func NewPostID(value string) (PostID, error) {
	if value == "" {
		return PostID{}, errors.New("post ID cannot be empty")
	}
	if !primitive.IsValidObjectID(value) {
		return PostID{}, errors.New("post ID must be a valid ObjectID")
	}
	return PostID{value: value}, nil
}

// GeneratePostID generates a new PostID.
func GeneratePostID() PostID {
	return PostID{value: primitive.NewObjectID().Hex()}
}

// Value returns the string value of the identifier.
func (p PostID) Value() string {
	return p.value
}

// String returns the string form of the identifier.
func (p PostID) String() string {
	return p.value
}

// IsZero indicates if the identifier is empty.
func (p PostID) IsZero() bool {
	return p.value == ""
}

// Equals compares two identifiers.
func (p PostID) Equals(other PostID) bool {
	return p.value == other.value
}
