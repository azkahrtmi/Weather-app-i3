package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// RunMigrations runs all database migrations in the migrations directory
func RunMigrations() error {
	log.Println("Running database migrations...")

	// Ensure DB connection is established
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}

	// Track migration status with a migrations table
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Read all migration files
	migrationsDir := "migrations"
	migrationFiles, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migration files: %v", err)
	}

	// Sort migration files by name to ensure proper order
	for i := 0; i < len(migrationFiles); i++ {
		for j := i + 1; j < len(migrationFiles); j++ {
			if filepath.Base(migrationFiles[i]) > filepath.Base(migrationFiles[j]) {
				migrationFiles[i], migrationFiles[j] = migrationFiles[j], migrationFiles[i]
			}
		}
	}

	// Apply each migration
	for _, file := range migrationFiles {
		version := strings.TrimSuffix(filepath.Base(file), ".sql")

		// Check if migration has already been applied
		var exists bool
		err := DB.Get(&exists, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %v", err)
		}

		if exists {
			log.Printf("Migration %s already applied, skipping", version)
			continue
		}

		// Read migration file
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", file, err)
		}

		// Split migration into up/down parts
		parts := strings.Split(string(content), "-- Down migration")
		upMigration := parts[0]

		// Begin transaction
		tx, err := DB.Beginx()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %v", err)
		}

		// Execute migration
		_, err = tx.Exec(upMigration)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to apply migration %s: %v", version, err)
		}

		// Record migration
		_, err = tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %v", version, err)
		}

		// Commit transaction
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("failed to commit transaction: %v", err)
		}

		log.Printf("Successfully applied migration: %s", version)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
