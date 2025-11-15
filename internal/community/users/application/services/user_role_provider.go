package services

import (
	"context"
	"log"
	"strings"

	"Gommunity/internal/community/users/domain/model/valueobjects"
	"Gommunity/internal/community/users/domain/repositories"
)

type UserRoleProviderService struct {
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
}

func NewUserRoleProviderService(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
) *UserRoleProviderService {
	return &UserRoleProviderService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *UserRoleProviderService) GetUserRoleByUserID(ctx context.Context, userID string) (string, error) {
	// Parse userID
	userIDVO, err := valueobjects.NewUserID(userID)
	if err != nil {
		log.Printf("Invalid userID format: %s, error: %v", userID, err)
		return "", err
	}

	// Find user
	user, err := s.userRepo.FindByUserID(ctx, userIDVO)
	if err != nil {
		log.Printf("Error finding user %s: %v", userID, err)
		return "", err
	}

	if user == nil {
		log.Printf("User not found: %s", userID)
		return "", nil
	}

	// Get role
	role, err := s.roleRepo.FindByID(ctx, user.RoleID())
	if err != nil {
		log.Printf("Error finding role for user %s: %v", userID, err)
		return "", err
	}

	if role == nil {
		log.Printf("Role not found for user %s", userID)
		return "", nil
	}

	// Convert to uppercase format: "teacher" -> "ROLE_TEACHER"
	roleName := "ROLE_" + strings.ToUpper(role.Name())
	log.Printf("Fetched role for user %s: %s (from DB: %s)", userID, roleName, role.Name())

	return roleName, nil
}
