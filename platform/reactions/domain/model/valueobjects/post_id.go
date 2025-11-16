package valueobjects

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PostID represents a reference to a post from the posts bounded context.
type PostID struct {
	value string `json:"value" bson:"post_id"`
}

// NewPostID creates and validates a PostID.
func NewPostID(value string) (PostID, error) {
	if value == "" {
		return PostID{}, errors.New("post ID cannot be empty")
	}
	if _, err := primitive.ObjectIDFromHex(value); err != nil {
		return PostID{}, errors.New("post ID must be a valid ObjectID")
	}
	return PostID{value: value}, nil
}

// Value returns the string value of the PostID.
func (p PostID) Value() string {
	return p.value
}

// String returns the string representation of the PostID.
func (p PostID) String() string {
	return p.value
}

// IsZero checks if the PostID is unset.
func (p PostID) IsZero() bool {
	return p.value == ""
}
