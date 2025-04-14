package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/bisre1921/billing-and-invoice-system/config"
	"github.com/bisre1921/billing-and-invoice-system/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateCompany godoc
// @Summary Creates a new company
// @Description Registers a new company in the system
// @Tags Company
// @Accept json
// @Produce json
// @Param company body models.Company true "Company Data"
// @Success 201 {object} map[string]interface{} "Company created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /company/create [post]
// @Security BearerAuth
func CreateCompany(c *gin.Context) {
	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idFromToken, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication token is required."})
		c.Abort()
		return
	}

	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()
	company.Owner = idFromToken.(string)

	result, err := config.DB.Collection("companies").InsertOne(context.Background(), company)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Company created successfully", "company_id": result.InsertedID})
}

// GetCompany godoc
// @Summary Get company by ID
// @Description Retrieves a single company by its ID if owned by the authenticated user
// @Tags Company
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} models.Company "Company retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Company not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /company/{id} [get]
// @Security BearerAuth
func GetCompany(c *gin.Context) {
	companyId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	var company models.Company
	err = config.DB.Collection("companies").FindOne(context.Background(), bson.M{"_id": companyId}).Decode(&company)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, company)
}
