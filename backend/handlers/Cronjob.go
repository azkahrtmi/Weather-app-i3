package handlers

import (
	"log"
	"time"

	"backend/config"
	"backend/models"
	"backend/services"
)


func CheckAndUpdateWeather() {
	db := config.GetDB() // Gunakan DB client dari config.go

	var locations []models.Location

	err := db.Select(&locations, "SELECT * FROM locations")
	if err != nil {
		log.Printf("Error fetching locations: %v", err)
		return
	}

	now := time.Now()

	for _, loc := range locations {
		if !loc.UpdateIntervalHours.Valid {
			continue 
		}

		interval := time.Duration(loc.UpdateIntervalHours.Int64) * time.Hour

		if loc.UpdatedAt.IsZero() || now.Sub(loc.UpdatedAt) >= interval {
			go UpdateLocationWeather(loc) // Jalankan secara asynchronous
		}
	}
}

func UpdateLocationWeather(loc models.Location) {
	log.Printf("Updating weather for location %s...", loc.Name)

	err := services.FetchAndSaveWeather(loc)
	if err != nil {
		log.Printf("Failed to update weather for %s: %v", loc.Name, err)
	}
}
