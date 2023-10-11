package user

import (
	"context"
	"database/sql"
	"github.com/sonngocme/words-reminder-be/db"
)

type Storage interface {
	SignUpUser(ctx context.Context, arg db.SignUpUserParams) (sql.Result, error)
}

type Service struct {
	storage Storage
}

var _ service = (*Service)(nil)

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) SignUpUser(ctx context.Context, arg db.SignUpUserParams) (sql.Result, error) {
	return s.storage.SignUpUser(ctx, arg)
}
