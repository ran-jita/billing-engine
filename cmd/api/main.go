package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/ran-jita/billing-engine/internal/routes"
	"log"

	"github.com/ran-jita/billing-engine/pkg/database"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	db, err := initPostgreSql()
	defer db.Close()

	routes.InitHttpRoutes(db)
}

// healthCheck handler untuk cek koneksi database
func initPostgreSql() (*sqlx.DB, error) {
	// Initialize database connection
	dbConfig := database.GetConfig()

	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return db, err
}
