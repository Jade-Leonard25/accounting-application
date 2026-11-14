package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"payroll-main-process-data/database"
	"payroll-main-process-data/models"
)

// CompanyHandler manages company CRUD operations
func CompanyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if database.DB == nil {
		http.Error(w, `{"error":"Database not initialized"}`, http.StatusInternalServerError)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	var companyID string
	if len(parts) > 2 {
		companyID = parts[2]
	}

	switch r.Method {
	case http.MethodGet:
		if companyID != "" {
			var company models.Company
			// Preload Admin, CEO, and Employees for complete context
			if result := database.DB.Preload("Admin").Preload("CEO").Preload("Employees").First(&company, "id = ?", companyID); result.Error != nil {
				http.Error(w, `{"error":"Company not found"}`, http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(company)
		} else {
			var companies []models.Company
			if result := database.DB.Preload("Admin").Preload("CEO").Preload("Employees").Find(&companies); result.Error != nil {
				http.Error(w, `{"error":"`+result.Error.Error()+`"}`, http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(companies)
		}

	case http.MethodPost:
		var company models.Company
		if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
			http.Error(w, `{"error":"Invalid JSON body"}`, http.StatusBadRequest)
			return
		}

		if result := database.DB.Create(&company); result.Error != nil {
			http.Error(w, `{"error":"`+result.Error.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  "success",
			"message": "Company created successfully",
			"company": company,
		})

	case http.MethodPut:
		if companyID == "" {
			http.Error(w, `{"error":"Company ID required for update"}`, http.StatusBadRequest)
			return
		}

		var updateData models.Company
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			http.Error(w, `{"error":"Invalid JSON body"}`, http.StatusBadRequest)
			return
		}

		var company models.Company
		if result := database.DB.First(&company, "id = ?", companyID); result.Error != nil {
			http.Error(w, `{"error":"Company not found"}`, http.StatusNotFound)
			return
		}

		database.DB.Model(&company).Updates(updateData)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  "success",
			"message": "Company updated successfully",
			"company": company,
		})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"Method not allowed"}`))
	}
}