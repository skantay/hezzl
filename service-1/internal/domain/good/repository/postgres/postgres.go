package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/skantay/hezzl/internal/domain/good/model"
	"github.com/skantay/hezzl/internal/domain/good/repository"
)

type Collection struct {
	Good  model.Good   `json:"good"`
	Goods []model.Good `json:"goods"`
}

type goodRepository struct {
	db *sql.DB
	nc *nats.Conn
}

func New(db *sql.DB, nc *nats.Conn) repository.GoodRepository {
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

func (g goodRepository) Create(ctx context.Context, good model.Good) (model.Good, error) {
	stmt := `INSERT INTO goods(project_id, name, description, priority, removed, created_at)
             VALUES($1, $2, $3, $4, $5, $6) RETURNING *;`

	var newGood model.Good
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
			return newGood, model.ErrGoodNotFound
		}
		return newGood, fmt.Errorf("db exec error: %w", err)
	}

	ec, err := nats.NewEncodedConn(g.nc, "json")
	if err != nil {
		return newGood, err
	}

	if err := ec.Publish("good", Collection{Good: newGood}); err != nil {
		return newGood, fmt.Errorf("publish error: %w", err)
	}

	return newGood, nil
}

func (g goodRepository) Delete(ctx context.Context, id, projectID int) (model.Good, error) {
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Good{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt := `UPDATE goods SET   
                 removed = $1
             WHERE id = $2 AND project_id = $3 RETURNING *;`

	var updatedGood model.Good

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
		return model.Good{}, model.ErrGoodNotFound
	}

	err = tx.Commit()
	if err != nil {
		return model.Good{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	ec, err := nats.NewEncodedConn(g.nc, "json")
	if err != nil {
		return model.Good{}, err
	}

	if err := ec.Publish("good", Collection{Good: updatedGood}); err != nil {
		return model.Good{}, fmt.Errorf("publish error: %w", err)
	}

	return updatedGood, nil
}

func (g goodRepository) UpdateNameDesc(ctx context.Context, good model.Good, emptyDesc bool) (model.Good, error) {
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Good{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var description sql.NullString

	if good.Description == "" && emptyDesc {
		description.Valid = false
	} else {
		description.String = good.Description
		description.Valid = true
	}

	stmt := `UPDATE goods SET   
                 name = $1,   
                 description = COALESCE($2, description)   
             WHERE id = $3 AND project_id = $4 RETURNING *;`

	var updatedGood model.Good

	if err := g.db.QueryRowContext(
		ctx,
		stmt,
		good.Name,
		description,
		good.ID,
		good.ProjectID,
	).Scan(
		&updatedGood.ID,
		&updatedGood.ProjectID,
		&updatedGood.Name,
		&updatedGood.Description,
		&updatedGood.Priority,
		&updatedGood.Removed,
		&updatedGood.CreatedAt,
	); err != nil {
		return model.Good{}, model.ErrGoodNotFound
	}

	err = tx.Commit()
	if err != nil {
		return model.Good{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	ec, err := nats.NewEncodedConn(g.nc, "json")
	if err != nil {
		return model.Good{}, err
	}

	if err := ec.Publish("good", Collection{Good: updatedGood}); err != nil {
		return model.Good{}, fmt.Errorf("publish error: %w", err)
	}

	return updatedGood, nil
}

func (g goodRepository) UpdatePriority(ctx context.Context, priority, id, projectID int) ([]model.Good, error) {
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
		if err != nil {
			tx.Rollback()
		}
	}()

	result := make([]model.Good, 0)

	var newPriority int

	stmt := `SELECT priority FROM goods WHERE id = $1 AND project_id = $2;`
	_ = g.db.QueryRow(stmt, id, projectID).Scan(&newPriority)

	stmt = `UPDATE goods SET priority = $1 WHERE priority = $2 AND project_id = $3;`
	_, _ = tx.ExecContext(ctx, stmt, newPriority, priority, projectID)

	stmt = `UPDATE goods SET priority = $1 WHERE id = $2 AND project_id = $3;`

	_, err = tx.ExecContext(ctx, stmt, priority, id, projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrGoodNotFound
		}

		return nil, fmt.Errorf("query error: %w", err)
	}

	stmt = `SELECT * FROM goods WHERE priority IN($1, $2)`

	rows, err := g.db.QueryContext(ctx, stmt, priority, newPriority)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var good model.Good

		if err := rows.Scan(
			&good.ID,
			&good.ProjectID,
			&good.Name,
			&good.Description,
			&good.Priority,
			&good.Removed,
			&good.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		result = append(result, good)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	if len(result) == 0 {
		return nil, model.ErrGoodNotFound
	}

	ec, err := nats.NewEncodedConn(g.nc, "json")
	if err != nil {
		return nil, err
	}

	if err := ec.Publish("good", Collection{Goods: result}); err != nil {
		return nil, fmt.Errorf("publish error: %w", err)
	}

	return result, nil
}

func (g goodRepository) Get(ctx context.Context, id int) (model.Good, error) {
	stmt := `SELECT * FROM goods WHERE id = $1`

	var good model.Good
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
			return model.Good{}, fmt.Errorf("good with id #%d %w", id, model.ErrGoodNotFound)
		}

		return model.Good{}, fmt.Errorf("query error: %w", err)
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
