package valueobjects

import (
	"errors"
	"strings"
)

// PostContent represents the markdown content of a post.
type PostContent struct {
	value string `json:"value" bson:"content"`
}

// NewPostContent validates post content.
// Content can be plain text or markdown - no format restrictions.
func NewPostContent(value string) (PostContent, error) {
	if strings.TrimSpace(value) == "" {
		return PostContent{}, errors.New("post content cannot be empty")
	}

	// Users can write anything they want - no format restrictions
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
