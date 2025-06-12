package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bisre1921/billing-and-invoice-system/config"
	"github.com/bisre1921/billing-and-invoice-system/controllers"
	"github.com/bisre1921/billing-and-invoice-system/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestGenerateInvoice(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/invoice/generate", controllers.GenerateInvoice)
	// Test cases
	testCases := []struct {
		name           string
		requestBody    models.Invoice
		expectedStatus int
		setupMock      func(mt *mtest.T)
	}{		{
			name: "Valid Invoice Generation",
			requestBody: models.Invoice{
				CustomerID:      primitive.NewObjectID().Hex(),
				CompanyID:       primitive.NewObjectID().Hex(),
				ReferenceNumber: "INV-2025-001",
				PaymentType:     "cash",
				Terms:           "Net 30",
				Items: []models.InvoiceItem{
					{
						ItemID:    primitive.NewObjectID().Hex(),
						ItemName:  "Test Product",
						Quantity:  2,
						UnitPrice: 49.99,
						Discount:  10,
					},
					{
						ItemID:    primitive.NewObjectID().Hex(),
						ItemName:  "Another Product",
						Quantity:  1,
						UnitPrice: 29.99,
						Discount:  0,
					},
				},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(mt *mtest.T) {
				mt.AddMockResponses(bson.D{
					{Key: "ok", Value: 1},
					{Key: "insertedId", Value: primitive.NewObjectID()},
				})
			},		},
	}
	// Run tests using MongoDB mock
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("GenerateInvoiceTests", func(mt *mtest.T) {
		config.DB = mt.DB

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Setup mock
				tc.setupMock(mt)

				// Create request
				jsonData, _ := json.Marshal(tc.requestBody)
				req, _ := http.NewRequest("POST", "/invoice/generate", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				// Perform request
				router.ServeHTTP(w, req)

				// Check response
				assert.Equal(t, tc.expectedStatus, w.Code)

				if tc.expectedStatus == http.StatusOK {
					var response map[string]interface{}
					err := json.Unmarshal(w.Body.Bytes(), &response)
					assert.NoError(t, err)
					assert.Contains(t, response, "message")
					assert.Equal(t, "Invoice generated successfully", response["message"])
					assert.Contains(t, response, "invoice")

					// Validate invoice data
					invoice := response["invoice"].(map[string]interface{})
					assert.NotNil(t, invoice["id"])
					assert.Equal(t, "Paid", invoice["status"])

					// Verify calculated amount
					var expectedTotal float64 = 0
					for _, item := range tc.requestBody.Items {
						discountAmount := item.UnitPrice * float64(item.Discount) / 100
						subtotal := float64(item.Quantity) * (item.UnitPrice - discountAmount)
						expectedTotal += subtotal
					}
					assert.Equal(t, expectedTotal, invoice["amount"])
				}
			})
		}	})
}

// End of test file
