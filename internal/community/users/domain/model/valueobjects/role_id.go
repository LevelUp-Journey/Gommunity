package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

type RoleID struct {
	value string `json:"value" bson:"role_id"`
}

func NewRoleID(value string) (RoleID, error) {
	if value == "" {
		return RoleID{}, errors.New("role ID cannot be empty")
	}

	// Validate UUID format
	if _, err := uuid.Parse(value); err != nil {
		return RoleID{}, errors.New("role ID must be a valid UUID")
	}

	return RoleID{value: value}, nil
}

func (r RoleID) Value() string {
	return r.value
}

func (r RoleID) String() string {
	return r.value
}

func (r RoleID) IsZero() bool {
	return r.value == ""
}
