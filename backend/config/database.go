package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectDatabase() {
	
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found (this is okay in production)")
	}

	
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("Environment variable POSTGRES_DSN tidak ditemukan")
	}

	
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Database connected!")
}


func GetDB() *sqlx.DB {
	return DB
}
