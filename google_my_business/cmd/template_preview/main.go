package main

import (
	"flag"
	"fmt"
	"google_my_business/config"
	"google_my_business/database"
	"google_my_business/google_my_business_api"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Parse command line arguments
	port := flag.Int("port", 8080, "Port to run the template preview server on")
	flag.Parse()

	// Print intro message
	fmt.Printf("Starting monthly report template preview server\n")
	fmt.Printf("Use this to preview and iterate on the HTML template design\n")
	fmt.Printf("Server will connect to database to load real client data\n")
	fmt.Printf("------------------------------------------------------\n")

	// Load system configuration
	cfg := config.ReadProperties()

	// Initialize database connection
	fmt.Printf("Connecting to database: %s@%s:%s/%s\n", cfg.DbUsername, cfg.DbAddress, cfg.DbPort, cfg.DbName)
	db := database.OpenDB(cfg.DbName, cfg.DbAddress, cfg.DbPort, cfg.DbUsername, cfg.DbPassword)
	if db == nil {
		log.Fatal("Failed to connect to database")
	}
	defer db.Close()

	// Test database connection
	err := db.Ping()
	if err != nil {
		log.Fatalf("Database connection test failed: %v", err)
	}

	fmt.Printf("Database connection successful!\n")
	fmt.Printf("------------------------------------------------------\n")

	// Start the preview server with database connection
	google_my_business_api.StartReportPreviewServer(*port, db)
}
