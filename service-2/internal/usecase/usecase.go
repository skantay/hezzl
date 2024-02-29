package usecase

import (
	"encoding/json"
	"fmt"

	"github.com/skantay/service-2/internal/entity"
	repository "github.com/skantay/service-2/internal/usecase/repository/clickhouse"
)

type Service struct {
	Good GoodUsecase
}

func NewService(good GoodUsecase) Service {
	return Service{good}
}

type GoodUsecase interface {
	Create(data []byte) error
}

type goodUsecase struct {
	repo repository.GoodRepository
}

func New(repo repository.GoodRepository) GoodUsecase {
	return goodUsecase{repo}
}

func (g goodUsecase) Create(data []byte) error {
	var collection entity.Collection

	err := json.Unmarshal(data, &collection)
	if err != nil {
		return fmt.Errorf("here error: %w", err)
	}

	return g.repo.Create(collection)
}
