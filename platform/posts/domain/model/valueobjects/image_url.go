package valueobjects

import (
	"errors"
	"net/url"
	"strings"
)

// ImageURL represents a validated image location.
type ImageURL struct {
	value string `json:"value" bson:"image_url"`
}

// NewImageURL validates and creates an ImageURL.
func NewImageURL(value string) (ImageURL, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ImageURL{}, errors.New("image URL cannot be empty")
	}

	parsed, err := url.Parse(trimmed)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ImageURL{}, errors.New("image URL must be an absolute URL")
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ImageURL{}, errors.New("image URL must use http or https scheme")
	}

	return ImageURL{value: trimmed}, nil
}

// Value returns the URL string.
func (i ImageURL) Value() string {
	return i.value
}

// String returns the URL string.
func (i ImageURL) String() string {
	return i.value
}

// IsZero indicates whether the URL is empty.
func (i ImageURL) IsZero() bool {
	return i.value == ""
}
