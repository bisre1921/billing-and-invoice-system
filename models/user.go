package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Password  string             `json:"password" bson:"password"`
	Address   string             `json:"address,omitempty" bson:"address,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type UpdateUserInput struct {
	Name      *string   `json:"name,omitempty"`
	Email     *string   `json:"email,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	Password  *string   `json:"password,omitempty"`
	Address   *string   `json:"address,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
