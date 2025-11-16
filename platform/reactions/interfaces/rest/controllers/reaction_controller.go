package controllers

import (
	"net/http"
	"strings"

	"Gommunity/platform/reactions/domain/model/commands"
	"Gommunity/platform/reactions/domain/model/entities"
	"Gommunity/platform/reactions/domain/model/queries"
	"Gommunity/platform/reactions/domain/model/valueobjects"
	"Gommunity/platform/reactions/domain/services"
	"Gommunity/platform/reactions/interfaces/rest/resources"
	"Gommunity/shared/infrastructure/middleware"

	"github.com/gin-gonic/gin"
)

// ReactionController handles HTTP requests for reactions.
type ReactionController struct {
	commandService services.ReactionCommandService
	queryService   services.ReactionQueryService
}

// NewReactionController builds a ReactionController.
func NewReactionController(
	commandService services.ReactionCommandService,
	queryService services.ReactionQueryService,
) *ReactionController {
	return &ReactionController{
		commandService: commandService,
		queryService:   queryService,
	}
}

// AddReaction godoc
// @Summary Add or update a reaction to a post
// @Description Allows a user to add a reaction to a post. If the user already reacted, the reaction type will be updated.
// @Tags reactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param post_id path string true "Post ID (ObjectID)"
// @Param request body resources.AddReactionResource true "Reaction payload"
// @Success 201 {object} resources.ReactionResource
// @Failure 400 {object} resources.ErrorResponse
// @Failure 401 {object} resources.ErrorResponse
// @Failure 404 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /api/v1/posts/{post_id}/reactions [post]
func (c *ReactionController) AddReaction(ctx *gin.Context) {
	postIDValue := ctx.Param("post_id")

	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, resources.ErrorResponse{Error: "authentication required"})
		return
	}

	var req resources.AddReactionResource
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "invalid request body"})
		return
	}

	postID, err := valueobjects.NewPostID(postIDValue)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "invalid post id"})
		return
	}

	userIDVO, err := valueobjects.NewUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "invalid user id"})
		return
	}

	reactionType, err := valueobjects.NewReactionType(req.ReactionType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: err.Error()})
		return
	}

	cmd, err := commands.NewAddReactionCommand(postID, userIDVO, reactionType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: err.Error()})
		return
	}

	_, err = c.commandService.HandleAdd(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(mapErrorToStatus(err), resources.ErrorResponse{Error: err.Error()})
		return
	}

	// Retrieve the created/updated reaction
	query, _ := queries.NewGetUserReactionOnPostQuery(postID, userIDVO)
	reaction, err := c.queryService.HandleGetUserReactionOnPost(ctx.Request.Context(), query)
	if err != nil || reaction == nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{Error: "unable to retrieve reaction"})
		return
	}

	ctx.JSON(http.StatusCreated, c.toResource(reaction))
}

// RemoveReaction godoc
// @Summary Remove a reaction from a post
// @Description Allows a user to remove their reaction from a post.
// @Tags reactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param post_id path string true "Post ID (ObjectID)"
// @Success 204 "Reaction removed"
// @Failure 400 {object} resources.ErrorResponse
// @Failure 401 {object} resources.ErrorResponse
// @Failure 404 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /api/v1/posts/{post_id}/reactions [delete]
func (c *ReactionController) RemoveReaction(ctx *gin.Context) {
	postIDValue := ctx.Param("post_id")

	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, resources.ErrorResponse{Error: "authentication required"})
		return
	}

	postID, err := valueobjects.NewPostID(postIDValue)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "invalid post id"})
		return
	}

	userIDVO, err := valueobjects.NewUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "invalid user id"})
		return
	}

	cmd, err := commands.NewRemoveReactionCommand(postID, userIDVO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: err.Error()})
		return
	}

	if err := c.commandService.HandleRemove(ctx.Request.Context(), cmd); err != nil {
		ctx.JSON(mapErrorToStatus(err), resources.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetReactionCountByPost godoc
// @Summary Get reaction counts for a post
// @Description Retrieves aggregated reaction counts grouped by reaction type for a specific post.
// @Tags reactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param post_id path string true "Post ID (ObjectID)"
// @Success 200 {object} resources.ReactionCountResource
// @Failure 400 {object} resources.ErrorResponse
// @Failure 401 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /api/v1/posts/{post_id}/reactions/count [get]
func (c *ReactionController) GetReactionCountByPost(ctx *gin.Context) {
	postIDValue := ctx.Param("post_id")

	postID, err := valueobjects.NewPostID(postIDValue)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "invalid post id"})
		return
	}

	query, _ := queries.NewGetReactionCountByPostQuery(postID)
	summary, err := c.queryService.HandleGetCountByPost(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{Error: "failed to retrieve reaction counts"})
		return
	}

	response := resources.ReactionCountResource{
		PostID:     postID.Value(),
		TotalCount: summary.TotalCount,
		Counts:     summary.Counts,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetUserReactionOnPost godoc
// @Summary Get current user's reaction on a post
// @Description Retrieves the authenticated user's reaction on a specific post.
// @Tags reactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param post_id path string true "Post ID (ObjectID)"
// @Success 200 {object} resources.ReactionResource
// @Failure 400 {object} resources.ErrorResponse
// @Failure 401 {object} resources.ErrorResponse
// @Failure 404 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /api/v1/posts/{post_id}/reactions/me [get]
func (c *ReactionController) GetUserReactionOnPost(ctx *gin.Context) {
	postIDValue := ctx.Param("post_id")

	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, resources.ErrorResponse{Error: "authentication required"})
		return
	}

	postID, err := valueobjects.NewPostID(postIDValue)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "invalid post id"})
		return
	}

	userIDVO, err := valueobjects.NewUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "invalid user id"})
		return
	}

	query, _ := queries.NewGetUserReactionOnPostQuery(postID, userIDVO)
	reaction, err := c.queryService.HandleGetUserReactionOnPost(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{Error: "failed to retrieve user reaction"})
		return
	}

	if reaction == nil {
		ctx.JSON(http.StatusNotFound, resources.ErrorResponse{Error: "no reaction found"})
		return
	}

	ctx.JSON(http.StatusOK, c.toResource(reaction))
}

func (c *ReactionController) toResource(reaction *entities.Reaction) resources.ReactionResource {
	return resources.ReactionResource{
		ReactionID:   reaction.ReactionID().Value(),
		PostID:       reaction.PostID().Value(),
		UserID:       reaction.UserID().Value(),
		ReactionType: reaction.ReactionType().Value(),
		CreatedAt:    reaction.CreatedAt(),
		UpdatedAt:    reaction.UpdatedAt(),
	}
}

func mapErrorToStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	lower := strings.ToLower(err.Error())
	switch {
	case strings.Contains(lower, "not found"):
		return http.StatusNotFound
	case strings.Contains(lower, "already"):
		return http.StatusConflict
	case strings.Contains(lower, "not allowed"):
		return http.StatusForbidden
	default:
		return http.StatusBadRequest
	}
}
