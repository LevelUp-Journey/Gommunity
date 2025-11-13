package repositories

import (
	"context"
	"errors"
	"log"

	"Gommunity/internal/community/users/domain/model/entities"
	"Gommunity/internal/community/users/domain/model/valueobjects"
	domain_repos "Gommunity/internal/community/users/domain/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type roleRepositoryImpl struct {
	collection *mongo.Collection
}

// NewRoleRepository creates a new RoleRepository implementation
func NewRoleRepository(collection *mongo.Collection) domain_repos.RoleRepository {
	return &roleRepositoryImpl{
		collection: collection,
	}
}

// roleDocument represents the MongoDB document structure for Role
type roleDocument struct {
	ID        string `bson:"_id"`
	RoleID    string `bson:"role_id"`
	Name      string `bson:"name"`
	CreatedAt int64  `bson:"created_at"`
}

// Save saves a new role to the database
func (r *roleRepositoryImpl) Save(ctx context.Context, role *entities.Role) error {
	doc := r.entityToDocument(role)

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("role already exists")
		}
		log.Printf("Error saving role to MongoDB: %v", err)
		return err
	}

	log.Printf("Role saved to MongoDB: %s", role.Name())
	return nil
}

// FindByID finds a role by role ID
func (r *roleRepositoryImpl) FindByID(ctx context.Context, roleID valueobjects.RoleID) (*entities.Role, error) {
	filter := bson.M{"role_id": roleID.Value()}

	var doc roleDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Printf("Error finding role by RoleID in MongoDB: %v", err)
		return nil, err
	}

	return r.documentToEntity(&doc)
}

// FindByName finds a role by name
func (r *roleRepositoryImpl) FindByName(ctx context.Context, name string) (*entities.Role, error) {
	filter := bson.M{"name": name}

	var doc roleDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Printf("Error finding role by name in MongoDB: %v", err)
		return nil, err
	}

	return r.documentToEntity(&doc)
}

// FindAll finds all roles
func (r *roleRepositoryImpl) FindAll(ctx context.Context) ([]*entities.Role, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding all roles in MongoDB: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var roles []*entities.Role
	for cursor.Next(ctx) {
		var doc roleDocument
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("Error decoding role document: %v", err)
			return nil, err
		}

		role, err := r.documentToEntity(&doc)
		if err != nil {
			log.Printf("Error converting document to entity: %v", err)
			return nil, err
		}

		roles = append(roles, role)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return roles, nil
}

// Helper methods for conversion between entity and document

func (r *roleRepositoryImpl) entityToDocument(role *entities.Role) *roleDocument {
	return &roleDocument{
		ID:        role.ID(),
		RoleID:    role.RoleID().Value(),
		Name:      role.Name(),
		CreatedAt: role.CreatedAt().Unix(),
	}
}

func (r *roleRepositoryImpl) documentToEntity(doc *roleDocument) (*entities.Role, error) {
	roleID, err := valueobjects.NewRoleID(doc.RoleID)
	if err != nil {
		return nil, err
	}

	role, err := entities.NewRole(roleID, doc.Name)
	if err != nil {
		return nil, err
	}

	return role, nil
}
