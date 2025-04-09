package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/bisre1921/billing-and-invoice-system/config"
	"github.com/bisre1921/billing-and-invoice-system/models"
	"github.com/gin-gonic/gin"
)

// AddItem godoc
// @Summary Add a new item
// @Description Business Owner or Employee adds a new item
// @Tags Item
// @Accept json
// @Produce json
// @Param item body models.Item true "Item Data"
// @Success 201 {object} models.GenericResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /item/add [post]
// @Security BearerAuth
func AddItem(c *gin.Context) {
	var item models.Item

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if item.CompanyID.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CompanyID is required"})
		return
	}

	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	result, err := config.DB.Collection("items").InsertOne(context.Background(), item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Item added successfully",
		"id":      result.InsertedID,
	})
}
