package main

import (
	"log"
	"net/http"

	"github.com/darrenjon/restaurant-ordering-system/internal/config"
	"github.com/darrenjon/restaurant-ordering-system/internal/database"
	"github.com/darrenjon/restaurant-ordering-system/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Load database configuration
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	// Create database manager
	dbManager, err := database.NewManager(dbConfig)
	if err != nil {
		log.Fatalf("Failed to create database manager: %v", err)
	}

	// Run auto migrations
	if err := dbManager.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run auto migrations: %v", err)
	}

	r := mux.NewRouter()

	// Restaurant info routes
	r.HandleFunc("/api/restaurant-info", handlers.GetRestaurantInfo(dbManager)).Methods("GET")
	r.HandleFunc("/api/restaurant-info", handlers.UpdateRestaurantInfo(dbManager)).Methods("PUT")
	r.HandleFunc("/api/restaurant-info/open", handlers.CheckRestaurantOpen(dbManager)).Methods("GET")

	// Add a simple health check route
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
