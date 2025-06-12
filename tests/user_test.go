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
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestRegisterUser(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/auth/register/user", controllers.RegisterUser)

	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		setupMock      func(mt *mtest.T)
	}{		{
			name: "Valid User Registration",
			requestBody: map[string]interface{}{
				"name":     "Test User",
				"email":    "testuser@example.com",
				"password": "securepassword123",
				"role":     "owner",			},
			expectedStatus: http.StatusCreated,
			setupMock: func(mt *mtest.T) {
				// Mock the FindOne call for existing user check - simulate no document found
				// First response: no document found (allows registration to proceed)
				// Second response: successful insert
				mt.AddMockResponses(
					mtest.CreateCursorResponse(0, "billing_invoice.users", mtest.FirstBatch),
					mtest.CreateSuccessResponse(),
				)
			},
		},
	}

	// Run tests using MongoDB mock
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("RegisterUserTests", func(mt *mtest.T) {
		config.DB = mt.DB

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Setup mock
				tc.setupMock(mt)

				// Create request
				jsonData, _ := json.Marshal(tc.requestBody)
				req, _ := http.NewRequest("POST", "/auth/register/user", bytes.NewBuffer(jsonData))
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
					assert.Equal(t, "User created successfully", response["message"])
					assert.Contains(t, response, "id")
				}
			})
		}
	})
}
