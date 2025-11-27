package valueobjects

import (
	"errors"
	"strings"
)

const (
	CommunityRoleMember = "member"
	CommunityRoleAdmin  = "admin"
	CommunityRoleOwner  = "owner"
)

// CommunityRole represents a user's privilege within a community from the posts perspective.
type CommunityRole struct {
	value string `json:"value" bson:"role"`
}

var allowedCommunityRoles = map[string]struct{}{
	CommunityRoleMember: {},
	CommunityRoleAdmin:  {},
	CommunityRoleOwner:  {},
}

// NewCommunityRole validates and creates a CommunityRole.
func NewCommunityRole(value string) (CommunityRole, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return CommunityRole{}, errors.New("role cannot be empty")
	}
	if _, ok := allowedCommunityRoles[normalized]; !ok {
		return CommunityRole{}, errors.New("invalid community role")
	}
	return CommunityRole{value: normalized}, nil
}

// Value returns the underlying string.
func (r CommunityRole) Value() string {
	return r.value
}

// IsZero indicates if the role is unset.
func (r CommunityRole) IsZero() bool {
	return r.value == ""
}

// IsMember indicates if the role is member.
func (r CommunityRole) IsMember() bool {
	return r.value == CommunityRoleMember
}

// IsAdmin indicates if the role is admin.
func (r CommunityRole) IsAdmin() bool {
	return r.value == CommunityRoleAdmin
}

// IsOwner indicates if the role is owner.
func (r CommunityRole) IsOwner() bool {
	return r.value == CommunityRoleOwner
}

// IsAdminOrOwner indicates admin or owner privileges.
func (r CommunityRole) IsAdminOrOwner() bool {
	return r.IsAdmin() || r.IsOwner()
}
