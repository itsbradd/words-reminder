package user

import (
	"context"
	"database/sql"
	"github.com/sonngocme/words-reminder-be/db"
	"golang.org/x/crypto/bcrypt"
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

func (s *Service) HashPassword(pass string) (string, error) {
	// TODO: Move Bcrypt hashing password to another pkg
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(bytes), err
}
