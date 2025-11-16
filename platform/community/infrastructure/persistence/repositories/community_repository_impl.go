package repositories

import (
	"context"
	"errors"
	"log"

	"Gommunity/platform/community/domain/model/entities"
	"Gommunity/platform/community/domain/model/valueobjects"
	domain_repos "Gommunity/platform/community/domain/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type communityRepositoryImpl struct {
	collection *mongo.Collection
}

// NewCommunityRepository creates a new CommunityRepository implementation
func NewCommunityRepository(collection *mongo.Collection) domain_repos.CommunityRepository {
	return &communityRepositoryImpl{
		collection: collection,
	}
}

// communityDocument represents the MongoDB document structure
type communityDocument struct {
	ID          string  `bson:"_id"`
	CommunityID string  `bson:"community_id"`
	OwnerID     string  `bson:"owner_id"`
	Name        string  `bson:"name"`
	Description string  `bson:"description"`
	IconURL     *string `bson:"icon_url"`
	BannerURL   *string `bson:"banner_url"`
	IsPrivate   bool    `bson:"is_private"`
	CreatedAt   int64   `bson:"created_at"`
	UpdatedAt   int64   `bson:"updated_at"`
}

// Save saves a new community to the database
func (r *communityRepositoryImpl) Save(ctx context.Context, community *entities.Community) error {
	doc := r.entityToDocument(community)

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("community already exists")
		}
		log.Printf("Error saving community to MongoDB: %v", err)
		return err
	}

	log.Printf("Community saved to MongoDB: %s", community.Name().Value())
	return nil
}

// Update updates an existing community in the database
func (r *communityRepositoryImpl) Update(ctx context.Context, community *entities.Community) error {
	filter := bson.M{"community_id": community.CommunityID().Value()}

	update := bson.M{
		"$set": bson.M{
			"name":        community.Name().Value(),
			"description": community.Description().Value(),
			"icon_url":    community.IconURL(),
			"banner_url":  community.BannerURL(),
			"is_private":  community.IsPrivate(),
			"updated_at":  community.UpdatedAt().Unix(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating community in MongoDB: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("community not found")
	}

	log.Printf("Community updated in MongoDB: %s", community.CommunityID().Value())
	return nil
}

// FindByID finds a community by community ID
func (r *communityRepositoryImpl) FindByID(ctx context.Context, communityID valueobjects.CommunityID) (*entities.Community, error) {
	filter := bson.M{"community_id": communityID.Value()}

	var doc communityDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Printf("Error finding community by ID in MongoDB: %v", err)
		return nil, err
	}

	return r.documentToEntity(&doc)
}

// FindByOwnerID finds all communities owned by a specific owner
func (r *communityRepositoryImpl) FindByOwnerID(ctx context.Context, ownerID valueobjects.OwnerID) ([]*entities.Community, error) {
	filter := bson.M{"owner_id": ownerID.Value()}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error finding communities by owner ID in MongoDB: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var communities []*entities.Community
	for cursor.Next(ctx) {
		var doc communityDocument
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("Error decoding community document: %v", err)
			return nil, err
		}

		community, err := r.documentToEntity(&doc)
		if err != nil {
			log.Printf("Error converting document to entity: %v", err)
			return nil, err
		}

		communities = append(communities, community)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return communities, nil
}

// FindAll finds all communities
func (r *communityRepositoryImpl) FindAll(ctx context.Context) ([]*entities.Community, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding all communities in MongoDB: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var communities []*entities.Community
	for cursor.Next(ctx) {
		var doc communityDocument
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("Error decoding community document: %v", err)
			return nil, err
		}

		community, err := r.documentToEntity(&doc)
		if err != nil {
			log.Printf("Error converting document to entity: %v", err)
			return nil, err
		}

		communities = append(communities, community)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return communities, nil
}

// Delete deletes a community by community ID
func (r *communityRepositoryImpl) Delete(ctx context.Context, communityID valueobjects.CommunityID) error {
	filter := bson.M{"community_id": communityID.Value()}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error deleting community from MongoDB: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("community not found")
	}

	log.Printf("Community deleted from MongoDB: %s", communityID.Value())
	return nil
}

// ExistsByID checks if a community exists by community ID
func (r *communityRepositoryImpl) ExistsByID(ctx context.Context, communityID valueobjects.CommunityID) (bool, error) {
	filter := bson.M{"community_id": communityID.Value()}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Error checking community existence in MongoDB: %v", err)
		return false, err
	}

	return count > 0, nil
}

// Helper methods for conversion between entity and document

func (r *communityRepositoryImpl) entityToDocument(community *entities.Community) *communityDocument {
	return &communityDocument{
		ID:          community.ID(),
		CommunityID: community.CommunityID().Value(),
		OwnerID:     community.OwnerID().Value(),
		Name:        community.Name().Value(),
		Description: community.Description().Value(),
		IconURL:     community.IconURL(),
		BannerURL:   community.BannerURL(),
		IsPrivate:   community.IsPrivate(),
		CreatedAt:   community.CreatedAt().Unix(),
		UpdatedAt:   community.UpdatedAt().Unix(),
	}
}

func (r *communityRepositoryImpl) documentToEntity(doc *communityDocument) (*entities.Community, error) {
	ownerID, err := valueobjects.NewOwnerID(doc.OwnerID)
	if err != nil {
		return nil, err
	}

	name, err := valueobjects.NewCommunityName(doc.Name)
	if err != nil {
		return nil, err
	}

	description, err := valueobjects.NewDescription(doc.Description)
	if err != nil {
		return nil, err
	}

	community, err := entities.NewCommunity(ownerID, name, description, doc.IconURL, doc.BannerURL, doc.IsPrivate)
	if err != nil {
		return nil, err
	}

	return community, nil
}
