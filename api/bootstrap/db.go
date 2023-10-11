package bootstrap

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sonngocme/words-reminder-be/db"
	"go.uber.org/fx"
)

func NewDBConn(lc fx.Lifecycle) (*db.Queries, error) {
	conn, err := sql.Open("mysql", "root:thisisverysecret@/words_reminder?parseTime=true")
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
