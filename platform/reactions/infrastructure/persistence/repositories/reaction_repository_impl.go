package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"Gommunity/platform/reactions/domain/model/entities"
	"Gommunity/platform/reactions/domain/model/valueobjects"
	domain_repositories "Gommunity/platform/reactions/domain/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type reactionRepositoryImpl struct {
	collection *mongo.Collection
}

// NewReactionRepository creates a MongoDB-backed ReactionRepository.
func NewReactionRepository(collection *mongo.Collection) domain_repositories.ReactionRepository {
	return &reactionRepositoryImpl{
		collection: collection,
	}
}

type reactionDocument struct {
	ID           string `bson:"_id"`
	ReactionID   string `bson:"reaction_id"`
	PostID       string `bson:"post_id"`
	UserID       string `bson:"user_id"`
	ReactionType string `bson:"reaction_type"`
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
}

// Save inserts a new reaction document.
func (r *reactionRepositoryImpl) Save(ctx context.Context, reaction *entities.Reaction) error {
	doc := r.entityToDocument(reaction)
	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("user already reacted to this post")
		}
		log.Printf("failed to insert reaction: %v", err)
		return err
	}
	return nil
}

// Update modifies an existing reaction document.
func (r *reactionRepositoryImpl) Update(ctx context.Context, reaction *entities.Reaction) error {
	filter := bson.M{"reaction_id": reaction.ReactionID().Value()}
	update := bson.M{
		"$set": bson.M{
			"reaction_type": reaction.ReactionType().Value(),
			"updated_at":    reaction.UpdatedAt().Unix(),
		},
	}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("failed to update reaction: %v", err)
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("reaction not found")
	}
	return nil
}

// FindByID retrieves a reaction by its identifier.
func (r *reactionRepositoryImpl) FindByID(ctx context.Context, reactionID valueobjects.ReactionID) (*entities.Reaction, error) {
	filter := bson.M{"reaction_id": reactionID.Value()}

	var doc reactionDocument
	if err := r.collection.FindOne(ctx, filter).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		log.Printf("failed to find reaction by id: %v", err)
		return nil, err
	}
	return r.documentToEntity(&doc)
}

// FindByPostAndUser retrieves a user's reaction to a specific post.
func (r *reactionRepositoryImpl) FindByPostAndUser(ctx context.Context, postID valueobjects.PostID, userID valueobjects.UserID) (*entities.Reaction, error) {
	filter := bson.M{
		"post_id": postID.Value(),
		"user_id": userID.Value(),
	}

	var doc reactionDocument
	if err := r.collection.FindOne(ctx, filter).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		log.Printf("failed to find reaction by post and user: %v", err)
		return nil, err
	}
	return r.documentToEntity(&doc)
}

// FindByPost retrieves all reactions for a specific post.
func (r *reactionRepositoryImpl) FindByPost(ctx context.Context, postID valueobjects.PostID) ([]*entities.Reaction, error) {
	filter := bson.M{"post_id": postID.Value()}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Printf("failed to find reactions by post: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var reactions []*entities.Reaction
	for cursor.Next(ctx) {
		var doc reactionDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		entity, err := r.documentToEntity(&doc)
		if err != nil {
			return nil, err
		}
		reactions = append(reactions, entity)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return reactions, nil
}

// CountByPost returns reaction counts grouped by type for a post.
func (r *reactionRepositoryImpl) CountByPost(ctx context.Context, postID valueobjects.PostID) (map[string]int, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "post_id", Value: postID.Value()}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$reaction_type"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("failed to count reactions by post: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	counts := make(map[string]int)
	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		counts[result.ID] = result.Count
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return counts, nil
}

// Delete removes a reaction by identifier.
func (r *reactionRepositoryImpl) Delete(ctx context.Context, reactionID valueobjects.ReactionID) error {
	filter := bson.M{"reaction_id": reactionID.Value()}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("failed to delete reaction: %v", err)
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("reaction not found")
	}
	return nil
}

// DeleteByPostAndUser removes a user's reaction from a specific post.
func (r *reactionRepositoryImpl) DeleteByPostAndUser(ctx context.Context, postID valueobjects.PostID, userID valueobjects.UserID) error {
	filter := bson.M{
		"post_id": postID.Value(),
		"user_id": userID.Value(),
	}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("failed to delete reaction by post and user: %v", err)
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("reaction not found")
	}
	return nil
}

// DeleteByPostIDs removes reactions for a list of posts
func (r *reactionRepositoryImpl) DeleteByPostIDs(ctx context.Context, postIDs []valueobjects.PostID) error {
	if len(postIDs) == 0 {
		return nil
	}

	values := make([]string, 0, len(postIDs))
	for _, id := range postIDs {
		values = append(values, id.Value())
	}

	filter := bson.M{
		"post_id": bson.M{"$in": values},
	}

	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Printf("failed to delete reactions by posts: %v", err)
		return err
	}

	return nil
}

func (r *reactionRepositoryImpl) entityToDocument(reaction *entities.Reaction) *reactionDocument {
	return &reactionDocument{
		ID:           reaction.ReactionID().Value(),
		ReactionID:   reaction.ReactionID().Value(),
		PostID:       reaction.PostID().Value(),
		UserID:       reaction.UserID().Value(),
		ReactionType: reaction.ReactionType().Value(),
		CreatedAt:    reaction.CreatedAt().Unix(),
		UpdatedAt:    reaction.UpdatedAt().Unix(),
	}
}

func (r *reactionRepositoryImpl) documentToEntity(doc *reactionDocument) (*entities.Reaction, error) {
	reactionID, err := valueobjects.NewReactionID(doc.ReactionID)
	if err != nil {
		return nil, err
	}
	postID, err := valueobjects.NewPostID(doc.PostID)
	if err != nil {
		return nil, err
	}
	userID, err := valueobjects.NewUserID(doc.UserID)
	if err != nil {
		return nil, err
	}
	reactionType, err := valueobjects.NewReactionType(doc.ReactionType)
	if err != nil {
		return nil, err
	}

	reaction := entities.ReconstructReaction(
		doc.ID,
		reactionID,
		postID,
		userID,
		reactionType,
		time.Unix(doc.CreatedAt, 0),
		time.Unix(doc.UpdatedAt, 0),
	)

	return reaction, nil
}
