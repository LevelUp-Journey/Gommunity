package valueobjects

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubscriptionID struct {
	value string `json:"value" bson:"subscription_id"`
}

func NewSubscriptionID(value string) (SubscriptionID, error) {
	if value == "" {
		return SubscriptionID{}, errors.New("subscription ID cannot be empty")
	}
	if !primitive.IsValidObjectID(value) {
		return SubscriptionID{}, errors.New("subscription ID must be a valid ObjectID")
	}
	return SubscriptionID{value: value}, nil
}

func GenerateSubscriptionID() SubscriptionID {
	return SubscriptionID{value: primitive.NewObjectID().Hex()}
}

func (s SubscriptionID) Value() string {
	return s.value
}

func (s SubscriptionID) String() string {
	return s.value
}

func (s SubscriptionID) IsZero() bool {
	return s.value == ""
}

func (s SubscriptionID) Equals(other SubscriptionID) bool {
	return s.value == other.value
}

func (s SubscriptionID) MarshalBSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s.value)), nil
}
