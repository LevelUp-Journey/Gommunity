package controllers

import (
	"net/http"

	"Gommunity/platform/users/domain/model/commands"
	"Gommunity/platform/users/domain/model/entities"
	"Gommunity/platform/users/domain/model/queries"
	"Gommunity/platform/users/domain/model/valueobjects"
	"Gommunity/platform/users/domain/services"
	"Gommunity/platform/users/interfaces/rest/resources"
	"Gommunity/shared/infrastructure/middleware"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	commandService services.UserCommandService
	queryService   services.UserQueryService
}

func NewUserController(
	commandService services.UserCommandService,
	queryService services.UserQueryService,
) *UserController {
	return &UserController{
		commandService: commandService,
		queryService:   queryService,
	}
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get a specific user by their user ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} resources.UserResource
// @Failure 400 {object} resources.ErrorResponse
// @Failure 401 {object} resources.ErrorResponse
// @Failure 404 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /api/v1/users/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	userIDParam := ctx.Param("id")

	userID, err := valueobjects.NewUserID(userIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: "Invalid user ID format",
		})
		return
	}

	query, err := queries.NewGetUserByIDQuery(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	user, err := c.queryService.HandleGetByID(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{
			Error: "Failed to retrieve user",
		})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusNotFound, resources.ErrorResponse{
			Error: "User not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, c.transformUserToResource(user))
}

// GetUserByUsername godoc
// @Summary Get user by username
// @Description Get a specific user by their username
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Success 200 {object} resources.UserResource
// @Failure 400 {object} resources.ErrorResponse
// @Failure 401 {object} resources.ErrorResponse
// @Failure 404 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /api/v1/users/username/{username} [get]
func (c *UserController) GetUserByUsername(ctx *gin.Context) {
	usernameParam := ctx.Param("username")

	username, err := valueobjects.NewUsername(usernameParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: "Invalid username format",
		})
		return
	}

	query, err := queries.NewGetUserByUsernameQuery(username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	user, err := c.queryService.HandleGetByUsername(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{
			Error: "Failed to retrieve user",
		})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusNotFound, resources.ErrorResponse{
			Error: "User not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, c.transformUserToResource(user))
}

// UpdateBannerURL godoc
// @Summary Update user banner URL
// @Description Update the banner URL for the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Param request body resources.UpdateBannerURLResource true "Banner URL update request"
// @Success 200 {object} resources.UserResource
// @Failure 400 {object} resources.ErrorResponse
// @Failure 401 {object} resources.ErrorResponse
// @Failure 403 {object} resources.ErrorResponse
// @Failure 404 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /api/v1/users/{id}/banner [put]
func (c *UserController) UpdateBannerURL(ctx *gin.Context) {
	userIDParam := ctx.Param("id")

	// Get authenticated user ID from context
	authenticatedUserID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, resources.ErrorResponse{
			Error: "Authentication required",
		})
		return
	}

	// Validate that the user is updating their own banner
	if authenticatedUserID != userIDParam {
		ctx.JSON(http.StatusForbidden, resources.ErrorResponse{
			Error: "You can only update your own banner",
		})
		return
	}

	var req resources.UpdateBannerURLResource
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	userID, err := valueobjects.NewUserID(userIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: "Invalid user ID format",
		})
		return
	}

	cmd := commands.NewUpdateBannerURLCommand(userID, req.BannerURL)

	if err := c.commandService.HandleUpdateBanner(ctx.Request.Context(), cmd); err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, resources.ErrorResponse{
				Error: "User not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{
			Error: "Failed to update banner URL",
		})
		return
	}

	// Retrieve updated user
	query, _ := queries.NewGetUserByIDQuery(userID)
	user, err := c.queryService.HandleGetByID(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{
			Error: "Failed to retrieve updated user",
		})
		return
	}

	ctx.JSON(http.StatusOK, c.transformUserToResource(user))
}

func (c *UserController) transformUserToResource(user *entities.User) resources.UserResource {
	// Note: Users BC no longer manages roles
	// Roles are managed per-community in Subscriptions BC
	// IAM roles (STUDENT, TEACHER, ADMIN) come from JWT token
	return resources.UserResource{
		ID:         user.ID(),
		UserID:     user.UserID().Value(),
		ProfileID:  user.ProfileID().Value(),
		Username:   user.Username().Value(),
		ProfileURL: user.ProfileURL(),
		BannerURL:  user.BannerURL(),
		UpdatedAt:  user.UpdatedAt(),
		CreatedAt:  user.CreatedAt(),
	}
}
