package user

import (
	"context"
	"database/sql"
	"github.com/sonngocme/words-reminder-be/db"
)

type mySQLStorage struct {
	db *db.Queries
}

var _ Storage = (*mySQLStorage)(nil)

func NewStorage(db *db.Queries) Storage {
	return &mySQLStorage{db: db}
}

func (s *mySQLStorage) SignUpUser(ctx context.Context, arg db.SignUpUserParams) (sql.Result, error) {
	return s.db.SignUpUser(ctx, arg)
}