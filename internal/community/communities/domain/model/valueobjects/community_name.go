package valueobjects

import (
	"errors"
	"strings"
)

type CommunityName struct {
	value string `json:"value" bson:"name"`
}

func NewCommunityName(value string) (CommunityName, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return CommunityName{}, errors.New("community name cannot be empty")
	}

	if len(trimmed) < 3 {
		return CommunityName{}, errors.New("community name must be at least 3 characters long")
	}

	if len(trimmed) > 100 {
		return CommunityName{}, errors.New("community name cannot exceed 100 characters")
	}

	return CommunityName{value: trimmed}, nil
}

func (c CommunityName) Value() string {
	return c.value
}

func (c CommunityName) String() string {
	return c.value
}

func (c CommunityName) IsZero() bool {
	return c.value == ""
}
