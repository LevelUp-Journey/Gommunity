package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"Gommunity/platform/posts/domain/model/entities"
	"Gommunity/platform/posts/domain/model/valueobjects"
	domain_repositories "Gommunity/platform/posts/domain/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type postRepositoryImpl struct {
	collection *mongo.Collection
}

// NewPostRepository creates a MongoDB-backed PostRepository.
func NewPostRepository(collection *mongo.Collection) domain_repositories.PostRepository {
	return &postRepositoryImpl{
		collection: collection,
	}
}

type postDocument struct {
	ID          string   `bson:"_id"`
	PostID      string   `bson:"post_id"`
	CommunityID string   `bson:"community_id"`
	AuthorID    string   `bson:"author_id"`
	PostType    string   `bson:"post_type"`
	Content     string   `bson:"content"`
	Images      []string `bson:"images"`
	CreatedAt   int64    `bson:"created_at"`
	UpdatedAt   int64    `bson:"updated_at"`
}

// Save inserts a new post document.
func (r *postRepositoryImpl) Save(ctx context.Context, post *entities.Post) error {
	doc := r.entityToDocument(post)
	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("post already exists")
		}
		log.Printf("failed to insert post: %v", err)
		return err
	}
	return nil
}

// FindByID retrieves a post by its identifier.
func (r *postRepositoryImpl) FindByID(ctx context.Context, postID valueobjects.PostID) (*entities.Post, error) {
	filter := bson.M{"post_id": postID.Value()}

	var doc postDocument
	if err := r.collection.FindOne(ctx, filter).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		log.Printf("failed to find post by id: %v", err)
		return nil, err
	}
	return r.documentToEntity(&doc)
}

// FindByCommunity retrieves posts belonging to a community.
func (r *postRepositoryImpl) FindByCommunity(ctx context.Context, communityID valueobjects.CommunityID, limit, offset *int) ([]*entities.Post, error) {
	filter := bson.M{"community_id": communityID.Value()}
	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	if limit != nil {
		findOptions.SetLimit(int64(*limit))
	}
	if offset != nil {
		findOptions.SetSkip(int64(*offset))
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Printf("failed to find posts by community: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*entities.Post
	for cursor.Next(ctx) {
		var doc postDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		entity, err := r.documentToEntity(&doc)
		if err != nil {
			return nil, err
		}
		posts = append(posts, entity)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// FindByCommunities retrieves posts from multiple communities, optionally filtered by post type
func (r *postRepositoryImpl) FindByCommunities(ctx context.Context, communityIDs []valueobjects.CommunityID, postType *valueobjects.PostType, limit, offset *int) ([]*entities.Post, error) {
	// Convert community IDs to strings
	communityIDStrings := make([]string, len(communityIDs))
	for i, id := range communityIDs {
		communityIDStrings[i] = id.Value()
	}

	// Build filter
	filter := bson.M{"community_id": bson.M{"$in": communityIDStrings}}

	// Add post type filter if provided
	if postType != nil {
		filter["post_type"] = postType.Value()
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	if limit != nil {
		findOptions.SetLimit(int64(*limit))
	}
	if offset != nil {
		findOptions.SetSkip(int64(*offset))
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Printf("failed to find posts by communities: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*entities.Post
	for cursor.Next(ctx) {
		var doc postDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		entity, err := r.documentToEntity(&doc)
		if err != nil {
			return nil, err
		}
		posts = append(posts, entity)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// FindByAuthorAndCommunity retrieves a publication constraint for a user.
func (r *postRepositoryImpl) FindByAuthorAndCommunity(ctx context.Context, authorID valueobjects.AuthorID, communityID valueobjects.CommunityID) (*entities.Post, error) {
	filter := bson.M{
		"author_id":    authorID.Value(),
		"community_id": communityID.Value(),
	}

	var doc postDocument
	if err := r.collection.FindOne(ctx, filter).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return r.documentToEntity(&doc)
}

// Delete removes a post by identifier.
func (r *postRepositoryImpl) Delete(ctx context.Context, postID valueobjects.PostID) error {
	filter := bson.M{"post_id": postID.Value()}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("failed to delete post: %v", err)
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("post not found")
	}
	return nil
}

// FindPostIDsByCommunity returns post IDs for a community (lightweight)
func (r *postRepositoryImpl) FindPostIDsByCommunity(ctx context.Context, communityID valueobjects.CommunityID) ([]valueobjects.PostID, error) {
	filter := bson.M{"community_id": communityID.Value()}
	projection := bson.M{"post_id": 1}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		log.Printf("failed to list post ids by community: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var ids []valueobjects.PostID
	for cursor.Next(ctx) {
		var doc postDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		postID, err := valueobjects.NewPostID(doc.PostID)
		if err != nil {
			return nil, err
		}
		ids = append(ids, postID)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

// DeleteByCommunity removes all posts for a community
func (r *postRepositoryImpl) DeleteByCommunity(ctx context.Context, communityID valueobjects.CommunityID) error {
	filter := bson.M{"community_id": communityID.Value()}
	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Printf("failed to delete posts by community: %v", err)
		return err
	}
	return nil
}

func (r *postRepositoryImpl) entityToDocument(post *entities.Post) *postDocument {
	return &postDocument{
		ID:          post.PostID().Value(),
		PostID:      post.PostID().Value(),
		CommunityID: post.CommunityID().Value(),
		AuthorID:    post.AuthorID().Value(),
		PostType:    post.PostType().Value(),
		Content:     post.Content().Value(),
		Images:      post.Images().URLs(),
		CreatedAt:   post.CreatedAt().Unix(),
		UpdatedAt:   post.UpdatedAt().Unix(),
	}
}

func (r *postRepositoryImpl) documentToEntity(doc *postDocument) (*entities.Post, error) {
	postID, err := valueobjects.NewPostID(doc.PostID)
	if err != nil {
		return nil, err
	}
	communityID, err := valueobjects.NewCommunityID(doc.CommunityID)
	if err != nil {
		return nil, err
	}
	authorID, err := valueobjects.NewAuthorID(doc.AuthorID)
	if err != nil {
		return nil, err
	}
	postType, err := valueobjects.NewPostType(doc.PostType)
	if err != nil {
		return nil, err
	}
	content, err := valueobjects.NewPostContent(doc.Content)
	if err != nil {
		return nil, err
	}
	images, err := valueobjects.NewPostImages(doc.Images)
	if err != nil {
		return nil, err
	}

	post := entities.ReconstructPost(
		doc.ID,
		postID,
		communityID,
		authorID,
		postType,
		content,
		images,
		time.Unix(doc.CreatedAt, 0),
		time.Unix(doc.UpdatedAt, 0),
	)

	return post, nil
}
