package repositories

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Gommunity/platform/subscriptions/domain/model/entities"
	"Gommunity/platform/subscriptions/domain/model/valueobjects"
	domain_repos "Gommunity/platform/subscriptions/domain/repositories"
)

type subscriptionRepositoryImpl struct {
	collection *mongo.Collection
}

// NewSubscriptionRepository creates a new SubscriptionRepository implementation
func NewSubscriptionRepository(collection *mongo.Collection) domain_repos.SubscriptionRepository {
	return &subscriptionRepositoryImpl{
		collection: collection,
	}
}

// subscriptionDocument represents the MongoDB document structure
type subscriptionDocument struct {
	ID             string `bson:"_id"`
	SubscriptionID string `bson:"subscription_id"`
	UserID         string `bson:"user_id"`
	CommunityID    string `bson:"community_id"`
	Role           string `bson:"role"`
	CreatedAt      int64  `bson:"created_at"`
	UpdatedAt      int64  `bson:"updated_at"`
}

// Save persists a subscription
func (r *subscriptionRepositoryImpl) Save(ctx context.Context, subscription *entities.Subscription) error {
	doc := r.toDocument(subscription)

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("subscription already exists")
		}
		return err
	}

	return nil
}

// FindByID retrieves a subscription by its ID
func (r *subscriptionRepositoryImpl) FindByID(ctx context.Context, id valueobjects.SubscriptionID) (*entities.Subscription, error) {
	filter := bson.M{"subscription_id": id.Value()}

	var doc subscriptionDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return r.toEntity(&doc)
}

// FindByUserAndCommunity retrieves a subscription by user ID and community ID
func (r *subscriptionRepositoryImpl) FindByUserAndCommunity(ctx context.Context, userID valueobjects.UserID, communityID valueobjects.CommunityID) (*entities.Subscription, error) {
	filter := bson.M{
		"user_id":      userID.Value(),
		"community_id": communityID.Value(),
	}

	var doc subscriptionDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return r.toEntity(&doc)
}

// FindAllByCommunityID retrieves all subscriptions for a specific community
func (r *subscriptionRepositoryImpl) FindAllByCommunityID(ctx context.Context, communityID valueobjects.CommunityID, limit, offset *int) ([]*entities.Subscription, error) {
	filter := bson.M{"community_id": communityID.Value()}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	if offset != nil && *offset > 0 {
		opts.SetSkip(int64(*offset))
	}

	if limit != nil && *limit > 0 {
		opts.SetLimit(int64(*limit))
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var subscriptions []*entities.Subscription
	for cursor.Next(ctx) {
		var doc subscriptionDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		subscription, err := r.toEntity(&doc)
		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, subscription)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// FindAllByUserID retrieves all subscriptions for a specific user
func (r *subscriptionRepositoryImpl) FindAllByUserID(ctx context.Context, userID valueobjects.UserID) ([]*entities.Subscription, error) {
	filter := bson.M{"user_id": userID.Value()}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var subscriptions []*entities.Subscription
	for cursor.Next(ctx) {
		var doc subscriptionDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		subscription, err := r.toEntity(&doc)
		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, subscription)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// CountByCommunityID returns the total number of subscriptions for a community
func (r *subscriptionRepositoryImpl) CountByCommunityID(ctx context.Context, communityID valueobjects.CommunityID) (int64, error) {
	filter := bson.M{"community_id": communityID.Value()}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// ExistsByUserAndCommunity checks if a subscription exists for a user in a community
func (r *subscriptionRepositoryImpl) ExistsByUserAndCommunity(ctx context.Context, userID valueobjects.UserID, communityID valueobjects.CommunityID) (bool, error) {
	filter := bson.M{
		"user_id":      userID.Value(),
		"community_id": communityID.Value(),
	}

	count, err := r.collection.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Delete removes a subscription
func (r *subscriptionRepositoryImpl) Delete(ctx context.Context, id valueobjects.SubscriptionID) error {
	filter := bson.M{"subscription_id": id.Value()}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("subscription not found")
	}

	return nil
}

// DeleteByUserAndCommunity removes a subscription by user and community
func (r *subscriptionRepositoryImpl) DeleteByUserAndCommunity(ctx context.Context, userID valueobjects.UserID, communityID valueobjects.CommunityID) error {
	filter := bson.M{
		"user_id":      userID.Value(),
		"community_id": communityID.Value(),
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("subscription not found")
	}

	return nil
}

// DeleteByCommunity removes all subscriptions for a community
func (r *subscriptionRepositoryImpl) DeleteByCommunity(ctx context.Context, communityID valueobjects.CommunityID) error {
	filter := bson.M{
		"community_id": communityID.Value(),
	}

	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

// toDocument converts an entity to a document
func (r *subscriptionRepositoryImpl) toDocument(subscription *entities.Subscription) *subscriptionDocument {
	return &subscriptionDocument{
		ID:             subscription.ID(),
		SubscriptionID: subscription.SubscriptionID().Value(),
		UserID:         subscription.UserID().Value(),
		CommunityID:    subscription.CommunityID().Value(),
		Role:           subscription.Role().Value(),
		CreatedAt:      subscription.CreatedAt().Unix(),
		UpdatedAt:      subscription.UpdatedAt().Unix(),
	}
}

// toEntity converts a document to an entity
func (r *subscriptionRepositoryImpl) toEntity(doc *subscriptionDocument) (*entities.Subscription, error) {
	subscriptionID, err := valueobjects.NewSubscriptionID(doc.SubscriptionID)
	if err != nil {
		return nil, err
	}

	userID, err := valueobjects.NewUserID(doc.UserID)
	if err != nil {
		return nil, err
	}

	communityID, err := valueobjects.NewCommunityID(doc.CommunityID)
	if err != nil {
		return nil, err
	}

	role, err := valueobjects.NewCommunityRole(doc.Role)
	if err != nil {
		return nil, err
	}

	return entities.ReconstructSubscription(
		doc.ID,
		subscriptionID,
		userID,
		communityID,
		role,
		time.Unix(doc.CreatedAt, 0),
		time.Unix(doc.UpdatedAt, 0),
	), nil
}
