package commands

import (
	"Gommunity/internal/community/users/domain/model/valueobjects"
	"errors"
)

type UpdateRoleCommand struct {
	userID valueobjects.UserID
	roleID valueobjects.RoleID
}

func NewUpdateRoleCommand(userID valueobjects.UserID, roleID valueobjects.RoleID) (UpdateRoleCommand, error) {
	if userID.IsZero() {
		return UpdateRoleCommand{}, errors.New("userID cannot be empty")
	}
	if roleID.IsZero() {
		return UpdateRoleCommand{}, errors.New("roleID cannot be empty")
	}

	return UpdateRoleCommand{
		userID: userID,
		roleID: roleID,
	}, nil
}

func (c UpdateRoleCommand) UserID() valueobjects.UserID {
	return c.userID
}

func (c UpdateRoleCommand) RoleID() valueobjects.RoleID {
	return c.roleID
}
