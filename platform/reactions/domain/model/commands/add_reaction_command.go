package commands

import (
	"errors"

	"Gommunity/platform/reactions/domain/model/valueobjects"
)

// AddReactionCommand represents the intent to add a reaction to a post.
type AddReactionCommand struct {
	postID       valueobjects.PostID
	userID       valueobjects.UserID
	reactionType valueobjects.ReactionType
}

// NewAddReactionCommand validates and builds an AddReactionCommand.
func NewAddReactionCommand(
	postID valueobjects.PostID,
	userID valueobjects.UserID,
	reactionType valueobjects.ReactionType,
) (AddReactionCommand, error) {
	if postID.IsZero() {
		return AddReactionCommand{}, errors.New("post ID is required")
	}
	if userID.IsZero() {
		return AddReactionCommand{}, errors.New("user ID is required")
	}
	if reactionType.IsZero() {
		return AddReactionCommand{}, errors.New("reaction type is required")
	}

	return AddReactionCommand{
		postID:       postID,
		userID:       userID,
		reactionType: reactionType,
	}, nil
}

// PostID returns the post identifier.
func (c AddReactionCommand) PostID() valueobjects.PostID {
	return c.postID
}

// UserID returns the user identifier.
func (c AddReactionCommand) UserID() valueobjects.UserID {
	return c.userID
}

// ReactionType returns the reaction type.
func (c AddReactionCommand) ReactionType() valueobjects.ReactionType {
	return c.reactionType
}
