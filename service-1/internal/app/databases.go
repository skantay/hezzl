package app

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/skantay/hezzl/config"

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/lib/pq"
)

func connectRedis(cfg config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Database.Redis.Host, cfg.Database.Redis.Port),
		Password: cfg.Database.Redis.Password,
		DB:       0,
	})
	if err := client.Ping(); err.Err() != nil {
		return nil, fmt.Errorf("redis client ping error: %w", err.Err())
	}

	return client, nil
}

func connectPostgres(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf(`
		user=%s
		password=%s
		dbname=%s
		host=%s
		port=%d
		sslmode=%s`,
		cfg.Database.Postgres.User,
		cfg.Database.Postgres.Password,
		cfg.Database.Postgres.DBName,
		cfg.Database.Postgres.Host,
		cfg.Database.Postgres.Port,
		cfg.Database.Postgres.SSLMode,
	))
	if err != nil {
		return nil, fmt.Errorf("sql open error: %w", err)
	}

	// Pinging the database
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping error: %w", err)
	}

	return db, nil
}

func migrateUp(db *sql.DB) error {
	// Migrate up file
	data, err := os.ReadFile("./migrations/setup.up.sql")
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	// Migrating up
	_, err = db.Exec(string(data))
	if err != nil {
		return fmt.Errorf("database setup migration error: %w", err)
	}

	return nil
}

func migrateDown(db *sql.DB) error {
	// Migrate down file
	data, err := os.ReadFile("./migrations/setup.down.sql")
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	// Migrating down
	_, err = db.Exec(string(data))
	if err != nil {
		return fmt.Errorf("database setup migration error: %w", err)
	}
	return nil
}
