package entities

import (
	"time"

	"Gommunity/platform/feed/domain/model/valueobjects"
)

// FeedItem represents a post item in a user's feed
type FeedItem struct {
	postID      valueobjects.PostID
	communityID valueobjects.CommunityID
	authorID    string
	content     string
	messageType string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewFeedItem creates a new FeedItem
func NewFeedItem(
	postID valueobjects.PostID,
	communityID valueobjects.CommunityID,
	authorID string,
	content string,
	messageType string,
	createdAt time.Time,
	updatedAt time.Time,
) *FeedItem {
	return &FeedItem{
		postID:      postID,
		communityID: communityID,
		authorID:    authorID,
		content:     content,
		messageType: messageType,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// Getters
func (f *FeedItem) PostID() valueobjects.PostID {
	return f.postID
}

func (f *FeedItem) CommunityID() valueobjects.CommunityID {
	return f.communityID
}

func (f *FeedItem) AuthorID() string {
	return f.authorID
}

func (f *FeedItem) Content() string {
	return f.content
}

func (f *FeedItem) MessageType() string {
	return f.messageType
}

func (f *FeedItem) CreatedAt() time.Time {
	return f.createdAt
}

func (f *FeedItem) UpdatedAt() time.Time {
	return f.updatedAt
}
