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

	// Restaurant info routes
	r.HandleFunc("/api/restaurant-info", handlers.GetRestaurantInfo(dbManager)).Methods("GET")
	r.HandleFunc("/api/restaurant-info", handlers.UpdateRestaurantInfo(dbManager)).Methods("PUT")
	r.HandleFunc("/api/restaurant-info/open", handlers.CheckRestaurantOpen(dbManager)).Methods("GET")

	// Category routes
	r.HandleFunc("/api/categories", handlers.GetCategories(dbManager)).Methods("GET")
	r.HandleFunc("/api/categories", handlers.CreateCategory(dbManager)).Methods("POST")
	r.HandleFunc("/api/categories/{id}", handlers.UpdateCategory(dbManager)).Methods("PUT")
	r.HandleFunc("/api/categories/{id}", handlers.DeleteCategory(dbManager)).Methods("DELETE")

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
