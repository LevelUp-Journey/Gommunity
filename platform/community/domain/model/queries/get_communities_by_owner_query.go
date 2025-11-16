package queries

import (
	"errors"

	"Gommunity/platform/community/domain/model/valueobjects"
)

type GetCommunitiesByOwnerQuery struct {
	ownerID valueobjects.OwnerID
}

func NewGetCommunitiesByOwnerQuery(ownerID valueobjects.OwnerID) (GetCommunitiesByOwnerQuery, error) {
	if ownerID.IsZero() {
		return GetCommunitiesByOwnerQuery{}, errors.New("ownerID cannot be empty")
	}

	return GetCommunitiesByOwnerQuery{
		ownerID: ownerID,
	}, nil
}

func (q GetCommunitiesByOwnerQuery) OwnerID() valueobjects.OwnerID {
	return q.ownerID
}
