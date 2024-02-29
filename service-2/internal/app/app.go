package app

import (
	"context"
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/nats-io/nats.go"
	"github.com/skantay/service-2/config"
	"github.com/skantay/service-2/internal/controller/nats/v"
	"github.com/skantay/service-2/internal/usecase"
	repository "github.com/skantay/service-2/internal/usecase/repository/clickhouse"
	"github.com/skantay/service-2/pkg/connClickhouse"
	"github.com/skantay/service-2/pkg/migrate"
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

	db, err := connClickhouse.ConnectClickhouse(cfg)
	if err != nil {
		return err
	}

	if err := migrate.MigrateUp(db); err != nil {
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

	repo := repository.New(db)

	usecaseGood := usecase.New(repo)

	service := usecase.NewService(usecaseGood)

	ctrl := v.New(nc, log, service)

	if err := ctrl.Serve(context.Background()); err != nil {
		return fmt.Errorf("Controller error: %w", err)
	}

	return nil
}
