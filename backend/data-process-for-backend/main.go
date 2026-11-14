package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"payroll-main-process-data/database"
	"payroll-main-process-data/models"

	"github.com/joho/godotenv"
)

var idRegex = regexp.MustCompile(`^[a-zA-Z0-9]{3}-[a-zA-Z0-9]{3}-[a-zA-Z0-9]{3}$`)

func employeeHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for your Vue frontend
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if database.DB == nil {
		log.Printf("ERROR: database.DB is nil!")
		http.Error(w, `{"error":"Database not initialized"}`, http.StatusInternalServerError)
		return
	}

	// Parse path segments to check if an ID is present (e.g., /api/employees/98b-01p-37x)
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	var empID string
	if len(parts) > 2 {
		empID = parts[2]
	}

	switch r.Method {
	case http.MethodGet:
		if empID != "" {
			// Get single employee
			var emp models.Employee
			if result := database.DB.First(&emp, "id = ?", empID); result.Error != nil {
				http.Error(w, `{"error":"Employee not found"}`, http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(emp)
		} else {
			// List all employees
			var employees []models.Employee
			if result := database.DB.Find(&employees); result.Error != nil {
				http.Error(w, `{"error":"`+result.Error.Error()+`"}`, http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(employees)
		}

	case http.MethodPost:
		var emp models.Employee
		if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
			http.Error(w, `{"error":"Invalid JSON body"}`, http.StatusBadRequest)
			return
		}

		if !idRegex.MatchString(emp.ID) {
			http.Error(w, `{"error":"Invalid ID format. Expected xxx-xxx-xxx"}`, http.StatusBadRequest)
			return
		}

		if result := database.DB.Create(&emp); result.Error != nil {
			http.Error(w, `{"error":"`+result.Error.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"status":   "success",
			"message":  "Employee record created successfully",
			"employee": emp,
		})

	case http.MethodPut:
		if empID == "" {
			http.Error(w, `{"error":"Employee ID required for update"}`, http.StatusBadRequest)
			return
		}

		var updateData models.Employee
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			http.Error(w, `{"error":"Invalid JSON body"}`, http.StatusBadRequest)
			return
		}

		var emp models.Employee
		if result := database.DB.First(&emp, "id = ?", empID); result.Error != nil {
			http.Error(w, `{"error":"Employee not found"}`, http.StatusNotFound)
			return
		}

		// Update fields
		database.DB.Model(&emp).Updates(updateData)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"status":   "success",
			"message":  "Employee record updated successfully",
			"employee": emp,
		})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"Method not allowed"}`))
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	database.ConnectDB()

	if database.DB != nil {
		database.DB.AutoMigrate(&models.Employee{})
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Go backend with Neon & GORM is running successfully!"))
	})
	
	// Register both base and trailing slash so all subpaths route properly to employeeHandler
	mux.HandleFunc("/api/employees", employeeHandler)
	mux.HandleFunc("/api/employees/", employeeHandler)

	loggedMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming Request: [%s] %s", r.Method, r.URL.Path)
		mux.ServeHTTP(w, r)
	})

	log.Println("Server running on port 8080...")
	if err := http.ListenAndServe(":8080", loggedMux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
