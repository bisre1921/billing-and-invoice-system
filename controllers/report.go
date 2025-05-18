package controllers

import (
	"context"
	"fmt"
	"net/http"
	// "strconv"
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
// @Description Generate a new report based on the provided type, title, and description. The content is auto-generated.
// @Tags Reports
// @Accept json
// @Produce json
// @Param report body models.GenerateReportRequest true "Report generation data"
// @Success 200 {object} map[string]interface{} "Report generated successfully"
// @Failure 400 {object} map[string]string "Invalid report input"
// @Failure 500 {object} map[string]string "Failed to generate report"
// @Router /report/generate [post]
func GenerateReport(c *gin.Context) {
	var req models.GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report input", "details": err.Error()})
		return
	}

	companyID := req.CompanyID
	reportType := req.Type

	var content string
	var err error

	switch reportType {
    case "Financial":
        fmt.Println("Generating Financial Report Content") 
        content, err = generateFinancialReportContent(companyID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate financial report content", "details": err.Error()})
            return
        }
    case "Inventory":
        fmt.Println("Generating Inventory Report Content") 
        content, err = generateInventoryReportContent(companyID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate inventory report content", "details": err.Error()})
            return
        }
    case "Customer Analysis":
        fmt.Println("Generating Customer Analysis Report Content") 
        content, err = generateCustomerAnalysisReportContent(companyID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate customer analysis report content", "details": err.Error()})
            return
        }
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report type"})
        return
    }

	report := models.Report{
		CompanyID:        companyID,
		Title:            req.Title,
		Description:      req.Description,
		CreatedBy:        req.CreatedBy,
		CreatedDate:      time.Now(),
		LastModifiedDate: time.Now(),
		Type:             reportType,
		Status:           "Generated",
		Content:          content,
	}

	res, err := config.DB.Collection("reports").InsertOne(context.Background(), report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save report to database"})
		return
	}

	report.ID = res.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusOK, gin.H{"message": "Report generated successfully", "report": report})
}

func generateFinancialReportContent(companyID string) (string, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	var company models.Company
	companyObjID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return "", fmt.Errorf("invalid company ID format: %w", err)
	}
	err = config.DB.Collection("companies").FindOne(context.Background(), bson.M{"_id": companyObjID}).Decode(&company)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve company: %w", err)
	}
	companyName := company.Name

	paidRevenuePipeline := []bson.M{
		{"$match": bson.M{
			"company_id": companyID,
			"status":     "Paid",
			"date": bson.M{
				"$gte": startOfMonth,
				"$lte": endOfMonth,
			},
		}},
		{"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
		}},
	}

	unpaidRevenuePipeline := []bson.M{
		{"$match": bson.M{
			"company_id": companyID,
			"status":     "Unpaid",
			"date": bson.M{
				"$gte": startOfMonth,
				"$lte": endOfMonth,
			},
		}},
		{"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
		}},
	}

	paidCursor, err := config.DB.Collection("invoices").Aggregate(context.Background(), paidRevenuePipeline)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve paid monthly revenue: %w", err)
	}
	defer paidCursor.Close(context.Background())
	var paidRevenue float64
	if paidCursor.Next(context.Background()) {
		var row bson.M
		if err := paidCursor.Decode(&row); err != nil {
			return "", fmt.Errorf("failed to decode paid monthly revenue: %w", err)
		}
		if total, ok := row["total"].(float64); ok {
			paidRevenue = total
		}
	}

	unpaidCursor, err := config.DB.Collection("invoices").Aggregate(context.Background(), unpaidRevenuePipeline)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve unpaid monthly revenue: %w", err)
	}
	defer unpaidCursor.Close(context.Background())
	var unpaidRevenue float64
	if unpaidCursor.Next(context.Background()) {
		var row bson.M
		if err := unpaidCursor.Decode(&row); err != nil {
			return "", fmt.Errorf("failed to decode unpaid monthly revenue: %w", err)
		}
		if total, ok := row["total"].(float64); ok {
			unpaidRevenue = total
		}
	}

	totalRevenue := paidRevenue + unpaidRevenue

	content := fmt.Sprintf(`
		Financial Report for %s - %s
		-----------------------------------
		Total Revenue: $%.2f
		Paid Revenue:  $%.2f
		Unpaid Revenue: $%.2f
	`, companyName, now.Format("January 2006"), totalRevenue, paidRevenue, unpaidRevenue)

	return content, nil
}

func generateInventoryReportContent(companyID string) (string, error) {
	
	var company models.Company
	companyObjID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return "", fmt.Errorf("invalid company ID format: %w", err)
	}
	err = config.DB.Collection("companies").FindOne(context.Background(), bson.M{"_id": companyObjID}).Decode(&company)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve company: %w", err)
	}
	companyName := company.Name
	fmt.Println("companyID: ", companyID)

	var items []models.Item
	itemCursor, err := config.DB.Collection("items").Find(context.Background(), bson.M{"company_id": companyID})
	fmt.Println("itemCursor: ", itemCursor)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve items cursor: %w", err)
	}
	defer itemCursor.Close(context.Background())
	if err := itemCursor.All(context.Background(), &items); err != nil {
		return "", fmt.Errorf("failed to decode items: %w", err)
	}

	
	fmt.Println("items: ", items)
	var employees []models.Employee
	employeeCursor, err := config.DB.Collection("employees").Find(context.Background(), bson.M{"company_id": companyID})
	fmt.Println("employeeCursor: ", employeeCursor)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve employees cursor: %w", err)
	}
	defer employeeCursor.Close(context.Background())
	if err := employeeCursor.All(context.Background(), &employees); err != nil {
		return "", fmt.Errorf("failed to decode employees: %w", err)
	}
	fmt.Println("employees: ", employees)


	var customers []models.Customer
	customerCursor, err := config.DB.Collection("customers").Find(context.Background(), bson.M{"company_id": companyID})
	fmt.Println("customerCursor: ", customerCursor)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve customers cursor: %w", err)
	}
	defer customerCursor.Close(context.Background())
	if err := customerCursor.All(context.Background(), &customers); err != nil {
		return "", fmt.Errorf("failed to decode customers: %w", err)
	}
	fmt.Println("customers: ", customers)


	content := fmt.Sprintf(`
		Inventory Report for %s - %s
		-----------------------------------
		Number of Items: %d
		Number of Employees: %d
		Number of Customers: %d

		--- Items ---
	`, companyName, time.Now().Format("January 2006"), len(items), len(employees), len(customers))

	for _, item := range items {
		content += fmt.Sprintf("- %s (ID: %s)\n", item.Name, item.ID.Hex())
	}

	content += "\n--- Employees ---\n"
	for _, employee := range employees {
		content += fmt.Sprintf("- %s (Email: %s)\n", employee.Name, employee.Email)
	}

	content += "\n--- Customers ---\n"
	for _, customer := range customers {
		content += fmt.Sprintf("- %s (Email: %s)\n", customer.Name, customer.Email)
	}

	return content, nil
}

func generateCustomerAnalysisReportContent(companyID string) (string, error) {
	
	var company models.Company
	companyObjID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return "", fmt.Errorf("invalid company ID format: %w", err)
	}
	err = config.DB.Collection("companies").FindOne(context.Background(), bson.M{"_id": companyObjID}).Decode(&company)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve company: %w", err)
	}
	companyName := company.Name

	var customers []models.Customer
	customerCursor, err := config.DB.Collection("customers").Find(context.Background(), bson.M{"company_id": companyID})
	if err != nil {
		return "", fmt.Errorf("failed to retrieve customers cursor: %w", err)
	}
	defer customerCursor.Close(context.Background())
	if err := customerCursor.All(context.Background(), &customers); err != nil {
		return "", fmt.Errorf("failed to decode customers: %w", err)
	}
	fmt.Println("customers: ", customers)
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	customerData := make(map[string]float64)

	for _, customer := range customers {
		var customerInvoices []models.Invoice
		invoiceCursor, err := config.DB.Collection("invoices").Find(context.Background(), bson.M{"company_id": companyID, "customer_id": customer.ID.Hex(), "status": "Paid", "date": bson.M{"$gte": startOfMonth, "$lte": endOfMonth}})
		if err != nil {
			return "", fmt.Errorf("failed to retrieve invoices for customer %s cursor: %w", customer.Name, err)
		}
		defer invoiceCursor.Close(context.Background())
		var totalSpent float64 = 0
		if err := invoiceCursor.All(context.Background(), &customerInvoices); err != nil {
			return "", fmt.Errorf("failed to decode invoices for customer %s: %w", customer.Name, err)
		}
		for _, invoice := range customerInvoices {
			totalSpent += invoice.Amount
		}
		customerData[customer.Name] = totalSpent
	}

	content := fmt.Sprintf(`
		Customer Analysis Report for %s - %s
		---------------------------------------
		Total Number of Customers: %d

		--- Customer Spending (Current Month) ---
	`, companyName, time.Now().Format("January 2006"), len(customers))

	for name, amount := range customerData {
		content += fmt.Sprintf("- %s: $%.2f\n", name, amount)
	}

	return content, nil
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	var report models.Report
	err = config.DB.Collection("reports").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"report": report})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve reports"})
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
	pdf.Ln(6)
	pdf.Cell(40, 10, "Description: "+report.Description)
	pdf.Ln(6)
	pdf.Cell(40, 10, "Created By: "+report.CreatedBy)
	pdf.Ln(6)
	pdf.Cell(40, 10, "Created Date: "+report.CreatedDate.Format(time.RFC3339))
	pdf.Ln(6)
	pdf.Cell(40, 10, "Last Modified Date: "+report.LastModifiedDate.Format(time.RFC3339))
	pdf.Ln(6)
	pdf.Cell(40, 10, "Type: "+report.Type)
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	const lineHeight = 5
	const margin = 10
	const pageWidth float64 = 210 - 2*margin
	y := pdf.GetY()

	lines := splitTextToLines(report.Content, pageWidth, pdf)
	for _, line := range lines {
		if y > 280 { 
			pdf.AddPage()
			y = margin
		}
		pdf.SetY(y)
		pdf.SetX(margin)
		pdf.Cell(0, lineHeight, line)
		y += lineHeight
	}

	err = pdf.Output(c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download report"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=report_"+report.ID.Hex()+".pdf")
	c.Header("Content-Type", "application/pdf")
}

func splitTextToLines(text string, width float64, pdf *gofpdf.Fpdf) []string {
	var lines []string
	currentLine := ""
	words := []rune(text)

	for _, word := range words {
		tempLine := currentLine + string(word)
		lineWidth := pdf.GetStringWidth(tempLine)
		if lineWidth < width {
			currentLine = tempLine
		} else {
			lines = append(lines, currentLine)
			currentLine = string(word)
		}
	}
	lines = append(lines, currentLine)
	return lines
}