package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Description  string             `json:"description" bson:"description"`
	SellingPrice float64            `json:"selling_price" bson:"selling_price"`
	BuyingPrice  float64            `json:"buying_price" bson:"buying_price"`
	Quantity     int                `json:"quantity" bson:"quantity"`
	Unit         string             `json:"unit" bson:"unit"`
	CompanyID    primitive.ObjectID `json:"company_id" bson:"company_id"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
