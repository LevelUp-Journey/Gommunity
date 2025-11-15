package queries

import (
	"errors"

	"Gommunity/platform/posts/domain/model/valueobjects"
)

// GetPostByIDQuery retrieves a post by its identifier.
type GetPostByIDQuery struct {
	postID valueobjects.PostID
}

// NewGetPostByIDQuery validates input and creates the query.
func NewGetPostByIDQuery(postID valueobjects.PostID) (GetPostByIDQuery, error) {
	if postID.IsZero() {
		return GetPostByIDQuery{}, errors.New("post ID is required")
	}
	return GetPostByIDQuery{postID: postID}, nil
}

// PostID returns the identifier for lookup.
func (q GetPostByIDQuery) PostID() valueobjects.PostID {
	return q.postID
}
