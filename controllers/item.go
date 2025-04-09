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

// UpdateItem godoc
// @Summary Update an existing item
// @Description Business Owner or Employee updates item info
// @Tags Item
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param item body models.Item true "Updated item data"
// @Success 200 {object} models.GenericResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /item/update/{id} [put]
// @Security BearerAuth
func UpdateItem(c *gin.Context) {
	id := c.Param("id")
	itemID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"name":          item.Name,
			"description":   item.Description,
			"selling_price": item.SellingPrice,
			"buying_price":  item.BuyingPrice,
			"unit":          item.Unit,
			"quantity":      item.Quantity,
			"updated_at":    item.UpdatedAt,
		},
	}

	result, err := config.DB.Collection("items").UpdateByID(context.Background(), itemID, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item updated successfully"})
}

// DeleteItem godoc
// @Summary Delete an existing item
// @Description Business Owner or Employee deletes an item
// @Tags Item
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} models.GenericResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /item/delete/{id} [delete]
// @Security BearerAuth
func DeleteItem(c *gin.Context) {
	id := c.Param("id")
	itemID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	result, err := config.DB.Collection("items").DeleteOne(context.Background(), bson.M{"_id": itemID})
	if err != nil || result.DeletedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}

// ListItems godoc
// @Summary View all items
// @Description Shows all items if they exist, or a message if none are found
// @Tags Item
// @Accept json
// @Produce json
// @Success 200 {array} models.Item
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /item/all [get]
// @Security BearerAuth
func ListItems(c *gin.Context) {
	cursor, err := config.DB.Collection("items").Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong while retrieving items"})
		return
	}
	defer cursor.Close(context.Background())

	var items []models.Item
	if err := cursor.All(context.Background(), &items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding item list"})
		return
	}

	if len(items) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No items found in the system"})
		return
	}

	c.JSON(http.StatusOK, items)
}
