package service

import (
	"database/sql"

	"github.com/nats-io/nats.go"
	"github.com/skantay/service-2/internal/domain/good/repository"
	"github.com/skantay/service-2/internal/domain/good/usecase"
)

type Service struct {
	GoodService usecase.GoodUsecase
	Nc          *nats.Conn
}

func New(db *sql.DB, nc *nats.Conn) Service {
	repoGood := repository.New(db)

	usecaseGood := usecase.New(repoGood)

	return Service{
		GoodService: usecaseGood,
		Nc:          nc,
	}
}
