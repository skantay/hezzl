package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator"
	"github.com/skantay/hezzl/config"
	"github.com/skantay/hezzl/internal/controller/api"
	"github.com/skantay/hezzl/internal/controller/mq/nats/v"
	"github.com/skantay/hezzl/internal/repository/postgres"
	cache "github.com/skantay/hezzl/internal/repository/redis"
	"github.com/skantay/hezzl/internal/usecase"
	"github.com/skantay/hezzl/pkg/migrate"
	psql "github.com/skantay/hezzl/pkg/postgres"
	rds "github.com/skantay/hezzl/pkg/redis"

	"github.com/nats-io/nats.go"
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
	db, err := psql.ConnectPostgres(cfg)
	if err != nil {
		return fmt.Errorf("postgres connection error: %w", err)
	}
	defer db.Close()

	// Migrating up
	if err := migrate.MigrateUp(db); err != nil {
		return fmt.Errorf("migration up error: %w", err)
	}

	/*


		//In case database needs to be dropped


		defer func() error {
			// Migrating down
			if err := migrate.MigrateDown(db); err != nil {
				return fmt.Errorf("migration down error: %w", err)
			}
			return nil
		}()


	*/

	// Connecting redis client
	client, err := rds.ConnectRedis(cfg)
	if err != nil {
		return fmt.Errorf("redis connection error: %w", err)
	}
	defer client.Close()

	// Connecting to nats
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%d", cfg.Nats.Host, cfg.Nats.Port))
	if err != nil {
		return err
	}

	natsI := v.New(nc)

	goodUsecase := usecase.NewGoodUsecase(
		postgres.New(db, natsI),
		cache.New(client))

	service := usecase.NewService(goodUsecase)

	validate := validator.New()

	// Init controller
	ctrl := api.New(service, log, cfg, validate)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info("Service started")

	go func() error {
		// Run controller
		if err := ctrl.Serve(ctx); err != nil {
			return fmt.Errorf("Controller error: %w", err)
		}

		return nil
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown

	log.Info("Server shut down")
	return nil
}
