package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	AppID     string             `bson:"app_id"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	PassHash  []byte             `bson:"pass_hash"`
	Phone     string             `bson:"phone"`
	CreatedAt time.Time          `bson:"created_at"`
}

type App struct {
	ID     int
	Name   string
	Secret string
}
