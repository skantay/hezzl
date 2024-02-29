package rds

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/skantay/hezzl/config"
)

func ConnectRedis(cfg config.Config) (*redis.Client, error) {
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
