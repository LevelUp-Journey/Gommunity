package resources

import "time"

// PostResource represents a post in responses.
// Note: Post type has been removed - all posts are messages.
// Only community owners and admins can create posts.
type PostResource struct {
	PostID      string    `json:"postId" example:"64c2f1e5b9d3a45f78901234"`
	CommunityID string    `json:"communityId" example:"550e8400-e29b-41d4-a716-446655440002"`
	AuthorID    string    `json:"authorId" example:"550e8400-e29b-41d4-a716-446655440003"`
	Content     string    `json:"content" example:"Hello students!\nRemember to submit your projects."`
	Images      []string  `json:"images" example:"https://example.com/image.png"`
	CreatedAt   time.Time `json:"createdAt" example:"2025-01-12T12:00:00Z"`
	UpdatedAt   time.Time `json:"updatedAt" example:"2025-01-12T12:05:00Z"`
}

// CreatePostResource represents the payload to create a post.
// Only community owners and admins can create posts.
type CreatePostResource struct {
	Content string   `json:"content" example:"Hello community!" binding:"required"`
	Images  []string `json:"images" binding:"omitempty,dive,url"`
}

// ErrorResponse represents an error payload.
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}
