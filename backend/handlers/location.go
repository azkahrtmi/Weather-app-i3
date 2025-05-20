package handlers

import (
	"backend/config"
	"backend/models"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SaveLocation(c *gin.Context) {
	var loc models.Location
	if err := c.ShouldBindJSON(&loc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var locationID int
	insertQuery := `INSERT INTO locations (name, latitude, longitude) VALUES ($1, $2, $3) RETURNING id`
	err := config.DB.QueryRow(insertQuery, loc.Name, loc.Latitude, loc.Longitude).Scan(&locationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert location"})
		return
	}

	// Ambil data cuaca saat ini
	weather, err := services.FetchWeatherData(loc.Latitude, loc.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch current weather", "detail": err.Error()})
		return
	}

	updateQuery := `
        UPDATE locations
        SET weather_summary = $1,
            temperature = $2,
            wind_speed = $3,
            updated_at = NOW()
        WHERE id = $4`
	_, err = config.DB.Exec(updateQuery,
		weather.Current.Summary,
		weather.Current.Temperature,
		weather.Current.Wind.Speed,
		locationID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location"})
		return
	}

	// Ambil data prediksi 1 hari ke depan
	forecastDay, forecastSummary, err := services.FetchOneDayForecast(loc.Latitude, loc.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch forecast", "detail": err.Error()})
		return
	}

	// Simpan ke tabel predictions
	insertPrediction := `INSERT INTO predictions (location_id, date, summary) VALUES ($1, $2, $3)`
	_, err = config.DB.Exec(insertPrediction, locationID, forecastDay, forecastSummary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert prediction", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Success to Add new Location and Forecast",
	})
}

func GetLocations(c *gin.Context) {
	var locations []models.Location

	err := config.DB.Select(&locations, "SELECT * FROM locations ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}

	c.JSON(http.StatusOK, locations)
}

func DeleteLocation(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	_, err = config.DB.Exec("DELETE FROM locations WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Location deleted successfully"})
}

func UpdateLocationName(c *gin.Context) {
	
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	
	var req struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	
	weather, err := services.FetchWeatherData(req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to fetch weather data",
			"detail": err.Error(),
		})
		return
	}

	
	updateQuery := `
		UPDATE locations
		SET name = $1,
			latitude = $2,
			longitude = $3,
			weather_summary = $4,
			temperature = $5,
			wind_speed = $6,
			updated_at = NOW()
		WHERE id = $7`
	_, err = config.DB.Exec(updateQuery,
		req.Name,
		req.Latitude,
		req.Longitude,
		weather.Current.Summary,
		weather.Current.Temperature,
		weather.Current.Wind.Speed,
		id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location and weather"})
		return
	}

	// Ambil prediksi cuaca 1 hari ke depan
	forecastDay, forecastSummary, err := services.FetchOneDayForecast(req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to fetch forecast",
			"detail": err.Error(),
		})
		return
	}

	
deleteOld := `DELETE FROM predictions WHERE location_id = $1 AND date = $2`
_, err = config.DB.Exec(deleteOld, id, forecastDay)
if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old prediction", "detail": err.Error()})
	return
}


insertPrediction := `INSERT INTO predictions (location_id, date, summary) VALUES ($1, $2, $3)`
_, err = config.DB.Exec(insertPrediction, id, forecastDay, forecastSummary)
if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert prediction", "detail": err.Error()})
	return
}


	c.JSON(http.StatusOK, gin.H{"message": "Location updated successfully with new weather and forecast"})
}


func GetPredictions(c *gin.Context) {
	var predictions []models.Prediction

	err := config.DB.Select(&predictions, "SELECT * FROM predictions ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch predictions"})
		return
	}

	c.JSON(http.StatusOK, predictions)
}

// GET /locations/:id/predictions
func GetPredictionsByLocation(c *gin.Context) {
	locationID := c.Param("id")
	var predictions []models.Prediction

	err := config.DB.Select(&predictions, "SELECT * FROM predictions WHERE location_id = $1 ORDER BY date DESC", locationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch predictions"})
		return
	}

	c.JSON(http.StatusOK, predictions)
}
