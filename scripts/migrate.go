package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "local"
	}

	envFile := fmt.Sprintf(".env.%s", env)

	if err := godotenv.Load(envFile); err != nil {
		fmt.Printf("No %s file found, fallback ke .env\n", envFile)
		_ = godotenv.Load(".env")
	}

	// Flag untuk menentukan aksi migrasi
	var up, down bool
	var steps int
	flag.BoolVar(&up, "up", false, "Migrate up")
	flag.BoolVar(&down, "down", false, "Migrate down")
	flag.IntVar(&steps, "steps", 0, "Number of steps to migrate")
	flag.Parse()

	// Validate migration direction
	if (up && down) || (!up && !down) {
		log.Fatalf("[ERROR] Please specify either --up or --down, but not both")
	}

	direction := "up"
	if down {
		direction = "down"
	}

	log.Printf("[CONFIG] Migration direction: %s, Steps: %d\n", direction, steps)

	// Build connection string
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSchema := os.Getenv("DB_SCHEMA")

	if dbHost == "" {
		dbHost = "localhost"
		log.Println("[CONFIG] Using default DB_HOST: localhost")
	}
	if dbPort == "" {
		dbPort = "5432"
		log.Println("[CONFIG] Using default DB_PORT: 5432")
	}

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table=%s_migrations",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
		dbSchema,
	)

	log.Printf("[DB] Attempting to create migration with host: %s, port: %s, database: %s, schema: %s\n", dbHost, dbPort, dbName, dbSchema)

	m, err := migrate.New("file://db/migrations", connectionString)
	if err != nil {
		log.Fatalf("[MIGRATION] Failed to create migration: %v\n", err)
	}

	// Jalankan sesuai argumen
	log.Println("[MIGRATION] Starting migration process...")
	if direction == "up" {
		log.Printf("[MIGRATION] Performing %d up migration step(s)\n", steps)
		if steps > 0 {
			err = m.Steps(steps)
		} else {
			err = m.Up()
		}
	} else {
		log.Printf("[MIGRATION] Performing %d down migration step(s)\n", steps)
		if steps > 0 {
			err = m.Steps(-steps)
		} else {
			err = m.Down()
		}
	}

	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("[MIGRATION] No changes to apply. Database is up to date.")
		} else {
			log.Fatalf("[MIGRATION] Migration failed with error: %v\n", err)
		}
	} else {
		log.Println("[MIGRATION] Migration completed successfully")
	}
}
