package queries

import (
	"errors"

	"Gommunity/platform/posts/domain/model/valueobjects"
)

// GetPostsByCommunityQuery retrieves posts for a community.
type GetPostsByCommunityQuery struct {
	communityID valueobjects.CommunityID
	limit       *int
	offset      *int
}

// NewGetPostsByCommunityQuery validates input and creates the query.
func NewGetPostsByCommunityQuery(communityID valueobjects.CommunityID) (GetPostsByCommunityQuery, error) {
	if communityID.IsZero() {
		return GetPostsByCommunityQuery{}, errors.New("community ID is required")
	}
	return GetPostsByCommunityQuery{communityID: communityID}, nil
}

// WithPagination sets pagination options.
func (q GetPostsByCommunityQuery) WithPagination(limit, offset int) GetPostsByCommunityQuery {
	if limit > 0 {
		q.limit = &limit
	}
	if offset >= 0 {
		q.offset = &offset
	}
	return q
}

// CommunityID returns the community identifier.
func (q GetPostsByCommunityQuery) CommunityID() valueobjects.CommunityID {
	return q.communityID
}

// Limit returns the optional limit.
func (q GetPostsByCommunityQuery) Limit() *int {
	return q.limit
}

// Offset returns the optional offset.
func (q GetPostsByCommunityQuery) Offset() *int {
	return q.offset
}
