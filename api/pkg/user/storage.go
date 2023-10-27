package user

import (
	"context"
	"github.com/itsbradd/words-reminder-be/db"
)

type mySQLStorage struct {
	db *db.Queries
}

var _ Storage = (*mySQLStorage)(nil)

func NewStorage(db *db.Queries) Storage {
	return &mySQLStorage{db: db}
}

func (s *mySQLStorage) CreateUser(ctx context.Context, arg db.CreateUserParams) (int64, error) {
	return s.db.CreateUser(ctx, arg)
}

func (s *mySQLStorage) SetUserRefreshToken(ctx context.Context, arg db.SetUserRefreshTokenParams) error {
	return s.db.SetUserRefreshToken(ctx, arg)
}

func (s *mySQLStorage) GetUserByUsername(ctx context.Context, username string) (db.User, error) {
	return s.db.GetUserByUsername(ctx, username)
}
