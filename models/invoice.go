package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CustomerID      string             `json:"customer_id" bson:"customer_id" binding:"required"`
	CompanyID       string             `json:"company_id" bson:"company_id" binding:"required"`
	ReferenceNumber string             `json:"reference_number" bson:"reference_number" binding:"required"`
	Date            time.Time          `json:"date" bson:"date"`
	Terms           string             `json:"terms" bson:"terms"`
	Status          string             `json:"status" bson:"status"`
	Amount          float64            `json:"amount" bson:"amount"`
	PaymentType     string             `json:"payment_type" bson:"payment_type" binding:"required"`
	DueDate         *time.Time         `json:"due_date,omitempty" bson:"due_date,omitempty"`
	PaymentDate     time.Time          `json:"payment_date,omitempty" bson:"payment_date,omitempty"`
	Items           []InvoiceItem      `json:"items" bson:"items"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

type InvoiceItem struct {
	ItemID    string  `json:"item_id" bson:"item_id"`
	ItemName  string  `json:"item_name" bson:"item_name"`
	Quantity  int     `json:"quantity" bson:"quantity"`
	UnitPrice float64 `json:"unit_price" bson:"unit_price"`
	Discount  float64 `json:"discount" bson:"discount"`
	Subtotal  float64 `json:"subtotal" bson:"subtotal"`
}

type UpdatePaymentStatusRequest struct {
	PaymentDate time.Time `json:"payment_date"`
}
