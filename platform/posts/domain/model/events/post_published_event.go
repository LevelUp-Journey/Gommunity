package events

import (
	"time"

	"Gommunity/platform/posts/domain/model/valueobjects"
)

// PostPublishedEvent represents the publication of a post.
type PostPublishedEvent struct {
	postID      valueobjects.PostID
	communityID valueobjects.CommunityID
	authorID    valueobjects.AuthorID
	postType    valueobjects.PostType
	occurredOn  time.Time
}

// NewPostPublishedEvent creates a new PostPublishedEvent.
func NewPostPublishedEvent(
	postID valueobjects.PostID,
	communityID valueobjects.CommunityID,
	authorID valueobjects.AuthorID,
	postType valueobjects.PostType,
) PostPublishedEvent {
	return PostPublishedEvent{
		postID:      postID,
		communityID: communityID,
		authorID:    authorID,
		postType:    postType,
		occurredOn:  time.Now(),
	}
}

// PostID returns the post identifier.
func (e PostPublishedEvent) PostID() valueobjects.PostID {
	return e.postID
}

// CommunityID returns the community identifier.
func (e PostPublishedEvent) CommunityID() valueobjects.CommunityID {
	return e.communityID
}

// AuthorID returns the author identifier.
func (e PostPublishedEvent) AuthorID() valueobjects.AuthorID {
	return e.authorID
}

// PostType returns the type of post.
func (e PostPublishedEvent) PostType() valueobjects.PostType {
	return e.postType
}

// OccurredOn returns the timestamp of the event.
func (e PostPublishedEvent) OccurredOn() time.Time {
	return e.occurredOn
}
