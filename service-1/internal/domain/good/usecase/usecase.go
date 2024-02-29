package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/skantay/hezzl/internal/domain/good/model"
	"github.com/skantay/hezzl/internal/domain/good/repository"
	"go.uber.org/zap"
)

type GoodUsecase interface {
	CreateGood(ctx context.Context, projectID int, name string) (model.Good, error)
	DeleteGood(ctx context.Context, id, projectID int) (model.Good, error)
	UpdateGood(ctx context.Context, id, projectID int, name, desc string, emptyDesc bool) (model.Good, error)
	GetGoods(ctx context.Context, limit, offset int) ([]model.Good, error)
	Reprioritiize(ctx context.Context, priority, id, projectID int) ([]model.Good, error)
}

type goodUsecase struct {
	repo  repository.GoodRepository
	cache repository.GoodCacheRepository
	log   *zap.Logger
}

func New(repo repository.GoodRepository, cache repository.GoodCacheRepository, log *zap.Logger) GoodUsecase {
	return goodUsecase{repo, cache, log}
}

func (g goodUsecase) CreateGood(ctx context.Context, projectID int, name string) (model.Good, error) {
	maxPriority, err := g.repo.GetMaxPriority(ctx)
	if err != nil {
		return model.Good{}, err
	}

	good, err := g.repo.Create(ctx, model.Good{
		Name:      name,
		ProjectID: projectID,
		Priority:  maxPriority + 1,
		Removed:   false,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return good, err
	}

	return good, nil
}

func (g goodUsecase) DeleteGood(ctx context.Context, id, projectID int) (model.Good, error) {
	key := fmt.Sprintf("goods_%d", id)
	_ = g.cache.Delete(ctx, key)

	deletedGood, err := g.repo.Delete(ctx, id, projectID)
	if err != nil {
		return model.Good{}, err
	}

	return deletedGood, nil
}

func (g goodUsecase) UpdateGood(ctx context.Context, id, projectID int, name, desc string, emptyDesc bool) (model.Good, error) {
	key := fmt.Sprintf("goods_%d", id)
	_ = g.cache.Delete(ctx, key)

	good, err := g.repo.UpdateNameDesc(ctx, model.Good{
		ID:          id,
		ProjectID:   projectID,
		Name:        name,
		Description: desc,
	}, emptyDesc)
	if err != nil {
		return model.Good{}, err
	}

	return good, nil
}

func (g goodUsecase) GetGoods(ctx context.Context, limit, offset int) ([]model.Good, error) {
	length := int(math.Abs(float64(limit - offset)))
	goods := make([]model.Good, 0, length)

	maxRows, err := g.repo.CountRows(ctx)
	if err != nil {
		return nil, err
	}

	if limit > maxRows {
		limit = maxRows - 1
	}

	for i := offset; limit >= 0; limit-- {

		var good model.Good

		key := fmt.Sprintf("goods_%d", i)

		good, err := g.cache.Get(ctx, key)
		if err != nil && !errors.Is(err, model.ErrGoodNotFound) {
			return goods, fmt.Errorf("redis error get: %w", err)
		}

		g.log.Sugar().Infof("%s getting from cache", key)

		if errors.Is(err, model.ErrGoodNotFound) {
			g.log.Sugar().Infof("in cache not found. %s getting from db", key)

			good, err = g.repo.Get(ctx, i)
			if err != nil {
				return goods, err
			}

			if err := g.cache.Create(ctx, good, key, time.Minute); err != nil {
				return goods, fmt.Errorf("redis error create: %w", err)
			}
		}

		goods = append(goods, good)
		i++
	}

	return goods, nil
}

func (g goodUsecase) Reprioritiize(ctx context.Context, priority, id, projectID int) ([]model.Good, error) {
	goods, err := g.repo.UpdatePriority(ctx, priority, id, projectID)
	if err != nil {
		return nil, err
	}

	for _, good := range goods {
		key := fmt.Sprintf("goods_%d", good.ID)
		_ = g.cache.Delete(ctx, key)
	}

	return goods, nil
}
