package handlers

import (
	"fmt"
	"net/http"
	"project/db" // Make sure to import the db package or correct the path
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

// Helper function to parse dates in dd/mm/yy format
func parseDate(dateStr string) (time.Time, error) {
	// Parse date in "dd/mm/yy" format
	return time.Parse("02/01/06", dateStr)
}

// GetUsersBetweenDates handles fetching users between two dates
func GetUsersBetweenDates(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both start_date and end_date are required"})
		return
	}

	// Parse start and end dates
	start, err := parseDate(startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid start_date format: %v", err)})
		return
	}

	end, err := parseDate(endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid end_date format: %v", err)})
		return
	}

	// Fetch users by date range
	users, err := FetchUsersByDateRange(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// FetchUsersByDateRange retrieves users whose date falls between the specified start and end dates
func FetchUsersByDateRange(startDate, endDate time.Time) ([]map[string]interface{}, error) {
	query := `SELECT id, name, email, registration_no, phone_no, date FROM users WHERE STR_TO_DATE(date, '%d/%m/%y') BETWEEN ? AND ?`
	rows, err := db.DB.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %v", err)
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var name, email, registrationNo, phoneNo, date string
		if err := rows.Scan(&id, &name, &email, &registrationNo, &phoneNo, &date); err != nil {
			return nil, fmt.Errorf("error scanning user data: %v", err)
		}

		user := map[string]interface{}{
			"id":              id,
			"name":            name,
			"email":           email,
			"registration_no": registrationNo,
			"phone_no":        phoneNo,
			"date":            date,
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return users, nil
}

// ExportUsersToExcel exports the users' data to an Excel file
func ExportUsersToExcel(c *gin.Context) {
	// Get the start and end dates
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both start_date and end_date are required"})
		return
	}

	// Parse start and end dates
	start, err := parseDate(startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid start_date format: %v", err)})
		return
	}

	end, err := parseDate(endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid end_date format: %v", err)})
		return
	}

	// Fetch users by date range
	users, err := FetchUsersByDateRange(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}

	// Create a new Excel file
	f := excelize.NewFile()

	// Create a new sheet
	sheetIndex, err := f.NewSheet("Users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error creating new sheet: %v", err)})
		return
	}

	// Set the active sheet to the newly created one
	f.SetActiveSheet(sheetIndex)

	// Set headers
	headers := []string{"ID", "Name", "Email", "Registration No", "Phone No", "Date"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Users", cell, header)
	}

	// Add user data to Excel
	for i, user := range users {
		f.SetCellValue("Users", fmt.Sprintf("A%d", i+2), user["id"])
		f.SetCellValue("Users", fmt.Sprintf("B%d", i+2), user["name"])
		f.SetCellValue("Users", fmt.Sprintf("C%d", i+2), user["email"])
		f.SetCellValue("Users", fmt.Sprintf("D%d", i+2), user["registration_no"])
		f.SetCellValue("Users", fmt.Sprintf("E%d", i+2), user["phone_no"])
		f.SetCellValue("Users", fmt.Sprintf("F%d", i+2), user["date"])
	}

	// Save to file and send the file
	fileName := "UsersData.xlsx"
	if err := f.SaveAs(fileName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving Excel file"})
		return
	}

	c.File(fileName)
}

// ExportUsersToPDF exports the users' data to a PDF file
// ExportUsersToPDF handles exporting the users' data to a PDF file
func ExportUsersToPDF(c *gin.Context) {
	// Get the start and end dates
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both start_date and end_date are required"})
		return
	}

	// Parse start and end dates in "dd/mm/yy" format
	start, err := time.Parse("02/01/06", startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid start_date format: %v", err)})
		return
	}

	end, err := time.Parse("02/01/06", endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid end_date format: %v", err)})
		return
	}

	// Fetch users by date range
	users, err := FetchUsersByDateRange(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}

	// Create a new PDF document
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	// Add title
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Users Data")

	// Add table headers
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(20, 10, "ID", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Name", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 10, "Email", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 10, "Registration No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Phone No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 10, "Date", "1", 0, "C", false, 0, "")
	pdf.Ln(-1) // Move to the next line

	// Add user data
	pdf.SetFont("Arial", "", 12)
	for _, user := range users {
		pdf.CellFormat(20, 10, fmt.Sprintf("%d", user["id"]), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%s", user["name"]), "1", 0, "L", false, 0, "")
		pdf.CellFormat(60, 10, fmt.Sprintf("%s", user["email"]), "1", 0, "L", false, 0, "")
		pdf.CellFormat(50, 10, fmt.Sprintf("%s", user["registration_no"]), "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%s", user["phone_no"]), "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%s", user["date"]), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	// Output the PDF directly to the response
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=UsersData.pdf")
	err = pdf.Output(c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating PDF"})
		return
	}
}
