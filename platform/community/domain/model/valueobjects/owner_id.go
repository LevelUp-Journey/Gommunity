package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

type OwnerID struct {
	value string `json:"value" bson:"owner_id"`
}

func NewOwnerID(value string) (OwnerID, error) {
	if value == "" {
		return OwnerID{}, errors.New("owner ID cannot be empty")
	}

	// Validate UUID format
	if _, err := uuid.Parse(value); err != nil {
		return OwnerID{}, errors.New("owner ID must be a valid UUID")
	}

	return OwnerID{value: value}, nil
}

func (o OwnerID) Value() string {
	return o.value
}

func (o OwnerID) String() string {
	return o.value
}

func (o OwnerID) IsZero() bool {
	return o.value == ""
}
