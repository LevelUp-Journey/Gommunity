package commands

import (
	"Gommunity/platform/users/domain/model/valueobjects"
	"errors"
	"net/url"
)

type UpdateBannerURLCommand struct {
	userID    valueobjects.UserID
	bannerURL string
}

func NewUpdateBannerURLCommand(userID valueobjects.UserID, bannerURL string) (UpdateBannerURLCommand, error) {
	if userID.IsZero() {
		return UpdateBannerURLCommand{}, errors.New("userID cannot be empty")
	}

	if bannerURL == "" {
		return UpdateBannerURLCommand{}, errors.New("bannerURL cannot be empty")
	}

	// Validate URL format
	parsedURL, err := url.ParseRequestURI(bannerURL)
	if err != nil {
		return UpdateBannerURLCommand{}, errors.New("bannerURL must be a valid URL")
	}

	// Validate URL scheme (must be http or https)
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return UpdateBannerURLCommand{}, errors.New("bannerURL must use http or https scheme")
	}

	return UpdateBannerURLCommand{
		userID:    userID,
		bannerURL: bannerURL,
	}, nil
}

func (c UpdateBannerURLCommand) UserID() valueobjects.UserID {
	return c.userID
}

func (c UpdateBannerURLCommand) BannerURL() string {
	return c.bannerURL
}
