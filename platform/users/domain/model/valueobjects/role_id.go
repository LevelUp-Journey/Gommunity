package valueobjects

import (
	"errors"
)

type RoleID struct {
	value string `json:"value" bson:"role_id"`
}

func NewRoleID(value string) (RoleID, error) {
	if value == "" {
		return RoleID{}, errors.New("role ID cannot be empty")
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
