package models

import (
	"time"
)

type User struct {
	ID        string     `json:"id" bson:"_id,omitempty"`
	Name      string     `json:"name" bson:"name"`
	Email     string     `json:"email" bson:"email"`
	Phone     string     `json:"phone,omitempty" bson:"phone,omitempty"`
	Password  string     `json:"-" bson:"password"`
	Address   string     `json:"address,omitempty" bson:"address,omitempty"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}
