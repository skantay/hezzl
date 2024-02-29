package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/skantay/hezzl/internal/entity"
)

type GoodCacheRepository interface {
	Create(ctx context.Context, good entity.Good, key string, duration time.Duration) error
	Delete(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (entity.Good, error)
}

type goodRepository struct {
	db *redis.Client
}

func New(db *redis.Client) GoodCacheRepository {
	return goodRepository{db}
}

func (g goodRepository) Create(ctx context.Context, good entity.Good, key string, duration time.Duration) error {
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

func (g goodRepository) Get(ctx context.Context, key string) (entity.Good, error) {
	data, err := g.db.Get(key).Result()

	if err == redis.Nil {
		return entity.Good{}, entity.ErrGoodNotFound
	} else if err != nil {
		return entity.Good{}, err
	}

	var goods entity.Good
	err = json.Unmarshal([]byte(data), &goods)
	if err != nil {
		return entity.Good{}, err
	}

	return goods, nil
}
