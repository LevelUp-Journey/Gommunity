package queries

import (
	"errors"

	"Gommunity/platform/reactions/domain/model/valueobjects"
)

// GetUserReactionOnPostQuery represents a request to check if a user reacted to a post.
type GetUserReactionOnPostQuery struct {
	postID valueobjects.PostID
	userID valueobjects.UserID
}

// NewGetUserReactionOnPostQuery validates and builds a GetUserReactionOnPostQuery.
func NewGetUserReactionOnPostQuery(
	postID valueobjects.PostID,
	userID valueobjects.UserID,
) (GetUserReactionOnPostQuery, error) {
	if postID.IsZero() {
		return GetUserReactionOnPostQuery{}, errors.New("post ID is required")
	}
	if userID.IsZero() {
		return GetUserReactionOnPostQuery{}, errors.New("user ID is required")
	}
	return GetUserReactionOnPostQuery{
		postID: postID,
		userID: userID,
	}, nil
}

// PostID returns the post identifier.
func (q GetUserReactionOnPostQuery) PostID() valueobjects.PostID {
	return q.postID
}

// UserID returns the user identifier.
func (q GetUserReactionOnPostQuery) UserID() valueobjects.UserID {
	return q.userID
}
