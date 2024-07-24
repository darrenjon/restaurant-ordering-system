package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/darrenjon/restaurant-ordering-system/internal/database"
	"github.com/darrenjon/restaurant-ordering-system/internal/models"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(db *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate required fields
		if user.Username == "" || user.Password == "" || user.Name == "" || user.Role == "" {
			http.Error(w, "Username, password, name, and role are required", http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)

		result := db.GetDB().Create(&user)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		user.Password = "" // Don't send the password back
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

func GetUsers(db *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []models.User
		result := db.GetDB().Find(&users)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		// Don't send the password back
		for i := range users {
			users[i].Password = ""
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func GetUser(db *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var user models.User
		result := db.GetDB().First(&user, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		user.Password = "" // Don't send the password back
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func UpdateUser(db *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var updatedUser models.User
		if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var user models.User
		result := db.GetDB().First(&user, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		if updatedUser.Username != "" {
			user.Username = updatedUser.Username
		}
		if updatedUser.Name != "" {
			user.Name = updatedUser.Name
		}
		if updatedUser.Role != "" {
			user.Role = updatedUser.Role
		}
		if updatedUser.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Error hashing password", http.StatusInternalServerError)
				return
			}
			user.Password = string(hashedPassword)
		}

		result = db.GetDB().Save(&user)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		user.Password = ""
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func DeleteUser(db *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		result := db.GetDB().Delete(&models.User{}, id)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		if result.RowsAffected == 0 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
	}
}
