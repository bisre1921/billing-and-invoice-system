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
// @Security BearerAuth
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
// @Security BearerAuth
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

// GetAllEmployees godoc
// @Summary Get all employees for a company
// @Description Business Owner views all employees for their company
// @Tags Employee
// @Accept json
// @Produce json
// @Param company_id query string true "Company ID"
// @Success 200 {array} models.Employee
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /employee/all [get]
// @Security BearerAuth
func GetAllEmployees(c *gin.Context) {
	companyIDParam := c.Query("company_id")
	if companyIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company ID is required"})
		return
	}

	companyID, err := primitive.ObjectIDFromHex(companyIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	filter := bson.M{"company_id": companyID}
	cursor, err := config.DB.Collection("employees").Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employees"})
		return
	}
	defer cursor.Close(context.Background())

	var employees []models.Employee
	if err := cursor.All(context.Background(), &employees); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse employee data"})
		return
	}

	c.JSON(http.StatusOK, employees)
}

// GetEmployee godoc
// @Summary Get an employee by ID
// @Description Business Owner views an employee by ID
// @Tags Employee
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} models.Employee
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /employee/{id} [get]
// @Security BearerAuth
func GetEmployee(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	filter := bson.M{"_id": objID}
	var employee models.Employee
	err = config.DB.Collection("employees").FindOne(context.Background(), filter).Decode(&employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employee"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

// UpdateEmployee godoc
// @Summary Update an employee
// @Description Business Owner updates an employee by ID
// @Tags Employee
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Param employee body models.Employee true "Employee Data"
// @Success 200 {object} models.GenericResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /employee/update/{id} [put]
// @Security BearerAuth
func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee.UpdatedAt = time.Now()

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": employee}
	result, err := config.DB.Collection("employees").UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Employee updated successfully",
	})
}
