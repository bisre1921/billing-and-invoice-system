package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID                     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name                   string             `json:"name" bson:"name"`
	Email                  string             `json:"email" bson:"email"`
	Phone                  string             `json:"phone" bson:"phone"`
	Address                string             `json:"address" bson:"address"`
	TIN                    string             `json:"tin" bson:"tin"`
	MaxCreditAmount        float64            `json:"max_credit_amount" bson:"max_credit_amount"`
	CurrentCreditAvailable float64            `json:"current_credit_available" bson:"current_credit_available"`
	CompanyID              primitive.ObjectID `json:"company_id" bson:"company_id"`
	CreatedAt              time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt              time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
