package services

import (
	"backend/config"
	"backend/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
)

type WeatherResponse struct {
	Current struct {
		Summary     string  `json:"summary"`
		Temperature float64 `json:"temperature"`
		Wind        struct {
			Speed float64 `json:"speed"`
		} `json:"wind"`
	} `json:"current"`
}


// ambil current weather
func FetchWeatherData(lat, lon float64) (WeatherResponse, error) {
	var weather WeatherResponse


	apiKey := os.Getenv("METEOSOURCE_API_KEY")
	if apiKey == "" {
		return weather, fmt.Errorf("API key not found in environment variable")
	}

	url := fmt.Sprintf(
		"https://www.meteosource.com/api/v1/free/point?lat=%f&lon=%f&sections=current&timezone=auto&units=metric&language=en&key=%s",
		lat, lon, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return weather, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return weather, fmt.Errorf("MeteoSource API error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return weather, err
	}

	err = json.Unmarshal(body, &weather)
	if err != nil {
		return weather, err
	}

	return weather, nil
}


type OneDayForecastResponse struct {
	Daily struct {
		Data []struct {
			Day     string `json:"day"`
			Summary string `json:"summary"`
		} `json:"data"`
	} `json:"daily"`
}

func FetchOneDayForecast(lat, lon float64) (string, string, error) {
	var forecast OneDayForecastResponse

	apiKey := os.Getenv("METEOSOURCE_API_KEY")
	if apiKey == "" {
		return "", "", fmt.Errorf("API key not set")
	}

	url := fmt.Sprintf(
		"https://www.meteosource.com/api/v1/free/point?lat=%f&lon=%f&sections=daily&timezone=auto&units=metric&language=en&key=%s",
		lat, lon, apiKey,
	)

	// Debug: URL yang dipanggil
	fmt.Println("Requesting URL:", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", "", fmt.Errorf("HTTP request error: %v", err)
	}
	defer resp.Body.Close()

	// Debug: Status code
	fmt.Println("Response status:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("API error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("read body error: %v", err)
	}

	// Debug: Raw JSON body
	fmt.Println("Raw JSON:\n", string(body))

	err = json.Unmarshal(body, &forecast)
	if err != nil {
		return "", "", fmt.Errorf("JSON unmarshal error: %v", err)
	}

	dataLen := len(forecast.Daily.Data)
	// Debug: jumlah data
	fmt.Println("Length of forecast data:", dataLen)

	if dataLen == 0 {
		return "", "", fmt.Errorf("no forecast data received")
	}

	// Gunakan data[1] jika tersedia, jika tidak fallback ke data[0]
	index := 0
	if dataLen > 1 {
		index = 1
	}
	day := forecast.Daily.Data[index].Day
	summary := forecast.Daily.Data[index].Summary

	// Debug: hasil akhir
	fmt.Printf("Forecast - Day: %s, Summary: %s\n", day, summary)

	return day, summary, nil
}


func FetchAndSaveWeather(loc models.Location) error {
    weather, err := FetchWeatherData(loc.Latitude, loc.Longitude)
    if err != nil {
        return err
    }

    // Jika tidak berubah, skip update
    if loc.WeatherSummary.Valid && loc.WeatherSummary.String == weather.Current.Summary &&
       loc.Temperature.Valid && math.Abs(loc.Temperature.Float64 - weather.Current.Temperature) < 0.01 &&
       loc.WindSpeed.Valid && math.Abs(loc.WindSpeed.Float64 - weather.Current.Wind.Speed) < 0.01 {
        log.Printf("No weather change for %s. Skipping update.", loc.Name)
        return nil
    }

    // Update ke DB karena berubah
    _, err = config.GetDB().Exec(`
        UPDATE locations 
        SET weather_summary = $1, temperature = $2, wind_speed = $3, updated_at = NOW()
        WHERE id = $4
    `, weather.Current.Summary, weather.Current.Temperature, weather.Current.Wind.Speed, loc.ID)

    return err
}
