package valueobjects

import (
	"errors"
	"strings"
)

// PostContent represents the markdown content of a post.
type PostContent struct {
	value string `json:"value" bson:"content"`
}

// NewPostContent validates markdown content ensuring it preserves line breaks.
func NewPostContent(value string) (PostContent, error) {
	if strings.TrimSpace(value) == "" {
		return PostContent{}, errors.New("post content cannot be empty")
	}

	if !containsLineBreak(value) {
		return PostContent{}, errors.New("post content must include at least one line break for markdown formatting")
	}

	return PostContent{value: value}, nil
}

// Value returns the raw markdown content.
func (c PostContent) Value() string {
	return c.value
}

// IsZero indicates if the content is empty.
func (c PostContent) IsZero() bool {
	return c.value == ""
}

func containsLineBreak(value string) bool {
	return strings.Contains(value, "\n") || strings.Contains(value, "\r")
}
