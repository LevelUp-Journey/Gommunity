package acl

import (
	"context"
	"time"
)

// PostData represents post data exposed to other bounded contexts
type PostData struct {
	PostID      string
	CommunityID string
	AuthorID    string
	Content     string
	MessageType string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PostsFacade exposes posts operations to other bounded contexts.
type PostsFacade interface {
	PostExists(ctx context.Context, postID string) (bool, error)
	GetAnnouncementsByCommunities(ctx context.Context, communityIDs []string, limit, offset *int) ([]*PostData, error)
}
