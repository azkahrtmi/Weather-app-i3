package models

import (
	"database/sql"
	"time"
)

type Location struct {
	ID             int       `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	Latitude       float64   `db:"latitude" json:"latitude"`
	Longitude      float64   `db:"longitude" json:"longitude"`
	WeatherSummary sql.NullString    `db:"weather_summary" json:"weather_summary"`
	Temperature    sql.NullFloat64   `db:"temperature" json:"temperature"`
	WindSpeed      sql.NullFloat64  `db:"wind_speed" json:"wind_speed"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
	UpdateIntervalHours sql.NullInt64 `db:"update_interval_hours" json:"-"`

}
