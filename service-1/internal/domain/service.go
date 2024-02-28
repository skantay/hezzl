package domain

import (
	"database/sql"

	"github.com/nats-io/nats.go"
	"github.com/skantay/hezzl/internal/domain/good/repository/postgres"
	"github.com/skantay/hezzl/internal/domain/good/repository/redisV1"
	"github.com/skantay/hezzl/internal/domain/good/usecase"
	"go.uber.org/zap"

	"github.com/go-redis/redis"
)

type Service struct {
	GoodService usecase.GoodUsecase
}

func New(db *sql.DB, cache *redis.Client, nc *nats.Conn, log *zap.Logger) Service {
	goodRepo := postgres.New(db, nc)

	goodCache := redisV1.New(cache)

	goodService := usecase.New(goodRepo, goodCache, log)

	return Service{
		GoodService: goodService,
	}
}
