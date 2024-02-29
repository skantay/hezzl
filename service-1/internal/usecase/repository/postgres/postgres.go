package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/skantay/hezzl/internal/controller/mq/nats/v"
	"github.com/skantay/hezzl/internal/entity"
)

type GoodRepository interface {
	Create(ctx context.Context, good entity.Good) (entity.Good, error)
	Delete(ctx context.Context, id, projectID int) (entity.Good, error)
	UpdateDesc(ctx context.Context, name string, id, projectID int) (entity.Good, error)
	UpdateName(ctx context.Context, desc string, id, projectID int) (entity.Good, error)
	UpdatePriority(ctx context.Context, priority, id, projectID int) ([]entity.Good, error)
	Get(ctx context.Context, id int) (entity.Good, error)
	GetMaxPriority(ctx context.Context) (int, error)
	CountRows(ctx context.Context) (int, error)
}

type Collection struct {
	Good  entity.Good   `json:"good"`
	Goods []entity.Good `json:"goods"`
}

type goodRepository struct {
	db *sql.DB
	nc v.NC
}

func New(db *sql.DB, nc v.NC) GoodRepository {
	return goodRepository{db, nc}
}

func (g goodRepository) GetMaxPriority(ctx context.Context) (int, error) {
	var maxPriority *int

	if err := g.db.QueryRowContext(ctx, "SELECT MAX(priority) FROM goods").Scan(&maxPriority); err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("get max priority error: %w", err)
	}
	if maxPriority == nil {
		return 0, nil
	}
	return *maxPriority, nil
}

func (g goodRepository) Create(ctx context.Context, good entity.Good) (entity.Good, error) {
	stmt := `INSERT INTO goods(project_id, name, description, priority, removed, created_at)
             VALUES($1, $2, $3, $4, $5, $6) RETURNING *;`

	var newGood entity.Good
	err := g.db.QueryRowContext(ctx, stmt,
		good.ProjectID,
		good.Name,
		good.Description,
		good.Priority,
		good.Removed,
		good.CreatedAt,
	).Scan(
		&newGood.ID,
		&newGood.ProjectID,
		&newGood.Name,
		&newGood.Description,
		&newGood.Priority,
		&newGood.Removed,
		&newGood.CreatedAt,
	)
	if err != nil {
		if err.Error() == `pq: insert or update on table "goods" violates foreign key constraint "goods_project_id_fkey"` {
			return newGood, entity.ErrProjectNotFound
		}
		return newGood, fmt.Errorf("trouble executing db: %w", err)
	}

	return newGood, g.nc.SendJSON(ctx, Collection{Good: newGood})
}

func (g goodRepository) Delete(ctx context.Context, id, projectID int) (entity.Good, error) {
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		return entity.Good{}, fmt.Errorf("trouble with starting a transaction: %w", err)
	}
	defer tx.Rollback()

	stmt := `UPDATE goods SET   
                 removed = $1
             WHERE id = $2 AND project_id = $3 RETURNING *;`

	var updatedGood entity.Good

	if err := g.db.QueryRowContext(
		ctx,
		stmt,
		true,
		id,
		projectID,
	).Scan(
		&updatedGood.ID,
		&updatedGood.ProjectID,
		&updatedGood.Name,
		&updatedGood.Description,
		&updatedGood.Priority,
		&updatedGood.Removed,
		&updatedGood.CreatedAt,
	); err != nil {
		return entity.Good{}, entity.ErrGoodNotFound
	}

	err = tx.Commit()
	if err != nil {
		return entity.Good{}, fmt.Errorf("trouble with committing a transaction: %w", err)
	}

	return updatedGood, g.nc.SendJSON(ctx, Collection{Good: updatedGood})
}

func (g goodRepository) UpdateName(ctx context.Context, name string, id, projectID int) (entity.Good, error) {
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		return entity.Good{}, fmt.Errorf("trouble with starting a transaction: %w", err)
	}
	defer tx.Rollback()

	stmt := `UPDATE goods SET   
                 name = $1  
             WHERE id = $2 AND project_id = $3 RETURNING *;`

	var updatedGood entity.Good

	if err := g.db.QueryRowContext(
		ctx,
		stmt,
		name,
		id,
		projectID,
	).Scan(
		&updatedGood.ID,
		&updatedGood.ProjectID,
		&updatedGood.Name,
		&updatedGood.Description,
		&updatedGood.Priority,
		&updatedGood.Removed,
		&updatedGood.CreatedAt,
	); err != nil {
		return entity.Good{}, entity.ErrGoodNotFound
	}

	err = tx.Commit()
	if err != nil {
		return entity.Good{}, fmt.Errorf("trouble with committing a transaction: %w", err)
	}

	return updatedGood, g.nc.SendJSON(ctx, Collection{Good: updatedGood})
}

func (g goodRepository) UpdateDesc(ctx context.Context, desc string, id, projectID int) (entity.Good, error) {
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		return entity.Good{}, fmt.Errorf("trouble with starting a transaction: %w", err)
	}
	defer tx.Rollback()

	stmt := `UPDATE goods SET
				description = $1   
             WHERE id = $2 AND project_id = $3 RETURNING *;`

	var updatedGood entity.Good

	if err := g.db.QueryRowContext(
		ctx,
		stmt,
		desc,
		id,
		projectID,
	).Scan(
		&updatedGood.ID,
		&updatedGood.ProjectID,
		&updatedGood.Name,
		&updatedGood.Description,
		&updatedGood.Priority,
		&updatedGood.Removed,
		&updatedGood.CreatedAt,
	); err != nil {
		return entity.Good{}, entity.ErrGoodNotFound
	}

	err = tx.Commit()
	if err != nil {
		return entity.Good{}, fmt.Errorf("trouble with committing a transaction: %w", err)
	}

	return updatedGood, g.nc.SendJSON(ctx, Collection{Good: updatedGood})
}

func (g goodRepository) UpdatePriority(ctx context.Context, priority, id, projectID int) ([]entity.Good, error) {
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("trouble with starting a transaction: %w", err)
	}
	defer tx.Rollback()

	result := make([]entity.Good, 0)

	stmt := `UPDATE goods SET priority = priority + 1 WHERE priority >= $1 AND project_id = $2;`
	_, err = tx.ExecContext(ctx, stmt, priority, projectID)
	if err != nil {
		return nil, fmt.Errorf("trouble with updating priorities: %w", err)
	}

	stmt = `UPDATE goods SET priority = $1 WHERE id = $2 AND project_id = $3;`
	_, err = tx.ExecContext(ctx, stmt, priority, id, projectID)
	if err != nil {
		return nil, entity.ErrGoodNotFound
	}

	stmt = `SELECT * FROM goods WHERE project_id = $1 ORDER BY priority;`

	rows, err := tx.QueryContext(ctx, stmt, projectID)
	if err != nil {
		return nil, fmt.Errorf("trouble with updating entities: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var good entity.Good
		if err := rows.Scan(
			&good.ID,
			&good.ProjectID,
			&good.Name,
			&good.Description,
			&good.Priority,
			&good.Removed,
			&good.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("trouble with scanning row: %w", err)
		}

		result = append(result, good)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("trouble with committing transaction: %w", err)
	}

	return result, g.nc.SendJSON(ctx, Collection{Goods: result})
}

func (g goodRepository) Get(ctx context.Context, id int) (entity.Good, error) {
	stmt := `SELECT * FROM goods WHERE id = $1`

	var good entity.Good
	err := g.db.QueryRowContext(ctx, stmt, id).Scan(
		&good.ID,
		&good.ProjectID,
		&good.Name,
		&good.Description,
		&good.Priority,
		&good.Removed,
		&good.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Good{}, fmt.Errorf("good with id #%d %w", id, entity.ErrGoodNotFound)
		}

		return entity.Good{}, fmt.Errorf("query error: %w", err)
	}

	return good, nil
}

func (g goodRepository) CountRows(ctx context.Context) (int, error) {
	stmt := `SELECT COUNT(id) FROM goods;`

	var count int

	if err := g.db.QueryRowContext(ctx, stmt).Scan(&count); err != nil {
		return 0, fmt.Errorf("query error: %w", err)
	}

	return count, nil
}