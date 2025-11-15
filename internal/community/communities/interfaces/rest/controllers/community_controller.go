package controllers

import (
	"log"
	"net/http"

	"Gommunity/internal/community/communities/domain/model/commands"
	"Gommunity/internal/community/communities/domain/model/entities"
	"Gommunity/internal/community/communities/domain/model/queries"
	"Gommunity/internal/community/communities/domain/model/valueobjects"
	"Gommunity/internal/community/communities/domain/services"
	"Gommunity/internal/community/communities/interfaces/rest/resources"
	"Gommunity/shared/infrastructure/middleware"

	"github.com/gin-gonic/gin"
)

type CommunityController struct {
	commandService services.CommunityCommandService
	queryService   services.CommunityQueryService
}

func NewCommunityController(
	commandService services.CommunityCommandService,
	queryService services.CommunityQueryService,
) *CommunityController {
	return &CommunityController{
		commandService: commandService,
		queryService:   queryService,
	}
}

// CreateCommunity godoc
// @Summary Create a new community
// @Description Create a new community (only for ROLE_TEACHER and ROLE_ADMIN)
// @Tags communities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body resources.CreateCommunityResource true "Community creation request"
// @Success 201 {object} resources.CommunityResource
// @Failure 400 {object} resources.ErrorResponse
// @Failure 401 {object} resources.ErrorResponse
// @Failure 403 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /communities [post]
func (c *CommunityController) CreateCommunity(ctx *gin.Context) {
	// Get authenticated user ID from context
	authenticatedUserID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, resources.ErrorResponse{
			Error: "Authentication required",
		})
		return
	}

	// Get user role from context
	role, err := middleware.GetRoleFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, resources.ErrorResponse{
			Error: "Role not found in token",
		})
		return
	}

	log.Printf("Creating community - User: %s, Role: '%s'", authenticatedUserID, role)

	// Verify user has ROLE_TEACHER or ROLE_ADMIN
	if role != "ROLE_TEACHER" && role != "ROLE_ADMIN" {
		log.Printf("Access denied - Role '%s' is not authorized to create communities", role)
		ctx.JSON(http.StatusForbidden, resources.ErrorResponse{
			Error: "Only teachers and admins can create communities",
		})
		return
	}

	var req resources.CreateCommunityResource
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Create value objects
	ownerID, err := valueobjects.NewOwnerID(authenticatedUserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: "Invalid owner ID",
		})
		return
	}

	name, err := valueobjects.NewCommunityName(req.Name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	description, err := valueobjects.NewDescription(req.Description)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Create command
	cmd, err := commands.NewCreateCommunityCommand(ownerID, name, description, req.IconURL, req.BannerURL)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Execute command
	communityID, err := c.commandService.HandleCreate(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{
			Error: "Failed to create community",
		})
		return
	}

	// Retrieve created community
	query, _ := queries.NewGetCommunityByIDQuery(*communityID)
	community, err := c.queryService.HandleGetByID(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{
			Error: "Failed to retrieve created community",
		})
		return
	}

	response := c.transformCommunityToResource(community)
	ctx.JSON(http.StatusCreated, response)
}

// GetCommunityByID godoc
// @Summary Get community by ID
// @Description Get a specific community by its ID
// @Tags communities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Community ID (UUID)"
// @Success 200 {object} resources.CommunityResource
// @Failure 400 {object} resources.ErrorResponse
// @Failure 401 {object} resources.ErrorResponse
// @Failure 404 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /communities/{id} [get]
func (c *CommunityController) GetCommunityByID(ctx *gin.Context) {
	communityIDParam := ctx.Param("id")

	communityID, err := valueobjects.NewCommunityID(communityIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: "Invalid community ID format",
		})
		return
	}

	query, err := queries.NewGetCommunityByIDQuery(communityID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	community, err := c.queryService.HandleGetByID(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{
			Error: "Failed to retrieve community",
		})
		return
	}

	if community == nil {
		ctx.JSON(http.StatusNotFound, resources.ErrorResponse{
			Error: "Community not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, c.transformCommunityToResource(community))
}

// GetMyCommunitiesAsOwner godoc
// @Summary Get my communities as owner
// @Description Get all communities owned by the authenticated user
// @Tags communities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} resources.CommunityResource
// @Failure 401 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /communities/my-communities [get]
func (c *CommunityController) GetMyCommunitiesAsOwner(ctx *gin.Context) {
	// Get authenticated user ID from context
	authenticatedUserID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, resources.ErrorResponse{
			Error: "Authentication required",
		})
		return
	}

	ownerID, err := valueobjects.NewOwnerID(authenticatedUserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: "Invalid owner ID",
		})
		return
	}

	query, err := queries.NewGetCommunitiesByOwnerQuery(ownerID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	communities, err := c.queryService.HandleGetByOwner(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{
			Error: "Failed to retrieve communities",
		})
		return
	}

	response := make([]resources.CommunityResource, 0, len(communities))
	for _, community := range communities {
		response = append(response, c.transformCommunityToResource(community))
	}

	ctx.JSON(http.StatusOK, response)
}

// GetAllCommunities godoc
// @Summary Get all communities
// @Description Get all communities (paginated)
// @Tags communities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} resources.CommunityResource
// @Failure 401 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /communities [get]
func (c *CommunityController) GetAllCommunities(ctx *gin.Context) {
	query := queries.NewGetAllCommunitiesQuery()

	communities, err := c.queryService.HandleGetAll(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{
			Error: "Failed to retrieve communities",
		})
		return
	}

	response := make([]resources.CommunityResource, 0, len(communities))
	for _, community := range communities {
		response = append(response, c.transformCommunityToResource(community))
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteCommunity godoc
// @Summary Delete community
// @Description Delete a community (only owner can delete)
// @Tags communities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Community ID (UUID)"
// @Success 204
// @Failure 400 {object} resources.ErrorResponse
// @Failure 401 {object} resources.ErrorResponse
// @Failure 403 {object} resources.ErrorResponse
// @Failure 404 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /communities/{id} [delete]
func (c *CommunityController) DeleteCommunity(ctx *gin.Context) {
	communityIDParam := ctx.Param("id")

	// Get authenticated user ID from context
	authenticatedUserID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, resources.ErrorResponse{
			Error: "Authentication required",
		})
		return
	}

	communityID, err := valueobjects.NewCommunityID(communityIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: "Invalid community ID format",
		})
		return
	}

	ownerID, err := valueobjects.NewOwnerID(authenticatedUserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: "Invalid owner ID",
		})
		return
	}

	cmd, err := commands.NewDeleteCommunityCommand(communityID, ownerID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := c.commandService.HandleDelete(ctx.Request.Context(), cmd); err != nil {
		if err.Error() == "community not found" {
			ctx.JSON(http.StatusNotFound, resources.ErrorResponse{
				Error: "Community not found",
			})
			return
		}
		if err.Error() == "only the owner can delete the community" {
			ctx.JSON(http.StatusForbidden, resources.ErrorResponse{
				Error: "Only the owner can delete the community",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{
			Error: "Failed to delete community",
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// UpdateCommunityInfo godoc
// @Summary Update community information
// @Description Update community name, description, icon and banner (only owner can update)
// @Tags communities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Community ID (UUID)"
// @Param request body resources.UpdateCommunityResource true "Community update request"
// @Success 200 {object} resources.CommunityResource
// @Failure 400 {object} resources.ErrorResponse
// @Failure 401 {object} resources.ErrorResponse
// @Failure 403 {object} resources.ErrorResponse
// @Failure 404 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /communities/{id} [put]
func (c *CommunityController) UpdateCommunityInfo(ctx *gin.Context) {
	communityIDParam := ctx.Param("id")

	// Get authenticated user ID from context
	authenticatedUserID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, resources.ErrorResponse{
			Error: "Authentication required",
		})
		return
	}

	communityID, err := valueobjects.NewCommunityID(communityIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: "Invalid community ID format",
		})
		return
	}

	var req resources.UpdateCommunityResource
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Validate that at least one field is provided
	if req.Name == nil && req.Description == nil && req.IconURL == nil && req.BannerURL == nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: "At least one field must be provided for update",
		})
		return
	}

	// First, verify the community exists and the user is the owner
	query, err := queries.NewGetCommunityByIDQuery(communityID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	community, err := c.queryService.HandleGetByID(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{
			Error: "Failed to retrieve community",
		})
		return
	}

	if community == nil {
		ctx.JSON(http.StatusNotFound, resources.ErrorResponse{
			Error: "Community not found",
		})
		return
	}

	// Verify the user is the owner
	if !community.IsOwner(authenticatedUserID) {
		ctx.JSON(http.StatusForbidden, resources.ErrorResponse{
			Error: "Only the owner can update community information",
		})
		return
	}

	// Use existing values if not provided in request
	name := community.Name()
	if req.Name != nil {
		name, err = valueobjects.NewCommunityName(*req.Name)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
				Error: err.Error(),
			})
			return
		}
	}

	description := community.Description()
	if req.Description != nil {
		description, err = valueobjects.NewDescription(*req.Description)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
				Error: err.Error(),
			})
			return
		}
	}

	// Create and execute command
	cmd, err := commands.NewUpdateCommunityInfoCommand(communityID, name, description, req.IconURL, req.BannerURL)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := c.commandService.HandleUpdateInfo(ctx.Request.Context(), cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{
			Error: "Failed to update community information",
		})
		return
	}

	// Retrieve updated community
	updatedCommunity, err := c.queryService.HandleGetByID(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{
			Error: "Failed to retrieve updated community",
		})
		return
	}

	response := c.transformCommunityToResource(updatedCommunity)
	ctx.JSON(http.StatusOK, response)
}

// UpdateCommunityPrivacy godoc
// @Summary Update community privacy status
// @Description Update the privacy status of a community (only owner can update)
// @Tags communities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Community ID (UUID)"
// @Param request body resources.UpdateCommunityPrivacyResource true "Community privacy update request"
// @Success 200 {object} resources.CommunityResource
// @Failure 400 {object} resources.ErrorResponse
// @Failure 401 {object} resources.ErrorResponse
// @Failure 403 {object} resources.ErrorResponse
// @Failure 404 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /communities/{id}/privacy [patch]
func (c *CommunityController) UpdateCommunityPrivacy(ctx *gin.Context) {
	communityIDParam := ctx.Param("id")

	// Get authenticated user ID from context
	authenticatedUserID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, resources.ErrorResponse{
			Error: "Authentication required",
		})
		return
	}

	communityID, err := valueobjects.NewCommunityID(communityIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: "Invalid community ID format",
		})
		return
	}

	var req resources.UpdateCommunityPrivacyResource
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// First, verify the community exists and the user is the owner
	query, err := queries.NewGetCommunityByIDQuery(communityID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	community, err := c.queryService.HandleGetByID(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{
			Error: "Failed to retrieve community",
		})
		return
	}

	if community == nil {
		ctx.JSON(http.StatusNotFound, resources.ErrorResponse{
			Error: "Community not found",
		})
		return
	}

	// Verify the user is the owner
	if !community.IsOwner(authenticatedUserID) {
		ctx.JSON(http.StatusForbidden, resources.ErrorResponse{
			Error: "Only the owner can update community privacy",
		})
		return
	}

	// Create and execute command
	cmd, err := commands.NewUpdateCommunityPrivacyCommand(communityID, req.IsPrivate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := c.commandService.HandleUpdatePrivacy(ctx.Request.Context(), cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{
			Error: "Failed to update community privacy",
		})
		return
	}

	// Retrieve updated community
	updatedCommunity, err := c.queryService.HandleGetByID(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{
			Error: "Failed to retrieve updated community",
		})
		return
	}

	response := c.transformCommunityToResource(updatedCommunity)
	ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) transformCommunityToResource(community *entities.Community) resources.CommunityResource {
	return resources.CommunityResource{
		ID:          community.ID(),
		CommunityID: community.CommunityID().Value(),
		OwnerID:     community.OwnerID().Value(),
		Name:        community.Name().Value(),
		Description: community.Description().Value(),
		IconURL:     community.IconURL(),
		BannerURL:   community.BannerURL(),
		IsPrivate:   community.IsPrivate(),
		CreatedAt:   community.CreatedAt(),
		UpdatedAt:   community.UpdatedAt(),
	}
}
