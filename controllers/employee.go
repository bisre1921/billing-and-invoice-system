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

// AddEmployee godoc
// @Summary Add a new employee
// @Description Business Owner adds a new employee
// @Tags Employee
// @Accept json
// @Produce json
// @Param employee body models.Employee true "Employee Data"
// @Success 201 {object} models.GenericResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /employee/add [post]
func AddEmployee(c *gin.Context) {
	var employee models.Employee

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate Company ID
	if employee.CompanyID.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CompanyID is required"})
		return
	}

	employee.CreatedAt = time.Now()
	employee.UpdatedAt = time.Now()

	result, err := config.DB.Collection("employees").InsertOne(context.Background(), employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add employee"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Employee added successfully",
		"id":      result.InsertedID,
	})
}

// DeleteEmployee godoc
// @Summary Delete an employee
// @Description Business Owner deletes an employee by ID
// @Tags Employee
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} models.GenericResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /employee/delete/{id} [delete]
func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	filter := bson.M{"_id": objID}
	res, err := config.DB.Collection("employees").DeleteOne(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employee"})
		return
	}

	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Employee deleted successfully",
		"id":      id,
	})
}
