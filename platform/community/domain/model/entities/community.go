package entities

import (
	"time"

	"Gommunity/platform/community/domain/model/valueobjects"

	"github.com/google/uuid"
)

type Community struct {
	communityID valueobjects.CommunityID   `bson:"_id"`
	ownerID     valueobjects.OwnerID       `bson:"owner_id"`
	name        valueobjects.CommunityName `bson:"name"`
	description valueobjects.Description   `bson:"description"`
	iconURL     *string                    `bson:"icon_url"`
	bannerURL   *string                    `bson:"banner_url"`
	isPrivate   bool                       `bson:"is_private"`
	createdAt   time.Time                  `bson:"created_at"`
	updatedAt   time.Time                  `bson:"updated_at"`
}

// NewCommunity creates a new Community entity
func NewCommunity(
	ownerID valueobjects.OwnerID,
	name valueobjects.CommunityName,
	description valueobjects.Description,
	iconURL *string,
	bannerURL *string,
	isPrivate bool,
) (*Community, error) {
	now := time.Now()
	communityID := uuid.New().String()

	communityIDVO, err := valueobjects.NewCommunityID(communityID)
	if err != nil {
		return nil, err
	}

	return &Community{
		communityID: communityIDVO,
		ownerID:     ownerID,
		name:        name,
		description: description,
		iconURL:     iconURL,
		bannerURL:   bannerURL,
		isPrivate:   isPrivate,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Getters
func (c *Community) ID() string {
	return c.communityID.Value()
}

func (c *Community) CommunityID() valueobjects.CommunityID {
	return c.communityID
}

func (c *Community) OwnerID() valueobjects.OwnerID {
	return c.ownerID
}

func (c *Community) Name() valueobjects.CommunityName {
	return c.name
}

func (c *Community) Description() valueobjects.Description {
	return c.description
}

func (c *Community) IconURL() *string {
	return c.iconURL
}

func (c *Community) BannerURL() *string {
	return c.bannerURL
}

func (c *Community) IsPrivate() bool {
	return c.isPrivate
}

func (c *Community) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Community) UpdatedAt() time.Time {
	return c.updatedAt
}

// Business methods
func (c *Community) UpdateInfo(name valueobjects.CommunityName, description valueobjects.Description, iconURL *string, bannerURL *string) {
	c.name = name
	c.description = description
	if iconURL != nil {
		c.iconURL = iconURL
	}
	if bannerURL != nil {
		c.bannerURL = bannerURL
	}
	c.updatedAt = time.Now()
}

func (c *Community) UpdateIconURL(iconURL string) {
	c.iconURL = &iconURL
	c.updatedAt = time.Now()
}

func (c *Community) UpdateBannerURL(bannerURL string) {
	c.bannerURL = &bannerURL
	c.updatedAt = time.Now()
}

func (c *Community) MakePrivate() {
	c.isPrivate = true
	c.updatedAt = time.Now()
}

func (c *Community) MakePublic() {
	c.isPrivate = false
	c.updatedAt = time.Now()
}

func (c *Community) UpdatePrivacy(isPrivate bool) {
	c.isPrivate = isPrivate
	c.updatedAt = time.Now()
}

// IsOwner checks if the given userID is the owner of the community
func (c *Community) IsOwner(userID string) bool {
	return c.ownerID.Value() == userID
}
