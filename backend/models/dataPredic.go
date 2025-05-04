package models

import "time"

type Prediction struct {
	ID         int       `db:"id"`
	LocationID int       `db:"location_id"`
	Date       string    `db:"date"` 
	Summary    string    `db:"summary"`
	CreatedAt  time.Time `db:"created_at"`
}
