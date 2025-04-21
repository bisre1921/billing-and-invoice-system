package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Report struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CompanyID        string             `json:"company_id" bson:"company_id"`
	Title            string             `json:"title"`
	Description      string             `json:"description"`
	CreatedBy        primitive.ObjectID `json:"created_by" bson:"created_by"`
	CreatedDate      time.Time          `json:"created_date"`
	LastModifiedDate time.Time          `json:"last_modified_date"`
	Type             string             `json:"type"`
	Status           string             `json:"status"`
	Content          string             `json:"content"`
	FilePath         string             `json:"file_path,omitempty"`
}
