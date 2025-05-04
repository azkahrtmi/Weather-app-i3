package scheduler

import (
	"log"

	"backend/handlers"

	"github.com/robfig/cron/v3"
)

func StartWeatherScheduler() {
	c := cron.New()

	_, err := c.AddFunc("@every 10s", func() {
		log.Println("Running weather update scheduler...")
		handlers.CheckAndUpdateWeather()
	})
	if err != nil {
		log.Fatalf("Failed to schedule job: %v", err)
	}

	c.Start()
	log.Println("Scheduler started")
}
