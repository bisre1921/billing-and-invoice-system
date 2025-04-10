package models

import "time"

type Invoice struct {
	ID              string        `json:"id" bson:"_id,omitempty"`
	CustomerID      string        `json:"customer_id" bson:"customer_id" binding:"required"`
	ReferenceNumber string        `json:"reference_number" bson:"reference_number" binding:"required"`
	Date            time.Time     `json:"date" bson:"date" binding:"required"`
	Terms           string        `json:"terms" bson:"terms"`
	Status          string        `json:"status" bson:"status"`
	Amount          float64       `json:"amount" bson:"amount"`
	DueDate         time.Time     `json:"due_date" bson:"due_date"`
	PaymentDate     time.Time     `json:"payment_date,omitempty" bson:"payment_date,omitempty"`
	Items           []InvoiceItem `json:"items" bson:"items"`
	CreatedAt       time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" bson:"updated_at"`
}

type InvoiceItem struct {
	ItemID    string  `json:"item_id" bson:"item_id"`
	ItemName  string  `json:"item_name" bson:"item_name"`
	Quantity  int     `json:"quantity" bson:"quantity"`
	UnitPrice float64 `json:"unit_price" bson:"unit_price"`
	Discount  float64 `json:"discount" bson:"discount"`
	Subtotal  float64 `json:"subtotal" bson:"subtotal"`
}
