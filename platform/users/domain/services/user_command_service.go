package services

import (
	"context"

	"Gommunity/platform/users/domain/model/commands"
)

type UserCommandService interface {
	HandleUpdateBanner(ctx context.Context, cmd commands.UpdateBannerURLCommand) error
	HandleUpdateRole(ctx context.Context, cmd commands.UpdateRoleCommand) error
}
