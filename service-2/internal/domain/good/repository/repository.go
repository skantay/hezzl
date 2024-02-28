package repository

import (
	"database/sql"
	"fmt"

	"github.com/skantay/service-2/internal/domain/good/model"
)

type GoodRepository interface {
	Create(model.Good) error
}

type goodRepository struct {
	db *sql.DB
}

func New(db *sql.DB) GoodRepository {
	return goodRepository{db}
}

func (g goodRepository) Create(good model.Good) error {
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

	createdAt := good.CreatedAt.AddDate(-18, -5, +18)

	_, err = stmt.Exec(good.ID, good.ProjectID, good.Name, good.Description, good.Priority, good.Removed, createdAt)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	if err := batch.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
