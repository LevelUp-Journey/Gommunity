package resources


import "time"

// ReactionResource represents a reaction in responses.
type ReactionResource struct {
	ReactionID   string    `json:"reactionId" example:"64c2f1e5b9d3a45f78901234"`
	PostID       string    `json:"postId" example:"64c2f1e5b9d3a45f78901235"`
	UserID       string    `json:"userId" example:"550e8400-e29b-41d4-a716-446655440000"`
	ReactionType string    `json:"reactionType" example:"like"`
	CreatedAt    time.Time `json:"createdAt" example:"2025-01-12T12:00:00Z"`
	UpdatedAt    time.Time `json:"updatedAt" example:"2025-01-12T12:05:00Z"`
}

// AddReactionResource represents the payload to add a reaction.
type AddReactionResource struct {
	ReactionType string `json:"reactionType" example:"like" binding:"required,oneof=like love haha wow sad angry"`
}

// ReactionCountResource represents aggregated reaction counts for a post.
type ReactionCountResource struct {
	PostID     string         `json:"postId" example:"64c2f1e5b9d3a45f78901235"`
	TotalCount int            `json:"totalCount" example:"42"`
	Counts     map[string]int `json:"counts" example:"like:25,love:10,haha:5,wow:2"`
}

// ErrorResponse represents an error payload.
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}
