package entities

import (
	"errors"
	"time"

	"Gommunity/platform/reactions/domain/model/valueobjects"
)

// Reaction represents a user's reaction to a post.
type Reaction struct {
	id           string
	reactionID   valueobjects.ReactionID
	postID       valueobjects.PostID
	userID       valueobjects.UserID
	reactionType valueobjects.ReactionType
	createdAt    time.Time
	updatedAt    time.Time
}

// NewReaction creates a new reaction aggregate.
func NewReaction(
	postID valueobjects.PostID,
	userID valueobjects.UserID,
	reactionType valueobjects.ReactionType,
) (*Reaction, error) {
	if postID.IsZero() {
		return nil, errors.New("post ID is required")
	}
	if userID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	if reactionType.IsZero() {
		return nil, errors.New("reaction type is required")
	}

	now := time.Now()
	reactionID := valueobjects.GenerateReactionID()

	return &Reaction{
		id:           reactionID.Value(),
		reactionID:   reactionID,
		postID:       postID,
		userID:       userID,
		reactionType: reactionType,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// ReconstructReaction rebuilds a reaction from persistence.
func ReconstructReaction(
	id string,
	reactionID valueobjects.ReactionID,
	postID valueobjects.PostID,
	userID valueobjects.UserID,
	reactionType valueobjects.ReactionType,
	createdAt time.Time,
	updatedAt time.Time,
) *Reaction {
	return &Reaction{
		id:           id,
		reactionID:   reactionID,
		postID:       postID,
		userID:       userID,
		reactionType: reactionType,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// ChangeReactionType updates the reaction type (business rule: user can change their reaction).
func (r *Reaction) ChangeReactionType(newType valueobjects.ReactionType) error {
	if newType.IsZero() {
		return errors.New("new reaction type is required")
	}
	r.reactionType = newType
	r.updatedAt = time.Now()
	return nil
}

// ID returns the persistence identifier.
func (r *Reaction) ID() string {
	return r.id
}

// ReactionID returns the aggregate identifier.
func (r *Reaction) ReactionID() valueobjects.ReactionID {
	return r.reactionID
}

// PostID returns the post identifier.
func (r *Reaction) PostID() valueobjects.PostID {
	return r.postID
}

// UserID returns the user identifier.
func (r *Reaction) UserID() valueobjects.UserID {
	return r.userID
}

// ReactionType returns the reaction type.
func (r *Reaction) ReactionType() valueobjects.ReactionType {
	return r.reactionType
}

// CreatedAt returns the creation timestamp.
func (r *Reaction) CreatedAt() time.Time {
	return r.createdAt
}

// UpdatedAt returns the update timestamp.
func (r *Reaction) UpdatedAt() time.Time {
	return r.updatedAt
}
