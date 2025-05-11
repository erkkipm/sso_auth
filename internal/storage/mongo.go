package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/erkkipm/sso_auth/internal/configs"
	"github.com/erkkipm/sso_auth/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	Collection *mongo.Collection
}

func NewStorage(ctx context.Context, cfg configs.MongoDB) (*Storage, error) {

	mongoURI := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.AuthDB,
	)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}
	return &Storage{Collection: client.Database(cfg.Database).Collection(cfg.Collection.Users)}, nil
}

// CreateUser ... Запись пользователя в Базу
func (s *Storage) CreateUser(ctx context.Context, u models.User) error {
	_, err := s.Collection.InsertOne(ctx, u)
	return err
}
func (s *Storage) GetUserByEmailAndApp(ctx context.Context, u models.User) (*models.User, error) {
	filter := bson.M{
		"email":  u.Email,
		"app_id": u.AppID,
	}
	var user models.User
	err := s.Collection.FindOne(ctx, filter).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil // пользователя нет — это не ошибка
	}
	if err != nil {
		return nil, err // реальная ошибка при запросе
	}
	return &user, nil
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
