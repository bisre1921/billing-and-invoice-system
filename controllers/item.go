package controllers

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/bisre1921/billing-and-invoice-system/config"
	"github.com/bisre1921/billing-and-invoice-system/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	// Validate item code
	if item.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item code is required"})
		return
	}

	// Validate item category
	if item.Category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item category is required"})
		return
	}

	// Check if item code already exists for this company
	var existingItem models.Item
	err := config.DB.Collection("items").FindOne(context.Background(),
		bson.M{
			"code":       item.Code,
			"company_id": item.CompanyID,
		}).Decode(&existingItem)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Item with this code already exists"})
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

	item.UpdatedAt = time.Now() // If code is being updated, check if new code already exists for this company
	if item.Code != "" {
		var existingItem models.Item
		err := config.DB.Collection("items").FindOne(context.Background(),
			bson.M{
				"code":       item.Code,
				"company_id": item.CompanyID,
				"_id":        bson.M{"$ne": itemID}, // exclude current item
			}).Decode(&existingItem)

		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Item with this code already exists"})
			return
		}
	}
	update := bson.M{
		"$set": bson.M{
			"code":          item.Code,
			"name":          item.Name,
			"description":   item.Description,
			"category":      item.Category,
			"selling_price": item.SellingPrice,
			"unit":          item.Unit,
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

// GetItemsByCompanyID godoc
// @Summary Get all items for a specific company
// @Description Retrieves all items belonging to a specific company.
// @Tags Item
// @Accept json
// @Produce json
// @Param company_id path string true "Company ID"
// @Success 200 {array} models.Item "Successfully retrieved items"
// @Failure 400 {object} models.ErrorResponse "Invalid company ID"
// @Failure 404 {object} models.ErrorResponse "No items found for the company"
// @Failure 500 {object} models.ErrorResponse "Failed to retrieve items"
// @Router /item/company/company_id [get]
func GetItemsByCompanyID(c *gin.Context) {
	companyIDParam := c.Param("company_id")
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
	cursor, err := config.DB.Collection("items").Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	defer cursor.Close(context.Background())

	var items []models.Item
	if err := cursor.All(context.Background(), &items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode items"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetItem godoc
// @Summary Get an item by ID
// @Description Retrieves a single item by its ID
// @Tags Item
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} models.Item "Item found"
// @Failure 400 {object} models.ErrorResponse "Invalid item ID"
// @Failure 404 {object} models.ErrorResponse "Item not found"
// @Failure 500 {object} models.ErrorResponse "Failed to fetch item"
// @Router /item/{id} [get]
// @Security BearerAuth
func GetItem(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var item models.Item
	err = config.DB.Collection("items").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// ImportItems godoc
// @Summary Import items from CSV file
// @Description Import multiple items from a CSV file
// @Tags Item
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file containing item data"
// @Param company_id formData string true "Company ID"
// @Success 201 {object} models.GenericResponse "Items imported successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid file format or data"
// @Failure 500 {object} models.ErrorResponse "Failed to import items"
// @Router /item/import [post]
// @Security BearerAuth
func ImportItems(c *gin.Context) {
	// Get company ID from form data
	companyID := c.PostForm("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company ID is required"})
		return
	}

	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Check file extension
	if file.Filename[len(file.Filename)-4:] != ".csv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only CSV files are allowed"})
		return
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	// Create CSV reader
	reader := csv.NewReader(src)

	// Read header
	header, err := reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read CSV header"})
		return
	}

	// Validate header
	expectedHeaders := []string{"code", "name", "description", "category", "selling_price", "unit"}
	if len(header) < len(expectedHeaders) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CSV format. Required columns: code, name, description, category, selling_price, unit"})
		return
	}

	var items []interface{}
	var duplicateCodes []string
	lineNumber := 2 // Start from 2 because line 1 is header

	// Read and process each row
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error reading line %d: %v", lineNumber, err)})
			return
		}

		// Check if item code already exists
		var existingItem models.Item
		err = config.DB.Collection("items").FindOne(context.Background(),
			bson.M{
				"code":       row[0],
				"company_id": objCompanyID,
			}).Decode(&existingItem)

		if err == nil {
			// Item code already exists
			duplicateCodes = append(duplicateCodes, row[0])
			lineNumber++
			continue
		}

		// Parse selling price
		sellingPrice, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Invalid selling price at line %d: %s", lineNumber, row[4]),
			})
			return
		}

		// Create new item
		item := models.Item{
			ID:           primitive.NewObjectID(),
			Code:         row[0],
			Name:         row[1],
			Description:  row[2],
			Category:     row[3],
			SellingPrice: sellingPrice,
			Unit:         row[5],
			CompanyID:    objCompanyID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		items = append(items, item)
		lineNumber++
	}

	if len(items) == 0 {
		if len(duplicateCodes) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("All items have duplicate codes: %v", duplicateCodes),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid items found in CSV"})
		return
	}

	// Insert all valid items
	_, err = config.DB.Collection("items").InsertMany(context.Background(), items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import items"})
		return
	}

	response := gin.H{
		"message": fmt.Sprintf("Successfully imported %d items", len(items)),
	}
	if len(duplicateCodes) > 0 {
		response["skipped_codes"] = duplicateCodes
		response["message"] = fmt.Sprintf("Imported %d items. Skipped %d items with duplicate codes",
			len(items), len(duplicateCodes))
	}

	c.JSON(http.StatusCreated, response)
}
