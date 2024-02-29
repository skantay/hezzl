package usecase

import (
	"encoding/json"
	"fmt"

	"github.com/skantay/service-2/internal/domain/good/model"
	"github.com/skantay/service-2/internal/domain/good/repository"
)

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
	var collection model.Collection

	err := json.Unmarshal(data, &collection)
	if err != nil {
		return fmt.Errorf("here error: %w", err)
	}

	return g.repo.Create(collection)
}
