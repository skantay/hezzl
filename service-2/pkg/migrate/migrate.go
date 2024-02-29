package migrate

import (
	"database/sql"
	"fmt"
	"os"
)

func MigrateUp(db *sql.DB) error {
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

func MigrateDown(db *sql.DB) error {
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
