package main

import (
	"log"
	"net/http"

	"payroll-main-process-data/controllers"
	"payroll-main-process-data/database"
	"payroll-main-process-data/models"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	database.ConnectDB()

	if database.DB != nil {
		// Passing models in order ensures foreign key constraints are created properly
		err := database.DB.AutoMigrate(
			&models.Admin{},
			&models.Company{},
			&models.Employee{},
		)
		if err != nil {
			log.Fatalf("Failed to auto-migrate database: %v", err)
		}
	}

	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Go backend with Neon & GORM is running successfully!"))
	})
	
	// Register Admin routes
	mux.HandleFunc("/api/admins", controllers.AdminHandler)
	mux.HandleFunc("/api/admins/", controllers.AdminHandler)

	// Register Company routes
	mux.HandleFunc("/api/companies", controllers.CompanyHandler)
	mux.HandleFunc("/api/companies/", controllers.CompanyHandler)

	// Register Employee routes
	mux.HandleFunc("/api/employees", controllers.EmployeeHandler)
	mux.HandleFunc("/api/employees/", controllers.EmployeeHandler)

	loggedMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming Request: [%s] %s", r.Method, r.URL.Path)
		mux.ServeHTTP(w, r)
	})

	log.Println("Server running on port 8080...")
	if err := http.ListenAndServe(":8080", loggedMux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}