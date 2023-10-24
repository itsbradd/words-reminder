package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sonngocme/words-reminder-be/db"
	"go.uber.org/fx"
	"os"
)

func formatDBSource() string {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	return fmt.Sprintf("%s:%s@/%s?parseTime=true", user, pass, dbName)
}

func NewDBConn(lc fx.Lifecycle) (*db.Queries, error) {
	conn, err := sql.Open("mysql", formatDBSource())
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	queries := db.New(conn)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})
	return queries, nil
}

func ProvideDBConn() fx.Option {
	return fx.Provide(NewDBConn)
}
