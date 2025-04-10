package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/bisre1921/billing-and-invoice-system/config"
	"github.com/bisre1921/billing-and-invoice-system/models"
	"github.com/gin-gonic/gin"
)

// GenerateInvoice godoc
// @Summary Generate a new invoice
// @Description Generate a new invoice for a customer with item list and auto-calculated total.
// @Tags Invoices
// @Accept json
// @Produce json
// @Param invoice body models.Invoice true "Invoice data"
// @Success 200 {object} map[string]interface{} "Invoice generated successfully"
// @Failure 400 {object} map[string]string "Invalid invoice input"
// @Failure 500 {object} map[string]string "Failed to generate invoice"
// @Router /invoice/generate [post]
func GenerateInvoice(c *gin.Context) {
	var invoice models.Invoice
	if err := c.ShouldBindJSON(&invoice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice input", "details": err.Error()})
		return
	}

	var total float64 = 0
	for i, item := range invoice.Items {
		discountAmount := item.UnitPrice * float64(item.Discount) / 100
		subtotal := float64(item.Quantity) * (item.UnitPrice - discountAmount)
		invoice.Items[i].Subtotal = subtotal
		total += subtotal
	}

	invoice.Amount = total
	invoice.CreatedAt = time.Now()
	invoice.UpdatedAt = time.Now()

	_, err := config.DB.Collection("invoices").InsertOne(context.Background(), invoice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate invoice"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice generated successfully", "invoice": invoice})
}
