package valueobjects

import (
	"slices"
	"strings"
)

// PostImages encapsulates the images attached to a post.
type PostImages struct {
	values []ImageURL `json:"values" bson:"images"`
}

// NewPostImages validates image URLs and removes duplicates.
func NewPostImages(imageURLs []string) (PostImages, error) {
	if len(imageURLs) == 0 {
		return PostImages{values: []ImageURL{}}, nil
	}

	seen := make(map[string]struct{}, len(imageURLs))
	values := make([]ImageURL, 0, len(imageURLs))

	for _, raw := range imageURLs {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}

		imageURL, err := NewImageURL(trimmed)
		if err != nil {
			return PostImages{}, err
		}

		if _, exists := seen[imageURL.Value()]; exists {
			continue
		}
		seen[imageURL.Value()] = struct{}{}
		values = append(values, imageURL)
	}

	return PostImages{values: values}, nil
}

// PostImagesFromValueObjects reconstructs PostImages from validated value objects.
func PostImagesFromValueObjects(images []ImageURL) PostImages {
	copied := slices.Clone(images)
	if copied == nil {
		copied = []ImageURL{}
	}
	return PostImages{values: copied}
}

// Values returns the underlying image value objects.
func (p PostImages) Values() []ImageURL {
	return slices.Clone(p.values)
}

// URLs returns the image URLs as strings.
func (p PostImages) URLs() []string {
	if len(p.values) == 0 {
		return []string{}
	}
	urls := make([]string, 0, len(p.values))
	for _, image := range p.values {
		urls = append(urls, image.Value())
	}
	return urls
}

// IsEmpty indicates whether there are no images.
func (p PostImages) IsEmpty() bool {
	return len(p.values) == 0
}
