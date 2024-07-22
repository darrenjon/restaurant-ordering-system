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

func GetMenuItems(db *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var menuItems []models.MenuItem
		result := db.GetDB().Preload("AddOns").Find(&menuItems)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(menuItems)
	}
}

func GetMenuItem(db *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
			return
		}

		var menuItem models.MenuItem
		result := db.GetDB().Preload("AddOns").First(&menuItem, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Menu item not found", http.StatusNotFound)
			} else {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(menuItem)
	}
}

func CreateMenuItem(db *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var menuItem models.MenuItem
		if err := json.NewDecoder(r.Body).Decode(&menuItem); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := db.GetDB().Create(&menuItem)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(menuItem)
	}
}

func UpdateMenuItem(db *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
			return
		}

		var updatedMenuItem models.MenuItem
		if err := json.NewDecoder(r.Body).Decode(&updatedMenuItem); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Ensure the ID in the URL matches the ID in the request body
		tx := db.GetDB().Begin()
		// get the existing menu item
		var existingMenuItem models.MenuItem
		result := db.GetDB().First(&existingMenuItem, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Menu item not found", http.StatusNotFound)
			} else {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Update fields
		existingMenuItem.Name = updatedMenuItem.Name
		existingMenuItem.Description = updatedMenuItem.Description
		existingMenuItem.Price = updatedMenuItem.Price
		existingMenuItem.ImageURL = updatedMenuItem.ImageURL
		existingMenuItem.IsAvailable = updatedMenuItem.IsAvailable
		existingMenuItem.CategoryID = updatedMenuItem.CategoryID
		// Delete existing add-ons
		if err := tx.Where("menu_item_id = ?", existingMenuItem.ID).Delete(&models.AddOn{}).Error; err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create new add-ons
		for i := range updatedMenuItem.AddOns {
			updatedMenuItem.AddOns[i].MenuItemID = existingMenuItem.ID
		}
		if err := tx.Create(&updatedMenuItem.AddOns).Error; err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Save the updated menu item
		if err := tx.Save(&existingMenuItem).Error; err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Preload the add-ons
		if err := db.GetDB().Preload("AddOns").First(&existingMenuItem, existingMenuItem.ID).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingMenuItem)
	}
}

func DeleteMenuItem(db *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
			return
		}
		// Start a transaction
		tx := db.GetDB().Begin()
		// Delete the add-ons first
		if err := tx.Where("menu_item_id = ?", id).Delete(&models.AddOn{}).Error; err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Delete the menu item
		result := tx.Delete(&models.MenuItem{}, id)
		if result.Error != nil {
			tx.Rollback()
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		// Check if the menu item was found
		if result.RowsAffected == 0 {
			tx.Rollback()
			http.Error(w, "Menu item not found", http.StatusNotFound)
			return
		}
		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Menu item deleted successfully"})
	}
}
