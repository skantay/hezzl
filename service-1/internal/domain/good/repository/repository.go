package repository

import (
	"context"
	"time"

	"github.com/skantay/hezzl/internal/domain/good/model"
)

type GoodRepository interface {
	Create(ctx context.Context, good model.Good) (model.Good, error)
	Delete(ctx context.Context, id, projectID int) (model.Good, error)
	UpdateNameDesc(ctx context.Context, good model.Good, emptyDesc bool) (model.Good, error)
	UpdatePriority(ctx context.Context, priority, id, projectID int) ([]model.Good, error)
	Get(ctx context.Context, id int) (model.Good, error)
	GetMaxPriority(ctx context.Context) (int, error)
	CountRows(ctx context.Context) (int, error)
}

type GoodCacheRepository interface {
	Create(ctx context.Context, good model.Good, key string, duration time.Duration) error
	Delete(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (model.Good, error)
}
