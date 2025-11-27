package entities

import (
	"errors"
	"time"

	"Gommunity/platform/posts/domain/model/valueobjects"
)

// Post represents a message publication made inside a community.
// All posts are messages - announcements have been removed.
type Post struct {
	id          string
	postID      valueobjects.PostID
	communityID valueobjects.CommunityID
	authorID    valueobjects.AuthorID
	content     valueobjects.PostContent
	images      valueobjects.PostImages
	createdAt   time.Time
	updatedAt   time.Time
}

// NewPost creates a new post aggregate.
// Only community owners and admins can create posts.
func NewPost(
	communityID valueobjects.CommunityID,
	authorID valueobjects.AuthorID,
	content valueobjects.PostContent,
	images valueobjects.PostImages,
) (*Post, error) {
	if communityID.IsZero() {
		return nil, errors.New("community ID is required")
	}
	if authorID.IsZero() {
		return nil, errors.New("author ID is required")
	}
	if content.IsZero() {
		return nil, errors.New("post content is required")
	}

	now := time.Now()
	postID := valueobjects.GeneratePostID()

	return &Post{
		id:          postID.Value(),
		postID:      postID,
		communityID: communityID,
		authorID:    authorID,
		content:     content,
		images:      images,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// ReconstructPost rebuilds a post from persistence.
func ReconstructPost(
	id string,
	postID valueobjects.PostID,
	communityID valueobjects.CommunityID,
	authorID valueobjects.AuthorID,
	content valueobjects.PostContent,
	images valueobjects.PostImages,
	createdAt time.Time,
	updatedAt time.Time,
) *Post {
	return &Post{
		id:          id,
		postID:      postID,
		communityID: communityID,
		authorID:    authorID,
		content:     content,
		images:      images,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// ID returns the persistence identifier.
func (p *Post) ID() string {
	return p.id
}

// PostID returns the aggregate identifier.
func (p *Post) PostID() valueobjects.PostID {
	return p.postID
}

// CommunityID returns the community identifier.
func (p *Post) CommunityID() valueobjects.CommunityID {
	return p.communityID
}

// AuthorID returns the author identifier.
func (p *Post) AuthorID() valueobjects.AuthorID {
	return p.authorID
}

// Content returns the post content.
func (p *Post) Content() valueobjects.PostContent {
	return p.content
}

// Images returns the images attached to the post.
func (p *Post) Images() valueobjects.PostImages {
	return p.images
}

// CreatedAt returns the creation timestamp.
func (p *Post) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt returns the update timestamp.
func (p *Post) UpdatedAt() time.Time {
	return p.updatedAt
}
