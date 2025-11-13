package valueobjects

import (
	"errors"
	"strings"
)

// Predefined role IDs (UUIDs)
const (
	UserRoleIDStr   = "550e8400-e29b-41d4-a716-446655440001"
	MemberRoleIDStr = "550e8400-e29b-41d4-a716-446655440002"
	AdminRoleIDStr  = "550e8400-e29b-41d4-a716-446655440003"
	OwnerRoleIDStr  = "550e8400-e29b-41d4-a716-446655440004"
)

type Role struct {
	value string `json:"value" bson:"role"`
}

var validRoles = map[string]bool{
	"user":   true,
	"member": true,
	"admin":  true,
	"owner":  true,
}

func NewRole(value string) (Role, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return Role{}, errors.New("role cannot be empty")
	}
	if !validRoles[normalized] {
		return Role{}, errors.New("invalid role: must be one of user, member, admin, owner")
	}
	return Role{value: normalized}, nil
}

func (r Role) Value() string {
	return r.value
}

func (r Role) String() string {
	return r.value
}

func (r Role) IsZero() bool {
	return r.value == ""
}

// Predefined roles for convenience
var (
	UserRole   = Role{value: "user"}
	MemberRole = Role{value: "member"}
	AdminRole  = Role{value: "admin"}
	OwnerRole  = Role{value: "owner"}
)
