package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/skantay/hezzl/config"
	"github.com/skantay/hezzl/internal/domain/good/model"
	"github.com/skantay/hezzl/internal/domain/good/repository"

	_ "github.com/lib/pq"
)

/*
In projects table:

1, 'test', timestamp;

In goods table:

1, 1, 'toy', 'an expensive toy', 5, 't', timestamp;

*/

// Test Create Function
func TestCreate(t *testing.T) {
	goodRepo, exec, close, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	testCases := []struct {
		name          string
		ctx           context.Context
		good          model.Good
		expectedError bool
		expectedID    int
	}{
		{
			name: "positive #1",
			ctx:  context.TODO(),
			good: model.Good{
				Name:        "Toy",
				ProjectID:   1,
				Description: "an expensive toy",
				Priority:    5,
				Removed:     false,
				CreatedAt:   time.Now(),
			},
			expectedError: false,
			expectedID:    2,
		},
		{
			name: "negative #2 - non existing project id",
			ctx:  context.TODO(),
			good: model.Good{
				ProjectID: 777,
			},
			expectedError: true,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			good, err := goodRepo.Create(tCase.ctx, tCase.good)

			if tCase.expectedError {
				if err == nil {
					t.Error("expected error")
				}
				t.Log(err)
			}

			t.Log(good)

			if good.ID != tCase.expectedID {
				t.Error("returning good did not match")
			}
		})
	}

	t.Cleanup(func() {
		down, err := os.ReadFile("migrations/setup.down.sql")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := exec(string(down)); err != nil {
			t.Fatal(err)
		}
	})
}

// -------------------

// Test Delete Function
// func TestDelete(t *testing.T) {
// 	goodRepo, exec, err := setup(t)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	testCases := []struct {
// 		name          string
// 		ctx           context.Context
// 		id            int
// 		projectID     int
// 		toDelete      bool
// 		expectedError bool
// 	}{
// 		{
// 			name:          "positive #1",
// 			ctx:           context.TODO(),
// 			id:            1,
// 			projectID:     1,
// 			toDelete:      true,
// 			expectedError: false,
// 		},
// 		{
// 			name:          "negative #2 - good does not exist",
// 			ctx:           context.TODO(),
// 			id:            2,
// 			projectID:     1,
// 			toDelete:      false,
// 			expectedError: false,
// 		},
// 		{
// 			name:          "negative #3 - projectId does not exist",
// 			ctx:           context.TODO(),
// 			id:            1,
// 			projectID:     2,
// 			toDelete:      false,
// 			expectedError: false,
// 		},
// 	}

// 	for _, tCase := range testCases {
// 		t.Run(tCase.name, func(t *testing.T) {
// 			isDeleted, err := goodRepo.Delete(tCase.ctx, tCase.id, tCase.projectID)
// 			if tCase.expectedError {
// 				if err == nil {
// 					t.Error("expected error")
// 				}
// 				t.Log(err)
// 			}

// 			if isDeleted != tCase.toDelete {
// 				t.Error("good had to be deleted, but did not")
// 			}
// 		})
// 	}

// 	t.Cleanup(func() {
// 		down, err := os.ReadFile("migrations/setup.down.sql")
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if _, err := exec(string(down)); err != nil {
// 			t.Fatal(err)
// 		}
// 		if err := goodRepo.Close(context.TODO()); err != nil {
// 			t.Fatal(err)
// 		}
// 	})
// }

// -------------------

// Test Update Function
func TestUpdate(t *testing.T) {
	goodRepo, exec, close, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	testCases := []struct {
		name          string
		ctx           context.Context
		good          model.Good
		expectedError bool
		err           error
	}{
		{
			name: "positive #1",
			ctx:  context.TODO(),
			good: model.Good{
				ID:          1,
				Name:        "New Toy",
				ProjectID:   1,
				Description: "new desc expensive toy",
				Priority:    5,
				Removed:     false,
				CreatedAt:   time.Now(),
			},
			expectedError: false,
		},
		{
			name: "positive #2 without description",
			ctx:  context.TODO(),
			good: model.Good{
				ID:          1,
				Name:        "Toy",
				ProjectID:   1,
				Description: "",
				Priority:    5,
				Removed:     false,
				CreatedAt:   time.Now(),
			},
			expectedError: false,
		},
		{
			name:          "invalid #3 non existing id",
			ctx:           context.TODO(),
			good:          model.Good{},
			expectedError: true,
			err:           model.ErrGoodNotFound,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			good, err := goodRepo.UpdateNameDesc(tCase.ctx, tCase.good)
			if tCase.expectedError {
				if err == nil {
					t.Error("expected error")
				}
				if !errors.Is(err, tCase.err) {
					t.Errorf("unexpected error: %v \n expected error: %v", err, tCase.err)
				}
				t.Log(err)
			}

			t.Log(good)

			if good.Description != tCase.good.Description && good.Name != tCase.good.Name {
				t.Error("good was not updated")
			}
		})
	}
	t.Cleanup(func() {
		down, err := os.ReadFile("migrations/setup.down.sql")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := exec(string(down)); err != nil {
			t.Fatal(err)
		}
	})
}

// -------------------

// Test Get Function
func TestGet(t *testing.T) {
	goodRepo, exec, close, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	testCases := []struct {
		name          string
		ctx           context.Context
		id            int
		good          model.Good
		expectedError bool
	}{
		{
			name: "positive #1",
			ctx:  context.Background(),
			id:   1,
			good: model.Good{
				ID:          1,
				ProjectID:   1,
				Name:        "toy",
				Description: "an expensive toy",
				Priority:    5,
				Removed:     false,
			},
			expectedError: false,
		},
		{
			name:          "negative #1 good id does not exist",
			ctx:           context.Background(),
			id:            5,
			good:          model.Good{},
			expectedError: true,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			good, err := goodRepo.Get(tCase.ctx, tCase.id)

			if tCase.expectedError {
				if err == nil {
					t.Error("expected error")
				}
				t.Log(err)
			}

			if good.Name != tCase.good.Name {
				t.Error("goods did not match")
			}
		})
	}
	t.Cleanup(func() {
		down, err := os.ReadFile("migrations/setup.down.sql")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := exec(string(down)); err != nil {
			t.Fatal(err)
		}
	})
}

// -------------------

// -------------------

// setup function - connects a database connection and creates a good repository level instance
func setup(t *testing.T) (repository.GoodRepository, func(string, ...any) (sql.Result, error), func() error, error) {
	cfg := config.Database{
		Postgres: config.Postgres{
			User:     "user",
			Password: "pass",
			Host:     "localhost",
			DBName:   "domain",
			Port:     5432,
			SSLMode:  "disable",
		},
	}

	db, err := sql.Open("postgres", fmt.Sprintf(`
		user=%s
		password=%s
		dbname=%s
		host=%s
		port=%d
		sslmode=%s`,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.SSLMode,
	))
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	down, err := os.ReadFile("migrations/setup.return errdown.sql")
	if err != nil {
		t.Fatal(err)
	}

	up, err := os.ReadFile("migrations/setup.up.sql")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := db.Exec(string(down)); err != nil {
		t.Fatal(err)
	}

	if _, err := db.Exec(string(up)); err != nil {
		t.Fatal(err)
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatal(err)
	}

	return goodRepository{db, db, nc}, db.Exec, db.Close, nil
}
