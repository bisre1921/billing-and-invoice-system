package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/bisre1921/billing-and-invoice-system/config"
	"github.com/bisre1921/billing-and-invoice-system/models"
	"github.com/gin-gonic/gin"
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
