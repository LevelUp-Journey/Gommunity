package commands

import (
	"Gommunity/platform/users/domain/model/valueobjects"
)

type UpdateBannerURLCommand struct {
	userID    valueobjects.UserID
	bannerURL string
}

func NewUpdateBannerURLCommand(userID valueobjects.UserID, bannerURL string) UpdateBannerURLCommand {
	return UpdateBannerURLCommand{
		userID:    userID,
		bannerURL: bannerURL,
	}
}

func (c UpdateBannerURLCommand) UserID() valueobjects.UserID {
	return c.userID
}

func (c UpdateBannerURLCommand) BannerURL() string {
	return c.bannerURL
}
