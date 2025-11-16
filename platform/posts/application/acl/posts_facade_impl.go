package acl


import (
	"context"

	"Gommunity/platform/posts/domain/model/queries"
	"Gommunity/platform/posts/domain/model/valueobjects"
	"Gommunity/platform/posts/domain/services"
	"Gommunity/platform/posts/interfaces/acl"
)

type postsFacadeImpl struct {
	queryService services.PostQueryService
}

// NewPostsFacade constructs the posts facade implementation.
func NewPostsFacade(queryService services.PostQueryService) acl.PostsFacade {
	return &postsFacadeImpl{
		queryService: queryService,
	}
}

// PostExists checks if a post exists by ID.
func (f *postsFacadeImpl) PostExists(ctx context.Context, postID string) (bool, error) {
	postIDVO, err := valueobjects.NewPostID(postID)
	if err != nil {
		return false, err
	}

	query, err := queries.NewGetPostByIDQuery(postIDVO)
	if err != nil {
		return false, err
	}

	post, err := f.queryService.HandleGetByID(ctx, query)
	if err != nil {
		return false, err
	}

	return post != nil, nil
}
