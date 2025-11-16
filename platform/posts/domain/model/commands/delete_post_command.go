package commands

import (
	"errors"

	"Gommunity/platform/posts/domain/model/valueobjects"
)

// DeletePostCommand represents the intent to remove a post.
type DeletePostCommand struct {
	postID      valueobjects.PostID
	requestedBy valueobjects.AuthorID
}

// NewDeletePostCommand builds a DeletePostCommand.
func NewDeletePostCommand(
	postID valueobjects.PostID,
	requestedBy valueobjects.AuthorID,
) (DeletePostCommand, error) {
	if postID.IsZero() {
		return DeletePostCommand{}, errors.New("post ID is required")
	}
	if requestedBy.IsZero() {
		return DeletePostCommand{}, errors.New("requesting user ID is required")
	}

	return DeletePostCommand{
		postID:      postID,
		requestedBy: requestedBy,
	}, nil
}

// PostID returns the post identifier.
func (c DeletePostCommand) PostID() valueobjects.PostID {
	return c.postID
}

// RequestedBy returns the identifier of the user performing the action.
func (c DeletePostCommand) RequestedBy() valueobjects.AuthorID {
	return c.requestedBy
}
