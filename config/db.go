package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// DB is a package-level variable that holds the global DB connection
var DB *sql.DB

// ConnectDB initializes the connection to the PostgreSQL database.
func ConnectDB() {
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	// host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	name := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("postgres://%s:%s@db:%s/%s?sslmode=disable",
		user, pass, port, name)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("Failed to open DB: %v", err))
	}

	// Verify connection
	if err = DB.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to ping DB: %v", err))
	}

	fmt.Println("Database connection established!")
}
