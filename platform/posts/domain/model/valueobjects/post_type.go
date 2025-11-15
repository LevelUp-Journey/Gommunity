package valueobjects

import (
	"errors"
	"strings"
)

// Predefined post types.
const (
	MessagePostType      = "message"
	AnnouncementPostType = "announcement"
)

var validPostTypes = map[string]bool{
	MessagePostType:      true,
	AnnouncementPostType: true,
}

// PostType represents the type of a post (message or announcement).
type PostType struct {
	value string `json:"value" bson:"post_type"`
}

// NewPostType validates and creates a PostType.
func NewPostType(value string) (PostType, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return PostType{}, errors.New("post type cannot be empty")
	}
	if !validPostTypes[normalized] {
		return PostType{}, errors.New("post type must be either message or announcement")
	}
	return PostType{value: normalized}, nil
}

// DefaultMessageType returns the default message type.
func DefaultMessageType() PostType {
	return PostType{value: MessagePostType}
}

// Value returns the string value of the type.
func (t PostType) Value() string {
	return t.value
}

// String returns the string representation of the type.
func (t PostType) String() string {
	return t.value
}

// IsZero indicates if the type is unset.
func (t PostType) IsZero() bool {
	return t.value == ""
}

// IsAnnouncement indicates if the type is announcement.
func (t PostType) IsAnnouncement() bool {
	return t.value == AnnouncementPostType
}

// IsMessage indicates if the type is a message.
func (t PostType) IsMessage() bool {
	return t.value == MessagePostType
}
