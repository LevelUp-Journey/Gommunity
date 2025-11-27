package queries

import (
	"errors"

	"Gommunity/platform/reactions/domain/model/valueobjects"
)

// GetReactionCountByPostQuery represents a request to get reaction counts for a post.
type GetReactionCountByPostQuery struct {
	postID valueobjects.PostID
}

// NewGetReactionCountByPostQuery validates and builds a GetReactionCountByPostQuery.
func NewGetReactionCountByPostQuery(postID valueobjects.PostID) (GetReactionCountByPostQuery, error) {
	if postID.IsZero() {
		return GetReactionCountByPostQuery{}, errors.New("post ID is required")
	}
	return GetReactionCountByPostQuery{postID: postID}, nil
}

// PostID returns the post identifier.
func (q GetReactionCountByPostQuery) PostID() valueobjects.PostID {
	return q.postID
}
