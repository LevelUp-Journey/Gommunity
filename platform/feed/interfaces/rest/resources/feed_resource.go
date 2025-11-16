package resources

import "time"

// FeedItemResource represents a feed item in the REST API
type FeedItemResource struct {
	PostID      string    `json:"postId" example:"507f1f77bcf86cd799439011"`
	CommunityID string    `json:"communityId" example:"507f1f77bcf86cd799439012"`
	AuthorID    string    `json:"authorId" example:"user123"`
	Content     string    `json:"content" example:"This is an important announcement"`
	MessageType string    `json:"messageType" example:"announcement"`
	CreatedAt   time.Time `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updatedAt" example:"2023-01-01T00:00:00Z"`
}

// FeedResponse represents the feed response with pagination info
type FeedResponse struct {
	Items []FeedItemResource `json:"items"`
	Total int                `json:"total" example:"10"`
}
