package redisV1

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/skantay/hezzl/internal/domain/good/model"
	"github.com/skantay/hezzl/internal/domain/good/repository"
)

type goodRepository struct {
	db *redis.Client
}

func New(db *redis.Client) repository.GoodCacheRepository {
	return goodRepository{db}
}

func (g goodRepository) Create(ctx context.Context, good model.Good, key string, duration time.Duration) error {
	err := g.db.Set(key, good, duration)
	if err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (g goodRepository) Delete(ctx context.Context, key string) error {
	err := g.db.Del(key)
	if err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (g goodRepository) Get(ctx context.Context, key string) (model.Good, error) {
	data, err := g.db.Get(key).Result()
	if err == redis.Nil {
		return model.Good{}, model.ErrGoodNotFound
	} else if err != nil {
		return model.Good{}, err
	}

	var goods model.Good
	err = json.Unmarshal([]byte(data), &goods)
	if err != nil {
		return model.Good{}, err
	}

	return goods, nil
}
