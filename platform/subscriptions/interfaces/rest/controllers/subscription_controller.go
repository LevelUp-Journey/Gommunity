package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Gommunity/platform/subscriptions/application/outboundservices/acl"
	"Gommunity/platform/subscriptions/domain/model/commands"
	"Gommunity/platform/subscriptions/domain/model/queries"
	"Gommunity/platform/subscriptions/domain/model/valueobjects"
	"Gommunity/platform/subscriptions/domain/services"
	"Gommunity/platform/subscriptions/interfaces/rest/resources"
)

type SubscriptionController struct {
	commandService       services.SubscriptionCommandService
	queryService         services.SubscriptionQueryService
	externalUsersService *acl.ExternalUsersService
}

func NewSubscriptionController(
	commandService services.SubscriptionCommandService,
	queryService services.SubscriptionQueryService,
	externalUsersService *acl.ExternalUsersService,
) *SubscriptionController {
	return &SubscriptionController{
		commandService:       commandService,
		queryService:         queryService,
		externalUsersService: externalUsersService,
	}
}

// @Summary Subscribe a user to a community
// @Description Subscribe a user to a community with a specific role. IMPORTANT: Self-subscriptions (following a community) always receive 'member' role regardless of requested role. In public communities, users can only subscribe themselves. In private communities, owner/admin can add users by username and assign any role.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body resources.SubscribeUserResource true "Subscription request"
// @Success 201 {object} resources.SubscriptionResource
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /subscriptions [post]
func (c *SubscriptionController) SubscribeUser(ctx *gin.Context) {
	var req resources.SubscribeUserResource
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the requesting user ID from JWT context
	requestedByValue, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	requestedByID, ok := requestedByValue.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID in context"})
		return
	}

	requestedBy, err := valueobjects.NewUserID(requestedByID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid requesting user ID"})
		return
	}

	// Determine the target user ID (either from userID or username)
	var targetUserID valueobjects.UserID

	if req.UserID != nil {
		// User ID provided directly
		targetUserID, err = valueobjects.NewUserID(*req.UserID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}
	} else if req.Username != nil {
		// Username provided - fetch user ID from Users BC
		userIDPtr, err := c.externalUsersService.FetchUserIDByUsername(ctx.Request.Context(), *req.Username)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found with the provided username"})
			return
		}
		targetUserID = *userIDPtr
	} else {
		// Neither provided - subscribe the requesting user themselves
		targetUserID = requestedBy
	}

	// Create value objects
	communityID, err := valueobjects.NewCommunityID(req.CommunityID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid community ID"})
		return
	}

	role, err := valueobjects.NewCommunityRole(req.Role)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create command
	cmd, err := commands.NewSubscribeUserCommand(targetUserID, communityID, role, requestedBy)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Execute command
	_, err = c.commandService.Handle(ctx.Request.Context(), cmd)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "community not found" || err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "user is already subscribed to this community" {
			statusCode = http.StatusConflict
		} else if err.Error() == "only community owner or admins can add users to private communities" ||
			err.Error() == "users can only subscribe themselves to public communities" {
			statusCode = http.StatusForbidden
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the created subscription
	query, _ := queries.NewGetSubscriptionByUserAndCommunityQuery(targetUserID, communityID)
	subscription, err := c.queryService.Handle(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve created subscription"})
		return
	}

	response := resources.SubscriptionResource{
		SubscriptionID: subscription.SubscriptionID().Value(),
		UserID:         subscription.UserID().Value(),
		CommunityID:    subscription.CommunityID().Value(),
		Role:           subscription.Role().Value(),
		CreatedAt:      subscription.CreatedAt(),
		UpdatedAt:      subscription.UpdatedAt(),
	}

	ctx.JSON(http.StatusCreated, response)
}

// @Summary Unsubscribe a user from a community
// @Description Remove a user's subscription from a community. Users can unsubscribe themselves, or owner/admin can remove users.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body resources.UnsubscribeUserResource true "Unsubscribe request"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /subscriptions [delete]
func (c *SubscriptionController) UnsubscribeUser(ctx *gin.Context) {
	var req resources.UnsubscribeUserResource
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the requesting user ID from JWT context
	requestedByValue, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	requestedByID, ok := requestedByValue.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID in context"})
		return
	}

	requestedBy, err := valueobjects.NewUserID(requestedByID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid requesting user ID"})
		return
	}

	// Create value objects
	userID, err := valueobjects.NewUserID(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	communityID, err := valueobjects.NewCommunityID(req.CommunityID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid community ID"})
		return
	}

	// Create command
	cmd, err := commands.NewUnsubscribeUserCommand(userID, communityID, requestedBy)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Execute command
	err = c.commandService.HandleUnsubscribe(ctx.Request.Context(), cmd)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "community not found" || err.Error() == "subscription not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "only community owner, admins, or the user themselves can remove subscriptions" ||
			err.Error() == "community owner cannot unsubscribe from their own community" {
			statusCode = http.StatusForbidden
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary Get subscription count for a community
// @Description Get the total number of subscriptions for a specific community
// @Tags subscriptions
// @Produce json
// @Param community_id path string true "Community ID"
// @Success 200 {object} resources.SubscriptionCountResource
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /subscriptions/communities/{community_id}/count [get]
func (c *SubscriptionController) GetSubscriptionCount(ctx *gin.Context) {
	communityIDStr := ctx.Param("community_id")

	communityID, err := valueobjects.NewCommunityID(communityIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid community ID"})
		return
	}

	// Create query
	query, err := queries.NewGetSubscriptionCountByCommunityQuery(communityID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Execute query
	count, err := c.queryService.HandleCount(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := resources.SubscriptionCountResource{
		CommunityID: communityID.Value(),
		Count:       count,
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Get all subscriptions for a community
// @Description Get all subscriptions for a specific community with optional pagination
// @Tags subscriptions
// @Produce json
// @Param community_id path string true "Community ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} resources.SubscriptionListResource
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /subscriptions/communities/{community_id} [get]
func (c *SubscriptionController) GetAllSubscriptionsByCommunity(ctx *gin.Context) {
	communityIDStr := ctx.Param("community_id")

	communityID, err := valueobjects.NewCommunityID(communityIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid community ID"})
		return
	}

	// Create query
	query, err := queries.NewGetAllSubscriptionsByCommunityQuery(communityID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle pagination parameters
	// Note: For simplicity, pagination is not implemented in this version
	// You can add strconv.Atoi() to parse limit/offset query params if needed

	// Execute query
	subscriptions, err := c.queryService.HandleAll(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Transform to resources
	subscriptionResources := make([]resources.SubscriptionResource, 0, len(subscriptions))
	for _, sub := range subscriptions {
		subscriptionResources = append(subscriptionResources, resources.SubscriptionResource{
			SubscriptionID: sub.SubscriptionID().Value(),
			UserID:         sub.UserID().Value(),
			CommunityID:    sub.CommunityID().Value(),
			Role:           sub.Role().Value(),
			CreatedAt:      sub.CreatedAt(),
			UpdatedAt:      sub.UpdatedAt(),
		})
	}

	response := resources.SubscriptionListResource{
		Subscriptions: subscriptionResources,
		Total:         len(subscriptionResources),
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Get subscription by user and community
// @Description Get a specific subscription for a user in a community
// @Tags subscriptions
// @Produce json
// @Param user_id path int true "User ID"
// @Param community_id path string true "Community ID"
// @Success 200 {object} resources.SubscriptionResource
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /subscriptions/users/{user_id}/communities/{community_id} [get]
func (c *SubscriptionController) GetSubscriptionByUserAndCommunity(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	communityIDStr := ctx.Param("community_id")

	userID, err := valueobjects.NewUserID(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	communityID, err := valueobjects.NewCommunityID(communityIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid community ID"})
		return
	}

	// Create query
	query, err := queries.NewGetSubscriptionByUserAndCommunityQuery(userID, communityID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Execute query
	subscription, err := c.queryService.Handle(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if subscription == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	response := resources.SubscriptionResource{
		SubscriptionID: subscription.SubscriptionID().Value(),
		UserID:         subscription.UserID().Value(),
		CommunityID:    subscription.CommunityID().Value(),
		Role:           subscription.Role().Value(),
		CreatedAt:      subscription.CreatedAt(),
		UpdatedAt:      subscription.UpdatedAt(),
	}

	ctx.JSON(http.StatusOK, response)
}
