package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"Gommunity/platform/posts/domain/model/commands"
	"Gommunity/platform/posts/domain/model/entities"
	"Gommunity/platform/posts/domain/model/queries"
	"Gommunity/platform/posts/domain/model/valueobjects"
	"Gommunity/platform/posts/domain/services"
	"Gommunity/platform/posts/interfaces/rest/resources"
	"Gommunity/shared/infrastructure/middleware"

	"github.com/gin-gonic/gin"
)

// PostController handles HTTP requests for posts.
type PostController struct {
	commandService services.PostCommandService
	queryService   services.PostQueryService
}

// NewPostController builds a PostController.
func NewPostController(
	commandService services.PostCommandService,
	queryService services.PostQueryService,
) *PostController {
	return &PostController{
		commandService: commandService,
		queryService:   queryService,
	}
}

// CreatePost godoc
// @Summary Publish a new post
// @Description Members can create messages while admins or owners can also create announcements.
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param community_id path string true "Community ID (UUID)"
// @Param request body resources.CreatePostResource true "Post payload"
// @Success 201 {object} resources.PostResource
// @Failure 400 {object} resources.ErrorResponse
// @Failure 401 {object} resources.ErrorResponse
// @Failure 403 {object} resources.ErrorResponse
// @Failure 404 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /api/v1/communities/{community_id}/posts [post]
func (c *PostController) CreatePost(ctx *gin.Context) {
	communityIDValue := ctx.Param("community_id")

	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, resources.ErrorResponse{Error: "authentication required"})
		return
	}

	var req resources.CreatePostResource
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "invalid request body"})
		return
	}

	communityID, err := valueobjects.NewCommunityID(communityIDValue)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "invalid community id"})
		return
	}

	authorID, err := valueobjects.NewAuthorID(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "invalid author id"})
		return
	}

	postType := valueobjects.DefaultMessageType()
	if req.Type != "" {
		postType, err = valueobjects.NewPostType(req.Type)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: err.Error()})
			return
		}
	}

	content, err := valueobjects.NewPostContent(req.Content)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: err.Error()})
		return
	}

	images, err := valueobjects.NewPostImages(req.Images)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: err.Error()})
		return
	}

	cmd, err := commands.NewCreatePostCommand(communityID, authorID, postType, content, images)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: err.Error()})
		return
	}

	postID, err := c.commandService.HandlePublish(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(mapErrorToStatus(err), resources.ErrorResponse{Error: err.Error()})
		return
	}

	getQuery, _ := queries.NewGetPostByIDQuery(*postID)
	post, err := c.queryService.HandleGetByID(ctx.Request.Context(), getQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{Error: "unable to retrieve created post"})
		return
	}

	ctx.JSON(http.StatusCreated, c.toResource(post))
}

// GetPostByID godoc
// @Summary Get post by ID
// @Description Retrieves a post using its identifier within a community.
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param community_id path string true "Community ID (UUID)"
// @Param post_id path string true "Post ID (ObjectID)"
// @Success 200 {object} resources.PostResource
// @Failure 400 {object} resources.ErrorResponse
// @Failure 401 {object} resources.ErrorResponse
// @Failure 404 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /api/v1/communities/{community_id}/posts/{post_id} [get]
func (c *PostController) GetPostByID(ctx *gin.Context) {
	communityIDValue := ctx.Param("community_id")
	postIDValue := ctx.Param("post_id")

	communityID, err := valueobjects.NewCommunityID(communityIDValue)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "invalid community id"})
		return
	}

	postID, err := valueobjects.NewPostID(postIDValue)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "invalid post id"})
		return
	}

	query, _ := queries.NewGetPostByIDQuery(postID)
	post, err := c.queryService.HandleGetByID(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{Error: "failed to retrieve post"})
		return
	}

	if post == nil || post.CommunityID().Value() != communityID.Value() {
		ctx.JSON(http.StatusNotFound, resources.ErrorResponse{Error: "post not found"})
		return
	}

	ctx.JSON(http.StatusOK, c.toResource(post))
}

// GetPostsByCommunity godoc
// @Summary List posts by community
// @Description Retrieves posts published inside a community.
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param community_id path string true "Community ID (UUID)"
// @Param limit query int false "Limit results" minimum(1) maximum(100)
// @Param offset query int false "Skip results" minimum(0)
// @Success 200 {array} resources.PostResource
// @Failure 400 {object} resources.ErrorResponse
// @Failure 401 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /api/v1/communities/{community_id}/posts [get]
func (c *PostController) GetPostsByCommunity(ctx *gin.Context) {
	communityIDValue := ctx.Param("community_id")

	communityID, err := valueobjects.NewCommunityID(communityIDValue)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "invalid community id"})
		return
	}

	query, _ := queries.NewGetPostsByCommunityQuery(communityID)

	limitParam := 0
	offsetParam := -1

	if limitStr := ctx.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "limit must be a positive number"})
			return
		}
		if limit > 100 {
			ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "limit cannot exceed 100"})
			return
		}
		limitParam = limit
	}

	if offsetStr := ctx.Query("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "offset must be zero or positive"})
			return
		}
		offsetParam = offset
	}

	if limitParam > 0 || offsetParam >= 0 {
		query = query.WithPagination(limitParam, offsetParam)
	}

	posts, err := c.queryService.HandleGetByCommunity(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resources.ErrorResponse{Error: "failed to list posts"})
		return
	}

	response := make([]resources.PostResource, 0, len(posts))
	for _, post := range posts {
		response = append(response, c.toResource(post))
	}

	ctx.JSON(http.StatusOK, response)
}

// DeletePost godoc
// @Summary Delete a post
// @Description Only community admins or owners can delete posts.
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param community_id path string true "Community ID (UUID)"
// @Param post_id path string true "Post ID (ObjectID)"
// @Success 204 "Post deleted"
// @Failure 400 {object} resources.ErrorResponse
// @Failure 401 {object} resources.ErrorResponse
// @Failure 403 {object} resources.ErrorResponse
// @Failure 404 {object} resources.ErrorResponse
// @Failure 500 {object} resources.ErrorResponse
// @Router /api/v1/communities/{community_id}/posts/{post_id} [delete]
func (c *PostController) DeletePost(ctx *gin.Context) {
	communityIDValue := ctx.Param("community_id")
	postIDValue := ctx.Param("post_id")
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, resources.ErrorResponse{Error: "authentication required"})
		return
	}

	if _, err := valueobjects.NewCommunityID(communityIDValue); err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "invalid community id"})
		return
	}

	postID, err := valueobjects.NewPostID(postIDValue)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "invalid post id"})
		return
	}

	requesterID, err := valueobjects.NewAuthorID(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: "invalid requester id"})
		return
	}

	cmd, err := commands.NewDeletePostCommand(postID, requesterID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resources.ErrorResponse{Error: err.Error()})
		return
	}

	if err := c.commandService.HandleDelete(ctx.Request.Context(), cmd); err != nil {
		ctx.JSON(mapErrorToStatus(err), resources.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *PostController) toResource(post *entities.Post) resources.PostResource {
	return resources.PostResource{
		PostID:      post.PostID().Value(),
		CommunityID: post.CommunityID().Value(),
		AuthorID:    post.AuthorID().Value(),
		Type:        post.PostType().Value(),
		Content:     post.Content().Value(),
		Images:      post.Images().URLs(),
		CreatedAt:   post.CreatedAt(),
		UpdatedAt:   post.UpdatedAt(),
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
	case strings.Contains(lower, "only") || strings.Contains(lower, "not allowed") || strings.Contains(lower, "must"):
		return http.StatusForbidden
	default:
		return http.StatusBadRequest
	}
}
