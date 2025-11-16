package services

import (
	"context"

	"Gommunity/platform/community/domain/model/commands"
	"Gommunity/platform/community/domain/model/valueobjects"
)

type CommunityCommandService interface {
	HandleCreate(ctx context.Context, cmd commands.CreateCommunityCommand) (*valueobjects.CommunityID, error)
	HandleDelete(ctx context.Context, cmd commands.DeleteCommunityCommand) error
	HandleUpdatePrivacy(ctx context.Context, cmd commands.UpdateCommunityPrivacyCommand) error
	HandleUpdateInfo(ctx context.Context, cmd commands.UpdateCommunityInfoCommand) error
}
