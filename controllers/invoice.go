package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bisre1921/billing-and-invoice-system/config"
	"github.com/bisre1921/billing-and-invoice-system/models"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/gomail.v2"
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

	res, err := config.DB.Collection("invoices").InsertOne(context.Background(), invoice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate invoice"})
		return
	}

	invoice.ID = res.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusOK, gin.H{"message": "Invoice generated successfully", "invoice": invoice})

}

// GetInvoice godoc
// @Summary Get invoice by ID
// @Description Retrieve a specific invoice by its unique identifier
// @Tags Invoices
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} map[string]interface{} "Invoice retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid invoice ID"
// @Failure 404 {object} map[string]string "Invoice not found"
// @Router /invoice/{id} [get]
func GetInvoice(c *gin.Context) {
	invoiceID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(invoiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	var invoice models.Invoice
	err = config.DB.Collection("invoices").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&invoice)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// GetInvoicesByCompanyID godoc
// @Summary Get all invoices for a specific company
// @Description Retrieve all invoices associated with a given company ID.
// @Tags Invoices
// @Produce json
// @Param company_id path string true "Company ID"
// @Success 200 {array} []models.Invoice "Invoices retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid company ID"
// @Failure 404 {object} map[string]string "No invoices found for this company"
// @Failure 500 {object} map[string]string "Failed to retrieve invoices"
// @Router /invoice/companies/{company_id} [get]
func GetInvoicesByCompanyID(c *gin.Context) {
	companyID := c.Param("company_id")

	var invoices []models.Invoice
	cursor, err := config.DB.Collection("invoices").Find(context.Background(), bson.M{"company_id": companyID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invoices"})
		return
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &invoices); err != nil {
		if err.Error() == "mongo: no documents in result" {
			c.JSON(http.StatusNotFound, gin.H{"message": "No invoices found for this company"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode invoices"})
		return
	}

	c.JSON(http.StatusOK, invoices)
}

// SendInvoice godoc
// @Summary Send an invoice via email
// @Description Send a generated invoice to the customer via email
// @Tags Invoices
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} map[string]string "Invoice sent successfully"
// @Failure 400 {object} map[string]string "Invalid ID or customer"
// @Failure 404 {object} map[string]string "Invoice or customer not found"
// @Failure 500 {object} map[string]string "Failed to send email"
// @Router /invoice/send/{id} [post]
func SendInvoice(c *gin.Context) {
	invoiceID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(invoiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	var invoice models.Invoice
	err = config.DB.Collection("invoices").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&invoice)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	customerID, err := primitive.ObjectIDFromHex(invoice.CustomerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}
	var customer models.Customer
	err = config.DB.Collection("customers").FindOne(context.Background(), bson.M{"_id": customerID}).Decode(&customer)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	emailBody := fmt.Sprintf(
		"Hello %s,\n\nHere is your invoice:\n\nReference: %s\nDate: %s\nAmount: $%.2f\nDue Date: %s\n\nItems:\n",
		customer.Name,
		invoice.ReferenceNumber,
		invoice.Date.Format("2006-01-02"),
		invoice.Amount,
		invoice.DueDate.Format("2006-01-02"),
	)
	for _, item := range invoice.Items {
		emailBody += fmt.Sprintf("- %s x%d @ %.2f (Discount: %.2f%%) → Subtotal: $%.2f\n",
			item.ItemName, item.Quantity, item.UnitPrice, item.Discount, item.Subtotal)
	}

	emailBody += "\nThank you for your business!\n"

	m := gomail.NewMessage()
	m.SetHeader("From", "bisrattewodros3@gmail.com")
	m.SetHeader("To", customer.Email)
	m.SetHeader("Subject", "Your Invoice")
	m.SetBody("text/plain", emailBody)

	d := gomail.NewDialer("smtp.gmail.com", 587, "bisrattewodros3@gmail.com", "xtkd pntw wrfq rdak")

	if err := d.DialAndSend(m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice sent successfully to " + customer.Email})
}

// DownloadInvoice godoc
// @Summary Download invoice as PDF
// @Description Download a specific invoice by ID
// @Tags Invoices
// @Produce application/pdf
// @Param id path string true "Invoice ID"
// @Success 200 {file} binary
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 404 {object} map[string]string "Invoice not found"
// @Failure 500 {object} map[string]string "Failed to generate or send PDF"
// @Router /invoice/download/{id} [get]
func DownloadInvoice(c *gin.Context) {
	invoiceID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(invoiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	var invoice models.Invoice
	err = config.DB.Collection("invoices").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&invoice)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Invoice")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Invoice ID: %s", invoice.ID))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("Reference #: %s", invoice.ReferenceNumber))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("Date: %s", invoice.Date.Format("2006-01-02")))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("Due Date: %s", invoice.DueDate.Format("2006-01-02")))
	pdf.Ln(6)

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Items:")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	for _, item := range invoice.Items {
		pdf.CellFormat(0, 8, fmt.Sprintf("- %s x%d @ %.2f each (%.0f%% off) = %.2f", item.ItemName, item.Quantity, item.UnitPrice, item.Discount, item.Subtotal), "", 1, "", false, 0, "")
	}

	pdf.Ln(6)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Total Amount: %.2f", invoice.Amount))

	err = pdf.Output(c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=invoice_"+invoice.ID.Hex()+".pdf")
	c.Header("Content-Type", "application/pdf")
}
