package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/bisre1921/billing-and-invoice-system/config"
	"github.com/bisre1921/billing-and-invoice-system/models"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GenerateReport godoc
// @Summary Generate a new report
// @Description Generate a new report
// @Tags Reports
// @Accept json
// @Produce json
// @Param report body models.Report true "Report data"
// @Success 200 {object} map[string]interface{} "Report generated successfully"
// @Failure 400 {object} map[string]string "Invalid report input"
// @Failure 500 {object} map[string]string "Failed to generate report"
// @Router /report/generate [post]
func GenerateReport(c *gin.Context) {
	var report models.Report
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(400, gin.H{"error": "Invalid report input", "details": err.Error()})
		return
	}

	report.CreatedDate = time.Now()
	report.LastModifiedDate = time.Now()

	res, err := config.DB.Collection("reports").InsertOne(context.Background(), report)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate report"})
		return
	}

	report.ID = res.InsertedID.(primitive.ObjectID)

	c.JSON(200, gin.H{"message": "Report generated successfully", "report": report})
}

// GetReport godoc
// @Summary Get a report by ID
// @Description Get a report by ID
// @Tags Reports
// @Accept json
// @Produce json
// @Param id path string true "Report ID"
// @Success 200 {object} map[string]interface{} "Report retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid report ID"
// @Failure 500 {object} map[string]string "Failed to get report"
// @Router /report/{id} [get]
func GetReport(c *gin.Context) {
	reportId := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(reportId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid report ID"})
		return
	}

	var report models.Report
	err = config.DB.Collection("reports").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&report)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get report"})
		return
	}

	c.JSON(200, gin.H{"report": report})
}

// GetReportsByCompanyID godoc
// @Summary Get all reports for a specific company
// @Description Get all reports for a specific company
// @Tags Reports
// @Accept json
// @Produce json
// @Param company_id path string true "Company ID"
// @Success 200 {array} models.Report "Reports retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid company ID"
// @Failure 500 {object} map[string]string "Failed to get reports"
// @Router /report/companies/{company_id} [get]
func GetReportsByCompanyID(c *gin.Context) {
	companyId := c.Param("company_id")

	var reports []models.Report
	cursor, err := config.DB.Collection("reports").Find(context.Background(), bson.M{"company_id": companyId})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get reports"})
		return
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &reports); err != nil {
		if err.Error() == "mongo: no documents in result" {
			c.JSON(http.StatusNotFound, gin.H{"message": "No reports found for this company"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode reports"})
		return
	}

	c.JSON(http.StatusOK, reports)

}

// downloadReport godoc
// @Summary Download a report by ID
// @Description Download a report by ID
// @Tags Reports
// @Accept json
// @Produce json
// @Param id path string true "Report ID"
// @Success 200 {object} map[string]interface{} "Report downloaded successfully"
// @Failure 400 {object} map[string]string "Invalid report ID"
// @Failure 500 {object} map[string]string "Failed to download report"
// @Router /report/download/{id} [get]
func DownloadReport(c *gin.Context) {
	reportId := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(reportId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	var report models.Report
	err = config.DB.Collection("reports").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&report)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Report")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Title: "+report.Title)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Description: "+report.Description)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Created By: "+report.CreatedBy.Hex())
	pdf.Ln(10)
	pdf.Cell(40, 10, "Created Date: "+report.CreatedDate.Format(time.RFC3339))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Last Modified Date: "+report.LastModifiedDate.Format(time.RFC3339))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Type: "+report.Type)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Status: "+report.Status)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Content: "+report.Content)
	pdf.Ln(10)

	err = pdf.Output(c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download report"})
		return
	}

	c.Header("Content-Deposition", "attachment; filename=report_"+report.ID.Hex()+".pdf")
	c.Header("Content-Type", "application/pdf")
}
