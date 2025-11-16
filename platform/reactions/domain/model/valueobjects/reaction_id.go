package valueobjects

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReactionID represents the unique identifier for a reaction.
type ReactionID struct {
	value string `json:"value" bson:"_id"`
}

// NewReactionID creates and validates a ReactionID.
func NewReactionID(value string) (ReactionID, error) {
	if value == "" {
		return ReactionID{}, errors.New("reaction ID cannot be empty")
	}
	if _, err := primitive.ObjectIDFromHex(value); err != nil {
		return ReactionID{}, errors.New("reaction ID must be a valid ObjectID")
	}
	return ReactionID{value: value}, nil
}

// GenerateReactionID creates a new unique ReactionID.
func GenerateReactionID() ReactionID {
	return ReactionID{value: primitive.NewObjectID().Hex()}
}

// Value returns the string value of the ReactionID.
func (r ReactionID) Value() string {
	return r.value
}

// String returns the string representation of the ReactionID.
func (r ReactionID) String() string {
	return r.value
}

// IsZero checks if the ReactionID is unset.
func (r ReactionID) IsZero() bool {
	return r.value == ""
}
