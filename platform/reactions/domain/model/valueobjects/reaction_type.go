package valueobjects

import (
	"errors"
	"strings"
)

// Predefined reaction types.
const (
	LikeReactionType  = "like"
	LoveReactionType  = "love"
	HahaReactionType  = "haha"
	WowReactionType   = "wow"
	SadReactionType   = "sad"
	AngryReactionType = "angry"
)

var validReactionTypes = map[string]bool{
	LikeReactionType:  true,
	LoveReactionType:  true,
	HahaReactionType:  true,
	WowReactionType:   true,
	SadReactionType:   true,
	AngryReactionType: true,
}

// ReactionType represents the type of reaction (like, love, etc.).
type ReactionType struct {
	value string `json:"value" bson:"reaction_type"`
}

// NewReactionType validates and creates a ReactionType.
func NewReactionType(value string) (ReactionType, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return ReactionType{}, errors.New("reaction type cannot be empty")
	}
	if !validReactionTypes[normalized] {
		return ReactionType{}, errors.New("reaction type must be one of: like, love, haha, wow, sad, angry")
	}
	return ReactionType{value: normalized}, nil
}

// DefaultLikeType returns the default like reaction type.
func DefaultLikeType() ReactionType {
	return ReactionType{value: LikeReactionType}
}

// Value returns the string value of the reaction type.
func (r ReactionType) Value() string {
	return r.value
}

// String returns the string representation of the reaction type.
func (r ReactionType) String() string {
	return r.value
}

// IsZero indicates if the reaction type is unset.
func (r ReactionType) IsZero() bool {
	return r.value == ""
}

// IsLike indicates if the reaction is a like.
func (r ReactionType) IsLike() bool {
	return r.value == LikeReactionType
}

// IsLove indicates if the reaction is a love.
func (r ReactionType) IsLove() bool {
	return r.value == LoveReactionType
}
