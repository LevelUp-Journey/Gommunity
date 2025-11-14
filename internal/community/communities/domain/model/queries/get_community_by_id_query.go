package queries

import (
	"errors"

	"Gommunity/internal/community/communities/domain/model/valueobjects"
)

type GetCommunityByIDQuery struct {
	communityID valueobjects.CommunityID
}

func NewGetCommunityByIDQuery(communityID valueobjects.CommunityID) (GetCommunityByIDQuery, error) {
	if communityID.IsZero() {
		return GetCommunityByIDQuery{}, errors.New("communityID cannot be empty")
	}

	return GetCommunityByIDQuery{
		communityID: communityID,
	}, nil
}

func (q GetCommunityByIDQuery) CommunityID() valueobjects.CommunityID {
	return q.communityID
}
