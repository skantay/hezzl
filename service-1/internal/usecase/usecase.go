package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/skantay/hezzl/internal/entity"
	"github.com/skantay/hezzl/internal/repository/postgres"
	cache "github.com/skantay/hezzl/internal/repository/redis"
)

type Service struct {
	Good GoodUsecase
}

func NewService(good GoodUsecase) Service {
	return Service{good}
}

type GoodUsecase interface {
	Create(ctx context.Context, projectID int, name string) (entity.Good, error)
	Delete(ctx context.Context, id, projectID int) (entity.Good, error)
	Update(ctx context.Context, id, projectID int, name, desc string, emptyDesc bool) (entity.Good, error)
	List(ctx context.Context, limit, offset int) ([]entity.Good, error)
	Reprioritiize(ctx context.Context, priority, id, projectID int) ([]entity.Good, error)
}

type goodUsecase struct {
	repo  postgres.GoodRepository
	cache cache.GoodCacheRepository
}

func NewGoodUsecase(repo postgres.GoodRepository, cache cache.GoodCacheRepository) GoodUsecase {
	return goodUsecase{
		repo:  repo,
		cache: cache,
	}
}

func (g goodUsecase) Create(ctx context.Context, projectID int, name string) (entity.Good, error) {
	maxPriority, err := g.repo.GetMaxPriority(ctx, projectID)
	if err != nil {
		return entity.Good{}, fmt.Errorf("repository get max priority error: %w", err)
	}

	good, err := g.repo.Create(ctx, entity.Good{
		Name:      name,
		ProjectID: projectID,
		Priority:  maxPriority + 1,
		Removed:   false,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return good, fmt.Errorf("trouble creating a good: %w", err)
	}

	return good, nil
}

func (g goodUsecase) Delete(ctx context.Context, id, projectID int) (entity.Good, error) {
	deleted, err := g.repo.Delete(ctx, id, projectID)
	if err != nil {
		return entity.Good{}, fmt.Errorf("trouble deleting a good: %w", err)
	}

	key := fmt.Sprintf("goods_%d", deleted.ID)

	return deleted, g.cache.Delete(ctx, key)
}

func (g goodUsecase) Update(ctx context.Context, id, projectID int, name, desc string, emptyDesc bool) (entity.Good, error) {
	updated, err := g.repo.UpdateName(ctx, name, id, projectID)
	if err != nil {
		return entity.Good{}, fmt.Errorf("trouble updating a good: %w", err)
	}

	if !emptyDesc {
		updated, err = g.repo.UpdateDesc(ctx, desc, id, projectID)
		if err != nil {
			return entity.Good{}, fmt.Errorf("trouble updating a good: %w", err)
		}
	}

	key := fmt.Sprintf("goods_%d", id)

	return updated, g.cache.Delete(ctx, key)
}

func (g goodUsecase) List(ctx context.Context, limit, offset int) ([]entity.Good, error) {
	length := int(math.Abs(float64(limit - offset)))

	goods := make([]entity.Good, 0, length)

	maxRows, err := g.repo.CountRows(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository count rows error: %w", err)
	}

	if limit > maxRows {
		limit = maxRows
	}

	for limit > 0 {
		var good entity.Good

		key := fmt.Sprintf("goods_%d", offset)

		good, err := g.cache.Get(ctx, key)
		if err != nil && !errors.Is(err, entity.ErrGoodNotFound) {
			return nil, fmt.Errorf("cache error get: %w", err)
		}

		if errors.Is(err, entity.ErrGoodNotFound) {
			good, err = g.repo.Get(ctx, offset)
			if err != nil {
				if errors.Is(err, entity.ErrGoodNotFound) {
					return goods, entity.ErrGoodNotFound
				}
				return nil, fmt.Errorf("repository error get: %w", err)
			}

			if err := g.cache.Create(ctx, good, key, time.Minute); err != nil {
				return nil, fmt.Errorf("cache error create: %w", err)
			}
		}

		goods = append(goods, good)
		offset++
		limit--
	}

	return goods, nil
}

func (g goodUsecase) Reprioritiize(ctx context.Context, priority, id, projectID int) ([]entity.Good, error) {
	goods, err := g.repo.UpdatePriority(ctx, priority, id, projectID)
	if err != nil {
		return nil, fmt.Errorf("trouble getting a good: %w", err)
	}

	for _, good := range goods {
		key := fmt.Sprintf("goods_%d", good.ID)
		_ = g.cache.Delete(ctx, key)
	}

	return goods, nil
}
