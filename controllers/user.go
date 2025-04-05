package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bisre1921/billing-and-invoice-system/config"
	"github.com/bisre1921/billing-and-invoice-system/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateUser godoc
// @Summary Update user details
// @Description Updates user details by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UpdateUserInput true "User Data"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user/update/{id} [patch]
// @Security BearerAuth
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	fmt.Println(userID)

	idFromToken, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication token is required."})
		c.Abort()
		return
	}

	if idFromToken != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		c.Abort()
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var user models.User
	if err := config.DB.Collection("users").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user); err != nil {
		fmt.Println(err)

		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input models.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentTime := time.Now()

	input.UpdatedAt = currentTime

	_, err = config.DB.Collection("users").UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": input})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}
