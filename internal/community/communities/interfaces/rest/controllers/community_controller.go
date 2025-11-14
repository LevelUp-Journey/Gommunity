package controllers

import (
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
// @Description Create a new community (only for ROLE_TEACHER)
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

	// Verify user has ROLE_TEACHER
	if role != "ROLE_TEACHER" && role != "ROLE_ADMIN" {
		ctx.JSON(http.StatusForbidden, resources.ErrorResponse{
			Error: "Only teachers can create communities",
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
	cmd, err := commands.NewCreateCommunityCommand(ownerID, name, description)
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

func (c *CommunityController) transformCommunityToResource(community *entities.Community) resources.CommunityResource {
	return resources.CommunityResource{
		ID:          community.ID(),
		CommunityID: community.CommunityID().Value(),
		OwnerID:     community.OwnerID().Value(),
		Name:        community.Name().Value(),
		Description: community.Description().Value(),
		LogoURL:     community.LogoURL(),
		BannerURL:   community.BannerURL(),
		IsActive:    community.IsActive(),
		CreatedAt:   community.CreatedAt(),
		UpdatedAt:   community.UpdatedAt(),
	}
}
