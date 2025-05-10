package models

import "time"

type User struct {
	ID        string    `bson:"_id,omitempty"`
	AppID     string    `bson:"app_id"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	Phone     string    `bson:"phone"`
	CreatedAt time.Time `bson:"created_at"`
}
