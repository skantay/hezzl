package repository

import (
	"database/sql"
	"fmt"

	"github.com/skantay/service-2/internal/domain/good/model"
)

type GoodRepository interface {
	Create(collection model.Collection) error
}

type goodRepository struct {
	db *sql.DB
}

func New(db *sql.DB) GoodRepository {
	return goodRepository{db}
}

func (g goodRepository) Create(collection model.Collection) error {
	batch, err := g.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer batch.Rollback()

	stmt, err := batch.Prepare("INSERT INTO default.goods(ID, ProjectID, Name, Description, Priority, Removed, CreatedAt) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	if collection.Good != (model.Good{}) {
		createdAt := collection.Good.CreatedAt.AddDate(-18, -5, +18)
		_, err = stmt.Exec(collection.Good.ID, collection.Good.ProjectID, collection.Good.Name, collection.Good.Description, collection.Good.Priority, collection.Good.Removed, createdAt)
		if err != nil {
			return fmt.Errorf("failed to execute statement for single good: %w", err)
		}
	}

	for _, good := range collection.Goods {
		createdAt := good.CreatedAt.AddDate(-18, -5, +18)
		_, err = stmt.Exec(good.ID, good.ProjectID, good.Name, good.Description, good.Priority, good.Removed, createdAt)
		if err != nil {
			return fmt.Errorf("failed to execute statement for collection of goods: %w", err)
		}
	}

	if err := batch.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
