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
)

// RegisterCustomer godoc
// @Summary Register a new customer
// @Description Business Owner or Employee adds a new customer
// @Tags Customer
// @Accept json
// @Produce json
// @Param customer body models.Customer true "Customer Data"
// @Success 201 {object} models.GenericResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /customer/register [post]
// @Security BearerAuth
func RegisterCustomer(c *gin.Context) {
	var customer models.Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate Company ID
	if customer.CompanyID.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CompanyID is required"})
		return
	}

	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	result, err := config.DB.Collection("customers").InsertOne(context.Background(), customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register customer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Customer registered successfully",
		"id":      result.InsertedID,
	})
}

// UpdateCustomer godoc
// @Summary Update a customer by ID
// @Description Business Owner or Employee updates customer details
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param customer body models.Customer true "Updated Customer Info"
// @Success 200 {object} models.GenericResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /customer/update/{id} [put]
// @Security BearerAuth
func UpdateCustomer(c *gin.Context) {
	customerID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var updatedCustomer models.Customer
	if err := c.ShouldBindJSON(&updatedCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"name":       updatedCustomer.Name,
			"email":      updatedCustomer.Email,
			"phone":      updatedCustomer.Phone,
			"address":    updatedCustomer.Address,
			"company_id": updatedCustomer.CompanyID,
			"updated_at": time.Now(),
		},
	}

	result, err := config.DB.Collection("customers").UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		update,
	)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Customer update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Customer updated successfully",
	})
}
