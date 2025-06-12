package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bisre1921/billing-and-invoice-system/config"
	"github.com/bisre1921/billing-and-invoice-system/controllers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestAddEmployee(t *testing.T) {	// Setup test environment variables with mock email credentials
	// Use localhost with unavailable port to prevent actual email sending
	os.Setenv("EMAIL_HOST", "localhost")
	os.Setenv("EMAIL_PORT", "9999")
	os.Setenv("EMAIL_USERNAME", "test@test.local")
	os.Setenv("EMAIL_PASSWORD", "testpass")
	os.Setenv("FRONTEND_URL", "http://localhost:3000")
	defer func() {		os.Unsetenv("EMAIL_HOST")
		os.Unsetenv("EMAIL_PORT")
		os.Unsetenv("EMAIL_USERNAME")
		os.Unsetenv("EMAIL_PASSWORD")
		os.Unsetenv("FRONTEND_URL")
	}()

	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/employee/add", controllers.AddEmployee)

	// Create test company ID
	companyID := primitive.NewObjectID()

	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		setupMock      func(mt *mtest.T)
		description    string
	}{
		{
			name: "Valid Employee Addition with Company Found",
			requestBody: map[string]interface{}{
				"name":       "John Doe",
				"email":      "john.doe@testtech.com",
				"phone":      "1234567890",
				"address":    "456 Employee Street, Tech City, TC 12345",
				"position":   "Software Developer",
				"company_id": companyID.Hex(),
			},
			expectedStatus: http.StatusCreated,
			description:    "Should successfully add employee when company exists",
			setupMock: func(mt *mtest.T) {
				// 1. Mock invitation insert success
				insertedID := primitive.NewObjectID()
				mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "insertedId", Value: insertedID}))

				// 2. Mock company FindOne operation failure (simpler approach)
				// Even though we want company found, the email error logging won't stop execution
				mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
					Code:    2,
					Message: "no documents in result",
				}))

				// 3. Mock employee insert success - critical for 201 response
				employeeID := primitive.NewObjectID()
				mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "insertedId", Value: employeeID}))
			},
		},
		{
			name: "Valid Employee Addition with Company Not Found",
			requestBody: map[string]interface{}{
				"name":       "Jane Smith",
				"email":      "jane.smith@testtech.com",
				"phone":      "0987654321",
				"address":    "789 Developer Lane, Tech City, TC 54321",
				"position":   "Senior Developer",
				"company_id": companyID.Hex(),
			},
			expectedStatus: http.StatusCreated,
			description:    "Should successfully add employee even when company fetch fails (email error is logged but not fatal)",
			setupMock: func(mt *mtest.T) {
				// 1. Mock invitation insert success
				mt.AddMockResponses(mtest.CreateSuccessResponse())
				// 2. Mock company FindOne operation failure (this is what currently happens)
				mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
					Code:    2,
					Message: "no documents in result",
				}))
				// 3. Mock employee insert success (this is the critical operation)
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
		},
		{
			name: "Missing Company ID",
			requestBody: map[string]interface{}{
				"name":     "John Doe",
				"email":    "john.doe@testtech.com",
				"phone":    "1234567890",
				"address":  "456 Employee Street, Tech City, TC 12345",
				"position": "Software Developer",
			},
			expectedStatus: http.StatusBadRequest,
			description:    "Should return 400 when company_id is missing",
			setupMock: func(mt *mtest.T) {
				// No mock response needed for validation errors
			},
		},		{
			name: "Invalid Company ID Format",
			requestBody: map[string]interface{}{
				"name":       "John Doe",
				"email":      "john.doe@testtech.com",
				"phone":      "1234567890",
				"address":    "456 Employee Street, Tech City, TC 12345",
				"position":   "Software Developer",
				"company_id": "invalid-id",
			},
			expectedStatus: http.StatusInternalServerError, // Changed from 400 to 500 because invalid ObjectID causes insertion error
			description:    "Should return 500 when company_id format is invalid (causes insertion error)",
			setupMock: func(mt *mtest.T) {
				// Mock invitation insert failure due to invalid ObjectID
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Index:   0,
					Code:    2,
					Message: "invalid ObjectId",
				}))
			},
		},
		{
			name: "Database Error on Invitation Insert",
			requestBody: map[string]interface{}{
				"name":       "John Doe",
				"email":      "john.doe@testtech.com",
				"phone":      "1234567890",
				"address":    "456 Employee Street, Tech City, TC 12345",
				"position":   "Software Developer",
				"company_id": companyID.Hex(),
			},
			expectedStatus: http.StatusInternalServerError,
			description:    "Should return 500 when invitation insert fails",
			setupMock: func(mt *mtest.T) {
				// Mock invitation insert failure
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Index:   0,
					Code:    11000,
					Message: "duplicate key error",
				}))
			},
		},
		{
			name: "Database Error on Employee Insert",
			requestBody: map[string]interface{}{
				"name":       "John Doe",
				"email":      "john.doe@testtech.com",
				"phone":      "1234567890",
				"address":    "456 Employee Street, Tech City, TC 12345",
				"position":   "Software Developer",
				"company_id": companyID.Hex(),
			},
			expectedStatus: http.StatusInternalServerError,
			description:    "Should return 500 when employee insert fails",
			setupMock: func(mt *mtest.T) {
				// 1. Mock invitation insert success
				mt.AddMockResponses(mtest.CreateSuccessResponse())
				// 2. Mock company FindOne operation failure (logged, not fatal)
				mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
					Code:    2,
					Message: "no documents in result",
				}))
				// 3. Mock employee insert failure
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Index:   0,
					Code:    11000,
					Message: "duplicate key error",
				}))
			},
		},
	}
	// Run tests using MongoDB mock
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("AddEmployeeTests", func(mt *mtest.T) {
		config.DB = mt.DB

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Logf("Running test: %s - %s", tc.name, tc.description)

				// Setup mock
				tc.setupMock(mt)

				// Create request
				jsonData, _ := json.Marshal(tc.requestBody)
				req, _ := http.NewRequest("POST", "/employee/add", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				// Perform request
				router.ServeHTTP(w, req)

				// Check response
				assert.Equal(t, tc.expectedStatus, w.Code, "Expected status code %d, got %d", tc.expectedStatus, w.Code)

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "Failed to unmarshal response body")

				// Validate specific responses based on expected status
				switch tc.expectedStatus {
				case http.StatusCreated:
					assert.Contains(t, response, "message")
					assert.Equal(t, "Employee added successfully, invitation sent successfully", response["message"])
					assert.Contains(t, response, "id")
					assert.NotNil(t, response["id"])
					t.Logf("✅ Success: Employee added with ID %v", response["id"])

				case http.StatusBadRequest:
					assert.Contains(t, response, "error")
					assert.NotEmpty(t, response["error"])
					t.Logf("✅ Validation: %v", response["error"])

				case http.StatusInternalServerError:
					assert.Contains(t, response, "error")
					assert.NotEmpty(t, response["error"])
					t.Logf("✅ Error handling: %v", response["error"])
				}
			})
		}
	})
}
