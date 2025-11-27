package acl

import (
	"context"

	"Gommunity/platform/posts/domain/model/queries"
	"Gommunity/platform/posts/domain/model/valueobjects"
	"Gommunity/platform/posts/domain/repositories"
	"Gommunity/platform/posts/domain/services"
	"Gommunity/platform/posts/interfaces/acl"
)

type postsFacadeImpl struct {
	queryService services.PostQueryService
	postRepo     repositories.PostRepository
}

// NewPostsFacade constructs the posts facade implementation.
func NewPostsFacade(queryService services.PostQueryService, postRepo repositories.PostRepository) acl.PostsFacade {
	return &postsFacadeImpl{
		queryService: queryService,
		postRepo:     postRepo,
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

// GetAnnouncementsByCommunities retrieves posts from multiple communities.
// Note: Announcements have been removed - all posts are messages now.
// This method returns all posts for backward compatibility with Feed BC.
func (f *postsFacadeImpl) GetAnnouncementsByCommunities(ctx context.Context, communityIDs []string, limit, offset *int) ([]*acl.PostData, error) {
	// Convert string IDs to value objects
	communityIDVOs := make([]valueobjects.CommunityID, 0, len(communityIDs))
	for _, id := range communityIDs {
		communityIDVO, err := valueobjects.NewCommunityID(id)
		if err != nil {
			continue // Skip invalid IDs
		}
		communityIDVOs = append(communityIDVOs, communityIDVO)
	}

	if len(communityIDVOs) == 0 {
		return []*acl.PostData{}, nil
	}

	// Get all posts from repository (no type filtering)
	posts, err := f.postRepo.FindByCommunities(ctx, communityIDVOs, limit, offset)
	if err != nil {
		return nil, err
	}

	// Transform to PostData
	result := make([]*acl.PostData, len(posts))
	for i, post := range posts {
		result[i] = &acl.PostData{
			PostID:      post.PostID().Value(),
			CommunityID: post.CommunityID().Value(),
			AuthorID:    post.AuthorID().Value(),
			Content:     post.Content().Value(),
			MessageType: "message", // All posts are messages now
			CreatedAt:   post.CreatedAt(),
			UpdatedAt:   post.UpdatedAt(),
		}
	}

	return result, nil
}
