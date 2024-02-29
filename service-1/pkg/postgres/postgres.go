package psql

import (
	"database/sql"
	"fmt"

	"github.com/skantay/hezzl/config"

	_ "github.com/lib/pq"
)

func ConnectPostgres(cfg config.Config) (*sql.DB, error) {
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
