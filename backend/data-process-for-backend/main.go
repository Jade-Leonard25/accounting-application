package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"payroll-main-process-data/database"
	"payroll-main-process-data/models"

	"github.com/joho/godotenv"
)

var idRegex = regexp.MustCompile(`^[a-zA-Z0-9]{3}-[a-zA-Z0-9]{3}-[a-zA-Z0-9]{3}$`)

func employeeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	sessionToken := r.Header.Get("X-Session-Token")
	log.Printf("Auth Header: %s | Session Token: %s", authHeader, sessionToken)

	var emp models.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		log.Printf("ERROR decoding JSON: %v", err)
		http.Error(w, `{"error":"Invalid JSON body"}`, http.StatusBadRequest)
		return
	}

	if !idRegex.MatchString(emp.ID) {
		log.Printf("ERROR: Invalid ID format received: %s", emp.ID)
		http.Error(w, `{"error":"Invalid ID format. Expected xxx-xxx-xxx"}`, http.StatusBadRequest)
		return
	}

	if database.DB == nil {
		log.Printf("ERROR: database.DB is nil!")
		http.Error(w, `{"error":"Database not initialized"}`, http.StatusInternalServerError)
		return
	}

	if result := database.DB.Create(&emp); result.Error != nil {
		log.Printf("ERROR saving to NeonDB: %v", result.Error)
		http.Error(w, `{"error":"`+result.Error.Error()+`"}`, http.StatusInternalServerError)
		return
	}



	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"status":   "success",
		"message":  "Employee record created and stored in NeonDB",
		"employee": emp,
	})
}

func main() {
	// Load environment variables (.env) if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// Initialize database connection using ConnectDB()
	database.ConnectDB()

	// Run auto-migration for the Employee model
	if database.DB != nil {
		database.DB.AutoMigrate(&models.Employee{})
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Go backend with Neon & GORM is running successfully!"))
	})
	
	mux.HandleFunc("/api/employees", employeeHandler)

	loggedMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming Request: [%s] %s", r.Method, r.URL.Path)
		mux.ServeHTTP(w, r)
	})

	log.Println("Server running on port 8080...")
	if err := http.ListenAndServe(":8080", loggedMux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}