package server

import (
	"html/template"
	"log"
	"rm_client_portal/database"
	"shared_templates"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ReportsListHandler - retrieve a list of reports for the authenticated client
func ReportsListHandler(c *gin.Context) {
	email := getEmailFromJWT(c)
	clients := database.GetClientsForUserEmail(email)

	if len(clients) == 0 {
		c.JSON(404, gin.H{
			"error": "No clients found for user",
		})
		return
	}

	// Check if a specific client_id is requested via query parameter
	requestedClientIDStr := c.Query("client_id")
	var clientID int

	if requestedClientIDStr != "" {
		// Validate that the user has access to the requested client
		requestedClientID, err := strconv.Atoi(requestedClientIDStr)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid client_id parameter",
			})
			return
		}

		// Check if user has access to this client
		hasAccess := false
		for _, client := range clients {
			if int(client.ID) == requestedClientID {
				hasAccess = true
				clientID = requestedClientID
				break
			}
		}

		if !hasAccess {
			c.JSON(403, gin.H{
				"error": "Access denied to requested client",
			})
			return
		}
	} else {
		// Default to first client if no specific client requested
		clientID = int(clients[0].ID)
	}

	reports, err := database.GetClientReportsByClientID(clientID)
	if err != nil {
		log.Printf("Error retrieving reports for client %d: %v", clientID, err)
		c.JSON(500, gin.H{
			"error": "Failed to retrieve reports",
		})
		return
	}

	// Return simplified report list (without full location data)
	var reportsList []gin.H
	for _, report := range reports {
		reportsList = append(reportsList, gin.H{
			"report_id":    report.ReportID,
			"period_start": report.PeriodStart.Format("2006-01-02"),
			"period_end":   report.PeriodEnd.Format("2006-01-02"),
			"generated_at": report.GeneratedAt.Format("2006-01-02 15:04:05"),
			"client_name":  report.ClientName,
		})
	}

	c.JSON(200, gin.H{
		"reports": reportsList,
	})
}

// ReportHTMLHandler - return HTML for a specific report
func ReportHTMLHandler(c *gin.Context) {

	// Get report ID from URL parameter
	reportIDStr := c.Param("id")
	log.Printf("Report ID parameter: %s", reportIDStr)

	reportID, err := strconv.ParseInt(reportIDStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing report ID: %v", err)
		c.JSON(400, gin.H{
			"error": "Invalid report ID",
		})
		return
	}

	// Get client ID from JWT
	email := getEmailFromJWT(c)
	clients := database.GetClientsForUserEmail(email)

	if len(clients) == 0 {
		log.Printf("No clients found for user %s", email)
		c.JSON(404, gin.H{
			"error": "No clients found for user",
		})
		return
	}

	// Try to find the report across all user's clients
	var report *shared_templates.ClientReportData

	for _, client := range clients {
		clientID := int(client.ID)
		foundReport, lookupErr := database.GetClientReportByIDForClient(reportID, clientID)
		if lookupErr == nil && foundReport != nil {
			// Found the report for this client
			report = foundReport
			break
		}
	}

	if report == nil {
		log.Printf("Report %d not found for user %s", reportID, email)
		c.JSON(404, gin.H{
			"error": "Report not found",
		})
		return
	}

	// Parse and execute the HTML template
	tmpl, err := template.New("report").Parse(shared_templates.MonthlyReportTemplate)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		c.JSON(500, gin.H{
			"error": "Failed to parse template",
		})
		return
	}

	// Set content type to HTML and CSP headers to allow inline scripts
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Header("Content-Security-Policy", "script-src 'unsafe-inline' 'unsafe-eval'; object-src 'none'; base-uri 'none';")
	c.Header("X-Frame-Options", "SAMEORIGIN")

	// Execute the template with the report data
	err = tmpl.Execute(c.Writer, report)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		c.JSON(500, gin.H{
			"error": "Failed to render report",
		})
		return
	}
}
