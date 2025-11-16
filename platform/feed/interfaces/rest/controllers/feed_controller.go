package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Gommunity/platform/feed/domain/model/queries"
	"Gommunity/platform/feed/domain/model/valueobjects"
	"Gommunity/platform/feed/domain/services"
	"Gommunity/platform/feed/interfaces/rest/resources"
)

type FeedController struct {
	queryService services.FeedQueryService
}

func NewFeedController(queryService services.FeedQueryService) *FeedController {
	return &FeedController{
		queryService: queryService,
	}
}

// GetUserFeed retrieves the feed for the authenticated user
// @Summary Get user feed
// @Description Retrieves announcements from all communities the user is subscribed to
// @Tags feed
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page" default(20)
// @Param offset query int false "Page offset" default(0)
// @Success 200 {object} resources.FeedResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /feed [get]
func (c *FeedController) GetUserFeed(ctx *gin.Context) {
	// Get user ID from JWT token
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID format"})
		return
	}

	// Parse pagination parameters
	limit := 20
	offset := 0

	if limitParam := ctx.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetParam := ctx.Query("offset"); offsetParam != "" {
		if parsedOffset, err := strconv.Atoi(offsetParam); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Create user ID value object
	userIDVO, err := valueobjects.NewUserID(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create and execute query
	query, err := queries.NewGetUserFeedQuery(userIDVO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query = query.WithPagination(limit, offset)

	feedItems, err := c.queryService.Handle(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve feed"})
		return
	}

	// Transform to resources
	items := make([]resources.FeedItemResource, len(feedItems))
	for i, item := range feedItems {
		items[i] = resources.FeedItemResource{
			PostID:      item.PostID().Value(),
			CommunityID: item.CommunityID().Value(),
			AuthorID:    item.AuthorID(),
			Content:     item.Content(),
			MessageType: item.MessageType(),
			CreatedAt:   item.CreatedAt(),
			UpdatedAt:   item.UpdatedAt(),
		}
	}

	response := resources.FeedResponse{
		Items: items,
		Total: len(items),
	}

	ctx.JSON(http.StatusOK, response)
}
