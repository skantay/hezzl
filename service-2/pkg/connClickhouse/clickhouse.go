package connClickhouse

import (
	"database/sql"
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/skantay/service-2/config"
)

func ConnectClickhouse(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("clickhouse://%s:%d?username=%s&password=%s&database=%s",
		cfg.Database.Clickhouse.Host,
		cfg.Database.Clickhouse.Port,
		cfg.Database.Clickhouse.User,
		cfg.Database.Clickhouse.Password,
		cfg.Database.Clickhouse.DB,
	)
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
