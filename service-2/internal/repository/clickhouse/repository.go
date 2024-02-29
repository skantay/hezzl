package repository

import (
	"database/sql"
	"fmt"

	"github.com/skantay/service-2/internal/entity"
)

type GoodRepository interface {
	Create(collection entity.Collection) error
}

type goodRepository struct {
	db *sql.DB
}

func New(db *sql.DB) GoodRepository {
	return goodRepository{db}
}

func (g goodRepository) Create(collection entity.Collection) error {
	batch, err := g.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer batch.Rollback()

	stmt, err := batch.Prepare("INSERT INTO default.goods(ID, ProjectID, Name, Description, Priority, Removed, CreatedAt) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, good := range collection.Goods {
		_, err = stmt.Exec(
			good.ID,
			good.ProjectID,
			good.Name,
			good.Description,
			good.Priority,
			good.Removed,
			good.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to execute statement for collection of goods: %w", err)
		}
	}

	if err := batch.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
