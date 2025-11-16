package commandservices

import (
	"context"
	"errors"

	"Gommunity/platform/users/domain/model/commands"
	"Gommunity/platform/users/domain/repositories"
	"Gommunity/platform/users/domain/services"
)

type userCommandServiceImpl struct {
	userRepository repositories.UserRepository
}

func NewUserCommandService(userRepository repositories.UserRepository) services.UserCommandService {
	return &userCommandServiceImpl{
		userRepository: userRepository,
	}
}

func (s *userCommandServiceImpl) HandleUpdateBanner(ctx context.Context, cmd commands.UpdateBannerURLCommand) error {
	// Find user
	user, err := s.userRepository.FindByUserID(ctx, cmd.UserID())
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	// Update banner URL
	bannerURL := cmd.BannerURL()
	user.UpdateProfile(user.Username(), user.ProfileURL(), &bannerURL)

	// Save updated user
	return s.userRepository.Update(ctx, user)
}

func (s *userCommandServiceImpl) HandleUpdateRole(ctx context.Context, cmd commands.UpdateRoleCommand) error {
	// Find user
	user, err := s.userRepository.FindByUserID(ctx, cmd.UserID())
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	// Update role
	user.UpdateRoleID(cmd.RoleID())

	// Save updated user
	return s.userRepository.Update(ctx, user)
}
