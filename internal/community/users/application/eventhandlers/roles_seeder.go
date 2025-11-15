package eventhandlers

import (
	"context"
	"log"

	"Gommunity/internal/community/users/domain/model/entities"
	"Gommunity/internal/community/users/domain/model/valueobjects"
	domain_repos "Gommunity/internal/community/users/domain/repositories"
)

// SeedRoles seeds the default roles if they don't exist
func SeedRoles(ctx context.Context, roleRepo domain_repos.RoleRepository) error {
	roleNames := []string{"student", "teacher", "admin", "member", "owner"}

	for _, name := range roleNames {
		roleID, err := valueobjects.NewRoleID(name)
		if err != nil {
			return err
		}

		// Check if role exists
		existing, err := roleRepo.FindByID(ctx, roleID)
		if err != nil {
			return err
		}

		if existing == nil {
			// Create role
			role, err := entities.NewRole(roleID, name)
			if err != nil {
				return err
			}

			if err := roleRepo.Save(ctx, role); err != nil {
				return err
			}

			log.Printf("Seeded role: %s", name)
		}
	}

	return nil
}
