package commands

import (
	"errors"

	"Gommunity/platform/reactions/domain/model/valueobjects"
)

// RemoveReactionCommand represents the intent to remove a user's reaction from a post.
type RemoveReactionCommand struct {
	postID valueobjects.PostID
	userID valueobjects.UserID
}

// NewRemoveReactionCommand validates and builds a RemoveReactionCommand.
func NewRemoveReactionCommand(
	postID valueobjects.PostID,
	userID valueobjects.UserID,
) (RemoveReactionCommand, error) {
	if postID.IsZero() {
		return RemoveReactionCommand{}, errors.New("post ID is required")
	}
	if userID.IsZero() {
		return RemoveReactionCommand{}, errors.New("user ID is required")
	}

	return RemoveReactionCommand{
		postID: postID,
		userID: userID,
	}, nil
}

// PostID returns the post identifier.
func (c RemoveReactionCommand) PostID() valueobjects.PostID {
	return c.postID
}

// UserID returns the user identifier.
func (c RemoveReactionCommand) UserID() valueobjects.UserID {
	return c.userID
}
