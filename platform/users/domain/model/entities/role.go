package entities

import (
	"time"

	"Gommunity/platform/users/domain/model/valueobjects"
)

type Role struct {
	id        string              `bson:"_id"`
	roleID    valueobjects.RoleID `bson:"role_id"`
	name      string              `bson:"name"`
	createdAt time.Time           `bson:"created_at"`
}

// NewRole creates a new Role entity
func NewRole(roleID valueobjects.RoleID, name string) (*Role, error) {
	now := time.Now()

	return &Role{
		id:        roleID.Value(),
		roleID:    roleID,
		name:      name,
		createdAt: now,
	}, nil
}

// Getters
func (r *Role) ID() string {
	return r.id
}

func (r *Role) RoleID() valueobjects.RoleID {
	return r.roleID
}

func (r *Role) Name() string {
	return r.name
}

func (r *Role) CreatedAt() time.Time {
	return r.createdAt
}
