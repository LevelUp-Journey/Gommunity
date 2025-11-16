package services

import (
	"Gommunity/platform/feed/domain/model/entities"
	"Gommunity/platform/feed/domain/model/queries"
	"context"
)

type FeedQueryService interface {
	Handle(ctx context.Context, query queries.GetUserFeedQuery) ([]*entities.FeedItem, error)
}
