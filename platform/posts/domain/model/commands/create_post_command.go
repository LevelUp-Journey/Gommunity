package commands

import (
	"errors"

	"Gommunity/platform/posts/domain/model/valueobjects"
)

// CreatePostCommand represents the intent to publish a new post.
// Only community owners and admins can create posts.
type CreatePostCommand struct {
	communityID valueobjects.CommunityID
	authorID    valueobjects.AuthorID
	content     valueobjects.PostContent
	images      valueobjects.PostImages
}

// NewCreatePostCommand validates and builds a CreatePostCommand.
func NewCreatePostCommand(
	communityID valueobjects.CommunityID,
	authorID valueobjects.AuthorID,
	content valueobjects.PostContent,
	images valueobjects.PostImages,
) (CreatePostCommand, error) {
	if communityID.IsZero() {
		return CreatePostCommand{}, errors.New("community ID is required")
	}
	if authorID.IsZero() {
		return CreatePostCommand{}, errors.New("author ID is required")
	}
	if content.IsZero() {
		return CreatePostCommand{}, errors.New("content is required")
	}

	return CreatePostCommand{
		communityID: communityID,
		authorID:    authorID,
		content:     content,
		images:      images,
	}, nil
}

// CommunityID returns the community identifier for the command.
func (c CreatePostCommand) CommunityID() valueobjects.CommunityID {
	return c.communityID
}

// AuthorID returns the author identifier.
func (c CreatePostCommand) AuthorID() valueobjects.AuthorID {
	return c.authorID
}

// Content returns the post content.
func (c CreatePostCommand) Content() valueobjects.PostContent {
	return c.content
}

// Images returns the associated images.
func (c CreatePostCommand) Images() valueobjects.PostImages {
	return c.images
}
