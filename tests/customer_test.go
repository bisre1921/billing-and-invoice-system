package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bisre1921/billing-and-invoice-system/config"
	"github.com/bisre1921/billing-and-invoice-system/controllers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestRegisterCustomer(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/customer/register", controllers.RegisterCustomer)

	// Create a test company ID
	companyID := primitive.NewObjectID()
	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		setupMock      func(mt *mtest.T)
	}{		{
			name: "Valid Customer Registration",
			requestBody: map[string]interface{}{
				"name":                     "Test Customer",
				"email":                    "customer@example.com",
				"phone":                    "1234567890",
				"address":                  "123 Test Street, Test City",
				"tin":                      "123456789",
				"max_credit_amount":        1000.0,
				"current_credit_available": 1000.0,
				"company_id":               companyID.Hex(),
			},
			expectedStatus: http.StatusCreated,
			setupMock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
		},
	}

	// Run tests using MongoDB mock
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("RegisterCustomerTests", func(mt *mtest.T) {
		config.DB = mt.DB

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Setup mock
				tc.setupMock(mt)

				// Create request
				jsonData, _ := json.Marshal(tc.requestBody)
				req, _ := http.NewRequest("POST", "/customer/register", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				// Perform request
				router.ServeHTTP(w, req)				// Check response
				assert.Equal(t, tc.expectedStatus, w.Code)

				if tc.expectedStatus == http.StatusCreated {
					var response map[string]interface{}
					err := json.Unmarshal(w.Body.Bytes(), &response)
					assert.NoError(t, err)
					assert.Contains(t, response, "message")
					assert.Equal(t, "Customer registered successfully", response["message"])
					assert.Contains(t, response, "id")
				}
			})
		}
	})
}
