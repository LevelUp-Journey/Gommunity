package valueobjects

import (
	"errors"
	"strings"
)

// Predefined role names for Community context
// Note: These are community-specific roles (different from IAM roles: STUDENT, TEACHER, ADMIN)
// Community roles are scoped per community - a user can have different roles in different communities
const (
	MemberRoleName = "member"
	AdminRoleName  = "admin"
	OwnerRoleName  = "owner"
)

// CommunityRole represents a role that a user has within a specific community
// Roles are scoped to communities - a user can have different roles in different communities
// Valid roles: member, admin, owner
type CommunityRole struct {
	value string `json:"value" bson:"role"`
}

var validCommunityRoles = map[string]bool{
	MemberRoleName: true,
	AdminRoleName:  true,
	OwnerRoleName:  true,
}

func NewCommunityRole(value string) (CommunityRole, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return CommunityRole{}, errors.New("community role cannot be empty")
	}
	if !validCommunityRoles[normalized] {
		return CommunityRole{}, errors.New("invalid community role: must be one of member, admin, owner")
	}
	return CommunityRole{value: normalized}, nil
}

func (r CommunityRole) Value() string {
	return r.value
}

func (r CommunityRole) String() string {
	return r.value
}

func (r CommunityRole) IsZero() bool {
	return r.value == ""
}

func (r CommunityRole) Equals(other CommunityRole) bool {
	return r.value == other.value
}

// IsAdminOrOwner checks if the role has admin or owner privileges
func (r CommunityRole) IsAdminOrOwner() bool {
	return r.value == AdminRoleName || r.value == OwnerRoleName
}

// IsAdmin checks if the role is admin
func (r CommunityRole) IsAdmin() bool {
	return r.value == AdminRoleName
}

// IsOwner checks if the role is owner
func (r CommunityRole) IsOwner() bool {
	return r.value == OwnerRoleName
}

// IsMember checks if the role is member
func (r CommunityRole) IsMember() bool {
	return r.value == MemberRoleName
}

// Predefined roles for convenience
var (
	MemberRole = CommunityRole{value: MemberRoleName}
	AdminRole  = CommunityRole{value: AdminRoleName}
	OwnerRole  = CommunityRole{value: OwnerRoleName}
)
