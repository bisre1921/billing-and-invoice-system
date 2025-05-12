package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"phone" bson:"phone"`
	Address   string             `json:"address" bson:"address"`
	Position  string             `json:"position" bson:"position"`
	CompanyID primitive.ObjectID `json:"company_id" bson:"company_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type EmployeeInvitationRequest struct {
	Email     string             `json:"email" bson:"email"`
	CompanyID primitive.ObjectID `json:"company_id" bson:"company_id"`
	Position  string             `json:"position" bson:"position"`
	Name	  string             `json:"name" bson:"name"`
	Address   string             `json:"address" bson:"address"`
	Phone     string             `json:"phone" bson:"phone"`
}

type Invitation struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Token	 string             `json:"token" bson:"token"`
	Email	 string             `json:"email" bson:"email"`
	CompanyID primitive.ObjectID `json:"company_id" bson:"company_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}