package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"payroll-main-process-data/database"
	"payroll-main-process-data/models"
)

// AdminHandler manages admin CRUD operations
func AdminHandler(w http.ResponseWriter, r *http.Request) {
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
	var adminID string
	if len(parts) > 2 {
		adminID = parts[2]
	}

	switch r.Method {
	case http.MethodGet:
		if adminID != "" {
			var admin models.Admin
			// Preload companies managed by this admin
			if result := database.DB.Preload("Companies").First(&admin, "id = ?", adminID); result.Error != nil {
				http.Error(w, `{"error":"Admin not found"}`, http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(admin)
		} else {
			var admins []models.Admin
			if result := database.DB.Preload("Companies").Find(&admins); result.Error != nil {
				http.Error(w, `{"error":"`+result.Error.Error()+`"}`, http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(admins)
		}

	case http.MethodPost:
		var admin models.Admin
		if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
			http.Error(w, `{"error":"Invalid JSON body"}`, http.StatusBadRequest)
			return
		}

		// In production, make sure to hash PasswordHash before saving!
		if result := database.DB.Create(&admin); result.Error != nil {
			http.Error(w, `{"error":"`+result.Error.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  "success",
			"message": "Admin created successfully",
			"admin":   admin,
		})

	case http.MethodPut:
		if adminID == "" {
			http.Error(w, `{"error":"Admin ID required for update"}`, http.StatusBadRequest)
			return
		}

		var updateData models.Admin
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			http.Error(w, `{"error":"Invalid JSON body"}`, http.StatusBadRequest)
			return
		}

		var admin models.Admin
		if result := database.DB.First(&admin, "id = ?", adminID); result.Error != nil {
			http.Error(w, `{"error":"Admin not found"}`, http.StatusNotFound)
			return
		}

		database.DB.Model(&admin).Updates(updateData)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  "success",
			"message": "Admin updated successfully",
			"admin":   admin,
		})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"Method not allowed"}`))
	}
}