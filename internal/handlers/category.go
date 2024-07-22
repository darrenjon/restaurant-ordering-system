package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/darrenjon/restaurant-ordering-system/internal/database"
	"github.com/darrenjon/restaurant-ordering-system/internal/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetCategories(db *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var categories []models.Category
		result := db.GetDB().Order("display_order").Find(&categories)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(categories)
	}
}

func CreateCategory(db *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var category models.Category
		if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := db.GetDB().Create(&category)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(category)
	}
}

func UpdateCategory(db *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		var updatedCategory models.Category
		if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var existingCategory models.Category
		result := db.GetDB().First(&existingCategory, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Category not found", http.StatusNotFound)
			} else {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Update only specific fields
		existingCategory.Name = updatedCategory.Name
		existingCategory.DisplayOrder = updatedCategory.DisplayOrder

		result = db.GetDB().Save(&existingCategory)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingCategory)
	}
}

func DeleteCategory(db *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		result := db.GetDB().Delete(&models.Category{}, id)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		if result.RowsAffected == 0 {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted successfully"})
	}
}
