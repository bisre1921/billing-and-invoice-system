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

func TestAddItem(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/item/add", controllers.AddItem)
	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		setupMock      func(mt *mtest.T)
	}{		{
			name: "Valid Item Creation",
			requestBody: map[string]interface{}{
				"name":          "Test Item",
				"description":   "A test item for unit testing",
				"code":          "TEST001",
				"category":      "Electronics",
				"selling_price": 29.99,
				"buying_price":  19.99,
				"quantity":      100,
				"unit":          "pcs",
				"company_id":    primitive.NewObjectID().Hex(),
				"discount":      10.0,
			},
			expectedStatus: http.StatusCreated,
			setupMock: func(mt *mtest.T) {
				// Mock the FindOne call for existing item check (should return no document found)
				mt.AddMockResponses(
					mtest.CreateCommandErrorResponse(mtest.CommandError{
						Code:    4,
						Message: "not found",
						Name:    "CommandNotFound",
					}),
					mtest.CreateSuccessResponse(),
				)
			},
		},
	}
	// Run tests using MongoDB mock
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("AddItemTests", func(mt *mtest.T) {
		config.DB = mt.DB

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Setup mock
				tc.setupMock(mt)

				// Create request
				jsonData, _ := json.Marshal(tc.requestBody)
				req, _ := http.NewRequest("POST", "/item/add", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				// Perform request
				router.ServeHTTP(w, req)

				// Check response
				assert.Equal(t, tc.expectedStatus, w.Code)

				if tc.expectedStatus == http.StatusCreated {
					var response map[string]interface{}
					err := json.Unmarshal(w.Body.Bytes(), &response)
					assert.NoError(t, err)
					assert.Contains(t, response, "message")
					assert.Equal(t, "Item added successfully", response["message"])
					assert.Contains(t, response, "id")
				}
			})
		}
	})
}

// End of test file
