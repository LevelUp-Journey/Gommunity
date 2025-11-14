package valueobjects

import (
	"errors"
	"strings"
)

type Description struct {
	value string `json:"value" bson:"description"`
}

func NewDescription(value string) (Description, error) {
	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return Description{}, errors.New("description cannot be empty")
	}

	if len(trimmed) < 10 {
		return Description{}, errors.New("description must be at least 10 characters long")
	}

	if len(trimmed) > 500 {
		return Description{}, errors.New("description cannot exceed 500 characters")
	}

	return Description{value: trimmed}, nil
}

func (d Description) Value() string {
	return d.value
}

func (d Description) String() string {
	return d.value
}

func (d Description) IsZero() bool {
	return d.value == ""
}
