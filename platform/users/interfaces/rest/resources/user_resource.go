package resources

import "time"

// UserResource represents a user in API responses
type UserResource struct {
	ID         string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID     string    `json:"userId" example:"32f05fbf-9793-4205-980e-d23716627750"`
	ProfileID  string    `json:"profileId" example:"a751deae-573e-42e3-851c-04b242d6536d"`
	Username   string    `json:"username" example:"johndoe"`
	Role       string    `json:"role" example:"user"`
	ProfileURL *string   `json:"profileUrl,omitempty" example:"https://example.com/profile.jpg"`
	BannerURL  *string   `json:"bannerUrl,omitempty" example:"https://example.com/banner.jpg"`
	UpdatedAt  time.Time `json:"updatedAt" example:"2025-11-13T17:02:46Z"`
	CreatedAt  time.Time `json:"createdAt" example:"2025-11-13T17:02:46Z"`
}

// UpdateBannerURLResource represents the request to update banner URL
type UpdateBannerURLResource struct {
	BannerURL string `json:"bannerUrl" binding:"required,url" example:"https://example.com/banner.jpg"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}
