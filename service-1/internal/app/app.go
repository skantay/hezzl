package app

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/skantay/hezzl/config"
	"github.com/skantay/hezzl/internal/controller/api/ginV1"
	"github.com/skantay/hezzl/internal/domain"
	"go.uber.org/zap"
)

func Run() error {
	// Config setup
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	// Zap logger setup
	configLog := zap.NewDevelopmentConfig()
	configLog.DisableStacktrace = true
	log, err := configLog.Build()
	if err != nil {
		return fmt.Errorf("zap logger error: %w", err)
	}
	defer log.Sync()

	// Connecting to postgres
	db, err := connectPostgres(cfg)
	if err != nil {
		return fmt.Errorf("postgres connection error: %w", err)
	}
	defer db.Close()

	// Migrating up
	if err := migrateUp(db); err != nil {
		return fmt.Errorf("migration up error: %w", err)
	}

	// In case database needs to be dropped
	// defer func() error {
	// 	// Migrating down
	// 	if err := migrateDown(db); err != nil {
	// 		return fmt.Errorf("migration down error: %w", err)
	// 	}
	// 	return nil
	// }()

	// Connecting redis client
	client, err := connectRedis(cfg)
	if err != nil {
		return fmt.Errorf("redis connection error: %w", err)
	}
	defer client.Close()

	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%d", cfg.Nats.Host, cfg.Nats.Port))
	if err != nil {
		return err
	}

	// Init service
	service := domain.New(db, client, nc, log)

	// Init controller
	ctrl := ginV1.New(service, log, cfg)

	// Run controller
	if err := ctrl.Run(); err != nil {
		return err
	}

	return nil
}
