package entities

import (
	"time"

	"Gommunity/internal/community/users/domain/model/valueobjects"

	"github.com/google/uuid"
)

type User struct {
	id         string                 `bson:"_id"`
	userID     valueobjects.UserID    `bson:"user_id"`
	profileID  valueobjects.ProfileID `bson:"profile_id"`
	username   valueobjects.Username  `bson:"username"`
	roleID     valueobjects.RoleID    `bson:"role_id"`
	profileURL *string                `bson:"profile_url"`
	bannerURL  *string                `bson:"banner_url"`
	updatedAt  time.Time              `bson:"updated_at"`
	createdAt  time.Time              `bson:"created_at"`
}

// NewUser creates a new User entity
func NewUser(
	userID valueobjects.UserID,
	profileID valueobjects.ProfileID,
	username valueobjects.Username,
	profileURL *string,
) (*User, error) {
	now := time.Now()

	roleID, _ := valueobjects.NewRoleID(valueobjects.UserRoleIDStr)
	return &User{
		id:         uuid.New().String(),
		userID:     userID,
		profileID:  profileID,
		username:   username,
		roleID:     roleID,
		profileURL: profileURL,
		bannerURL:  nil,
		updatedAt:  now,
		createdAt:  now,
	}, nil
}

// Getters
func (u *User) ID() string {
	return u.id
}

func (u *User) UserID() valueobjects.UserID {
	return u.userID
}

func (u *User) ProfileID() valueobjects.ProfileID {
	return u.profileID
}

func (u *User) Username() valueobjects.Username {
	return u.username
}

func (u *User) RoleID() valueobjects.RoleID {
	return u.roleID
}

func (u *User) ProfileURL() *string {
	return u.profileURL
}

func (u *User) BannerURL() *string {
	return u.bannerURL
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdateProfile updates user profile information
func (u *User) UpdateProfile(username valueobjects.Username, profileURL *string, bannerURL *string) {
	u.username = username
	u.profileURL = profileURL
	if bannerURL != nil {
		u.bannerURL = bannerURL
	}
	u.updatedAt = time.Now()
}

func (u *User) UpdateRoleID(roleID valueobjects.RoleID) {
	u.roleID = roleID
	u.updatedAt = time.Now()
}
