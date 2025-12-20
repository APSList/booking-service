package lib

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ======== TYPES ========

// A type alias for the connection pool
type Database = pgxpool.Pool

// ======== METHODS ========

// GetDatabase returns a database pool to connect to the database asynchronously
func GetDatabase(logger Logger) *Database {

	// Debug: izpiši environment spremenljivke
	fmt.Println("DATABASE_URL:", os.Getenv("DATABASE_URL"))

	// Uporabi environment variable DATABASE_URL, ki mora vsebovati Supabase connection string
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		logger.Fatal("DATABASE_URL environment variable is not set")
		os.Exit(1)
	}

	// Create a connection pool to the database using pgxpool
	dbPool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		logger.Fatal("Unable to connect to database: ", err)
		os.Exit(1)
	}

	// Testiraj povezavo
	if err := dbPool.Ping(context.Background()); err != nil {
		logger.Fatal("Database ping failed: ", err)
		os.Exit(1)
	}

	logger.Info("Connected to the database successfully.")
	return dbPool
}
