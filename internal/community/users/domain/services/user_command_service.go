package services

import (
	"context"

	"Gommunity/internal/community/users/domain/model/commands"
)

type UserCommandService interface {
	HandleUpdateBanner(ctx context.Context, cmd commands.UpdateBannerURLCommand) error
}
