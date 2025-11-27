package entities

import (
	"errors"
	"time"

	"Gommunity/platform/subscriptions/domain/model/valueobjects"
)

// Subscription represents a user's subscription to a community with an associated role
// This is the aggregate root for the Subscriptions bounded context
type Subscription struct {
	id             string                      `bson:"_id"`
	subscriptionID valueobjects.SubscriptionID `bson:"subscription_id"`
	userID         valueobjects.UserID         `bson:"user_id"`
	communityID    valueobjects.CommunityID    `bson:"community_id"`
	role           valueobjects.CommunityRole  `bson:"role"`
	createdAt      time.Time                   `bson:"created_at"`
	updatedAt      time.Time                   `bson:"updated_at"`
}

// NewSubscription creates a new Subscription aggregate
func NewSubscription(
	userID valueobjects.UserID,
	communityID valueobjects.CommunityID,
	role valueobjects.CommunityRole,
) (*Subscription, error) {
	if userID.IsZero() {
		return nil, errors.New("user ID cannot be zero")
	}
	if communityID.IsZero() {
		return nil, errors.New("community ID cannot be empty")
	}
	if role.IsZero() {
		return nil, errors.New("role cannot be empty")
	}

	now := time.Now()
	subscriptionID := valueobjects.GenerateSubscriptionID()

	return &Subscription{
		id:             subscriptionID.Value(),
		subscriptionID: subscriptionID,
		userID:         userID,
		communityID:    communityID,
		role:           role,
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

// ReconstructSubscription reconstructs a Subscription from persistence
func ReconstructSubscription(
	id string,
	subscriptionID valueobjects.SubscriptionID,
	userID valueobjects.UserID,
	communityID valueobjects.CommunityID,
	role valueobjects.CommunityRole,
	createdAt time.Time,
	updatedAt time.Time,
) *Subscription {
	return &Subscription{
		id:             id,
		subscriptionID: subscriptionID,
		userID:         userID,
		communityID:    communityID,
		role:           role,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

// ID returns the MongoDB document ID
func (s *Subscription) ID() string {
	return s.id
}

// SubscriptionID returns the subscription ID
func (s *Subscription) SubscriptionID() valueobjects.SubscriptionID {
	return s.subscriptionID
}

// UserID returns the user ID
func (s *Subscription) UserID() valueobjects.UserID {
	return s.userID
}

// CommunityID returns the community ID
func (s *Subscription) CommunityID() valueobjects.CommunityID {
	return s.communityID
}

// Role returns the community-specific role
func (s *Subscription) Role() valueobjects.CommunityRole {
	return s.role
}

// CreatedAt returns the creation timestamp
func (s *Subscription) CreatedAt() time.Time {
	return s.createdAt
}

// UpdatedAt returns the last update timestamp
func (s *Subscription) UpdatedAt() time.Time {
	return s.updatedAt
}

// UpdateRole updates the role for this subscription
func (s *Subscription) UpdateRole(newRole valueobjects.CommunityRole) error {
	if newRole.IsZero() {
		return errors.New("new role cannot be empty")
	}

	s.role = newRole
	s.updatedAt = time.Now()
	return nil
}

// HasAdminOrOwnerRole checks if the subscription has admin or owner privileges
func (s *Subscription) HasAdminOrOwnerRole() bool {
	return s.role.IsAdminOrOwner()
}
