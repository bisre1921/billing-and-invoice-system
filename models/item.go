package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Code         string             `json:"code" bson:"code"`
	Name         string             `json:"name" bson:"name"`
	Description  string             `json:"description" bson:"description"`
	Category     string             `json:"category" bson:"category"`
	SellingPrice float64            `json:"selling_price" bson:"selling_price"`
	Unit         string             `json:"unit" bson:"unit"`
	CompanyID    primitive.ObjectID `json:"company_id" bson:"company_id"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
