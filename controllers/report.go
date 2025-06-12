package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/bisre1921/billing-and-invoice-system/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SalesReportRequest struct {
	CompanyID   string   `json:"company_id" binding:"required"`
	DateRange   string   `json:"date_range" binding:"required"` // e.g. today, last_7_days, last_month, last_3_months, custom
	CustomStart *string  `json:"custom_start,omitempty"`        // for custom range
	CustomEnd   *string  `json:"custom_end,omitempty"`
	Statuses    []string `json:"statuses"`   // ["Paid", "Unpaid"]
	Categories  []string `json:"categories"` // e.g. ["Electronics", ...]
}

type SalesReportItem struct {
	InvoiceID    string    `json:"invoice_id"`
	Date         time.Time `json:"date"`
	Status       string    `json:"status"`
	CustomerName string    `json:"customer_name"`
	Items        []struct {
		Name      string  `json:"name"`
		Category  string  `json:"category"`
		Quantity  int     `json:"quantity"`
		UnitPrice float64 `json:"unit_price"`
		Subtotal  float64 `json:"subtotal"`
	} `json:"items"`
	TotalAmount float64 `json:"total_amount"`
}

// GetSalesReport godoc
// @Summary Get sales report
// @Description Returns sales data for a company over a given period, with filters for status and item category.
// @Tags Reports
// @Accept json
// @Produce json
// @Param filters body SalesReportRequest true "Report filters"
// @Success 200 {array} SalesReportItem
// @Failure 400 {object} map[string]string
// @Router /report/sales [post]
func GetSalesReport(c *gin.Context) {
	var req SalesReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Date range
	var start, end time.Time
	now := time.Now()
	switch req.DateRange {
	case "today":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = start.Add(24 * time.Hour)
	case "last_7_days":
		end = now
		start = end.AddDate(0, 0, -7)
	case "last_month":
		end = now
		start = end.AddDate(0, -1, 0)
	case "last_3_months":
		end = now
		start = end.AddDate(0, -3, 0)
	case "custom":
		if req.CustomStart == nil || req.CustomEnd == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Custom start and end dates required"})
			return
		}
		var err error
		start, err = time.Parse("2006-01-02", *req.CustomStart)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid custom_start format (expected YYYY-MM-DD)"})
			return
		}
		end, err = time.Parse("2006-01-02", *req.CustomEnd)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid custom_end format (expected YYYY-MM-DD)"})
			return
		}
		end = end.Add(24 * time.Hour)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_range value"})
		return
	}

	filter := bson.M{
		"company_id": req.CompanyID,
		"date":       bson.M{"$gte": start, "$lt": end},
	}
	if len(req.Statuses) > 0 {
		filter["status"] = bson.M{"$in": req.Statuses}
	}

	cursor, err := config.DB.Collection("invoices").Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invoices"})
		return
	}
	defer cursor.Close(context.Background())

	var report []SalesReportItem
	for cursor.Next(context.Background()) {
		var inv struct {
			ID         primitive.ObjectID `bson:"_id"`
			Date       time.Time          `bson:"date"`
			Status     string             `bson:"status"`
			CustomerID string             `bson:"customer_id"`
			Items      []struct {
				ItemID    string  `bson:"item_id"`
				ItemName  string  `bson:"item_name"`
				Category  string  `bson:"category"`
				Quantity  int     `bson:"quantity"`
				UnitPrice float64 `bson:"unit_price"`
				Subtotal  float64 `bson:"subtotal"`
			} `bson:"items"`
			Amount float64 `bson:"amount"`
		}
		if err := cursor.Decode(&inv); err != nil {
			continue
		}

		// Filter by item category if needed
		var filteredItems []struct {
			Name      string  `json:"name"`
			Category  string  `json:"category"`
			Quantity  int     `json:"quantity"`
			UnitPrice float64 `json:"unit_price"`
			Subtotal  float64 `json:"subtotal"`
		}
		for _, it := range inv.Items {
			if len(req.Categories) == 0 || contains(req.Categories, it.Category) {
				filteredItems = append(filteredItems, struct {
					Name      string  `json:"name"`
					Category  string  `json:"category"`
					Quantity  int     `json:"quantity"`
					UnitPrice float64 `json:"unit_price"`
					Subtotal  float64 `json:"subtotal"`
				}{
					Name: it.ItemName, Category: it.Category, Quantity: it.Quantity, UnitPrice: it.UnitPrice, Subtotal: it.Subtotal,
				})
			}
		}
		if len(filteredItems) == 0 {
			continue
		}

		// Fetch customer name (optional, can be optimized)
		customerName := ""
		if inv.CustomerID != "" {
			var cust struct {
				Name string `bson:"name"`
			}
			_ = config.DB.Collection("customers").FindOne(context.Background(), bson.M{"_id": inv.CustomerID}).Decode(&cust)
			customerName = cust.Name
		}

		report = append(report, SalesReportItem{
			InvoiceID:    inv.ID.Hex(),
			Date:         inv.Date,
			Status:       inv.Status,
			CustomerName: customerName,
			Items:        filteredItems,
			TotalAmount:  inv.Amount,
		})
	}

	// Store the report in MongoDB as JSON string
	reportJSON, err := json.Marshal(report)
	if err == nil {
		_, _ = StoreSalesReport(req.CompanyID, "Sales Report", "Auto-generated sales report", "system", string(reportJSON))
	}

	c.JSON(http.StatusOK, report)
}

// Store a generated report in MongoDB
func StoreSalesReport(companyID, title, description, createdBy string, content string) (primitive.ObjectID, error) {
	if content == "" || content == "null" || content == "[]" {
		return primitive.NilObjectID, nil
	}
	report := bson.M{
		"company_id":         companyID,
		"title":              title,
		"description":        description,
		"created_by":         createdBy,
		"created_date":       time.Now(),
		"last_modified_date": time.Now(),
		"type":               "sales",
		"status":             "Generated",
		"content":            content,
	}
	res, err := config.DB.Collection("reports").InsertOne(context.Background(), report)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

// ListReports godoc
// @Summary List all stored reports
// @Description Fetch all stored reports (basic info only)
// @Tags Reports
// @Produce json
// @Success 200 {array} object
// @Failure 500 {object} map[string]string
// @Router /report/all [get]
func ListReports(c *gin.Context) {
	cursor, err := config.DB.Collection("reports").Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reports"})
		return
	}
	defer cursor.Close(context.Background())
	var reports []bson.M
	for cursor.Next(context.Background()) {
		var r bson.M
		if err := cursor.Decode(&r); err == nil {
			reports = append(reports, bson.M{
				"id":           r["_id"],
				"company_id":   r["company_id"],
				"title":        r["title"],
				"description":  r["description"],
				"created_by":   r["created_by"],
				"created_date": r["created_date"],
				"type":         r["type"],
				"status":       r["status"],
			})
		}
	}
	c.JSON(http.StatusOK, reports)
}

// GetReportDetails godoc
// @Summary Get report details
// @Description Fetch details of a specific report
// @Tags Reports
// @Produce json
// @Param id path string true "Report ID"
// @Success 200 {object} object
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /report/{id} [get]
func GetReportDetails(c *gin.Context) {
	reportID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}
	var report bson.M
	err = config.DB.Collection("reports").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&report)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}
	c.JSON(http.StatusOK, report)
}

// DeleteReport godoc
// @Summary Delete a report
// @Description Delete a generated report by ID
// @Tags Reports
// @Param id path string true "Report ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /report/{id} [delete]
func DeleteReport(c *gin.Context) {
	reportID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}
	_, err = config.DB.Collection("reports").DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete report"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Report deleted successfully"})
}

// DownloadReportCSV godoc
// @Summary Download a report as CSV
// @Description Download a generated report in CSV format
// @Tags Reports
// @Produce text/csv
// @Param id path string true "Report ID"
// @Success 200 {file} binary
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /report/download/{id} [get]
func DownloadReportCSV(c *gin.Context) {
	reportID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}
	var report bson.M
	err = config.DB.Collection("reports").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&report)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}
	c.Header("Content-Disposition", "attachment; filename=report.csv")
	c.Header("Content-Type", "text/csv")
	c.Writer.Write([]byte(report["content"].(string)))
}

func contains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
