package resources

import "time"

// CommunityResource represents a community in API responses
type CommunityResource struct {
	ID          string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CommunityID string    `json:"communityId" example:"550e8400-e29b-41d4-a716-446655440001"`
	OwnerID     string    `json:"ownerId" example:"550e8400-e29b-41d4-a716-446655440002"`
	Name        string    `json:"name" example:"Data Science Community"`
	Description string    `json:"description" example:"A community for data science enthusiasts to share knowledge and collaborate"`
	IconURL     *string   `json:"iconUrl,omitempty" example:"https://example.com/icon.jpg"`
	BannerURL   *string   `json:"bannerUrl,omitempty" example:"https://example.com/banner.jpg"`
	IsPrivate   bool      `json:"isPrivate" example:"false"`
	CreatedAt   time.Time `json:"createdAt" example:"2025-11-13T17:02:46Z"`
	UpdatedAt   time.Time `json:"updatedAt" example:"2025-11-13T17:02:46Z"`
}

// CreateCommunityResource represents the request to create a community
type CreateCommunityResource struct {
	Name        string  `json:"name" binding:"required,min=3,max=100" example:"Data Science Community"`
	Description string  `json:"description" binding:"required,min=10,max=500" example:"A community for data science enthusiasts to share knowledge and collaborate"`
	IconURL     *string `json:"iconUrl,omitempty" binding:"omitempty,url" example:"https://example.com/icon.jpg"`
	BannerURL   *string `json:"bannerUrl,omitempty" binding:"omitempty,url" example:"https://example.com/banner.jpg"`
}

// UpdateCommunityResource represents the request to update community information
type UpdateCommunityResource struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=3,max=100" example:"Updated Community Name"`
	Description *string `json:"description,omitempty" binding:"omitempty,min=10,max=500" example:"Updated community description"`
	IconURL     *string `json:"iconUrl,omitempty" binding:"omitempty,url" example:"https://example.com/icon.jpg"`
	BannerURL   *string `json:"bannerUrl,omitempty" binding:"omitempty,url" example:"https://example.com/banner.jpg"`
}

// UpdateCommunityPrivacyResource represents the request to update community privacy status
type UpdateCommunityPrivacyResource struct {
	IsPrivate bool `json:"isPrivate" binding:"required" example:"true"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}
