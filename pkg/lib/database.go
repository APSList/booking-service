package lib

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ======== TYPES ========

// A type alias for the connection pool
type Database = pgxpool.Pool

// ======== METHODS ========

func GetDatabase(logger Logger) *pgxpool.Pool { // Assuming your *Database is a wrapper or just the pool
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		logger.Fatal("DATABASE_URL environment variable is not set")
		os.Exit(1)
	}

	// 1. Parse the connection string into a config object
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		logger.Fatal("Unable to parse DATABASE_URL: ", err)
		os.Exit(1)
	}

	// 2. Fix for Supabase/PgBouncer:
	// This mode tells pgx to describe the statement but not give it a name,
	// avoiding collisions in transaction-mode poolers.
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeDescribeExec

	// 3. Create the pool using the modified config
	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Fatal("Unable to connect to database: ", err)
		os.Exit(1)
	}

	// Test connection
	if err := dbPool.Ping(context.Background()); err != nil {
		logger.Fatal("Database ping failed: ", err)
		os.Exit(1)
	}

	logger.Info("Connected to the database successfully using DescribeExec mode.")
	return dbPool
}
