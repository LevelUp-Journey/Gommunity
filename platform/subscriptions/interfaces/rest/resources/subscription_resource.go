package resources

import "time"

// SubscriptionResource represents a subscription in the REST API
type SubscriptionResource struct {
	SubscriptionID string    `json:"subscription_id" example:"507f1f77bcf86cd799439011"`
	UserID         string    `json:"user_id" example:"507f1f77bcf86cd799439013"`
	CommunityID    string    `json:"community_id" example:"507f1f77bcf86cd799439012"`
	Role           string    `json:"role" example:"member"`
	CreatedAt      time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt      time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// SubscribeUserResource represents the request to subscribe a user to a community
type SubscribeUserResource struct {
	UserID      *string `json:"user_id,omitempty" example:"507f1f77bcf86cd799439013"`
	Username    *string `json:"username,omitempty" example:"john_doe" validate:"omitempty,min=3,max=50"`
	CommunityID string  `json:"community_id" example:"507f1f77bcf86cd799439012" validate:"required"`
	Role        string  `json:"role" example:"member" validate:"required,oneof=member admin owner"`
}

// UnsubscribeUserResource represents the request to unsubscribe a user from a community
type UnsubscribeUserResource struct {
	UserID      string `json:"user_id" example:"507f1f77bcf86cd799439013" validate:"required"`
	CommunityID string `json:"community_id" example:"507f1f77bcf86cd799439012" validate:"required"`
}

// SubscriptionCountResource represents the subscription count for a community
type SubscriptionCountResource struct {
	CommunityID string `json:"community_id" example:"507f1f77bcf86cd799439012"`
	Count       int64  `json:"count" example:"150"`
}

// SubscriptionListResource represents a list of subscriptions
type SubscriptionListResource struct {
	Subscriptions []SubscriptionResource `json:"subscriptions"`
	Total         int                    `json:"total" example:"150"`
}
