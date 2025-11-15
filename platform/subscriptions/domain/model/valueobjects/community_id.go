package valueobjects

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CommunityID represents a Community ID from the Communities bounded context
type CommunityID struct {
	value string `json:"value" bson:"community_id"`
}

func NewCommunityID(value string) (CommunityID, error) {
	if value == "" {
		return CommunityID{}, errors.New("community ID cannot be empty")
	}
	if !primitive.IsValidObjectID(value) {
		return CommunityID{}, errors.New("community ID must be a valid ObjectID")
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

func (c CommunityID) Equals(other CommunityID) bool {
	return c.value == other.value
}

func (c CommunityID) MarshalBSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, c.value)), nil
}
