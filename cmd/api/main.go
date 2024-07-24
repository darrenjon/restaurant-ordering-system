package main

import (
	"net/http"

	"github.com/gorilla/mux"
	gormlogger "gorm.io/gorm/logger"

	"github.com/darrenjon/restaurant-ordering-system/internal/config"
	"github.com/darrenjon/restaurant-ordering-system/internal/database"
	"github.com/darrenjon/restaurant-ordering-system/internal/handlers"
	"github.com/darrenjon/restaurant-ordering-system/internal/logger"
)

func main() {
	// Load database configuration
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to load database config: %v", err)
	}

	// Create database manager
	dbManager, err := database.NewManager(dbConfig)
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to create database manager: %v", err)
	}

	// Run auto migrations
	if err := dbManager.AutoMigrate(); err != nil {
		logger.ErrorLogger.Fatalf("Failed to run auto migrations: %v", err)
	}

	// Set log mode to Info after initialization
	dbManager.SetLogMode(gormlogger.Info)

	r := mux.NewRouter()

	r.HandleFunc("/api/auth/login", handlers.Login(dbManager)).Methods("POST")
	r.HandleFunc("/api/auth/logout", handlers.Logout).Methods("POST")

	// User routes
	r.HandleFunc("/api/users", handlers.CreateUser(dbManager)).Methods("POST")
	// test auth middleware
	// r.HandleFunc("/api/users", middleware.AuthMiddleware(dbConfig)(handlers.GetUsers(dbManager))).Methods("GET")
	r.HandleFunc("/api/users", handlers.GetUsers(dbManager)).Methods("GET")
	r.HandleFunc("/api/users/{id}", handlers.GetUser(dbManager)).Methods("GET")
	r.HandleFunc("/api/users/{id}", handlers.UpdateUser(dbManager)).Methods("PUT")
	r.HandleFunc("/api/users/{id}", handlers.DeleteUser(dbManager)).Methods("DELETE")

	// Restaurant info routes
	r.HandleFunc("/api/restaurant-info", handlers.GetRestaurantInfo(dbManager)).Methods("GET")
	r.HandleFunc("/api/restaurant-info", handlers.UpdateRestaurantInfo(dbManager)).Methods("PUT")
	r.HandleFunc("/api/restaurant-info/open", handlers.CheckRestaurantOpen(dbManager)).Methods("GET")

	// Category routes
	r.HandleFunc("/api/categories", handlers.GetCategories(dbManager)).Methods("GET")
	r.HandleFunc("/api/categories", handlers.CreateCategory(dbManager)).Methods("POST")
	r.HandleFunc("/api/categories/{id}", handlers.UpdateCategory(dbManager)).Methods("PUT")
	r.HandleFunc("/api/categories/{id}", handlers.DeleteCategory(dbManager)).Methods("DELETE")

	// Menu item routes
	r.HandleFunc("/api/menu-items", handlers.GetMenuItems(dbManager)).Methods("GET")
	r.HandleFunc("/api/menu-items/{id}", handlers.GetMenuItem(dbManager)).Methods("GET")
	r.HandleFunc("/api/menu-items", handlers.CreateMenuItem(dbManager)).Methods("POST")
	r.HandleFunc("/api/menu-items/{id}", handlers.UpdateMenuItem(dbManager)).Methods("PUT")
	r.HandleFunc("/api/menu-items/{id}", handlers.DeleteMenuItem(dbManager)).Methods("DELETE")

	// Add a simple health check route
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	logger.InfoLogger.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.ErrorLogger.Fatalf("Error starting server: %v", err)
	}
}
