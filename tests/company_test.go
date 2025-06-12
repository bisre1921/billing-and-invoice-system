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

func TestCreateCompany(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Mock authentication middleware
	router.Use(func(c *gin.Context) {
		// Mock user ID from JWT token
		c.Set("userID", primitive.NewObjectID().Hex())
		c.Next()
	})

	router.POST("/company/create", controllers.CreateCompany)

	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		setupMock      func(mt *mtest.T)
	}{
		{
			name: "Valid Company Creation",
			requestBody: map[string]interface{}{
				"name":    "Test Tech Solutions",
				"email":   "info@testtech.com",
				"address": "123 Business District, Tech City, TC 12345",
			},
			expectedStatus: http.StatusCreated,
			setupMock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
		},
		{
			name: "Database Error",
			requestBody: map[string]interface{}{
				"name":    "Test Tech Solutions",
				"email":   "info@testtech.com",
				"address": "123 Business District, Tech City, TC 12345",
			},
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Index:   1,
					Code:    11000,
					Message: "E11000 duplicate key error",
				}))
			},
		},		{
			name: "Company with Missing Name",
			requestBody: map[string]interface{}{
				"email":   "info@testtech.com",
				"address": "123 Business District, Tech City, TC 12345",
			},
			expectedStatus: http.StatusCreated, // Controller doesn't validate required fields
			setupMock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
		},
		{
			name: "Company with Missing Email",
			requestBody: map[string]interface{}{
				"name":    "Test Tech Solutions",
				"address": "123 Business District, Tech City, TC 12345",
			},
			expectedStatus: http.StatusCreated, // Controller doesn't validate required fields
			setupMock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
		},
		{
			name: "Company with Empty Name",
			requestBody: map[string]interface{}{
				"name":    "",
				"email":   "info@testtech.com",
				"address": "123 Business District, Tech City, TC 12345",
			},
			expectedStatus: http.StatusCreated, // Controller doesn't validate required fields
			setupMock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
		},
	}

	// Run tests using MongoDB mock
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("CreateCompanyTests", func(mt *mtest.T) {
		config.DB = mt.DB

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Setup mock
				tc.setupMock(mt)

				// Create request
				jsonData, _ := json.Marshal(tc.requestBody)
				req, _ := http.NewRequest("POST", "/company/create", bytes.NewBuffer(jsonData))
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
					assert.Equal(t, "Company created successfully", response["message"])
					assert.Contains(t, response, "company_id")
					assert.NotNil(t, response["company_id"])
				}
			})
		}
	})
}

func TestCreateCompanyWithoutAuth(t *testing.T) {
	// Setup without authentication middleware
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/company/create", controllers.CreateCompany)

	// Test case for missing authentication
	requestBody := map[string]interface{}{
		"name":    "Test Tech Solutions",
		"email":   "info@testtech.com",
		"address": "123 Business District, Tech City, TC 12345",
	}

	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/company/create", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "Authentication token is required.", response["error"])
}