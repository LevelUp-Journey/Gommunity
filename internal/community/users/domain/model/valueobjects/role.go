package valueobjects

import (
	"errors"
	"strings"
)

// Predefined role IDs (UUIDs)
const (
	StudentRoleIDStr = "550e8400-e29b-41d4-a716-446655440001"
	TeacherRoleIDStr = "550e8400-e29b-41d4-a716-446655440002"
	AdminRoleIDStr   = "550e8400-e29b-41d4-a716-446655440003"
	MemberRoleIDStr  = "550e8400-e29b-41d4-a716-446655440004"
	OwnerRoleIDStr   = "550e8400-e29b-41d4-a716-446655440005"

	// Deprecated - for backward compatibility
	UserRoleIDStr = StudentRoleIDStr
)

type Role struct {
	value string `json:"value" bson:"role"`
}

var validRoles = map[string]bool{
	"student": true,
	"teacher": true,
	"admin":   true,
	"member":  true,
	"owner":   true,
}

func NewRole(value string) (Role, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return Role{}, errors.New("role cannot be empty")
	}
	if !validRoles[normalized] {
		return Role{}, errors.New("invalid role: must be one of student, teacher, admin, member, owner")
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
	StudentRole = Role{value: "student"}
	TeacherRole = Role{value: "teacher"}
	AdminRole   = Role{value: "admin"}
	MemberRole  = Role{value: "member"}
	OwnerRole   = Role{value: "owner"}

	// Deprecated - for backward compatibility
	UserRole = StudentRole
)
