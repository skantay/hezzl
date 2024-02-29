package app

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/nats-io/nats.go"
	"github.com/skantay/service-2/config"
	"github.com/skantay/service-2/internal/controller"
	service "github.com/skantay/service-2/internal/domain"
	"go.uber.org/zap"
)

func Run() error {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		return err
	}

	configLog := zap.NewDevelopmentConfig()
	configLog.DisableStacktrace = true
	log, err := configLog.Build()
	if err != nil {
		return fmt.Errorf("zap logger error: %w", err)
	}
	defer log.Sync()

	db, err := connectClickhouse(cfg)
	if err != nil {
		return err
	}

	if err := migrateUp(db); err != nil {
		return err
	}

	// In case
	// defer func() error {
	// 	if err := migrateDown(db); err != nil {
	// 		return err
	// 	}

	// 	return nil
	// }()

	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%d",
		cfg.Nats.Host,
		cfg.Nats.Port))
	if err != nil {
		return err
	}
	defer nc.Close()

	service := service.New(db, nc)

	controller := controller.New(service, log)

	if err := controller.Run(); err != nil {
		return err
	}

	return nil
}

func connectClickhouse(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("clickhouse://%s:%d?username=%s&password=%s&database=%s",
		cfg.Database.Clickhouse.Host,
		cfg.Database.Clickhouse.Port,
		cfg.Database.Clickhouse.User,
		cfg.Database.Clickhouse.Password,
		cfg.Database.Clickhouse.DB,
	)
	conn, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
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
