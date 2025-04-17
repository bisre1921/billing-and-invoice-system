// backend/models/activity.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Activity struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CompanyID  string             `json:"company_id"`
	Timestamp  time.Time          `json:"timestamp"`
	Type       string             `json:"type"`
	Message    string             `json:"message"`
	InvoiceID  string             `json:"invoice_id,omitempty"`
	CustomerID string             `json:"customer_id,omitempty"`
	EmployeeID string             `json:"employee_id,omitempty"`
	UserID     string             `json:"user_id,omitempty"`
}
