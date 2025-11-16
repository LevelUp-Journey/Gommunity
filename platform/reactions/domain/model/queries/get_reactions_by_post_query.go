package queries

import (
	"errors"

	"Gommunity/platform/reactions/domain/model/valueobjects"
)

// GetReactionsByPostQuery represents a request to retrieve all reactions for a post.
type GetReactionsByPostQuery struct {
	postID valueobjects.PostID
}

// NewGetReactionsByPostQuery validates and builds a GetReactionsByPostQuery.
func NewGetReactionsByPostQuery(postID valueobjects.PostID) (GetReactionsByPostQuery, error) {
	if postID.IsZero() {
		return GetReactionsByPostQuery{}, errors.New("post ID is required")
	}
	return GetReactionsByPostQuery{postID: postID}, nil
}

// PostID returns the post identifier.
func (q GetReactionsByPostQuery) PostID() valueobjects.PostID {
	return q.postID
}
