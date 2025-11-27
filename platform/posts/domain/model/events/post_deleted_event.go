package events

import (
	"time"

	"Gommunity/platform/posts/domain/model/valueobjects"
)

// PostDeletedEvent represents the deletion of a post.
type PostDeletedEvent struct {
	postID     valueobjects.PostID
	requester  valueobjects.AuthorID
	occurredOn time.Time
}

// NewPostDeletedEvent creates a new PostDeletedEvent.
func NewPostDeletedEvent(
	postID valueobjects.PostID,
	requester valueobjects.AuthorID,
) PostDeletedEvent {
	return PostDeletedEvent{
		postID:     postID,
		requester:  requester,
		occurredOn: time.Now(),
	}
}

// PostID returns the post identifier.
func (e PostDeletedEvent) PostID() valueobjects.PostID {
	return e.postID
}

// Requester returns the user who requested the deletion.
func (e PostDeletedEvent) Requester() valueobjects.AuthorID {
	return e.requester
}

// OccurredOn returns the event timestamp.
func (e PostDeletedEvent) OccurredOn() time.Time {
	return e.occurredOn
}
