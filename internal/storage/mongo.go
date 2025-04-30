package storage

import (
	"context"
	"github.com/erkkipm/sso_auth/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	Collection *mongo.Collection
}

func NewStorage(uri, db, coll string) (*Storage, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &Storage{Collection: client.Database(db).Collection(coll)}, nil
}

func (s *Storage) CreateUser(ctx context.Context, u models.User) error {
	_, err := s.Collection.InsertOne(ctx, u)
	return err
}

func (s *Storage) FindUser(ctx context.Context, appID, email string) (models.User, error) {
	var user models.User
	err := s.Collection.FindOne(ctx, bson.M{"app_id": appID, "email": email}).Decode(&user)
	return user, err
}

func (s *Storage) UpdatePassword(ctx context.Context, appID, email, newHash string) error {
	_, err := s.Collection.UpdateOne(ctx,
		bson.M{"app_id": appID, "email": email},
		bson.M{"$set": bson.M{"password": newHash}},
	)
	return err
}
