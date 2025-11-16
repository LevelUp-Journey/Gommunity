package repositories

import (
	"context"
	"errors"
	"log"

	"Gommunity/platform/users/domain/model/entities"
	"Gommunity/platform/users/domain/model/valueobjects"
	domain_repos "Gommunity/platform/users/domain/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepositoryImpl struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new UserRepository implementation
func NewUserRepository(collection *mongo.Collection) domain_repos.UserRepository {
	return &userRepositoryImpl{
		collection: collection,
	}
}

// userDocument represents the MongoDB document structure
type userDocument struct {
	ID         string  `bson:"_id"`
	UserID     string  `bson:"user_id"`
	ProfileID  string  `bson:"profile_id"`
	Username   string  `bson:"username"`
	ProfileURL *string `bson:"profile_url"`
	BannerURL  *string `bson:"banner_url"`
	UpdatedAt  int64   `bson:"updated_at"`
	CreatedAt  int64   `bson:"created_at"`
}

// Save saves a new user to the database
func (r *userRepositoryImpl) Save(ctx context.Context, user *entities.User) error {
	doc := r.entityToDocument(user)

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("user already exists")
		}
		log.Printf("Error saving user to MongoDB: %v", err)
		return err
	}

	log.Printf("User saved to MongoDB: %s", user.UserID().Value())
	return nil
}

// Update updates an existing user in the database
func (r *userRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	filter := bson.M{"user_id": user.UserID().Value()}

	update := bson.M{
		"$set": bson.M{
			"username":    user.Username().Value(),
			"profile_url": user.ProfileURL(),
			"banner_url":  user.BannerURL(),
			"updated_at":  user.UpdatedAt().Unix(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating user in MongoDB: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	log.Printf("User updated in MongoDB: %s", user.UserID().Value())
	return nil
}

// FindByUserID finds a user by user ID
func (r *userRepositoryImpl) FindByUserID(ctx context.Context, userID valueobjects.UserID) (*entities.User, error) {
	filter := bson.M{"user_id": userID.Value()}

	var doc userDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Printf("Error finding user by UserID in MongoDB: %v", err)
		return nil, err
	}

	return r.documentToEntity(&doc)
}

// FindByProfileID finds a user by profile ID
func (r *userRepositoryImpl) FindByProfileID(ctx context.Context, profileID valueobjects.ProfileID) (*entities.User, error) {
	filter := bson.M{"profile_id": profileID.Value()}

	var doc userDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Printf("Error finding user by ProfileID in MongoDB: %v", err)
		return nil, err
	}

	return r.documentToEntity(&doc)
}

// ExistsByUserID checks if a user exists by user ID
func (r *userRepositoryImpl) ExistsByUserID(ctx context.Context, userID valueobjects.UserID) (bool, error) {
	filter := bson.M{"user_id": userID.Value()}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Error checking user existence in MongoDB: %v", err)
		return false, err
	}

	return count > 0, nil
}

// FindByUsername finds a user by username
func (r *userRepositoryImpl) FindByUsername(ctx context.Context, username valueobjects.Username) (*entities.User, error) {
	filter := bson.M{"username": username.Value()}

	var doc userDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Printf("Error finding user by username in MongoDB: %v", err)
		return nil, err
	}

	return r.documentToEntity(&doc)
}

// Delete deletes a user by user ID
func (r *userRepositoryImpl) Delete(ctx context.Context, userID valueobjects.UserID) error {
	filter := bson.M{"user_id": userID.Value()}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error deleting user from MongoDB: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}

	log.Printf("User deleted from MongoDB: %s", userID.Value())
	return nil
}

// Helper methods for conversion between entity and document

func (r *userRepositoryImpl) entityToDocument(user *entities.User) *userDocument {
	return &userDocument{
		ID:         user.ID(),
		UserID:     user.UserID().Value(),
		ProfileID:  user.ProfileID().Value(),
		Username:   user.Username().Value(),
		ProfileURL: user.ProfileURL(),
		BannerURL:  user.BannerURL(),
		UpdatedAt:  user.UpdatedAt().Unix(),
		CreatedAt:  user.CreatedAt().Unix(),
	}
}

func (r *userRepositoryImpl) documentToEntity(doc *userDocument) (*entities.User, error) {
	userID, err := valueobjects.NewUserID(doc.UserID)
	if err != nil {
		return nil, err
	}

	profileID, err := valueobjects.NewProfileID(doc.ProfileID)
	if err != nil {
		return nil, err
	}

	username, err := valueobjects.NewUsername(doc.Username)
	if err != nil {
		return nil, err
	}

	user, err := entities.NewUser(userID, profileID, username, doc.ProfileURL)
	if err != nil {
		return nil, err
	}

	return user, nil
}
