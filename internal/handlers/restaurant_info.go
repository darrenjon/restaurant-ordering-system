package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/darrenjon/restaurant-ordering-system/internal/database"
	"github.com/darrenjon/restaurant-ordering-system/internal/models"
	"gorm.io/gorm"
)

func GetRestaurantInfo(db *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var info models.RestaurantInfo
		result := db.GetDB().Order("updated_at desc").First(&info)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Restaurant info not found", http.StatusNotFound)
			} else {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		json.NewEncoder(w).Encode(info)
	}
}

func UpdateRestaurantInfo(db *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var info models.RestaurantInfo
		if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// First, try to get the existing restaurant info
		var existingInfo models.RestaurantInfo
		result := db.GetDB().First(&existingInfo)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// If no record exists, create a new one
				result = db.GetDB().Create(&info)
			} else {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			// If a record exists, update it
			info.ID = existingInfo.ID // Ensure we're updating the existing record
			info.CreatedAt = existingInfo.CreatedAt
			result = db.GetDB().Save(&info)
		}

		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Restaurant info updated successfully"})
	}
}

func CheckRestaurantOpen(db *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var info models.RestaurantInfo
		result := db.GetDB().First(&info)
		if result.Error != nil {
			http.Error(w, "Restaurant info not found", http.StatusNotFound)
			return
		}

		isOpen := info.OpeningHours.IsOpen(time.Now())
		json.NewEncoder(w).Encode(map[string]bool{"isOpen": isOpen})
	}
}
